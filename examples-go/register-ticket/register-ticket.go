package main

/*
This API would accept JSON string as POST body and
create a Ticket in Zen Desk/Fresh Desk
*/

import (
	"log"
	"net/http"
	"os"
)

//Default values , this can be overridden by setting ENV variables
var (
	endPoint    = "https://helmion.zendesk.com/api/v2/users.json"
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

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	println("request status for ticket creation :" + resp.Status)

	defer resp.Body.Close()

	w.Header().Set("content-type", "application/json")
	w.Write([]byte(resp.Status))
}

// func main() {
// 	println("staritng app..")
// 	http.HandleFunc("/", Handler)
// 	http.ListenAndServe(":8085", nil)
// }

func init() {
	getEnvDatabaseConfig()
}
