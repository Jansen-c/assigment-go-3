package repository

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll(productSearchName string) ([]models.Product, error) {
	query := "SELECT id, name, price, stock FROM product"

	args := []any{} // Creates an empty slice of any type, e.g. can append like this: append(items, "string", 1, true)
	if productSearchName != "" {
		query += " WHERE p.name ILIKE $1"
		args = append(args, "%"+productSearchName+"%")
	}

	rows, err := repo.db.Query(query, args...) // ... = Unpacks the slice into individual arguments
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.Id, &p.Name, &p.Price, &p.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO product (name, price, stock) VALUES ($1, $2, $3) RETURNING id"
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock).Scan(&product.Id)
	return err
}

// GetByID - ambil product by ID
func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := "SELECT id, name, price, stock FROM product WHERE id = $1"

	var p models.Product
	err := repo.db.QueryRow(query, id).Scan(&p.Id, &p.Name, &p.Price, &p.Stock)
	if err == sql.ErrNoRows {
		return nil, errors.New("product tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE product SET name = $1, price = $2, stock = $3 WHERE id = $4"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.Id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product tidak ditemukan")
	}

	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM product WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("product tidak ditemukan")
	}

	return err
}
