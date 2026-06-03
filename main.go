package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/users/register", userRegisterHandler)
	mux.HandleFunc("/health-check", helthCheckHandler)

	server := http.Server{Addr: ":8080", Handler: mux}
	server.ListenAndServe()
}

func userRegisterHandler(resWriter http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		fmt.Fprintf(resWriter, "Invalid Method")
	}
}

func helthCheckHandler(resWriter http.ResponseWriter, _ *http.Request){
	fmt.Fprintf(resWriter, "Everything is Good")
}