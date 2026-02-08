package models

type Category struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
}

type Response struct {
	Results []Category `json:"results"`
	Count   int        `json:"count"`
}
