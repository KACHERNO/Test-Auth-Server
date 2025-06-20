package models

type Dbinfo struct {
	Current_database string `json:"current_database"`
	Version          string `json:"version"`
}
