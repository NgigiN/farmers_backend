package routes

import (
	"farm-backend/internal/config"
	activityHandlers "farm-backend/internal/handlers/activities"
	animalHandlers "farm-backend/internal/handlers/animals"
	authHandlers "farm-backend/internal/handlers/auth"
	inputHandlers "farm-backend/internal/handlers/inputs"
	plantHandlers "farm-backend/internal/handlers/plants"
	summaryHandlers "farm-backend/internal/handlers/summaries"
	"farm-backend/internal/middleware"
	animalServices "farm-backend/internal/services/animals"
	plantServices "farm-backend/internal/services/plants"
	summaryServices "farm-backend/internal/services/summaries"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize Auth Middleware
	authMiddleware := middleware.AuthMiddleware(cfg)

	// Plant Services
	plantService := plantServices.NewPlantService(db)
	landService := plantServices.NewLandService(db)
	seasonService := plantServices.NewSeasonService(db)

	activityService := plantServices.NewActivityService(db)
	inputService := plantServices.NewInputService(db)

	// Animal Services
	animalTypeService := animalServices.NewAnimalTypeService(db)
	herdService := animalServices.NewHerdService(db)
	animalService := animalServices.NewAnimalService(db)
	infrastructureService := animalServices.NewInfrastructureService(db)

	// Summary Services
	costCategoryService := summaryServices.NewCostCategoryService(db)
	revenueService := summaryServices.NewRevenueService(db)
	analysisService := summaryServices.NewAnalysisService(db)

	// Auth Service
	authService := plantServices.NewAuthService(db, cfg)

	// Auth Handler
	authHandler := authHandlers.NewAuthHandler(authService)

	// Plant Handlers
	plantHandler := plantHandlers.NewPlantHandler(plantService)
	landHandler := plantHandlers.NewLandHandler(landService)
	seasonHandler := plantHandlers.NewSeasonHandler(seasonService)

	// Unified Handlers
	activityHandler := activityHandlers.NewActivityHandler(activityService)
	inputHandler := inputHandlers.NewInputHandler(inputService)

	// Animal Handlers
	animalTypeHandler := animalHandlers.NewAnimalTypeHandler(animalTypeService)
	herdHandler := animalHandlers.NewHerdHandler(herdService)
	animalHandler := animalHandlers.NewAnimalHandler(animalService)
	infrastructureHandler := animalHandlers.NewInfrastructureHandler(infrastructureService)

	// Summary Handlers
	costCategoryHandler := summaryHandlers.NewCostCategoryHandler(costCategoryService)
	revenueHandler := summaryHandlers.NewRevenueHandler(revenueService)
	analysisHandler := summaryHandlers.NewAnalysisHandler(analysisService)

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
	}

	protected := api.Group("")
	protected.Use(authMiddleware)
	{
		// Plant Routes
		plants := protected.Group("/plants")
		{
			plants.POST("", plantHandler.CreatePlant)
			plants.GET("", plantHandler.ListPlants)
			plants.GET("/:id", plantHandler.GetPlant)
			plants.PUT("/:id", plantHandler.UpdatePlant)
			plants.DELETE("/:id", plantHandler.DeletePlant)
		}

		lands := protected.Group("/lands")
		{
			lands.POST("", landHandler.CreateLand)
			lands.GET("", landHandler.ListLands)
			lands.GET("/:id", landHandler.GetLand)
			lands.PUT("/:id", landHandler.UpdateLand)
			lands.DELETE("/:id", landHandler.DeleteLand)
		}

		seasons := protected.Group("/seasons")
		{
			seasons.POST("", seasonHandler.CreateSeason)
			seasons.GET("", seasonHandler.ListSeasons)
			seasons.GET("/:id", seasonHandler.GetSeason)
			seasons.PUT("/:id", seasonHandler.UpdateSeason)
			seasons.DELETE("/:id", seasonHandler.DeleteSeason)
		}

		// Animal Routes
		animalTypes := protected.Group("/animal-types")
		{
			animalTypes.POST("", animalTypeHandler.CreateAnimalType)
			animalTypes.GET("", animalTypeHandler.ListAnimalTypes)
			animalTypes.GET("/:id", animalTypeHandler.GetAnimalType)
			animalTypes.PUT("/:id", animalTypeHandler.UpdateAnimalType)
			animalTypes.DELETE("/:id", animalTypeHandler.DeleteAnimalType)
		}

		herds := protected.Group("/herds")
		{
			herds.POST("", herdHandler.CreateHerd)
			herds.GET("", herdHandler.ListHerds)
			herds.GET("/:id", herdHandler.GetHerd)
			herds.PUT("/:id", herdHandler.UpdateHerd)
			herds.DELETE("/:id", herdHandler.DeleteHerd)
		}

		animals := protected.Group("/animals")
		{
			animals.POST("", animalHandler.CreateAnimal)
			animals.GET("", animalHandler.ListAnimals)
			animals.GET("/:id", animalHandler.GetAnimal)
			animals.PUT("/:id", animalHandler.UpdateAnimal)
			animals.DELETE("/:id", animalHandler.DeleteAnimal)
		}

		infrastructure := protected.Group("/infrastructure")
		{
			infrastructure.POST("", infrastructureHandler.CreateInfrastructure)
			infrastructure.GET("", infrastructureHandler.ListInfrastructures)
			infrastructure.GET("/:id", infrastructureHandler.GetInfrastructure)
			infrastructure.PUT("/:id", infrastructureHandler.UpdateInfrastructure)
			infrastructure.DELETE("/:id", infrastructureHandler.DeleteInfrastructure)
		}

		// Unified Routes (Plant & Animal)
		activities := protected.Group("/activities")
		{
			activities.POST("", activityHandler.CreateActivity)
			activities.GET("", activityHandler.ListActivities) // plant|animal
			activities.GET("/:id", activityHandler.GetActivity)
			activities.PUT("/:id", activityHandler.UpdateActivity)
			activities.DELETE("/:id", activityHandler.DeleteActivity)
		}

		inputs := protected.Group("/inputs")
		{
			inputs.POST("", inputHandler.CreateInput)
			inputs.GET("", inputHandler.ListInputs) // plant|animal
			inputs.GET("/:id", inputHandler.GetInput)
			inputs.PUT("/:id", inputHandler.UpdateInput)
			inputs.DELETE("/:id", inputHandler.DeleteInput)
		}

		// Summary Routes
		costCategories := protected.Group("/cost-categories")
		{
			costCategories.POST("", costCategoryHandler.CreateCostCategory)
			costCategories.GET("", costCategoryHandler.ListCostCategories) // type=plant|animal&category=input|activity
			costCategories.GET("/:id", costCategoryHandler.GetCostCategory)
			costCategories.PUT("/:id", costCategoryHandler.UpdateCostCategory)
			costCategories.DELETE("/:id", costCategoryHandler.DeleteCostCategory)
		}

		revenue := protected.Group("/revenue")
		{
			revenue.POST("", revenueHandler.CreateRevenue)
			revenue.GET("", revenueHandler.ListRevenues) // source=plant|animal&start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
			revenue.GET("/:id", revenueHandler.GetRevenue)
			revenue.PUT("/:id", revenueHandler.UpdateRevenue)
			revenue.DELETE("/:id", revenueHandler.DeleteRevenue)
		}

		analysis := protected.Group("/analysis")
		{
			analysis.GET("/total-costs", analysisHandler.GetTotalCosts)
			analysis.GET("/total-revenue", analysisHandler.GetTotalRevenue)
			analysis.GET("/profit", analysisHandler.GetProfit)
			analysis.GET("/cost-breakdown", analysisHandler.GetCostBreakdown)
			analysis.GET("/revenue-breakdown", analysisHandler.GetRevenueBreakdown)
			analysis.GET("/monthly-summary", analysisHandler.GetMonthlySummary) // year=YYYY
		}
	}

	return router
}
