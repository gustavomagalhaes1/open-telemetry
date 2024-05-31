package weather

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type WeatherResponse struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func CelsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

func CelsiusToKelvin(celsius float64) float64 {
	return celsius + 273
}

func GetWeather(city string) (*WeatherResponse, error) {
	cityEncoded := url.QueryEscape(city)

	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=e1fece5bce574041a9f130048241703&q=%s&aqi=no", cityEncoded)

	log.Println(url)
	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}
	
	resp, err := client.Get(url)
	log.Println(resp)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var weather WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		return nil, err
	}
	return &weather, nil
}
