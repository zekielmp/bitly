package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/zekielmp/Bitly/internal/config"
	"gorm.io/gorm"
)

// Server struct
type Server struct {
	config *config.Config
	db     *gorm.DB
	logger *zerolog.Logger
}

// New creates an instance of Server
func New(cfg *config.Config, db *gorm.DB, logger *zerolog.Logger) *Server {
	return &Server{
		config: cfg,
		db:     db,
		logger: logger,
	}
}

// SetupRoute initializes routes for Server
func (s *Server) SetupRoute() *gin.Engine {
	router := gin.New()

	/* Add middlewares*/
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corMiddleware())
	// router.Use(s.adminMiddleware())

	/* Add route */
	router.GET("/health", s.healthCheck)

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", s.register)
			auth.POST("/login", s.login)
			auth.POST("/refresh", s.refreshToken)
			auth.POST("/logout", s.logout)

		}

	}

	return router
}
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) corMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type,Authorization, X-PIN")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	}
}

// func (s *Server) cors(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET,POST, PUT,DELETE,OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization, X-PIN")

// 		if r.Method == "OPTIONS" {
// 			r.Response.StatusCode = 204
// 			return
// 		}

// 	})
// }
