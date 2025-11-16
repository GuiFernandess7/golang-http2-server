package main

import (
	"fmt"
	"log"
	"net/http"
	mw "apiproject/internal/api/middlewares"
	"crypto/tls"
	"path/filepath"
	"time"
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

	rl := mw.NewRateLimiter(5, time.Minute)

	middlewares := mw.Compression(
		mw.ResponseTimeMiddleware(
			mw.SecurityHeaders(
				mw.Cors(mux),
			),
		),
	)

	server := &http.Server{
		Addr: port,
		Handler: rl.Middleware(middlewares),
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server Listening on port:", port)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("error starting server", err)
	}
}
