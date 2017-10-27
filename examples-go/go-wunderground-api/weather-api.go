package main

/*
This API will collect Weather Data by consuming
Wunderground API and return summary for given date and city

---INPUT---
{
  "city": "birmingham",
  "country": "",
  "date": "20170101"
}
---OUTPUT---
{
"fog": "0",
"rain": "1",
"maxtempm": "7",
"mintempm": "0",
"tornado": "0",
"maxpressurem": "1025",
"minpressurem": "1014",
"maxwspdm": "28",
"minwspdm": "7"
}

*/

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
)

const (
	ENV_WU_API_URL    = "WU_API_URL"
	ENV_WU_API_AC_URL = "WU_API_AC_URL"
	ENV_WU_API_KEY    = "WU_API_KEY"
)

//Default values , this can be overridden by setting ENV variables
var (
	apiURL          = "http://api.wunderground.com/api"
	autocompleteURL = "http://autocomplete.wunderground.com"
	apiKey          = "5045721b27fff489"
)

func init() {
	getEnvConfig()
}

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
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write([]byte(weatherDataJSON))

}

func getCityUniqueLink(city string, country string) (string, error) {

	//Query AutoCmplete API
	if len(country) == 0 {
		country = "GB" //default country to United kingdom if blank
	}
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

// GetLocalConditions returns weather summary for given date
func GetWeatherConditions(link string, dateString string) (string, error) {

	//form API URL for Historical Data
	historicalDataURL := apiURL + "/" + apiKey + "/history_" + url.QueryEscape(dateString) + link + ".json"

	println("historicalDataURL being queried : ", historicalDataURL)
	repsonse, err := http.Get(historicalDataURL)
	if err != nil {
		return "", err
	}
	defer repsonse.Body.Close()
	println("response Status for weatherAPI :", repsonse.Status)

	var historicalData HistoricalData
	err = json.NewDecoder(repsonse.Body).Decode(&historicalData)
	if err != nil {
		println(err)
		return "", err
	}

	//get summary for the given Date
	//dailySummary := historicalData.History.DailySummary[0]

	//marshal to JSON
	historicalDataJSON, err := json.Marshal(historicalData)
	if err != nil {
		println(err)
		return "", err
	}

	return string(historicalDataJSON), nil
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

//Model for WeatherAPI

type WeatherAPIInput struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Date    string `json:"date"`
}

type HistoricalData struct {
	Response Response `json:"response"`
	History  History  `json:"history"`
}

type Response struct {
	Version string `json:"version"`
}

type History struct {
	DailySummary []DailySummary `json:"dailysummary"`
}

type DailySummary struct {
	Fog          string `json:"fog"`
	Rain         string `json:"rain"`
	Maxtempm     string `json:"maxtempm"`
	Mintempm     string `json:"mintempm"`
	Tornado      string `json:"tornado"`
	Maxpressurem string `json:"maxpressurem"`
	Minpressurem string `json:"minpressurem"`
	Maxwspdm     string `json:"maxwspdm"`
	Minwspdm     string `json:"minwspdm"`
}

// func main() {
// 	println("staritng app..")
// 	http.HandleFunc("/", Handler)
// 	http.ListenAndServe(":8084", nil)
// }
