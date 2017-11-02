package main

/*
This API would echo same input JSON

INPUT == OUTPUT

*/
import (
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Message struct {
}

func Handler(w http.ResponseWriter, r *http.Request) {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	secret, err := clientset.Core().Secrets("default").Get("echo", metav1.GetOptions{})
	println(string(secret.Data["content"]))

	if err != nil {
		println("Error getting secret %s: %v", "echo", err)
		return
	}

	//read the content
	tokenContent := string(secret.Data["content"])

	//Read Post body
	// responseData, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(w, err.Error(), 400)
	// 	return
	// }
	// defer r.Body.Close()

	// write the content back in repsonse
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(tokenContent))
}

// func main() {
// 	println("staritng app..")
// 	http.HandleFunc("/", Handler)
// 	http.ListenAndServe(":8085", nil)
// }
