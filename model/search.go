package model

type SearchData struct {
	Source  string        `json:"source"`
	Query   string        `json:"q"`
	Expires int64         `json:"expires"`
	Results SearchResults `json:"results"`
}

type SearchResults struct {
	Users    []User    `json:"users"`
	Featured []Feature `json:"featured"`
}
