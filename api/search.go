package api

import (
	"encoding/json"
	"github.com/jlauser/gocache/model"
	"net/http"
	"time"
)

func (app *Application) SearchHandler(w http.ResponseWriter, r *http.Request) {
	// set defaults
	buildData := true
	var data model.SearchData
	now := time.Now().Unix()
	search := "*"
	// get the search term if passed
	queryParams := r.URL.Query()
	term := queryParams.Get("term")
	if term != "" {
		search = term
	}

	// Check the cache
	cachedData, ok := app.cacheGet(search)
	if ok {
		err := json.Unmarshal([]byte(cachedData), &data)
		if err != nil {
			http.Error(w, "Error retrieving from cache", http.StatusInternalServerError)
		}
		// Is the cache expired?
		if data.Expires == 0 || data.Expires > now {
			data.Source = "cache"
			buildData = false
		} else {
			buildData = true
		}
	}
	// build from origin
	if buildData {
		usersList, _ := app.CSV.Find("users", search)
		featuredList, _ := app.CSV.Find("content_featured", search)
		data = model.SearchData{
			Source:  "origin",
			Expires: time.Now().Add(time.Second * time.Duration(app.Config.Redis.DefaultCacheSeconds)).Unix(),
			Results: model.SearchResults{
				Users:    model.UsersFromList(usersList.([][]string)),
				Featured: model.FeaturesFromList(featuredList.([][]string)),
			},
		}
	}
	// encode for response
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
	// store in cache
	if buildData {
		jsonString := string(jsonData)
		app.cacheInsert(term, jsonString)
	}
}
