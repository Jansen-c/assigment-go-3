package service

import (
	"kasir-api/models"
	"kasir-api/repository"
)

type TransactionService struct {
	repo *repository.TransactionRepository
}

func NewTransactionService(repo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}

func (s *TransactionService) GetTodaysReport() (*models.TodaysReport, error) {
	return s.repo.GetTodaysReport()
}

func (s *TransactionService) GetReport(startDate string, endDate string) (*models.TodaysReport, error) {
	return s.repo.GetReport(startDate, endDate)
}
