package api

import (
	"encoding/json"
	"github.com/jlauser/gocache/model"
	"net/http"
)

func (app *Application) UsersHandler(w http.ResponseWriter, r *http.Request) {
	search := "*"
	queryParams := r.URL.Query()
	employeeId := queryParams.Get("id")
	if employeeId != "" {
		search = employeeId
	}
	result, ok := app.CSV.Find("users", search)
	if ok {
		users := model.UsersFromList(result.([][]string))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		}
	}
}
