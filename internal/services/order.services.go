package services

import (
	"errors"
	"fmt"

	"github.com/zekielmp/Bitly/internal/dto"
	"github.com/zekielmp/Bitly/internal/models"
	"gorm.io/gorm"
)

const (
	dateFormat = "2006-01-02T15:04:05Z"
)

type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{
		db: db,
	}
}

func (s *OrderService) CreateOrder(userID uint) (*dto.OrderResponse, error) {
	var orderRes *dto.OrderResponse

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var cart models.Cart
		if err := tx.Preload("CartItems.Product").Where("user_id = ?", userID).First(&cart).Error; err != nil {
			return errors.New("cart not found")
		}
		if len(cart.CartItems) == 0 {
			return errors.New("cart is empty")
		}

		//calculate and validate stock
		var totalAmount float64
		var orderItems []models.OrderItem

		for i := range cart.CartItems {
			cartItem := &cart.CartItems[i]

			if cartItem.Product.Stock < cartItem.Quantity {
				return fmt.Errorf("insufficient stock for product: %s", cartItem.Product.Name)
			}

			itemTotal := float64(cartItem.Quantity) * cartItem.Product.Price
			totalAmount += itemTotal

			orderItems = append(orderItems, models.OrderItem{
				ProductID: cartItem.ProductID,
				Quantity:  cartItem.Quantity,
				Price:     cartItem.Product.Price,
			})

			//Update product stock
			cartItem.Product.Stock -= cartItem.Quantity
			if err := tx.Save(&cartItem.Product).Error; err != nil {
				return err
			}

			//create order
			order := models.Order{
				UserID:      userID,
				Status:      string(models.OrderStatusPending),
				TotalAmount: totalAmount,
				OrderItems:  orderItems,
			}
			if err := tx.Create(&order).Error; err != nil {
				return err
			}

			//clear cart
			if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
				return err
			}
			res, err := s.getOrderResponse(tx, order.ID)
			if err != nil {
				return err
			}
			orderRes = res
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return orderRes, nil
}

func (s *OrderService) getOrderResponse(tx *gorm.DB, orderID uint) (*dto.OrderResponse, error) {
	var order models.Order
	if err := tx.Preload("OrderItems.Product.Category").First(&order, orderID).Error; err != nil {
		return nil, err
	}
	res := s.orderResponse(&order)
	return &res, nil
}

func (s *OrderService) orderResponse(order *models.Order) dto.OrderResponse {
	orderItems := make([]dto.OrderItemsResponse, len(order.OrderItems))
	for i := range order.OrderItems {
		item := order.OrderItems[i]
		orderItems = append(orderItems, dto.OrderItemsResponse{
			ID: item.ID,
			Product: dto.ProductResponse{
				ID:          item.Product.ID,
				CategoryID:  item.Product.CategoryID,
				Name:        item.Product.Name,
				Description: item.Product.Description,
				Price:       item.Product.Price,
				Stock:       item.Product.Stock,
				SKU:         item.Product.SKU,
				IsActive:    item.Product.IsActive,
				Category: dto.CategoryResponse{
					ID:          item.Product.Category.ID,
					Name:        item.Product.Category.Name,
					Description: item.Product.Category.Description,
					IsActive:    item.Product.Category.IsActive,
				},
			},
			Quantity: item.Quantity,
			Price:    item.Price,
		})
	}

	return dto.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Status:      string(order.Status),
		TotalAmount: order.TotalAmount,
		OrderItems:  orderItems,
		CreatedAt:   order.CreatedAt.Format(dateFormat),
	}
}
