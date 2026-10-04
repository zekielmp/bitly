package server

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zekielmp/Bitly/internal/dto"
	"github.com/zekielmp/Bitly/internal/utils"
)

func (s *Server) getCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	res, err := s.cart.GetCart(userID)
	if err != nil {
		utils.NotFoundResponse(c, "Cart not found", err)
		return
	}

	utils.SuccessResponse(c, "Cart retrieved successfully", res)

}

func (s *Server) addToCart(c *gin.Context) {
	user := c.GetUint("user_id")
	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid request data", err)
		return
	}
	res, err := s.cart.AddToCart(user, &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to add item to cart", err)
		return
	}
	utils.SuccessResponse(c, "Product added to cart successfully", res)
}

func (s *Server) updateCartItem(c *gin.Context) {
	user := c.GetUint("user_id") /*user_id*/

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid cart item ID", err)
		return
	}

	var req dto.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	res, err := s.cart.UpdateCartItem(user, uint(id), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update cart item", err)
		return
	}

	utils.SuccessResponse(c, "Cart item updated successfully", res)
}

func (s *Server) removeFromCart(c *gin.Context) {
	user := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid cart item ID", err)
		return
	}

	if err := s.cart.RemoveFromCart(user, uint(id)); err != nil {
		utils.InternalServerErrorResponse(c, "Failed to remove item from cart", err)
		return
	}
	utils.SuccessResponse(c, "Item removed successfully", nil)
}
