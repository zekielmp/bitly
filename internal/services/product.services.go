package services

import (
	"time"

	"github.com/zekielmp/Bitly/internal/dto"
	"github.com/zekielmp/Bitly/internal/models"
	"gorm.io/gorm"
)

type ProductServices struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductServices {
	return &ProductServices{
		db: db,
	}
}

func (p *ProductServices) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	if err := p.db.Create(&category).Error; err != nil {
		return nil, err
	}
	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
	}, nil

}

func (p *ProductServices) GetCategories() ([]dto.CategoryResponse, error) {
	var categories []models.Category

	if err := p.db.Where("is_active = ?", true).Find(&categories).Error; err != nil {
		return nil, err
	}

	response := make([]dto.CategoryResponse, len(categories))
	for i := range categories {
		response[i] = dto.CategoryResponse{
			ID:          categories[i].ID,
			Name:        categories[i].Name,
			Description: categories[i].Description,
			IsActive:    categories[i].IsActive,
		}
	}
	return response, nil
}

func (p ProductServices) GetCategory(category *models.Category, id uint) (*dto.CategoryResponse, error) {
	if err := p.db.First(&category, id).Error; err == nil {
		return nil, err
	}
	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
	}, nil
}

func (p *ProductServices) UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {

	var category models.Category
	if err := p.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	category.Name = req.Name
	category.Description = req.Description
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}
	if err := p.db.Save(&category).Error; err != nil {
		return nil, err
	}
	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
	}, nil
}
func (p *ProductServices) DeleteCategory(id uint) error {
	return p.db.Delete(models.Category{}, id).Error
}

// func (p *ProductServices) AddProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
// 	product := models.Product{
// 		Name:        req.Name,
// 		CategoryID:  req.CategoryID,
// 		Description: req.Description,
// 		Price:       req.Price,
// 		Stock:       req.Stock,
// 		SKU:         req.SKU,
// 		CreatedAt:   time.Now(),
// 	}
// 	var c models.Category
// 	category, err := p.GetCategory(&c, req.CategoryID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := p.db.Create(&product).Error; err != nil {
// 		return nil, err
// 	}

// 	return &dto.ProductResponse{
// 		ID: product.ID,
// 		// CategoryID:  category.ID,
// 		Name:        product.Name,
// 		Description: product.Description,
// 		Price:       product.Price,
// 		Stock:       product.Stock,
// 		SKU:         product.SKU,
// 		Category: dto.CategoryResponse{
// 			ID:          category.ID,
// 			Name:        category.Name,
// 			Description: category.Description,
// 			IsActive:    category.IsActive,
// 		},
// 	}, nil

// }
