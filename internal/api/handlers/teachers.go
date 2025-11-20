package handlers

import (
	"net/http"
	"strings"
	"encoding/json"
	"strconv"
	"fmt"
	"sync"

	models "apiproject/internal/models"
	sqlconnect "apiproject/internal/repository/sqlconnect"
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
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	defer db.Close()
	var newTeachers []models.Teacher
	err = json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO teachers (first_name, last_name, email, class, subject)
		VALUES (?, ?, ?, ?, ?)
	`
	stmt, err := db.Prepare(query)
	if err != nil {
		http.Error(w, "Error setting SQL statement", http.StatusInternalServerError)
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(
			newTeacher.FirstName,
			newTeacher.LastName,
			newTeacher.Email,
			newTeacher.Class,
			newTeacher.Subject,
		)
		if err != nil {
			http.Error(w, "Error inserting data into database", http.StatusInternalServerError)
			return
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Error getting last insert id", http.StatusInternalServerError)
			return
		}
		newTeacher.ID = int(lastID)
		addedTeachers[i] = newTeacher
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
