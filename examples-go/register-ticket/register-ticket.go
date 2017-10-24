package main

import (
	"log"
	"net/http"
	"os"
)

//Default values
var (
	ENV_API_END_POINT string = "https://helmion.freshdesk.com/api/v2/tickets"
	ENV_API_KEY       string = "MbzSRhpf0gBjQEdDLirp"
	ENV_API_PASSWORD  string = "X"
)

type TicketDetails struct {
	Email       string `json:"email"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	Priority    int    `json:"priority"`
	Name        string `json:"name"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	println("Opening HTTP Request...", ENV_API_END_POINT)
	req, err := http.NewRequest("POST", ENV_API_END_POINT, r.Body)
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(ENV_API_KEY, ENV_API_PASSWORD)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	println("reuqest staus" + resp.Status)

	defer resp.Body.Close()

	w.Header().Set("content-type", "application/json")
	w.Write([]byte(resp.Status))
}

func main() {
	println("staritng app..")
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8085", nil)
}

func init() {
	getEnvDatabaseConfig()
}

func getEnvDatabaseConfig() {
	log.Print("[CONFIG] Reading Env variables")
	endPoint := os.Getenv(ENV_API_END_POINT)
	apiKey := os.Getenv(ENV_API_KEY)
	apiPassword := os.Getenv(ENV_API_PASSWORD)

	if len(endPoint) > 0 {
		ENV_API_END_POINT = endPoint
	}
	if len(apiKey) > 0 {
		ENV_API_KEY = apiKey
	}
	if len(apiPassword) > 0 {
		ENV_API_PASSWORD = apiPassword
	}

}
