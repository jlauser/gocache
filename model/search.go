package model

type SearchResults struct {
	Users    []User    `json:"users"`
	Featured []Feature `json:"featured"`
}
