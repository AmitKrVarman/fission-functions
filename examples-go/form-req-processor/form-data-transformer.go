package main

/*
This API serves a web hook for TYPE FORM
and transform the input JSON to Zen Desk equivalent
request along side Weather API Input Data

*/

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	TV_CLAIM          = "TV claim"
	STORM_SURGE_CLAIM = "Storm surge claim"
	BOOLEAN_FEILD     = "boolean"
	TEXT_FEILD        = "text"
	DATE_FEILD        = "date"
	PARA_START        = "<p>"
	PARA_END          = "</p>"
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
	reader := ioutil.NopCloser(bytes.NewBuffer(buf))

	if r.Body == nil {
		http.Error(w, "Please send a valid JSON", 400)
		return
	}
	err := json.NewDecoder(reader).Decode(&typeFormdata)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	transformedData := transformData(typeFormdata)

	transformedDataJSON, err := json.Marshal(&transformedData)
	if err != nil {
		println(err)
		return
	}

	println("ticket details - ", string(transformedDataJSON))

	w.Header().Set("content-type", "application/json")
	w.Write([]byte(string(transformedDataJSON)))
}

//transform Type form data
func transformData(typeFormdata TypeFormData) TranformedData {
	//populate Fresh Desk Specific Struct
	var ticketDetails TicketDetails
	var transformedData TranformedData
	ticketDetails.Ticket.Subject = typeFormdata.FormResponse.Definition.Title
	ticketDetails.Ticket.Email = typeFormdata.FormResponse.Hidden.Email
	ticketDetails.Ticket.Phone = typeFormdata.FormResponse.Hidden.Phone

	//default data
	ticketDetails.Ticket.Status = "new"
	ticketDetails.Ticket.Priority = "normal"
	ticketDetails.Ticket.Type = "incident"
	ticketBody := ""

	//populate Descripton
	for i, field := range typeFormdata.FormResponse.Definition.Fields {
		answer := typeFormdata.FormResponse.Answers[i]
		ticketBody += PARA_START + field.Title

		switch field.Type {
		case BOOLEAN_FEILD:
			ticketBody += strconv.FormatBool(answer.Boolean) + PARA_END
		case DATE_FEILD:
			ticketBody += " : " + answer.Date + PARA_END
		default:
			ticketBody += " : " + answer.Text + PARA_END
		}

		ticketDetails.Ticket.Comment.Body = ticketBody

		//for Storm Surge Claims - work around as ID is not returned by TYPE form
		if strings.Compare(field.Title, STORM_SURGE_CLAIM) >= 0 {
			if strings.Contains(field.Title, "Where did the incident happen? (City/town name)") {
				transformedData.WeatherAPIInput.City = answer.Text
			}
			if strings.Contains(field.Title, "When did the incident happen?") {
				transformedData.WeatherAPIInput.Date = strings.Replace(answer.Date, "-", "", 2)
			}
		}

	}
	//update ticket specific data
	transformedData.TicketDetails = ticketDetails

	return transformedData
}

//Model for TYPE FORM

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

//Model for ZEN DESK

type TicketDetails struct {
	Ticket struct {
		Type     string `json:"type"`
		Subject  string `json:"subject"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Priority string `json:"priority"`
		Status   string `json:"status"`
		Comment  struct {
			Body string `json:"body"`
		} `json:"comment"`
		CustomFields []struct {
			ID    int    `json:"id"`
			Value string `json:"value"`
		} `json:"custom_fields"`
	} `json:"ticket"`
}

type WeatherAPIInput struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Date    string `json:"date"`
}

type TranformedData struct {
	TicketDetails   TicketDetails   `json:"ticketDetails"`
	WeatherAPIInput WeatherAPIInput `json:"weatherAPIInput"`
}
