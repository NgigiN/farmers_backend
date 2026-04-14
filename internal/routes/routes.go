package routes

import (
	"net/http"
	"strings"
	"time"

	"farm-backend/internal/config"
	activityHandlers "farm-backend/internal/handlers/activities"
	animalHandlers "farm-backend/internal/handlers/animals"
	authHandlers "farm-backend/internal/handlers/auth"
	inputHandlers "farm-backend/internal/handlers/inputs"
	plantHandlers "farm-backend/internal/handlers/plants"
	summaryHandlers "farm-backend/internal/handlers/summaries"
	userHandlers "farm-backend/internal/handlers/users"
	"farm-backend/internal/middleware"
	animalServices "farm-backend/internal/services/animals"
	authService "farm-backend/internal/services/auth"
	plantServices "farm-backend/internal/services/plants"
	summaryServices "farm-backend/internal/services/summaries"
	userServices "farm-backend/internal/services/users"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.LoggingMiddleware())

	// Build CORS allowed-origins from config.
	// Set ALLOWED_ORIGINS in .env, e.g.: ALLOWED_ORIGINS=http://localhost:3000,https://app.mysite.com
	allowedOrigins := []string{"http://localhost:3000"} // safe default
	if cfg.AllowedOrigins != "" {
		allowedOrigins = strings.Split(cfg.AllowedOrigins, ",")
	}
	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Auth Middleware
	authMiddleware := middleware.AuthMiddleware(cfg)
	// Rate limiter: max 10 requests per minute per IP on auth endpoints
	authRateLimiter := middleware.AuthRateLimiter(10, time.Minute)

	// ── Services ────────────────────────────────────────────────────────────
	// Auth
	svcAuth := authService.NewService(db, cfg)

	// Plant
	plantService := plantServices.NewPlantService(db)
	landService := plantServices.NewLandService(db)
	seasonService := plantServices.NewSeasonService(db)
	activityService := plantServices.NewActivityService(db)
	inputService := plantServices.NewInputService(db)

	// Animal
	animalTypeService := animalServices.NewAnimalTypeService(db)
	herdService := animalServices.NewHerdService(db)
	animalService := animalServices.NewAnimalService(db)
	infrastructureService := animalServices.NewInfrastructureService(db)

	// Summary
	costCategoryService := summaryServices.NewCostCategoryService(db)
	costService := summaryServices.NewCostService(db)
	revenueService := summaryServices.NewRevenueService(db)
	analysisService := summaryServices.NewAnalysisService(db)
	userService := userServices.NewUserService(db)

	// ── Handlers ────────────────────────────────────────────────────────────
	authHandler := authHandlers.NewAuthHandler(svcAuth)

	plantHandler := plantHandlers.NewPlantHandler(plantService)
	landHandler := plantHandlers.NewLandHandler(landService)
	seasonHandler := plantHandlers.NewSeasonHandler(seasonService)

	activityHandler := activityHandlers.NewActivityHandler(activityService)
	inputHandler := inputHandlers.NewInputHandler(inputService)

	animalTypeHandler := animalHandlers.NewAnimalTypeHandler(animalTypeService)
	herdHandler := animalHandlers.NewHerdHandler(herdService)
	animalHandler := animalHandlers.NewAnimalHandler(animalService)
	infrastructureHandler := animalHandlers.NewInfrastructureHandler(infrastructureService)

	costCategoryHandler := summaryHandlers.NewCostCategoryHandler(costCategoryService)
	costHandler := summaryHandlers.NewCostHandler(costService)
	revenueHandler := summaryHandlers.NewRevenueHandler(revenueService)
	analysisHandler := summaryHandlers.NewAnalysisHandler(analysisService)
	userHandler := userHandlers.NewUserHandler(userService)

	// ── Routes ──────────────────────────────────────────────────────────────
	api := router.Group("/api/v1")

	// Health check (public)
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Auth (public + rate-limited)
	auth := api.Group("/auth")
	auth.Use(authRateLimiter)
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(authMiddleware)
	{
		// ── Profile ──
		profile := protected.Group("/profile")
		{
			profile.GET("", userHandler.GetProfile)
			profile.PUT("", userHandler.UpdateProfile)
			profile.PUT("/password", userHandler.ChangePassword)
		}

		// ── Plants ──
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

		// ── Animals ──
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

		// ── Shared (plant & animal) ──
		// ?source_type=plant|animal
		activities := protected.Group("/activities")
		{
			activities.POST("", activityHandler.CreateActivity)
			activities.GET("", activityHandler.ListActivities)
			activities.GET("/:id", activityHandler.GetActivity)
			activities.PUT("/:id", activityHandler.UpdateActivity)
			activities.DELETE("/:id", activityHandler.DeleteActivity)
		}

		// ?source_type=plant|animal
		inputs := protected.Group("/inputs")
		{
			inputs.POST("", inputHandler.CreateInput)
			inputs.GET("", inputHandler.ListInputs)
			inputs.GET("/:id", inputHandler.GetInput)
			inputs.PUT("/:id", inputHandler.UpdateInput)
			inputs.DELETE("/:id", inputHandler.DeleteInput)
		}

		// ── Costs ──
		// ?type=plant|animal&category=input|activity
		costCategories := protected.Group("/cost-categories")
		{
			costCategories.POST("", costCategoryHandler.CreateCostCategory)
			costCategories.GET("", costCategoryHandler.ListCostCategories)
			costCategories.GET("/:id", costCategoryHandler.GetCostCategory)
			costCategories.PUT("/:id", costCategoryHandler.UpdateCostCategory)
			costCategories.DELETE("/:id", costCategoryHandler.DeleteCostCategory)
		}

		// ── Revenue ──
		// ?source=plant|animal  or  ?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
		revenue := protected.Group("/revenue")
		{
			revenue.POST("", revenueHandler.CreateRevenue)
			revenue.GET("", revenueHandler.ListRevenues)
			revenue.GET("/:id", revenueHandler.GetRevenue)
			revenue.PUT("/:id", revenueHandler.UpdateRevenue)
			revenue.DELETE("/:id", revenueHandler.DeleteRevenue)
		}

		// ── Analytics ──
		analysis := protected.Group("/analytics")
		{
			analysis.GET("/total-costs", analysisHandler.GetTotalCosts)
			analysis.GET("/total-costs-by-season", costHandler.GetTotalCostsBySeason)
			analysis.GET("/total-revenue", analysisHandler.GetTotalRevenue)
			analysis.GET("/profit", analysisHandler.GetProfit)
			analysis.GET("/cost-breakdown", analysisHandler.GetCostBreakdown)
			analysis.GET("/revenue-breakdown", analysisHandler.GetRevenueBreakdown)
			analysis.GET("/monthly-summary", analysisHandler.GetMonthlySummary) // ?year=YYYY
		}
	}

	return router
}
