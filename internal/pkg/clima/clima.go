package clima

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const weatherAPIKey = "8a140ace6b42089b39b3450624b0f495"

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust,omitempty"`
}

type Clouds struct {
	All int `json:"all"`
}

type Sys struct {
	Type    int    `json:"type"`
	ID      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int    `json:"sunrise"`
	Sunset  int    `json:"sunset"`
}

type WeatherData struct {
	Coord      Coord     `json:"coord"`
	Weather    []Weather `json:"weather"`
	Base       string    `json:"base"`
	Main       Main      `json:"main"`
	Visibility int       `json:"visibility"`
	Wind       Wind      `json:"wind"`
	Clouds     Clouds    `json:"clouds"`
	Dt         int       `json:"dt"`
	Sys        Sys       `json:"sys"`
	Timezone   int       `json:"timezone"`
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Cod        int       `json:"cod"`
}

type Temperatura struct {
	Temp_F float64
	Temp_C float64
	Temp_K float64
}

func converterCelsiusParaFahrenheit(celsius float64) float64 {
	return (celsius * 1.8) + 32
}

func converterCelsiusParaKelvin(celsius float64) float64 {
	return celsius + 273.15
}

func BuscarTemperatura(cidade string) (Temperatura, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", cidade, weatherAPIKey)

	resp, err := http.Get(url)
	if err != nil {
		return Temperatura{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Temperatura{}, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var temperatureData map[string]interface{}
	if err := decoder.Decode(&temperatureData); err != nil {
		return Temperatura{}, err
	}

	if mainData, ok := temperatureData["main"].(map[string]interface{}); ok {
		if temperature, ok := mainData["temp"].(float64); ok {
			return Temperatura{
				Temp_F: converterCelsiusParaFahrenheit(temperature),
				Temp_C: temperature,
				Temp_K: converterCelsiusParaKelvin(temperature),
			}, nil
		}
	}

	return Temperatura{}, fmt.Errorf("error: temperature not found in response data")
}
