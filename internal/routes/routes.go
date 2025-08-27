package routes

import (
	"LojaGin/internal/middleware"
	"LojaGin/internal/modules/auth"
	"LojaGin/internal/modules/cart"
	"LojaGin/internal/modules/category"
	"LojaGin/internal/modules/product"
	"LojaGin/internal/modules/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupAuthRoutes(api *gin.RouterGroup, db *gorm.DB) {
	userRepo := user.NewRepository(db)
	authService := auth.NewService(userRepo)
	authHandler := auth.NewHandler(authService)

	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}
}

func setupUserRoutes(api *gin.RouterGroup, db *gorm.DB) {
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	userRoutes := api.Group("/users").Use(middleware.AuthMiddleware())
	{
		userRoutes.GET("/", userHandler.GetAllUsers)
		userRoutes.GET("/:id", userHandler.GetUserByID)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
	}
}

func setupCategoryRoutes(api *gin.RouterGroup, db *gorm.DB) {
	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService)

	categoryRoutes := api.Group("/categories")
	{
		// Rotas públicas
		categoryRoutes.GET("/", categoryHandler.GetAllCategories)
		categoryRoutes.GET("/:id", categoryHandler.GetCategoryByID)

		// Rotas protegidas
		protectedCategoryRoutes := categoryRoutes.Group("/").Use(middleware.AuthMiddleware())
		protectedCategoryRoutes.POST("/", categoryHandler.CreateCategory)
		protectedCategoryRoutes.PUT("/:id", categoryHandler.UpdateCategory)
		protectedCategoryRoutes.DELETE("/:id", categoryHandler.DeleteCategory)
	}
}

func setupProductRoutes(api *gin.RouterGroup, db *gorm.DB) {
	productRepo := product.NewRepository(db)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)

	productRoutes := api.Group("/products")
	{
		// Rotas públicas
		productRoutes.GET("/", productHandler.GetAllProducts)
		productRoutes.GET("/:id", productHandler.GetProductByID)

		// Rotas protegidas
		protectedProductRoutes := productRoutes.Group("/").Use(middleware.AuthMiddleware())
		protectedProductRoutes.POST("/", productHandler.CreateProduct)
		protectedProductRoutes.PUT("/:id", productHandler.UpdateProduct)
		protectedProductRoutes.DELETE("/:id", productHandler.DeleteProduct)
	}
}

func setupCartRoutes(api *gin.RouterGroup, db *gorm.DB) {
	cartRepo := cart.NewRepository(db)
	cartService := cart.NewService(cartRepo)
	cartHandler := cart.NewHandler(cartService)

	cartRoutes := api.Group("/cart").Use(middleware.AuthMiddleware())
	{
		cartRoutes.GET("/", cartHandler.GetCart)
		cartRoutes.POST("/add", cartHandler.AddToCart)
		cartRoutes.POST("/remove", cartHandler.RemoveFromCart)
		cartRoutes.DELETE("/clear", cartHandler.ClearCart)
		cartRoutes.POST("/checkout", cartHandler.Checkout)
	}
}

func SetupAPI(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api")

	setupAuthRoutes(api, db)
	setupUserRoutes(api, db)
	setupCategoryRoutes(api, db)
	setupProductRoutes(api, db)
	setupCartRoutes(api, db)
}
