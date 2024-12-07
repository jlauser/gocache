package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/jlauser/gocache/model"
	"net/http"
	"time"
)

func (s *Server) SearchHandler(w http.ResponseWriter, r *http.Request) {

	// set defaults
	buildData := true
	var data model.SearchData
	now := time.Now().Unix()
	search := ""

	// check for search term passed in route of querystring
	query := chi.URLParam(r, "q")
	if query != "" {
		search = query
	} else {
		queryParams := r.URL.Query()
		query = queryParams.Get("q")
		if query != "" {
			search = query
		}
	}

	if search == "" {
		// No search passed
		data = model.SearchData{
			Source:  "",
			Expires: -1,
			Query:   search,
			Results: model.SearchResults{
				Users:    nil,
				Featured: nil,
			},
		}
		buildData = false
	} else {
		// Check the cache
		cachedData, ok := s.cacheGet(search)
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
			usersList, _ := s.CSV.Find("users", search)
			featuredList, _ := s.CSV.Find("content_featured", search)
			data = model.SearchData{
				Source:  "origin",
				Query:   search,
				Expires: time.Now().Add(time.Second * time.Duration(s.Config.Redis.DefaultCacheSeconds)).Unix(),
				Results: model.SearchResults{
					Users:    model.UsersFromList(usersList.([][]string)),
					Featured: model.FeaturesFromList(featuredList.([][]string)),
				},
			}
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
		s.cacheInsert(query, jsonString)
	}
}
