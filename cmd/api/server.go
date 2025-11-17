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

func execsHandler(w http.ResponseWriter, r *http.Request){
	switch r.Method {
		case http.MethodGet:
			w.Write([]byte("Hello GET Method on Execs Route"))
		case http.MethodPost:
			fmt.Println("Query: ", r.URL.Query())

			err := r.ParseForm()
			if err != nil {
				return
			}
			fmt.Println("Form data from POST: ", r.Form)
	}
}

func main() {
	port := ":8080"

	cert := filepath.Join("cmd", "config", "cert.pem")
	key := filepath.Join("cmd", "config", "key.pem")

	mux := http.NewServeMux()

	mux.HandleFunc("/execs", execsHandler)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	rl := mw.NewRateLimiter(5, time.Minute)

	hppOptions := mw.HPPOptions{
		CheckQuery: 			  true,
		CheckBody:  		      true,
		CheckBodyForContentType:  "application/x-www-form-urlencoded",
		Whitelist:				  []string{"sortBy", "sortOrder", "name", "age", "class"},
	}

	middlewares := mw.Compression(
		mw.ResponseTimeMiddleware(
			mw.SecurityHeaders(
				mw.Cors(mux),
			),
		),
	)

	secureMux := mw.HPP(hppOptions)(rl.Middleware(middlewares))

	server := &http.Server{
		Addr: port,
		Handler: secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server Listening on port:", port)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("error starting server", err)
	}
}
