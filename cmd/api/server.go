package main

import (
	"fmt"
	"log"
	"net/http"
	mw "apiproject/internal/api/middlewares"
	router "apiproject/internal/api/router"
	sqlconnect "apiproject/internal/repository/sqlconnect"
	//utils "apiproject/pkg/utils"
	"crypto/tls"
	"path/filepath"
)

func main() {
	_, err := sqlconnect.ConnectDB("test_database")
	if err != nil {
		fmt.Println("Error:----: ", err)
		return
	}

	port := ":8080"

	cert := filepath.Join("cmd", "config", "cert.pem")
	key := filepath.Join("cmd", "config", "key.pem")

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// hppOptions := mw.HPPOptions{
	// 	CheckQuery: 			  true,
	// 	CheckBody:  		      true,
	// 	CheckBodyForContentType:  "application/x-www-form-urlencoded",
	// 	Whitelist:				  []string{"sortBy", "sortOrder", "name", "age", "class"},
	// }
	//rl := mw.NewRateLimiter(5, time.Minute)
	// secureMux := utils.ApplyMiddlewares(
	// 	mux,
	// 	mw.HPP(hppOptions),
	// 	mw.Compression,
	// 	mw.SecurityHeaders,
	// 	mw.ResponseTimeMiddleware,
	// 	rl.Middleware,
	// 	mw.Cors,
	// )
	secureMux := mw.SecurityHeaders(router.Router())

	server := &http.Server{
		Addr: port,
		Handler: secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server Listening on port:", port)
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("error starting server", err)
	}
}
