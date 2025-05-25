package main

import (
	"log"
	"os"

	// "github.com/AryaJayadi/MedTrace_api/cmd/fabric" // No longer directly used for global setup
	"github.com/AryaJayadi/MedTrace_api/internal/auth"
	"github.com/AryaJayadi/MedTrace_api/internal/handlers"
	"github.com/AryaJayadi/MedTrace_api/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Global constants for Fabric setup are removed, as this is now dynamic.

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},                                                                // Consider making this configurable
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization}, // Added Authorization for JWT
	}))

	// Removed global Fabric initialization: orgConfig, orgSetup, network, contract

	// Chaincode and Channel names can still be read from ENV or use defaults from auth package
	// This is handled within AuthMiddleware now, but good to be aware.
	chaincodeName := os.Getenv("CHAINCODE_NAME")
	if chaincodeName == "" {
		chaincodeName = auth.DefaultChaincodeName
	}

	channelName := os.Getenv("CHANNEL_NAME")
	if channelName == "" {
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
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
