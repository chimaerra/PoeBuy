package models

type LivesearchNewItem struct {
	New []string `json:"new"`
}

type LivesearchAuthStatus struct {
	Auth bool `json:"auth"`
}
