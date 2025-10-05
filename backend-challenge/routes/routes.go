package routes

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	ginZap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"oolio.com/kart/configs"
	"oolio.com/kart/constants"
	"oolio.com/kart/controllers"
	"oolio.com/kart/docs"
	"oolio.com/kart/middlewares"
	"oolio.com/kart/repositories"
	"oolio.com/kart/services"
	"time"
)

// InitializeRoutes initializes the routes for the application.
func InitializeRoutes() *gin.Engine {
	if configs.ReleaseEnv == constants.ProdMode {
		gin.SetMode(gin.ReleaseMode)
	}

	initSwagger()

	router := gin.New()

	router.Use(cors.Default())

	router.Use(ginZap.RecoveryWithZap(configs.Logger, true))
	router.Use(ginZap.Ginzap(configs.Logger, time.RFC3339, false))

	pprof.Register(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	kartRouter := router.Group("/api")
	kartRouter.GET("/health", controllers.HealthCheckController.HealthCheck)

	pool, err := repositories.Pool()
	if err != nil {
		configs.Logger.Fatal("Failed to get database connection pool", zap.Error(err))
	}

	couponRepo := repositories.NewCouponRepositoryImpl(pool)

	if err := services.InitializeCouponService(couponRepo); err != nil {
		configs.Logger.Fatal("Failed to initialize coupon service", zap.Any("error", err))
	}

	productRepository := repositories.NewProductRepositoryImpl(pool)
	orderRepository := repositories.NewOrderRepositoryImpl(pool)

	productService := services.NewProductServiceImpl(productRepository)
	orderService := services.NewOrderServiceImpl(orderRepository, productRepository, services.CouponServiceImpl)

	productController := controllers.NewProductController(productService)
	orderController := controllers.NewOrderController(orderService)

	product := kartRouter.Group("/product")
	product.GET("", productController.GetProducts)
	product.GET("/:productId", productController.GetProductById)

	kartRouter.POST("/order", middlewares.APIKeyMiddleware(), orderController.PlaceOrder)

	return router
}

func initSwagger() {
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = "Kart API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "E-commerce API for managing products and orders"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", configs.Host, configs.Port)
}
