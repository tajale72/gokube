package router

import (
	"encoding/json"
	"fmt"
	"gokube/config"

	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func (h *Handlers) WeatherHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	w.Header().Set("Content-Type", "application/json")

	latStr := r.URL.Query().Get("lat")
	longStr := r.URL.Query().Get("long")

	h.logger.Print(
		"weather request received ",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("lat", latStr),
		zap.String("long", longStr),
		zap.String("remote_addr", r.RemoteAddr),
	)

	if err := validateCoordinates(latStr, longStr); err != nil {
		h.logger.Print(
			"weather request validation failed ",
			zap.String("lat", latStr),
			zap.String("long", longStr),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)

		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	forecastURL, err := fetchForecastURL(client, latStr, longStr)
	if err != nil {
		h.logger.Print(
			"failed fetching forecast url ",
			zap.String("lat", latStr),
			zap.String("long", longStr),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)

		writeJSONError(w, http.StatusFailedDependency, err.Error())
		return
	}

	h.logger.Print(
		"forecast url resolved",
		zap.String("lat", latStr),
		zap.String("long", longStr),
		zap.String("forecast_url", forecastURL),
	)

	forecastData, err := fetchForecastData(client, forecastURL)
	if err != nil {
		h.logger.Print(
			"failed fetching forecast data ",
			zap.String("forecast_url", forecastURL),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)

		writeJSONError(w, http.StatusFailedDependency, err.Error())
		return
	}

	todayPeriod := forecastData.Properties.Periods[0]

	response := ResponseData{
		ShortForecast:    todayPeriod.ShortForecast,
		Temperature:      todayPeriod.Temperature,
		TemperatureUnit:  todayPeriod.Unit,
		Characterization: getTemperatureCharacterization(todayPeriod.Temperature),
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Print(
			"failed encoding weather response ",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return
	}

	h.logger.Print(
		"weather request completed ",
		zap.String("lat", latStr),
		zap.String("long", longStr),
		zap.String("forecast", response.ShortForecast),
		zap.Int("temperature", response.Temperature),
		zap.String("temperature_unit", response.TemperatureUnit),
		zap.String("characterization", response.Characterization),
		zap.Duration("duration", time.Since(start)),
	)
}

func validateCoordinates(latStr, longStr string) error {
	if latStr == "" || longStr == "" {
		return fmt.Errorf("missing required parameters: lat and long")
	}

	if _, err := strconv.ParseFloat(latStr, 64); err != nil {
		return fmt.Errorf("invalid latitude format")
	}

	if _, err := strconv.ParseFloat(longStr, 64); err != nil {
		return fmt.Errorf("invalid longitude format")
	}

	return nil
}

func fetchForecastData(client *http.Client, forecastURL string) (NWSForecastResponse, error) {
	var forecastData NWSForecastResponse
	var lastErr error

	for attempt := 1; attempt <= 3; attempt++ {
		req, err := http.NewRequest(http.MethodGet, forecastURL, nil)
		if err != nil {
			return forecastData, fmt.Errorf("failed creating NWS forecast request")
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed retrieving weather details from NWS grid")

			if attempt < 3 {
				time.Sleep(time.Duration(attempt) * 200 * time.Millisecond)
			}

			continue
		}

		func() {
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				lastErr = fmt.Errorf("NWS forecast endpoint returned status %d", resp.StatusCode)
				return
			}

			if err := json.NewDecoder(resp.Body).Decode(&forecastData); err != nil {
				lastErr = fmt.Errorf("failed decoding weather periods")
				return
			}

			if len(forecastData.Properties.Periods) == 0 {
				lastErr = fmt.Errorf("no weather periods found in NWS response")
				return
			}

			lastErr = nil
		}()

		if lastErr == nil {
			return forecastData, nil
		}

		if resp != nil && !isRetryableStatus(resp.StatusCode) {
			return forecastData, lastErr
		}

		if attempt < 3 {
			time.Sleep(time.Duration(attempt) * 200 * time.Millisecond)
		}
	}

	return forecastData, lastErr
}

func getTemperatureCharacterization(temp int) string {
	switch {
	case temp <= 20:
		return "extremely cold"
	case temp <= 50:
		return "cold"
	case temp >= 100:
		return "extremely hot"
	case temp >= 85:
		return "hot"
	default:
		return "moderate"
	}
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func fetchForecastURL(client *http.Client, latStr, longStr string) (string, error) {
	pointsURL := fmt.Sprintf(
		"%s/points/%s,%s",
		config.GetString("WEATHER_API", "https://api.weather.gov"),
		latStr,
		longStr,
	)

	var lastErr error

	for attempt := 1; attempt <= 3; attempt++ {
		req, err := http.NewRequest(http.MethodGet, pointsURL, nil)
		if err != nil {
			return "", fmt.Errorf("failed creating NWS points request")
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed calling NWS points metadata endpoint")
			time.Sleep(time.Duration(attempt) * 200 * time.Millisecond)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var pointsData NWSPointsResponse

			if err := json.NewDecoder(resp.Body).Decode(&pointsData); err != nil {
				return "", fmt.Errorf("failed decoding NWS metadata payload")
			}

			if pointsData.Properties.ForecastURL == "" {
				return "", fmt.Errorf("missing forecast URL in NWS metadata response")
			}

			return pointsData.Properties.ForecastURL, nil
		}

		if !isRetryableStatus(resp.StatusCode) {
			return "", fmt.Errorf("NWS points endpoint returned status %d", resp.StatusCode)
		}

		lastErr = fmt.Errorf("NWS points endpoint returned status %d", resp.StatusCode)

		time.Sleep(time.Duration(attempt) * 200 * time.Millisecond)
	}

	return "", lastErr
}

func isRetryableStatus(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests ||
		statusCode == http.StatusInternalServerError ||
		statusCode == http.StatusBadGateway ||
		statusCode == http.StatusServiceUnavailable ||
		statusCode == http.StatusGatewayTimeout
}
