package api

import (
	"encoding/json"
	"github.com/jlauser/gocache/model"
	"net/http"
)

func (app *Application) SearchHandler(w http.ResponseWriter, r *http.Request) {
	search := "*"
	queryParams := r.URL.Query()
	term := queryParams.Get("term")
	if term != "" {
		search = term
	}
	usersList, _ := app.CSV.Find("users", search)
	featuredList, _ := app.CSV.Find("content_featured", search)
	results := model.SearchResults{
		Users:    model.UsersFromList(usersList.([][]string)),
		Featured: model.FeaturesFromList(featuredList.([][]string)),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)

	}
}
