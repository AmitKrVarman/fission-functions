package main

/*
This API would echo same input JSON

INPUT == OUTPUT

*/
import (
	"io/ioutil"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	//Read Post body
	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer r.Body.Close()

	//Write Post body
	w.Header().Set("content-type", "application/json")
	w.Write(responseData)
}

// func main() {
// 	println("staritng app..")
// 	http.HandleFunc("/", Handler)
// 	http.ListenAndServe(":8085", nil)
// }
