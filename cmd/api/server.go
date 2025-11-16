package main

import (
	"fmt"
	"log"
	"net/http"
	"apiproject/internal/api/middlewares"
	"crypto/tls"
	"path/filepath"
)

func main() {
	port := ":8080"

	cert := filepath.Join("cmd", "config", "cert.pem")
	key := filepath.Join("cmd", "config", "key.pem")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "Hello Server!")
	})

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	server := &http.Server{
		Addr: port,
		Handler: middlewares.SecurityHeaders(mux),
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server Listening on port:", port)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("error starting server", err)
	}
}
