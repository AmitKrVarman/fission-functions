package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	//Marhsal TYPE FORM DATA to TypeFormData struct
	var typeFormdata TypeFormData
	buf, _ := ioutil.ReadAll(r.Body)
	readerOne := ioutil.NopCloser(bytes.NewBuffer(buf))

	if r.Body == nil {
		http.Error(w, "Please send a valid JSON", 400)
		return
	}
	err := json.NewDecoder(readerOne).Decode(&typeFormdata)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	println(typeFormdata.EventID,
		typeFormdata.FormResponse.Hidden.Email,
		typeFormdata.FormResponse.Definition.Title)

	//Populate Weather API calls
	var weatherAPIInput WeatherAPIInput
	weatherAPIInput.Country = "GB" // defaulted to UK

	//populate Fresh Desk Specific Struct
	var ticketDetails TicketDetails
	ticketDetails.Email = typeFormdata.FormResponse.Hidden.Email
	ticketDetails.Subject = typeFormdata.FormResponse.Definition.Title
	ticketDetails.Name = typeFormdata.FormResponse.Hidden.Name
	ticketDetails.Phone = typeFormdata.FormResponse.Hidden.Phone
	//ticketDetails.Policy = typeFormdata.FormResponse.Hidden.Policy

	ticketDetails.Status = 2   //will be updated based on weather data
	ticketDetails.Priority = 1 //will be updated based on weather data
	ticketDetails.Description = ""

	//populate Descripton
	for i := 0; i < len(typeFormdata.FormResponse.Definition.Fields); i++ {
		ticketDetails.Description = ticketDetails.Description +
			" <p>" + typeFormdata.FormResponse.Definition.Fields[i].Title

		if typeFormdata.FormResponse.Definition.Fields[i].Type == "boolean" {
			if typeFormdata.FormResponse.Answers[i].Boolean {
				ticketDetails.Description = ticketDetails.Description + " : YES </p>"
			} else {
				ticketDetails.Description = ticketDetails.Description + " : NO </p>"
			}
		}
		if typeFormdata.FormResponse.Definition.Fields[i].Type == "date" {
			ticketDetails.Description = ticketDetails.Description + " : " + typeFormdata.FormResponse.Answers[i].Date + "</p>"
		}

		if strings.Contains(typeFormdata.FormResponse.Definition.Fields[i].Title,
			"Where did the incident happen? (City/town name)") {
			weatherAPIInput.City = typeFormdata.FormResponse.Answers[i].Text
		}
		if strings.Contains(typeFormdata.FormResponse.Definition.Fields[i].Title,
			"When did the incident happen?") {
			weatherAPIInput.Date = strings.Replace(typeFormdata.FormResponse.Answers[i].Date, "-", "", 2)
			println(weatherAPIInput.Date)
		}
	}

	println("Data being processed are -", ticketDetails.Email,
		ticketDetails.Subject,
		ticketDetails.Name,
		weatherAPIInput.City,
		weatherAPIInput.Date)

	weatherJSON, err := json.Marshal(&weatherAPIInput)
	if err != nil {
		println(err)
		return
	}

	//Call Wundergroud Fission function to Get Wind Speed
	weatherAPIResp, err := http.Post("http://fission.landg.madeden.net/get-weather-data",
		"application/json", bytes.NewBuffer(weatherJSON))

	println("response Status for weatherAPI :", weatherAPIResp.Status)
	var historicalData HistoricalData
	err = json.NewDecoder(weatherAPIResp.Body).Decode(&historicalData)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	// historicalDataJSON, err := json.Marshal(&historicalData)
	// if err != nil {
	// 	println(err)
	// 	return
	// }
	// println("historicalDataJSON - ", string(historicalDataJSON))

	//update ticket based on wind speed
	windSpeed := historicalData.History.DailySummary[0].Maxwspdm
	ticketDetails.CustomFields.Weather = "WIND Speed (mph): " + windSpeed
	n, _ := strconv.ParseFloat(windSpeed, 10)
	if n >= 20 {
		ticketDetails.Priority = 4
		ticketDetails.Status = 2
		ticketDetails.Subject = ticketDetails.Subject + " :  ACCEPTED"
	} else {
		ticketDetails.Priority = 1
		ticketDetails.Status = 4
		ticketDetails.Subject = ticketDetails.Subject + " :  DECLINED"
	}

	ticketJSON, err := json.Marshal(&ticketDetails)
	if err != nil {
		println(err)
		return
	}
	println("ticket details - ", string(ticketJSON))

	//Call Create Ticket API
	freshDeskResp, err := http.Post("http://fission.landg.madeden.net/register-ticket",
		"application/json", bytes.NewBuffer(ticketJSON))

	println("response Status for registration request:", freshDeskResp.Status)

	w.Header().Set("content-type", "application/json")
	w.Write([]byte(freshDeskResp.Status))

}

// func main() {
// 	println("staritng app..")
// 	http.HandleFunc("/", Handler)
// 	http.ListenAndServe(":8083", nil)
// }

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

//Model for FRESH DESK

type TypeFormData struct {
	EventID      string       `json:"event_id"`
	EventType    string       `json:"event_type"`
	FormResponse FormResponse `json:"form_response"`
}

type FormResponse struct {
	FormID      string     `json:"form_id"`
	Token       string     `json:"token"`
	SubmittedAt time.Time  `json:"submitted_at"`
	Hidden      Hidden     `json:"hidden"`
	Definition  Definition `json:"definition"`
	Answers     []Answers  `json:"answers"`
}

type Hidden struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Policy string `json:"policy"`
}

type Definition struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Fields []Fields `json:"fields"`
}

type Fields struct {
	ID                      string `json:"id"`
	Title                   string `json:"title"`
	Type                    string `json:"type"`
	AllowMultipleSelections bool   `json:"allow_multiple_selections"`
	AllowOtherChoice        bool   `json:"allow_other_choice"`
}

type Answers struct {
	Type    string `json:"type"`
	Text    string `json:"text,omitempty"`
	Field   Field  `json:"field"`
	FileURL string `json:"file_url,omitempty"`
	Date    string `json:"date,omitempty"`
	Boolean bool   `json:"boolean,omitempty"`
	Number  int    `json:"number,omitempty"`
}

type Field struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type TicketDetails struct {
	Email       string `json:"email"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	Priority    int    `json:"priority"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	//Policy       string       `json:"policy"`
	CustomFields CustomFields `json:"custom_fields"`
}

type CustomFields struct {
	Weather string `json:"weather"`
}
