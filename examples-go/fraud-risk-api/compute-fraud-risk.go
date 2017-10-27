package main

/*
This API will collect Weather Data by consuming
Wunderground API and return summary for given date and city

--- INPUT ---

Historical Weather Data received from Wunderground API
for any given date and city

--- OUTPUT ---
{
	"RiskScore" : 70
	"Description" : "Stormy weather identified"
}

*/
import (
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// do some fraud analysis

}
