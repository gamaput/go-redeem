package route

import (
	"fmt"
	"log"
	"time"

	"github.com/gamaput/go-redeem/controller"
	"github.com/gamaput/go-redeem/middleware"
	"github.com/gamaput/go-redeem/repository"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes : all the routes are defined here
func SetupRoutes(db *gorm.DB) {
	httpRouter := gin.Default()

	httpRouter.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	// Initialize  casbin adapter
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize casbin adapter: %v", err))
	}

	// Load model configuration file and policy store adapter
	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		panic(fmt.Sprintf("failed to create casbin enforcer: %v", err))
	}

	//add policy
	if hasPolicy := enforcer.HasPolicy("admin", "report", "read"); !hasPolicy {
		enforcer.AddPolicy("admin", "report", "read")
	}
	if hasPolicy := enforcer.HasPolicy("admin", "report", "write"); !hasPolicy {
		enforcer.AddPolicy("admin", "report", "write")
	}
	if hasPolicy := enforcer.HasPolicy("user", "report", "read"); !hasPolicy {
		enforcer.AddPolicy("user", "report", "read")
	}

	userRepository := repository.NewUserRepository(db)
	productRepository := repository.NewProductRepository(db)
	redeemCodeRepository := repository.NewRedeemCodeRepository(db)
	prizeCodeRepository := repository.NewPrizeRepository(db)

	if err := userRepository.Migrate(); err != nil {
		log.Fatal("User migrate err", err)
	}

	userController := controller.NewUserController(userRepository)
	productController := controller.NewProductController(productRepository)
	redeemController := controller.NewRedeemCodeController(redeemCodeRepository, prizeCodeRepository)
	prizeController := controller.NewPrizeController(prizeCodeRepository)

	apiRoutes := httpRouter.Group("/api")

	{
		apiRoutes.POST("/register", userController.AddUser(enforcer))
		apiRoutes.POST("/signin", userController.SignInUser)
		apiRoutes.GET("/logout", userController.Logout)
		apiRoutes.POST("/redeem", redeemController.RedeemCode)
		apiRoutes.GET("/rand-prize", prizeController.GetRandomPrize)
	}

	userProtectedRoutes := apiRoutes.Group("/users", middleware.AuthorizeJWT())
	{
		userProtectedRoutes.GET("/", middleware.Authorize("report", "read", enforcer), userController.GetAllUser)
		userProtectedRoutes.POST("/add", middleware.Authorize("report", "write", enforcer), userController.AddUser(enforcer))
		userProtectedRoutes.GET("/:user", middleware.Authorize("report", "read", enforcer), userController.GetUser)

		userProtectedRoutes.PATCH("/:user", middleware.Authorize("report", "write", enforcer), userController.UpdateUser)

		userProtectedRoutes.DELETE("/:user", middleware.Authorize("report", "write", enforcer), userController.DeleteUser)

	}

	productProductedRoutes := apiRoutes.Group("/products", middleware.AuthorizeJWT())
	{
		productProductedRoutes.GET("/", middleware.Authorize("report", "read", enforcer), productController.GetAllProducts)
		productProductedRoutes.POST("/add", middleware.Authorize("report", "read", enforcer), productController.CreateProduct(enforcer))
		productProductedRoutes.GET("/:product", middleware.Authorize("report", "read", enforcer), productController.GetProductByID)
		productProductedRoutes.PATCH("/:product", middleware.Authorize("report", "write", enforcer), productController.UpdateProduct)
		productProductedRoutes.DELETE("/:product", middleware.Authorize("report", "write", enforcer), productController.DeleteProduct)

	}
	redeemCodeRoutes := apiRoutes.Group("/voucher", middleware.AuthorizeJWT())
	{
		redeemCodeRoutes.GET("/", middleware.Authorize("report", "read", enforcer), redeemController.GetAllRedeems)
		redeemCodeRoutes.GET("/generate-code", middleware.Authorize("report", "write", enforcer), redeemController.GenerateCode)
		// redeemCodeRoutes.POST("/redeem", middleware.Authorize("report", "read", enforcer), redeemController.RedeemCode)
	}
	prizeCodeRoutes := apiRoutes.Group("/prizes", middleware.AuthorizeJWT())
	{
		prizeCodeRoutes.POST("/add", middleware.Authorize("report", "write", enforcer), prizeController.CreatePrize)
		prizeCodeRoutes.GET("/", middleware.Authorize("report", "write", enforcer), prizeController.GetAllPrizes)
		prizeCodeRoutes.DELETE("/:prize", middleware.Authorize("report", "write", enforcer), prizeController.DeletePrize)
		prizeCodeRoutes.PATCH("/:prize", middleware.Authorize("report", "write", enforcer), prizeController.UpdatePrize)
		prizeCodeRoutes.GET("/:prize", middleware.Authorize("report", "write", enforcer), prizeController.GetPrizeByID)
	}
	httpRouter.Run(":8081")

}
