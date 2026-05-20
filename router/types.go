package router

type ResponseData struct {
	ShortForecast    string `json:"short_forecast"`
	Temperature      int    `json:"temperature"`
	TemperatureUnit  string `json:"temperature_unit"`
	Characterization string `json:"characterization"`
}

// NWSPointsResponse parses the metadata endpoint response
type NWSPointsResponse struct {
	Properties struct {
		ForecastURL string `json:"forecast"`
	} `json:"properties"`
}

// NWSForecastResponse parses the actual gridded forecast endpoint response
type NWSForecastResponse struct {
	Properties struct {
		Periods []struct {
			Number        int    `json:"number"`
			Temperature   int    `json:"temperature"`
			Unit          string `json:"temperatureUnit"`
			ShortForecast string `json:"shortForecast"`
		} `json:"periods"`
	} `json:"properties"`
}
