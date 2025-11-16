package main

import (
	"fmt"
	"log"
	"net/http"
	"apiproject/internal/api/middlewares"
)

func main() {
	string = ":8080"

	cert := ""
	key := ""

	mux = http.NewServeMux()

	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "Hello Server!")
	})

	tlsConfig := tls.Config{
		MinVersion: tls.VersionTLS12
	}

	server := &http.Server{
		Addr: port,
		Handler: middlewares,
		TLSConfig: tlsConfig
	}

	fmt.Println("Server Listening on port:", port)
	err := server.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("error starting server", err)
	}
}
