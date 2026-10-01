package server

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zekielmp/Bitly/internal/dto"
	"github.com/zekielmp/Bitly/internal/utils"
)

func (s *Server) createCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	res, err := s.product.CreateCategory(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create category", err)
		return
	}

	utils.CreatedResponse(c, "Category created successfully", res)
}

func (s *Server) getCategories(c *gin.Context) {
	res, err := s.product.GetCategories()
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

	res, err := s.product.UpdateCategory(uint(id), &req)
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
		return
	}

	if err := s.product.DeleteCategory(uint(id)); err != nil {
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

	res, err := s.product.AddProduct(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to add product", err)
		return
	}

	utils.CreatedResponse(c, "Product added successfully", res)
}

func (s *Server) getProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	res, meta, err := s.product.GetProducts(page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get products", err)
		return
	}

	utils.PaginatedSuccessResponse(c, "Products retrieved successfully", res, *meta)

}

func (s *Server) getProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid Product ID", err)
		return
	}

	res, err := s.product.GetProduct(uint(id))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get product", err)
		return
	}

	utils.SuccessResponse(c, "Product retrieved successfully", res)
}

func (s *Server) updateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid Product ID", err)
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	res, err := s.product.UpdateProduct(&req, uint(id))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update product", err)
		return
	}
	utils.SuccessResponse(c, "Product updated successfully", res)
}

func (s *Server) deleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid Product ID", err)
		return
	}

	if err := s.product.DeleteProduct(uint(id)); err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete product", err)
		return
	}

	utils.SuccessResponse(c, "Product deleted successfully", nil)
}
