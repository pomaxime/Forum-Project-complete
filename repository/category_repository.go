package repository

import "forum/models"

type CategoryRepo struct {
    categories []models.Category
}

func NewCategoryRepo() *CategoryRepo {
    return &CategoryRepo{categories: models.DefaultCategories()}
}

func (r *CategoryRepo) IsValid(value string) bool {
    for _, category := range r.categories {
        if category.Value == value {
            return true
        }
    }
    return false
}

func (r *CategoryRepo) GetAll() []models.Category {
    return r.categories
}
