package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// func main() {
// 	println("staritng app..")
// 	http.HandleFunc("/", Handler)
// 	http.ListenAndServe(":8083", nil)
// }

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
		ticketDetails.Description +=
			" <p>" + typeFormdata.FormResponse.Definition.Fields[i].Title

		if typeFormdata.FormResponse.Definition.Fields[i].Type == "boolean" {
			if typeFormdata.FormResponse.Answers[i].Boolean {
				ticketDetails.Description += " : YES </p>"
			} else {
				ticketDetails.Description += " : NO </p>"
			}
		} else if typeFormdata.FormResponse.Definition.Fields[i].Type == "date" {
			ticketDetails.Description += " : " + typeFormdata.FormResponse.Answers[i].Date + "</p>"
		} else {
			ticketDetails.Description += " : " + typeFormdata.FormResponse.Answers[i].Text + "</p>"
		}

	}

	println("Data being processed are -", ticketDetails.Email,
		ticketDetails.Subject,
		ticketDetails.Name)

	ticketJSON, err := json.Marshal(&ticketDetails)
	if err != nil {
		println(err)
		return
	}

	println("ticket details - ", string(ticketJSON))

	w.Header().Set("content-type", "application/json")
	w.Write([]byte(string(ticketJSON)))
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
	Email        string       `json:"email"`
	Subject      string       `json:"subject"`
	Description  string       `json:"description"`
	Status       int          `json:"status"`
	Priority     int          `json:"priority"`
	Name         string       `json:"name"`
	Phone        string       `json:"phone"`
	CustomFields CustomFields `json:"custom_fields"`
}

type CustomFields struct {
	Weather string `json:"weather"`
}
