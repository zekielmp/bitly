package server

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zekielmp/Bitly/internal/utils"
)

func (s *Server) createOrder(c *gin.Context) {
	user := c.GetUint("user_id")

	res, err := s.order.CreateOrder(user)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to create order", err)
		return
	}
	utils.CreatedResponse(c, "Order created successfully", res)
}

func (s *Server) getOrders(c *gin.Context) {
	user := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	res, meta, err := s.order.GetOrders(user, page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve orders", err)
		return
	}
	utils.PaginatedSuccessResponse(c, "Orders retrieved successfully", res, *meta)

}
func (s *Server) getOrder(c *gin.Context) {
	user := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid order ID ", err)
		return
	}

	res, err := s.order.GetOrder(user, uint(id))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Orders not found", err)
		return
	}
	utils.SuccessResponse(c, "Order retrieved successfully", res)
}
