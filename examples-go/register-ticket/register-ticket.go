package main

/*
This API would accept JSON string as POST body and
create a Ticket in Zen Desk/Fresh Desk

INPUT - Zen Desk Create Ticket compliant JSON

OUTPUT - Ticket Meta Data JSON from Response Object
{
	"id": 133382282992,
	"ticket_id": 39,
	"created_at": "2017-10-25T18:32:55Z",
	"author_id": 115428050612,
	"metadata": {
		"system": {
		"ip_address": "2.122.25.146",
		"location": "Solihull, M2, United Kingdom",
		"latitude": 52.41669999999999,
		"longitude": -1.783299999999997
	},
	"custom": {}
	}
}
*/
import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

//Default values , this can be overridden by setting ENV variables
var (
	endPoint    = "https://landg.zendesk.com/api/v2/tickets.json"
	apiKey      = "aguaparachocolate@gmail.com/token"
	apiPassword = "3OdKmqWsCL8uKR0QSH0jqmdVWx44CagRnRS01mnL"
)

func getEnvDatabaseConfig() {
	log.Print("[CONFIG] Reading Env variables")
	endPointFromENV := os.Getenv("ENV_HELPDESK_API_EP")
	apiKeyFromENV := os.Getenv("ENV_HELPDESK_API_KEY")
	apiPasswordFromENV := os.Getenv("ENV_HELPDESK_API_PASSWORD")

	if len(endPointFromENV) > 0 {
		log.Print("[CONFIG] Setting Env variables", endPointFromENV)
		endPoint = endPointFromENV
	}
	if len(apiKeyFromENV) > 0 {
		apiKey = apiKeyFromENV
	}
	if len(apiPasswordFromENV) > 0 {
		apiPassword = apiPasswordFromENV
	}

}

func Handler(w http.ResponseWriter, r *http.Request) {

	println("Opening HTTP Request...", endPoint)
	req, err := http.NewRequest("POST", endPoint, r.Body)
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(apiKey, apiPassword)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	client := &http.Client{}

	zendeskAPIResp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	println("request status for ticket creation :" + zendeskAPIResp.Status)
	var ticketResponse TicketResponse
	err = json.NewDecoder(zendeskAPIResp.Body).Decode(&ticketResponse)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer zendeskAPIResp.Body.Close()

	//marshal response to JSON
	ticketAuditData := ticketResponse.Audit
	ticketResponseJSON, err := json.Marshal(&ticketAuditData)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write([]byte(ticketResponseJSON))
}

// func main() {
// 	println("staritng app..")
// 	http.HandleFunc("/", Handler)
// 	http.ListenAndServe(":8085", nil)
// }

func init() {
	getEnvDatabaseConfig()
}

type TicketResponse struct {
	Ticket struct {
		URL        string      `json:"url"`
		ID         int         `json:"id"`
		ExternalID interface{} `json:"external_id"`

		CreatedAt    time.Time   `json:"created_at"`
		UpdatedAt    time.Time   `json:"updated_at"`
		DueAt        interface{} `json:"due_at"`
		TicketFormID int64       `json:"ticket_form_id"`
	} `json:"ticket"`
	Audit struct {
		ID        int64     `json:"id"`
		TicketID  int       `json:"ticket_id"`
		CreatedAt time.Time `json:"created_at"`
		AuthorID  int64     `json:"author_id"`
		Metadata  struct {
			System struct {
				IPAddress string  `json:"ip_address"`
				Location  string  `json:"location"`
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"system"`
			Custom struct {
			} `json:"custom"`
		} `json:"metadata"`
	} `json:"audit"`
}
