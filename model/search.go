package model

type SearchData struct {
	Source  string        `json:"source"`
	Results SearchResults `json:"results"`
}

type SearchResults struct {
	Users    []User    `json:"users"`
	Featured []Feature `json:"featured"`
}
