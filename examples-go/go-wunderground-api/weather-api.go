package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const (
	ENV_WU_API_URL    = "WU_API_URL"
	ENV_WU_API_AC_URL = "WU_API_AC_URL"
	ENV_WU_API_KEY    = "WU_API_KEY"
)

//Default values , this can be removed and set as ENV variables
var (
	apiURL          = "http://api.wunderground.com/api"
	autocompleteURL = "http://autocomplete.wunderground.com"
	apiKey          = "5045721b27fff489"
)

func init() {
	getEnvConfig()
}

type InputData struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Date    string `json:"date"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	var inputData InputData

	if r.Body == nil {
		http.Error(w, "Please send a valid JSON", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&inputData)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	println("Query Weather Data for", inputData.City, inputData.Country, inputData.Date)

	//use Wundergroud AutoComplete API to get unique city link
	link, err := getCityUniqueLink(inputData.City, inputData.Country)

	//use Wundergroud API to retrieve Historical Data
	weatherDataJSON, err := GetWeatherConditions(link, inputData.Date)

	w.Write([]byte(weatherDataJSON))

}

func getCityUniqueLink(city string, country string) (string, error) {

	//Query AutoCmplete API
	autocompleteURL = autocompleteURL + "/aq?query=" + url.QueryEscape(city) + "&c=" + url.QueryEscape(country)
	println("autocompleteURL : ", autocompleteURL)
	acResp, err := http.Get(autocompleteURL)
	if err != nil {
		return "", err
	}
	defer acResp.Body.Close()

	acObj := autocomplete{}
	err = json.NewDecoder(acResp.Body).Decode(&acObj)

	if err != nil {
		return "", err
	}

	if len(acObj.Results) == 0 {
		println("No result found for ", city)
		return "", errors.New("No results found")
	}

	link := acObj.Results[0].Link
	println("link --- ", link)

	return link, err
}

// GetLocalConditions returns the conditions at the first place listed by
// Wunderground's autocomplete API.
func GetWeatherConditions(link string, dateString string) (string, error) {

	//form API URL for Historical Data
	historicalDataURL := apiURL + "/" + apiKey + "/history_" + url.QueryEscape(dateString) + link + ".json"

	println("historicalDataURL being queried : ", historicalDataURL)
	repsonse, err := http.Get(historicalDataURL)
	if err != nil {
		return "", err
	}
	defer repsonse.Body.Close()

	weatherData, err := ioutil.ReadAll(repsonse.Body)

	println("weather details JSON ", string(weatherData))

	return string(weatherData), nil
}

type autocomplete struct {
	Results []autocompleteResult `json:"RESULTS"`
}

type autocompleteResult struct {
	Link string `json:"l"`
}

type displaylocation struct {
	Name string `json:"full"`
}

// func main() {
// 	println("staritng app..")
// 	http.HandleFunc("/", Handler)
// 	http.ListenAndServe(":8084", nil)
// }

func getEnvConfig() {
	println("[CONFIG] Reading Env variables")
	apiURL := os.Getenv(ENV_WU_API_URL)
	autocompleteURL := os.Getenv(ENV_WU_API_AC_URL)
	apiKey := os.Getenv(ENV_WU_API_KEY)

	if len(apiURL) > 0 {
		println("[CONFIG] Missing Wundergroud API URL - Set ENV  ", ENV_WU_API_URL)
	}

	if len(autocompleteURL) > 0 {
		println("[CONFIG] Missing Wundergroud Auto Complete API URL - Set ENV  ", ENV_WU_API_AC_URL)
	}

	if len(apiKey) > 0 {
		println("[CONFIG] Missing Wundergroud Auto Complete API URL - Set ENV  ", ENV_WU_API_KEY)
	}
}
