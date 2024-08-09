package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, trend setters!")
}

func main() {
	http.HandleFunc("/", helloHandler)
	http.ListenAndServe(":8082", nil)
}
