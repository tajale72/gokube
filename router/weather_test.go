package router

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}

func withWeatherBaseURL(t *testing.T, baseURL string) {
	t.Helper()

	oldBaseURL := weatherBaseURL
	weatherBaseURL = baseURL

	t.Cleanup(func() {
		weatherBaseURL = oldBaseURL
	})
}

func newWeatherTestServer(t *testing.T, forecastStatus int, forecastBody string) *httptest.Server {
	t.Helper()

	var server *httptest.Server

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/points/32.7767,-96.7970":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"properties": {
					"forecast": "` + server.URL + `/forecast"
				}
			}`))

		case "/forecast":
			w.WriteHeader(forecastStatus)
			_, _ = w.Write([]byte(forecastBody))

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	t.Cleanup(server.Close)

	return server
}

func TestValidateCoordinates(t *testing.T) {
	tests := []struct {
		name    string
		lat     string
		long    string
		wantErr bool
	}{
		{"valid coordinates", "32.7767", "-96.7970", false},
		{"missing latitude", "", "-96.7970", true},
		{"missing longitude", "32.7767", "", true},
		{"invalid latitude", "abc", "-96.7970", true},
		{"invalid longitude", "32.7767", "xyz", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCoordinates(tt.lat, tt.long)

			if tt.wantErr && err == nil {
				t.Fatal("expected error but got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil error but got %v", err)
			}
		})
	}
}

func TestGetTemperatureCharacterization(t *testing.T) {
	tests := []struct {
		name     string
		temp     int
		expected string
	}{
		{"extremely cold below boundary", 10, "extremely cold"},
		{"extremely cold boundary", 20, "extremely cold"},
		{"cold lower boundary", 21, "cold"},
		{"cold upper boundary", 50, "cold"},
		{"moderate lower boundary", 51, "moderate"},
		{"moderate upper boundary", 84, "moderate"},
		{"hot lower boundary", 85, "hot"},
		{"hot upper boundary", 99, "hot"},
		{"extremely hot boundary", 100, "extremely hot"},
		{"extremely hot above boundary", 105, "extremely hot"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getTemperatureCharacterization(tt.temp)

			if got != tt.expected {
				t.Fatalf("expected %q but got %q", tt.expected, got)
			}
		})
	}
}

func TestWriteJSONError(t *testing.T) {
	rr := httptest.NewRecorder()

	writeJSONError(rr, http.StatusBadRequest, "invalid request")

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d but got %d", http.StatusBadRequest, rr.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed decoding response: %v", err)
	}

	if response["error"] != "invalid request" {
		t.Fatalf("expected error message %q but got %q", "invalid request", response["error"])
	}
}

func TestFetchForecastURL(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantURL    string
		wantErr    bool
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			body:       `{"properties":{"forecast":"http://example.com/forecast"}}`,
			wantURL:    "http://example.com/forecast",
			wantErr:    false,
		},
		{
			name:       "non 200 response",
			statusCode: http.StatusInternalServerError,
			body:       `{}`,
			wantErr:    true,
		},
		{
			name:       "invalid json",
			statusCode: http.StatusOK,
			body:       `invalid-json`,
			wantErr:    true,
		},
		{
			name:       "missing forecast url",
			statusCode: http.StatusOK,
			body:       `{"properties":{}}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.body))
			}))
			t.Cleanup(server.Close)

			withWeatherBaseURL(t, server.URL)

			got, err := fetchForecastURL(&http.Client{}, "32.7767", "-96.7970")

			if tt.wantErr && err == nil {
				t.Fatal("expected error but got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil error but got %v", err)
			}

			if got != tt.wantURL {
				t.Fatalf("expected URL %q but got %q", tt.wantURL, got)
			}
		})
	}
}

