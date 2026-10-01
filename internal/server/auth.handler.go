package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zekielmp/Bitly/internal/dto"
	"github.com/zekielmp/Bitly/internal/utils"
)

func (s *Server) register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid request data", err)
		return
	}

	response, err := s.auth.Register(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Registration failed", err)
		return
	}
	utils.CreatedResponse(c, "User registered successfully", response)

}

func (s *Server) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid request data", err)
		return
	}

	response, err := s.auth.Login(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "Login failed", err)
		return
	}
	utils.SuccessResponse(c, "Login successful", response)
}

func (s *Server) refreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	response, err := s.auth.RefreshToken(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Token refresh failed", err)
		return
	}
	utils.SuccessResponse(c, "Token refreshed successfully", response)
}

func (s *Server) logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	if err := s.auth.Logout(req.RefreshToken); err != nil {
		utils.InternalServerErrorResponse(c, "Logout Failed", err)
		return
	}
	utils.SuccessResponse(c, "Logout successful", nil)
}

func (s *Server) getprofile(c *gin.Context) {
	userID := c.GetUint("user_id")
	log.Println(userID)

	profile, err := s.user.GetProfile(userID)
	if err != nil {
		utils.NotFoundResponse(c, "User not found", err)
		return
	}
	utils.SuccessResponse(c, "Profile retrieved successfully", profile)
}

func (s *Server) updateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid request data", err)
		return
	}

	profile, err := s.user.UpdateProfile(userID, &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update profile", err)
		return
	}
	utils.SuccessResponse(c, "Profile retrieved successfully", profile)
}
