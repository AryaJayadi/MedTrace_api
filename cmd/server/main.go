package main

import (
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	OrgName      = "Org1"
	mspID        = "Org1MSP"
	cryptoPath   = "../../test-network/organizations/peerOrganizations/org1.example.com"
	certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts"
	keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "dns:///localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	orgConfig := fabric.OrgSetup{
		OrgName:      OrgName,
		MSPID:        mspID,
		cryptoPath:   cryptoPath,
		certPath:     certPath,
		keyPath:      keyPath,
		tlsCertPath:  tlsCertPath,
		peerEndpoint: peerEndpoint,
		gatewayPeer:  gatewayPeer,
		gateway:      nil,
	}
	setup := fabric.Initialize()

	// Setup Fabric Gateway and Contract
	contract, err := services.SetupFabricContract()
	if err != nil {
		log.Fatal("Failed to connect to Fabric:", err)
	}

	// Initialize services
	batchService := services.NewBatchService(contract)
	// drugService := services.NewDrugService(contract) // similarly

	// Initialize handlers
	batchHandler := handlers.NewBatchHandler(batchService)
	// drugHandler := handlers.NewDrugHandler(drugService)

	// Define routes
	// For clarity, we skip versioning here. In real apps consider /api/v1.
	e.POST("/batches", batchHandler.CreateBatch)
	e.GET("/batches/:id", batchHandler.GetBatch)
	// e.POST("/drugs", drugHandler.CreateDrug)
	// e.GET("/drugs/:id", drugHandler.GetDrug)

	// Start the server on port 8080
	e.Logger.Fatal(e.Start(":8080"))
}
