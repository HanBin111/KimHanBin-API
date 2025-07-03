package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"issue-api/models"
	"issue-api/store"

	"github.com/gorilla/mux"
)

type IssueRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UserID      *uint  `json:"userId"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": msg,
		"code":  code,
	})
}

func CreateIssue(w http.ResponseWriter, r *http.Request) {
	var req IssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "잘못된 요청입니다")
		return
	}
	if req.Title == "" {
		respondWithError(w, http.StatusBadRequest, "제목은 필수입니다")
		return
	}

	issue := models.Issue{
		Title:       req.Title,
		Description: req.Description,
	}
	result, err := store.GetStore().CreateIssue(issue, req.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func GetIssues(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}
	issues := store.GetStore().GetAllIssues(statusPtr)
	json.NewEncoder(w).Encode(map[string]interface{}{"issues": issues})
}

func GetIssueByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	issue, err := store.GetStore().GetIssueByID(uint(id))
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	json.NewEncoder(w).Encode(issue)
}

func UpdateIssue(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req IssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "잘못된 요청입니다")
		return
	}

	update := models.Issue{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}

	issue, err := store.GetStore().UpdateIssue(uint(id), update, req.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	json.NewEncoder(w).Encode(issue)
}
