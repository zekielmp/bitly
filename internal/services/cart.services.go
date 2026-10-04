package services

import (
	"errors"

	"github.com/zekielmp/Bitly/internal/dto"
	"github.com/zekielmp/Bitly/internal/models"
	"gorm.io/gorm"
)

type CartService struct {
	db *gorm.DB
}

func NewCartServices(db *gorm.DB) *CartService {
	return &CartService{
		db: db,
	}
}

func (s *CartService) GetCart(userID uint) (*dto.CartResponse, error) {
	var cart models.Cart
	err := s.db.Preload("CartItems.Product.Category").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}

	return s.cartResponse(&cart), nil

}

func (s *CartService) AddToCart(userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error) {

	/*Check if product exists*/
	var product models.Product
	if err := s.db.First(&product, req.ProductID).Error; err != nil {
		return nil, errors.New("product not found")
	}

	if product.Stock < req.Quantity {
		return nil, errors.New("insuffient stock")
	}

	/*Get or create cart*/
	var cart models.Cart
	if err := s.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		cart = models.Cart{UserID: userID}
		if err = s.db.Create(&cart).Error; err != nil {
			return nil, err
		}
	}

	/*Check if item already exists in cart*/
	var carItem models.CartItem
	if err := s.db.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&carItem).Error; err != nil {
		/*Create new cart item*/
		carItem = models.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		s.db.Create(&carItem)
	} else {
		/*update existign cart item*/
		carItem.Quantity += req.Quantity
		if carItem.Quantity > product.Stock {
			return nil, errors.New("insufficient stock")
		}
		s.db.Save(&carItem)
	}
	return s.GetCart(userID)
}

func (s *CartService) UpdateCartItem(userID, itemID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error) {
	var cartItem models.CartItem
	if err := s.db.Joins("JOIN carts ON cart_items.cart_id = cart.id").Where("cart_items.id = ? AND cart.user_id = ?", itemID, userID).First(&cartItem).Error; err != nil {
		return nil, errors.New("cart item not found")
	}

	var product models.Product
	if err := s.db.First(&product, cartItem.ProductID).Error; err != nil {
		return nil, errors.New("product not found")
	}
	if product.Stock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	cartItem.Quantity = req.Quantity
	if err := s.db.Save(&cartItem).Error; err != nil {
		return nil, err
	}
	return s.GetCart(userID)
}

func (s *CartService) RemoveFromCart(userID, itemsID uint) error {
	return s.db.Joins("JOIN carts ON cart_items.cart_id = carts.id").Where("cart_items.id = ? AND carts.user_id = ?", itemsID, userID).Delete(&models.CartItem{}).Error
}

func (s *CartService) cartResponse(cart *models.Cart) *dto.CartResponse {
	cartItems := make([]dto.CartItemsResponse, len(cart.CartItems))
	var total float64
	for i := range cart.CartItems {
		subtotal := float64(cart.CartItems[i].Quantity) * cart.CartItems[i].Product.Price
		total += subtotal
		cartItems[i] = dto.CartItemsResponse{
			ID: cart.CartItems[i].ID,
			Product: dto.ProductResponse{
				ID:          cart.CartItems[i].Product.ID,
				CategoryID:  cart.CartItems[i].Product.CategoryID,
				Name:        cart.CartItems[i].Product.Name,
				Description: cart.CartItems[i].Product.Description,
				Price:       cart.CartItems[i].Product.Price,
				Stock:       cart.CartItems[i].Product.Stock,
				SKU:         cart.CartItems[i].Product.SKU,
				IsActive:    cart.CartItems[i].Product.IsActive,
				Category: dto.CategoryResponse{
					ID:          cart.CartItems[i].Product.Category.ID,
					Name:        cart.CartItems[i].Product.Category.Name,
					Description: cart.CartItems[i].Product.Description,
					IsActive:    cart.CartItems[i].Product.IsActive,
				},
			},
			Quantity: cart.CartItems[i].Quantity,
			Subtotal: subtotal,
		}
	}

	return &dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CartItems: cartItems,
		Total:     total,
	}
}
