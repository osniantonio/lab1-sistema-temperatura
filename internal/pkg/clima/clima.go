package clima

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const weatherAPIKey = "8a140ace6b42089b39b3450624b0f495"

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
