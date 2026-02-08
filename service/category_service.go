package service

import (
	"kasir-api/models"
	"kasir-api/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) Create(data *models.Category) error {
	return s.repo.Create(data)
}

func (s *CategoryService) GetByID(id string) (*models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) Update(category *models.Category) error {
	return s.repo.Update(category)
}

func (s *CategoryService) Delete(id string) error { //? kayaknya yang perlu (id <type>) type di params, cuma kalau belum ada variable itu di-defined sebelumnya. beda kalau di atasnya ada id := ...
	return s.repo.Delete(id)
}
