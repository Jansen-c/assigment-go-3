package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handler"
	"kasir-api/models"
	"kasir-api/repository"
	"kasir-api/service"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

type PostResponse struct {
	Results []models.Category `json:"results"`
	Count   int               `json:"count"`
}

var category = []models.Category{
	{Id: "1", Name: "category pertama", Description: "desc1"},
	{Id: "2", Name: "category kedua", Description: "desc2"},
	{Id: "3", Name: "category ketiga", Description: "desc3"},
	{Id: "4", Name: "category keempat", Description: "desc4"},
}

func getAllCategories() []models.Category { // type must be defined exactly like this []Product instead of just Product. imagine def add(a: int, b: int) -> int: in python
	return category
}

func postNewCategory(name string, description string) models.Category { // type must be defined exactly like this []Product instead of just Product
	latestId := category[len(category)-1].Id
	// fmt.Println(latestId, "latestId")
	lastestIdInt, _ := strconv.Atoi(latestId)

	newCat := models.Category{
		// Id:          fmt.Sprintf("%d", len(category)+1), // sht. what a pain just to use template literal. %d, %f, %s has literraly the same style as C. Ffs why complicated on things like this, my god.
		Id:          fmt.Sprintf("%d", lastestIdInt+1), // copy the behaviour of latest id + 1. not so much like sql.
		Name:        name,
		Description: description,
	}
	category = append(category, newCat)
	return newCat
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	//? mirip ternary dengan else empty, cara bacanya return dari os.Stat lgs di-destructure, kalau ada (pake ":" bukan "?") jalanin yang pertama, kalau gada ga jalan
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		viper.ReadInConfig() //? emang harus pake "_ = viper" ... ?
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	defer db.Close() //?

	//? Dependency injection for what purposes?
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	//? Dependency injection for category
	categoryRepo := repository.NewCategoryRepository(db) // Same connection pool as above, i thought we make a new ones for every repo.
	// fmt.Println(categoryRepo, "categoryRepo")
	categoryService := service.NewCategoryService(categoryRepo)
	// fmt.Println(categoryService, "categoryService")
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Dependency injection for transaction
	// Transaction
	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepo)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)           // oh bisa langsung post gini, shorthand daripada kayak di bawah.
	http.HandleFunc("/api/report/hari-ini", transactionHandler.HandleReportToday) // Endpoint report today, sesuai yang diminta docs ass 3.
	http.HandleFunc("/api/report", transactionHandler.HandleReport)

	http.HandleFunc("/api/product", func(w http.ResponseWriter, r *http.Request) { // end trailing '/' is a must so ones with route params works.
		// controller.ProductsHandler(w, r)
		productHandler.HandleProducts(w, r)
		return
	})

	// --- /products with trailing "/"
	http.HandleFunc("/api/product/", func(w http.ResponseWriter, r *http.Request) { // end trailing '/' is a must so ones with route params works.
		// controller.ProductHandler(w, r)
		productHandler.HandleProductByID(w, r)
		return
	})

	// --- /categories without trailing "/"
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) { // end trailing '/' is a must so ones with route params works.
		categoryHandler.HandleCategory(w, r)
		return
	})

	// --- /categories/ with trailin "/". why not {id}? because there is no support on default go.
	//? just in you wondering again why need to separate one with no / and other one with / is go treat as different handleFunc.
	//? just like why in express id params :id always put at the very bottom. although i prefer something like ?id= for less code, but assignment asks this way.
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		categoryHandler.HandleCategoryByID(w, r)
		return

	})

	// --- health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{ //? this weird ahh syntax again.
			"status":  "OK",
			"message": "API Running",
		})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Jansen Go Second Assignment</title>
			<style>
				body { font-family: sans-serif; }
				ul { list-style: none; padding-left: 0; }
				li { margin: 0.5em 0; }
				code { background-color: #eee; padding: 0.2em 0.4em; border-radius: 3px; }
			</style>
		</head>
		<body>
			<h1>Available API Routes</h1>
			<h6>Third assignment. 
			<br>Fix inserting transaction details to db in transaction_repository.go, summary report, and summary report by date</h6>
			<ul>
				<li>GET <a href="/api/report/hari-ini">/api/report/hari-ini</a></li>
				<li>GET <a href="/api/report?start_date=2026-02-1&end_date=2026-02-17">/api/report?start_date=2026-02-1&end_date=2026-02-17</a></li>
			</ul>
		</body>
		</html>
		`
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
	})

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running on", addr)

	err = http.ListenAndServe(addr, nil) // timpa err atas jadi ga perlu re-declare lagi, toh sama sama error, beda message doang nanti.
	if err != nil {
		fmt.Println("Failed to start server")
	}
}
