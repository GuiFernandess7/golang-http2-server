package handlers

import (
	"net/http"
	"strings"
	"encoding/json"
	"strconv"
	"fmt"
	"sync"

	models "apiproject/internal/models"
)

var (
	teachers 	= make(map[int]models.Teacher)
	mutex 		= &sync.Mutex{}
	nextID 		= 1
)

func init() {
	teachers[nextID] = models.Teacher {
		ID: nextID,
		FirstName: "John",
		LastName: "Cena",
		Class: "A",
		Subject: "Fight",
	}
	nextID++
	teachers[nextID] = models.Teacher {
		ID: nextID,
		FirstName: "Jake",
		LastName: "Peralta",
		Class: "B",
		Subject: "Investigation",
	}
	nextID++
	teachers[nextID] = models.Teacher {
		ID: nextID,
		FirstName: "Robert",
		LastName: "Greene",
		Class: "B",
		Subject: "Biology",
	}
	nextID++
}

func TeacherHandler(w http.ResponseWriter, r *http.Request){
	switch r.Method {
		case http.MethodGet:
			getTeachersHandlers(w, r)
		case http.MethodPost:
			createTeacherHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
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

		teacherList := make([]models.Teacher, 0, len(teachers))
		for _, teacher := range teachers {
			if (firstName == "" || teacher.FirstName == firstName) &&
				(lastName == "" || teacher.LastName == lastName) {
				teacherList = append(teacherList, teacher)
			}
		}

		response := struct {
			Status string   `json:"status"`
			Count  int      `json:"count"`
			Data   []models.Teacher `json:"data"`
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
	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	addedTeachers := make([]models.Teacher, len(newTeachers))
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
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	json.NewEncoder(w).Encode(response)
}
