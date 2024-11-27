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
	var data model.SearchData
	cachedData, ok := app.cacheGet(search)
	if !ok {
		usersList, _ := app.CSV.Find("users", search)
		featuredList, _ := app.CSV.Find("content_featured", search)
		data = model.SearchData{
			Source: "origin",
			Results: model.SearchResults{
				Users:    model.UsersFromList(usersList.([][]string)),
				Featured: model.FeaturesFromList(featuredList.([][]string)),
			},
		}

	} else {
		err := json.Unmarshal([]byte(cachedData), &data)
		if err != nil {
			http.Error(w, "Error retrieving from cache", http.StatusInternalServerError)
		}
		data.Source = "cache"
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
	app.cacheInsert(term, string(jsonData))
}
