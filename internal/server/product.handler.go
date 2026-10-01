package server

import (
	"strconv"

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

	service := services.NewProductService(s.db)

	res, err := service.CreateCategory(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create category", err)
		return
	}

	utils.CreatedResponse(c, "Category created successfully", res)
}

func (s *Server) getCategories(c *gin.Context) {
	service := services.NewProductService(s.db)
	res, err := service.GetCategories()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get products", err)
		return
	}
	utils.SuccessResponse(c, "Product retrieved Successfully", res)
}

func (s *Server) updateCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid Category ID", err)
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	service := services.NewProductService(s.db)
	res, err := service.UpdateCategory(uint(id), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update category", err)
		return
	}
	utils.SuccessResponse(c, "Category updated successfully", res)
}

func (s *Server) deleteCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid Category ID", err)
	}

	service := services.NewProductService(s.db)
	if err := service.DeleteCategory(uint(id)); err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete category", err)
		return
	}

	utils.SuccessResponse(c, "Category deleted successfully", nil)
}

func (s *Server) addProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	service := services.NewProductService(s.db)
	res, err := service.AddProduct(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to add product", err)
		return
	}

	utils.CreatedResponse(c, "Product added successfully", res)
}
