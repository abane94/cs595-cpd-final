package main

import (
	"net/http"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := "Hello world"

	w.Write([]byte(message))
}

func main() {
	http.HandleFunc("/", sayHello)
	if err := http.ListenAndServe(":5000", nil); err != nil {
		panic(err)
	}
}
