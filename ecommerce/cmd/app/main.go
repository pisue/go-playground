package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Ecommerce Service Started on :8081")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Ecommerce Service!"))
	})
	http.ListenAndServe(":8081", nil)
}
