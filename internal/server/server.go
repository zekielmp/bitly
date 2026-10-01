package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/zekielmp/Bitly/internal/config"
	"github.com/zekielmp/Bitly/internal/services"
	"gorm.io/gorm"
)

// Server struct
type Server struct {
	config  *config.Config
	db      *gorm.DB
	logger  *zerolog.Logger
	auth    *services.AuthService
	user    *services.UserService
	product *services.ProductServices
	upload  *services.UploadService
}

// New creates an instance of Server
func New(cfg *config.Config,
	db *gorm.DB, logger *zerolog.Logger,
	auth *services.AuthService,
	product *services.ProductServices,
	user *services.UserService,
	upload *services.UploadService) *Server {
	return &Server{
		config:  cfg,
		db:      db,
		logger:  logger,
		auth:    auth,
		user:    user,
		product: product,
		upload:  upload,
	}
}

// SetupRoute initializes routes for Server
func (s *Server) SetupRoute() *gin.Engine {
	router := gin.New()

	/* Add middlewares*/
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corMiddleware())

	/* Add route */
	router.GET("/health", s.healthCheck)

	router.Static("/uploads", "./uploads")

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", s.register)
			auth.POST("/login", s.login)
			auth.POST("/refresh", s.refreshToken)
			auth.POST("/logout", s.logout)

		}
		protected := api.Group("/")
		protected.Use(s.authMiddleware())
		{
			/*user routes*/
			users := protected.Group("/users")
			{
				users.GET("/profile", s.getprofile)
				users.PUT("/profile", s.updateProfile)
			}
			/*category routes*/
			category := protected.Group("/categories")
			{
				category.POST("/", s.adminMiddleware(), s.createCategory)
				category.PUT("/:id", s.adminMiddleware(), s.updateCategory)
				category.DELETE("/:id", s.adminMiddleware(), s.deleteCategory)
			}
			/*product routes*/
			product := protected.Group("/products")
			{
				product.POST("/", s.adminMiddleware(), s.addProduct)
				product.PUT("/:id", s.adminMiddleware(), s.updateProduct)
				product.DELETE("/:id", s.adminMiddleware(), s.deleteProduct)
				product.POST("/:id/images", s.adminMiddleware(), s.uploadProductImages)

			}
		}
		{
		}

		/*public routes*/
		public := api.Group("/public")
		{
			public.GET("/categories", s.getCategories)
			public.GET("/products", s.getProducts)
			public.GET("/products/:id", s.getProduct)
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
