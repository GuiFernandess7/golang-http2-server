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
	"strconv"
)

type Teacher struct {
	ID 			int		`json:"id,omitempty"`
	FirstName 	string	`json:"first_name,omitempty"`
	LastName 	string	`json:"last_name,omitempty"`
	Class 		string  `json:"class,omitempty"`
	Subject		string	`json:"subject,omitempty"`
}

var (
	teachers 	= make(map[int]Teacher)
	mutex 		= &sync.Mutex{}
	nextID 		= 1
)

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
	nextID++
}

func getTeachersHandlers(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusForbidden)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")
	fmt.Println(idStr)

	if idStr == "" {
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
			Status string   `json:"status"`
			Count  int      `json:"count"`
			Data   []Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teacherList),
			Data:   teacherList,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return

	} else {
		idInt, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid parameter", http.StatusForbidden)
			return
		}

		teacher, exists := teachers[idInt]
		if !exists {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(teacher)
		return
	}
}

func createTeacherHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	var newTeachers []Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	addedTeachers := make([]Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		newTeacher.ID = nextID
		teachers[nextID] = newTeacher
		addedTeachers[i] = newTeacher
		nextID++
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string   `json:"status"`
		Count  int      `json:"count"`
		Data   []Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
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

	mux.HandleFunc("/teachers/", func(w http.ResponseWriter, r *http.Request) {
	    switch r.Method {
	    case http.MethodGet:
	        getTeachersHandlers(w, r)
	    case http.MethodPost:
	        createTeacherHandler(w, r)
	    default:
	        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	    }
	})

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