func TestFetchForecastURL_InvalidRequestURL(t *testing.T) {
	withWeatherBaseURL(t, "://bad-url")

	_, err := fetchForecastURL(&http.Client{}, "32.7767", "-96.7970")

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err.Error() != "failed creating NWS points request" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFetchForecastURL_ClientDoError(t *testing.T) {
	withWeatherBaseURL(t, "http://127.0.0.1:1")

	_, err := fetchForecastURL(&http.Client{}, "32.7767", "-96.7970")

	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestFetchForecastData(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantErr    bool
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			body: `{
				"properties": {
					"periods": [
						{
							"shortForecast": "Sunny",
							"temperature": 90,
							"temperatureUnit": "F"
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name:       "non 200 response",
			statusCode: http.StatusInternalServerError,
			body:       `{}`,
			wantErr:    true,
		},
		{
			name:       "invalid json",
			statusCode: http.StatusOK,
			body:       `invalid-json`,
			wantErr:    true,
		},
		{
			name:       "empty periods",
			statusCode: http.StatusOK,
			body:       `{"properties":{"periods":[]}}`,
			wantErr:    true,
		},
		{
			name:       "missing periods",
			statusCode: http.StatusOK,
			body:       `{"properties":{}}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.body))
			}))
			t.Cleanup(server.Close)

			got, err := fetchForecastData(&http.Client{}, server.URL)

			if tt.wantErr && err == nil {
				t.Fatal("expected error but got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil error but got %v", err)
			}

			if !tt.wantErr {
				period := got.Properties.Periods[0]

				if period.ShortForecast != "Sunny" {
					t.Fatalf("expected Sunny but got %s", period.ShortForecast)
				}

				if period.Temperature != 90 {
					t.Fatalf("expected 90 but got %d", period.Temperature)
				}

				if period.Unit != "F" {
					t.Fatalf("expected F but got %s", period.Unit)
				}
			}
		})
	}
}

func TestFetchForecastData_InvalidURL(t *testing.T) {
	_, err := fetchForecastData(&http.Client{}, "://bad-url")

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err.Error() != "failed creating NWS forecast request" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFetchForecastData_ClientDoError(t *testing.T) {
	_, err := fetchForecastData(&http.Client{}, "http://127.0.0.1:1")

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	if err.Error() != "failed retrieving weather details from NWS grid" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWeatherHandler(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		pointsStatus   int
		pointsBody     string
		forecastStatus int
		forecastBody   string
		expectedStatus int
	}{
		{
			name:           "invalid coordinates",
			url:            "/weather?lat=abc&long=xyz",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "forecast url error",
			url:            "/weather?lat=32.7767&long=-96.7970",
			pointsStatus:   http.StatusInternalServerError,
			pointsBody:     `{}`,
			expectedStatus: http.StatusFailedDependency,
		},
		{
			name:           "forecast data error",
			url:            "/weather?lat=32.7767&long=-96.7970",
			pointsStatus:   http.StatusOK,
			forecastStatus: http.StatusInternalServerError,
			forecastBody:   `{}`,
			expectedStatus: http.StatusFailedDependency,
		},
		{
			name:           "success",
			url:            "/weather?lat=32.7767&long=-96.7970",
			pointsStatus:   http.StatusOK,
			forecastStatus: http.StatusOK,
			forecastBody: `{
				"properties": {
					"periods": [
						{
							"shortForecast": "Sunny",
							"temperature": 90,
							"temperatureUnit": "F"
						}
					]
				}
			}`,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var server *httptest.Server

			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/points/32.7767,-96.7970":
					w.WriteHeader(tt.pointsStatus)

					if tt.pointsBody != "" {
						_, _ = w.Write([]byte(tt.pointsBody))
						return
					}

					_, _ = w.Write([]byte(`{
						"properties": {
							"forecast": "` + server.URL + `/forecast"
						}
					}`))

				case "/forecast":
					w.WriteHeader(tt.forecastStatus)
					_, _ = w.Write([]byte(tt.forecastBody))

				default:
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			t.Cleanup(server.Close)

			withWeatherBaseURL(t, server.URL)

			h := &Handlers{
				logger: testLogger(),
			}

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rr := httptest.NewRecorder()

			h.WeatherHandler(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected status %d but got %d. body=%s", tt.expectedStatus, rr.Code, rr.Body.String())
			}
		})
	}
}

func TestWeatherHandler_EncodeError(t *testing.T) {
	server := newWeatherTestServer(
		t,
		http.StatusOK,
		`{
			"properties": {
				"periods": [
					{
						"shortForecast": "Sunny",
						"temperature": 90,
						"temperatureUnit": "F"
					}
				]
			}
		}`,
	)

	withWeatherBaseURL(t, server.URL)

	h := &Handlers{
		logger: testLogger(),
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/weather?lat=32.7767&long=-96.7970",
		nil,
	)

	h.WeatherHandler(&failingResponseWriter{}, req)
}

type failingResponseWriter struct{}

func (f *failingResponseWriter) Header() http.Header {
	return http.Header{}
}

func (f *failingResponseWriter) Write([]byte) (int, error) {
	return 0, errors.New("forced write failure")
}

func (f *failingResponseWriter) WriteHeader(int) {}
