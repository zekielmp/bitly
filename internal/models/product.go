package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a product category in the system.
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	// Products []Product `json:"-"`
}

// Product represents a product in the system.
type Product struct {
	gorm.Model
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Price       float64        `json:"price" gorm:"not null"`
	Stock       int            `json:"stock" gorm:"default:0" `
	SKU         string         `json:"sku" gorm:"uniqueIndex;not null"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	// Category    Category       `json:"category" `
	// Images      []ProductImage `json:"images" `
	// OrdrerItems []OrderItem    `json:"-"`
	// CartItems   []CartItem     `json:"-"`
}

// ProductImage represents an image associated with a product.
type ProductImage struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ProductID uint           `json:"product_id" gorm:"not null"`
	URL       string         `json:"url" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	// Product Product `json:"-"`
}
