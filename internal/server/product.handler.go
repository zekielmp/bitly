package server

import (
	"github.com/gin-gonic/gin"
	"github.com/zekielmp/Bitly/internal/dto"
	"github.com/zekielmp/Bitly/internal/services"
	"github.com/zekielmp/Bitly/internal/utils"
)

func (s *Server) createCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	productService := services.NewProductService(s.db)

	res, err := productService.CreateCategory(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to create category", err)
		return
	}

	utils.CreatedResponse(c, "Category created successfully", res)
}

func (s *Server) addProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	productService := services.NewProductService(s.db)
	res, err := productService.AddProduct(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to add product", err)
		return
	}

	utils.CreatedResponse(c, "Product added successfully", res)
}
