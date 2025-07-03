package router

import (
	"issue-api/handlers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/issue", handlers.CreateIssue).Methods("POST")
	r.HandleFunc("/issues", handlers.GetIssues).Methods("GET")
	r.HandleFunc("/issue/{id:[0-9]+}", handlers.GetIssueByID).Methods("GET")
	r.HandleFunc("/issue/{id:[0-9]+}", handlers.UpdateIssue).Methods("PATCH")
	return r
}
