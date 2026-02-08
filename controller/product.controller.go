package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	// "strings"
	"kasir-api/models"
	// "assignment-go-1/models"
)

var product = []models.Product{
	{Id: "1", Name: "product pertama", Price: "10000", Stock: 10},
	{Id: "2", Name: "product kedua", Price: "20000", Stock: 20},
	{Id: "3", Name: "product ketiga", Price: "30000", Stock: 30},
	{Id: "4", Name: "product keempat", Price: "40000", Stock: 40},
}

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		result := models.ProductResponse{
			Results: product,
			Count:   len(product),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result) // Using json here
		return
	} else if r.Method == "POST" {
		var newProduct models.Product                      // think it like let newCategory: Category in typescript.
		err := json.NewDecoder(r.Body).Decode(&newProduct) //? how come newCategory suddenly have values? value is saved in err right?

		if err != nil { // early exit if anything seems off
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		// fmt.Println(newCategory.Name, newCategory.Description, "newCategory in POST", r.Body) // desc in json but retrieve it based on whats on struct above.
		// createdCategory := postNewCategory(newCategory.Name, newCategory.Description)

		latestId := product[len(product)-1].Id
		lastestIdInt, _ := strconv.Atoi(latestId)

		//TODO is this the way how to make new 'object'?
		newProd := models.Product{
			Id:   fmt.Sprintf("%d", lastestIdInt+1), // copy the behaviour of latest id + 1. not so much like sql.
			Name: newProduct.Name,
		}
		product = append(product, newProd)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // 201

		json.NewEncoder(w).Encode(newProd)

		return
	} else {
		// default to error
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func ProductHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/products/"):]
	if idStr == "" {
		http.Error(w, "ID parameter is missing for put method", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" && idStr != "" {

		for _, val := range product {
			if val.Id == idStr {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(val)
				return
			}
		}

		http.Error(w, "No product corresponds with this specific id", http.StatusNotFound)

	} else if r.Method == "PUT" {

		var newProduct models.Product
		err := json.NewDecoder(r.Body).Decode(&newProduct)

		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		for idx, val := range product {
			if val.Id == idStr {

				product[idx].Name = newProduct.Name
				// product[idx].Description = newProduct.Description

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(product[idx])
				return
			}
		}
	} else if r.Method == "POST" {
		var newProduct models.Product
		err := json.NewDecoder(r.Body).Decode(&newProduct)

		if err != nil { // early exit if anything seems off
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		latestId := product[len(product)-1].Id
		lastestIdInt, _ := strconv.Atoi(latestId)

		// this is how we make object, and then append it manually.
		newProd := models.Product{
			Id:    fmt.Sprintf("%d", lastestIdInt+1), // copy the behaviour of latest id + 1. not so much like sql.
			Name:  newProduct.Name,
			Price: newProduct.Price,
			Stock: newProduct.Stock,
		}
		product = append(product, newProd)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // 201

		json.NewEncoder(w).Encode(newProd)

		return
	} else if r.Method == "DELETE" {
		if idStr == "" {
			http.Error(w, "ID parameter is missing for delete method", http.StatusBadRequest)
			return
		}

		for idx, val := range product {
			if val.Id == idStr {
				arrBaru := append(product[:idx], product[idx+1:]...)
				product = arrBaru

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{
					"status":  "success",
					"message": "product with specific id deleted",
				})
				return
			}
		}
		http.Error(w, "No product found with specific id for delete method", http.StatusNotFound)
	} else {
		// default to error
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}
