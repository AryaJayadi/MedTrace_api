package main

import (
	"log"
	"os"
	"strings"

	"github.com/AryaJayadi/MedTrace_api/internal/auth"
	"github.com/AryaJayadi/MedTrace_api/internal/handlers"
	"github.com/AryaJayadi/MedTrace_api/internal/services"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Info: Error loading .env file (this is not fatal if environment variables are set directly):", err)
	} else {
		log.Println("Successfully loaded .env file")
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsEnv != "" {
		allowedOrigins = append(allowedOrigins, strings.Split(allowedOriginsEnv, ",")...)
	} else {
		allowedOrigins = []string{"http://localhost:5173"}
		log.Println("ALLOWED_ORIGINS not set in environment, using default:", allowedOrigins)
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowedOrigins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	chaincodeName := os.Getenv("CHAINCODE_NAME")
	if chaincodeName == "" {
		log.Println("CHAINCODE_NAME not set in environment, using default from auth package.")
		chaincodeName = auth.DefaultChaincodeName
	}

	channelName := os.Getenv("CHANNEL_NAME")
	if channelName == "" {
		log.Println("CHANNEL_NAME not set in environment, using default from auth package.")
		channelName = auth.DefaultChannelName
	}
	log.Printf("Using Chaincode: %s, Channel: %s", chaincodeName, channelName)

	// Services are instantiated without a contract. The contract will be passed per method.
	organizationService := services.NewOrganizationService() // Adjusted constructor
	batchService := services.NewBatchService()               // Adjusted constructor
	drugService := services.NewDrugService()                 // Adjusted constructor
	transferService := services.NewTransferService()         // Adjusted constructor
	ledgerService := services.NewLedgerService()             // Adjusted constructor

	// Handlers are instantiated with services.
	organizationHandler := handlers.NewOrganizationHandler(organizationService)
	batchHandler := handlers.NewBatchHandler(batchService)
	drugHandler := handlers.NewDrugHandler(drugService)
	transferHandler := handlers.NewTransferHandler(transferService)
	ledgerHandler := handlers.NewLedgerHandler(ledgerService)

	// --- Public Routes ---
	e.POST("/login", auth.LoginHandler)
	e.POST("/logout", auth.LogoutHandler) // Or GET, but POST is often preferred for logout
	e.POST("/refresh", auth.RefreshTokenHandler)

	// --- Protected Route Groups ---
	// These groups will use the AuthMiddleware to ensure a valid JWT and set up the Fabric context.

	orgGroup := e.Group("/organizations", auth.AuthMiddleware)
	orgGroup.GET("", organizationHandler.GetOrganizations)
	orgGroup.GET("/:id", organizationHandler.GetOrganizationByID)

	batchesGroup := e.Group("/batches", auth.AuthMiddleware)
	batchesGroup.POST("", batchHandler.CreateBatch)
	batchesGroup.GET("", batchHandler.GetAllBatches)
	batchesGroup.GET("/:id/exists", batchHandler.BatchExists)
	batchesGroup.GET("/:id", batchHandler.GetBatchByID)
	batchesGroup.PATCH("/batches", batchHandler.UpdateBatch) // Note: Route was /batches, should it be /:id?

	ledgerGroup := e.Group("/ledger", auth.AuthMiddleware)
	ledgerGroup.POST("/init", ledgerHandler.InitLedger)

	drugsGroup := e.Group("/drugs", auth.AuthMiddleware)
	drugsGroup.POST("", drugHandler.CreateDrug)
	drugsGroup.GET("/my", drugHandler.GetMyDrugs)
	drugsGroup.GET("/my/available", drugHandler.GetMyAvailDrugs)
	drugsGroup.GET("/:drugID", drugHandler.GetDrug)
	drugsGroup.GET("/batch/:batchID", drugHandler.GetDrugByBatch)

	transferGroup := e.Group("/transfers", auth.AuthMiddleware)
	transferGroup.POST("", transferHandler.CreateTransfer)
	transferGroup.GET("/my", transferHandler.GetMyTransfers)
	transferGroup.GET("/my/outgoing", transferHandler.GetMyOutTransfer)
	transferGroup.GET("/my/incoming", transferHandler.GetMyInTransfer)
	transferGroup.POST("/accept", transferHandler.AcceptTransfer)
	transferGroup.POST("/reject", transferHandler.RejectTransfer)
	transferGroup.GET("/:id", transferHandler.GetTransfer)

	port := os.Getenv("API_PORT")
	if port == "" {
		log.Println("API_PORT not set in environment, using default 8080")
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
