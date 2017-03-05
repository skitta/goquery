package openweather

import (
	"fmt"
	"net/http"

	simplejson "github.com/bitly/go-simplejson"
)

const (
	appid   string = "19e52752071f563f3d7add251cfe4fe8"
	baseURL string = "http://api.openweathermap.org/data/2.5/forecast/daily?"
)

// GetDaily return a report string fetch from openweather api
func GetDaily(city string) (info, des string, err error) {
	queryURL := baseURL + "APPID=" + appid + "&q=" + city + "&units=metric&cnt=1&mode=json"

	res, err := http.Get(queryURL)
	if err != nil {
		return
	}

	defer res.Body.Close()
	if res.StatusCode == 200 {
		data, _ := simplejson.NewFromReader(res.Body)
		weatherJSON := data.Get("list").GetIndex(0).Get("weather").GetIndex(0)
		tempJSON := data.Get("list").GetIndex(0).Get("temp")

		info, _ = weatherJSON.Get("main").String()
		detail, _ := weatherJSON.Get("description").String()
		maxTemp, _ := tempJSON.Get("max").Float64()
		minTemp, _ := tempJSON.Get("min").Float64()

		des = fmt.Sprintf("%s today with %2.2f°C/%2.2f°C", detail, maxTemp, minTemp)
	}

	return
}
