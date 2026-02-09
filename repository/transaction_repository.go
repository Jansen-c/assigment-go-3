package repository

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	var (
		res *models.Transaction
	)

	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	// Create an empty slice first. we will append this below.
	details := make([]models.TransactionDetail, 0)
	// Loop per item
	for _, item := range items {
		var productName string
		var productID, price, stock int
		// get product dapet pricing
		err := tx.QueryRow("SELECT id, name, price, stock FROM product WHERE id=$1", item.ProductID).Scan(&productID, &productName, &price, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}

		if err != nil {
			return nil, err
		}

		// Count subtotal then add to totalAmount
		subtotal := item.Quantity * price
		totalAmount += subtotal

		// kurangi jumlah stok
		_, err = tx.Exec("UPDATE product SET stock = stock - $1 WHERE id = $2", item.Quantity, productID)
		if err != nil {
			return nil, err
		}

		// Append to details
		details = append(details, models.TransactionDetail{
			ProductID:   productID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// insert transaction
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING ID", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// for i := range details {
	// 	details[i].TransactionID = transactionID
	// 	_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
	// 		transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	//* Task 3 no. 1: Refactor this, instead of exec per loop. Create loop, then append everything into arr, finally insert all at once as exec.
	// Insert transaction details. commented below is deprecated.
	// for i, detail := range details {
	// 	details[i].TransactionID = transactionID
	// 	_, err := tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2,$3,$4)", transactionID, detail.ProductID, detail.Quantity, detail.Subtotal)
	// 	if err != nil {
	if len(details) > 0 {
		query := "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES "
		// args := make([]interface{}, 0, len(details)*4) // hard to read, use any with unlimited size like below first.
		args := []any{}

		for i, detail := range details {
			details[i].TransactionID = transactionID
			if i > 0 {
				query += ","
			}
			query += fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
			args = append(args, transactionID, detail.ProductID, detail.Quantity, detail.Subtotal)
		}

		if _, err := tx.Exec(query, args...); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	res = &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}

	return res, nil
}

func (repo *TransactionRepository) GetTodaysReport() (*models.TodaysReport, error) {
	totalCountQuery := "SELECT COUNT(*) stock FROM transactions WHERE created_at >= CURRENT_DATE AND created_at < CURRENT_DATE + INTERVAL '1 day'"
	totalRevenueQuery := "SELECT COALESCE(SUM(total_amount), 0) FROM transactions WHERE created_at >= CURRENT_DATE AND created_at < CURRENT_DATE + INTERVAL '1 day'" // Coalesce ensure no null value, instead it will return 0]
	bestSellerProductQuery := `
		SELECT 
			p.id, 
			p.name,
			td.quantity 
		FROM transaction_details td 
		JOIN transactions t ON t.id = td.transaction_id
		JOIN product p ON p.id = td.product_id 
		WHERE t.created_at >= CURRENT_DATE AND t.created_at < CURRENT_DATE + INTERVAL '1 day' 
		GROUP BY p.id, p.name, td.quantity
		ORDER BY td.quantity DESC
		LIMIT 1
	`

	var totalTransactions int
	err := repo.db.QueryRow(totalCountQuery).Scan(&totalTransactions)
	if err != nil {
		return nil, err
	}

	var totalRevenue float64
	err = repo.db.QueryRow(totalRevenueQuery).Scan(&totalRevenue)
	if err != nil {
		return nil, err
	}

	var bestSeller models.BestSellerProduct
	err = repo.db.QueryRow(bestSellerProductQuery).Scan(&bestSeller.ProductID, &bestSeller.ProductName, &bestSeller.Qty)
	if err != nil {
		// If no sales today, handle gracefully
		if err == sql.ErrNoRows {
			bestSeller = models.BestSellerProduct{}
		} else {
			return nil, err
		}
	}

	report := &models.TodaysReport{
		TotalTransactions: totalTransactions,
		TotalRevenue:      totalRevenue,
		BestSeller:        bestSeller,
	}

	// defer repo.db.Close() // jangan asal close sialan

	return report, nil

}

func (repo *TransactionRepository) GetReport(startDate string, endDate string) (*models.TodaysReport, error) {
	totalCountQuery := fmt.Sprintf("SELECT COUNT(*) stock FROM transactions WHERE created_at >= '%s 00:00:00' AND created_at <=  '%s 23:59:59'", startDate, endDate)
	totalRevenueQuery := fmt.Sprintf("SELECT COALESCE(SUM(total_amount), 0) FROM transactions WHERE created_at >= '%s 00:00:00' AND created_at <= '%s 23:59:59'", startDate, endDate)
	bestSellerProductQuery := fmt.Sprintf(`
		SELECT 
			p.id, 
			p.name,
			COUNT(*) AS quantity 
		FROM transaction_details td 
		JOIN transactions t ON t.id = td.transaction_id
		JOIN product p ON p.id = td.product_id 
		WHERE t.created_at >= '%s 00:00:00' AND t.created_at <= '%s 23:59:59' 
		GROUP BY p.id, p.name 
		ORDER BY quantity DESC 
		LIMIT 1
	`, startDate, endDate)

	fmt.Println(totalCountQuery, "totalCountQuery")
	var totalTransactions int
	err := repo.db.QueryRow(totalCountQuery).Scan(&totalTransactions)
	if err != nil {
		return nil, err
	}

	var totalRevenue float64
	err = repo.db.QueryRow(totalRevenueQuery).Scan(&totalRevenue)
	if err != nil {
		return nil, err
	}

	var bestSeller models.BestSellerProduct
	err = repo.db.QueryRow(bestSellerProductQuery).Scan(&bestSeller.ProductID, &bestSeller.ProductName, &bestSeller.Qty)
	if err != nil {
		// If no sales today, handle gracefully
		if err == sql.ErrNoRows {
			bestSeller = models.BestSellerProduct{}
		} else {
			return nil, err
		}
	}

	report := &models.TodaysReport{
		TotalTransactions: totalTransactions,
		TotalRevenue:      totalRevenue,
		BestSeller:        bestSeller,
	}

	return report, nil

}
