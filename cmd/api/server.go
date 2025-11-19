package main

import (
	"fmt"
	"log"
	"net/http"
	mw "apiproject/internal/api/middlewares"
	"crypto/tls"
	"path/filepath"
	"sync"
	"encoding/json"
	"strings"
)

type Teacher struct {
	ID 			int
	FirstName 	string
	LastName 	string
	Class 		string
	Subject		string
}

var teachers = make(map[int]Teacher)
var mutex = &sync.Mutex{}
var nextID = 1

func init() {
	teachers[nextID] = Teacher {
		ID: nextID,
		FirstName: "John",
		LastName: "Cena",
		Class: "A",
		Subject: "Fight",
	}
	nextID++
	teachers[nextID] = Teacher {
		ID: nextID,
		FirstName: "Jake",
		LastName: "Peralta",
		Class: "B",
		Subject: "Investigation",
	}
	nextID++
	teachers[nextID] = Teacher {
		ID: nextID,
		FirstName: "Robert",
		LastName: "Greene",
		Class: "B",
		Subject: "Biology",
	}
}

func getTeachers(w http.ResponseWriter, r *http.Request){
	fmt.Println(r.Method)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusForbidden)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")
	fmt.Println(idStr)

	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")

	teacherList := make([]Teacher, 0, len(teachers))
	for _, teacher := range teachers {
		if (firstName == "" || teacher.FirstName == firstName) &&
			(lastName == "" || teacher.LastName == lastName) {
		teacherList = append(teacherList, teacher)
		}
	}

	response := struct {
		Status string `json:"status"`
		Count int  	  `json:"count"`
		Data []Teacher `json:"data"`
	}{
		Status: "success",
		Count: len(teacherList),
		Data: teacherList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type Middleware func(http.Handler) http.Handler

func applyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

func main() {
	port := ":8080"

	cert := filepath.Join("cmd", "config", "cert.pem")
	key := filepath.Join("cmd", "config", "key.pem")

	mux := http.NewServeMux()

	mux.HandleFunc("/teachers/", getTeachers)

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
	// secureMux := applyMiddlewares(
	// 	mux,
	// 	mw.HPP(hppOptions),
	// 	mw.Compression,
	// 	mw.SecurityHeaders,
	// 	mw.ResponseTimeMiddleware,
	// 	rl.Middleware,
	// 	mw.Cors,
	// )
	secureMux := mw.SecurityHeaders(mux)

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
