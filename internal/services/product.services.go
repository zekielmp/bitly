package services

import (
	"fmt"
	"time"

	"github.com/zekielmp/Bitly/internal/dto"
	"github.com/zekielmp/Bitly/internal/models"
	"github.com/zekielmp/Bitly/internal/utils"
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
	return p.db.Delete(&models.Category{}, id).Error
}

func (p *ProductServices) AddProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := models.Product{
		Name:        req.Name,
		CategoryID:  req.CategoryID,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		SKU:         req.SKU,
	}

	if err := p.db.Create(&product).Error; err != nil {
		return nil, err
	}

	return p.GetProduct(product.ID)

}

func (p *ProductServices) GetProducts(page, limit int) ([]dto.ProductResponse, *utils.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	var products []models.Product
	var total int64

	p.db.Model(&models.Product{}).Where("is_active = ?", true).Count(&total)
	if err := p.db.Preload("Category").Preload("Images").Where("is_active = ?", true).Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, nil, err
	}
	response := make([]dto.ProductResponse, len(products))
	for i := range products {
		response[i] = p.productResponse(&products[i])
	}

	totalPages := int((total*int64(limit) - 1) / int64(limit))
	meta := &utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
	return response, meta, nil

}

func (p *ProductServices) GetProduct(id uint) (*dto.ProductResponse, error) {
	var product models.Product
	if err := p.db.Preload("Category").Preload("Images").First(&product, id).Error; err != nil {
		return nil, err
	}
	res := p.productResponse(&product)
	return &res, nil
}

func (p *ProductServices) UpdateProduct(req *dto.UpdateProductRequest, id uint) (*dto.ProductResponse, error) {
	var product models.Product
	if err := p.db.First(&product, id).Error; err != nil {
		return nil, err
	}

	/*assign value to product models using req dto*/
	product.CategoryID = req.CategoryID
	product.Name = req.Name
	product.Description = req.Description
	product.Stock = req.Stock
	product.Price = req.Price
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := p.db.Save(&product).Error; err != nil {
		return nil, err
	}
	return p.GetProduct(id)
}

func (p *ProductServices) AddProductImage(productID uint, url, altText string) error {

	var count int64
	p.db.Model(&models.ProductImage{}).Where("product_id = ?", productID).Count(&count)

	image := models.ProductImage{
		ProductID: productID,
		URL:       url,
		AltText:   altText,
		IsPrimary: count == 0, /*First image is primary*/
	}

	return p.db.Create(&image).Error

}

func (p *ProductServices) DeleteProduct(id uint) error {
	return p.db.Delete(&models.Product{}, id).Error
}

func (p *ProductServices) AddProductReview(productID uint, userID uint, req *dto.ProductReviewRequest) error {
	var count int64
	p.db.Model(&models.ProductReviews{}).Where("product_id = ? AND user_id = ?", productID, userID).Count(&count)
	if count > 0 {
		return fmt.Errorf("user has already reviewed this product")
	}
	if err := p.db.Preload("User").Where("id = ?", userID).Find(&models.User{}).First(&models.User{}, userID).Error; err != nil {
		return fmt.Errorf("failed to find user")
	}

	review := models.ProductReviews{
		ProductID: productID,
		User: models.User{
			ID: userID,
		},
		Data:      req.Data,
		Rating:    req.Rating,
		CreatedAt: time.Now(),
	}
	return p.db.Create(&review).Error
}

func (p *ProductServices) productResponse(product *models.Product) dto.ProductResponse {
	images := make([]dto.ProductImageResponse, len(product.Images))
	for i := range product.Images {
		images[i] = dto.ProductImageResponse{
			ID:        product.Images[i].ID,
			URL:       product.Images[i].URL,
			AltText:   product.Images[i].AltText,
			IsPrimary: product.Images[i].IsPrimary,
		}
	}

	p.db.Preload("User").Where("product_id = ?", product.ID).Find(&product.ProductReviews)
	reviews := make([]dto.ProductReviewResponse, len(product.ProductReviews))
	for i := range product.ProductReviews {
		reviews[i] = dto.ProductReviewResponse{
			ID: product.ProductReviews[i].ID,
			User: dto.UserResponse{
				ID:        product.ProductReviews[i].User.ID,
				FirstName: product.ProductReviews[i].User.FirstName,
				LastName:  product.ProductReviews[i].User.LastName,
			},
			Data:   product.ProductReviews[i].Data,
			Rating: product.ProductReviews[i].Rating,
		}
	}

	return dto.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		SKU:         product.SKU,
		Category: dto.CategoryResponse{
			ID:          product.Category.ID,
			Name:        product.Category.Name,
			Description: product.Category.Description,
			IsActive:    product.Category.IsActive,
		},
		Images:  images,
		Reviews: reviews,
	}
}
