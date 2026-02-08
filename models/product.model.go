package models

type Product struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
	Stock int    `json:"stock"`
}

type ProductResponse struct {
	Results []Product `json:"results"`
	Count   int       `json:"count"`
}
