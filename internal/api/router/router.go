package router

import (
	"net/http"
	handlers "apiproject/internal/api/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/teachers/", handlers.TeacherHandler)
	return mux
}
