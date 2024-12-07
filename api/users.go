package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/jlauser/gocache/model"
	"net/http"
)

func (s *Server) UsersGetHandler(w http.ResponseWriter, r *http.Request) {
	search := "*"

	// check for user id
	employeeId := chi.URLParam(r, "id")
	if employeeId != "" {
		search = employeeId
	} else {
		queryParams := r.URL.Query()
		employeeId = queryParams.Get("id")
		if employeeId != "" {
			search = employeeId
		}
	}
	result, ok := s.CSV.Find("users", search)
	if ok {
		users := model.UsersFromList(result.([][]string))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		}
	}
}

func (s *Server) UserPostHandler(w http.ResponseWriter, r *http.Request) {
	employeeId := chi.URLParam(r, "id")
	if employeeId == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
}
