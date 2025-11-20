package router

import (
	"net/http"
)

func router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/teachers/", handlers.TeacherHandler)
	return mux
}
