package main

import (
	"log"
	"os"

	"github.com/AryaJayadi/MedTrace_api/cmd/fabric"
	"github.com/AryaJayadi/MedTrace_api/internal/handlers"
	"github.com/AryaJayadi/MedTrace_api/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	OrgName      = "Org1"
	mspID        = "Org1MSP"
	cryptoPath   = "../../../MedTrace_network/organizations/peerOrganizations/org1.medtrace.com"
	certPath     = cryptoPath + "/users/User1@org1.medtrace.com/msp/signcerts/User1@org1.medtrace.com-cert.pem"
	keyPath      = cryptoPath + "/users/User1@org1.medtrace.com/msp/keystore"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.medtrace.com/tls/ca.crt"
	peerEndpoint = "dns:///localhost:7051"
	gatewayPeer  = "peer0.org1.medtrace.com"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	orgConfig := fabric.OrgSetup{
		OrgName:      OrgName,
		MSPID:        mspID,
		CertPath:     certPath,
		KeyPath:      keyPath,
		TLSCertPath:  tlsCertPath,
		PeerEndpoint: peerEndpoint,
		GatewayPeer:  gatewayPeer,
	}
	orgSetup, err := fabric.Initialize(orgConfig)
	if err != nil {
		log.Fatalf("Failed to initialize Fabric Org: %v", err)
	}

	chaincodeName := "medtrace_cc"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "medtrace"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := orgSetup.Gateway.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	organizationService := services.NewOrganizationService(contract)
	batchService := services.NewBatchService(contract)
	drugService := services.NewDrugService(contract)
	transferService := services.NewTransferService(contract)
	ledgerService := services.NewLedgerService(contract)

	organizationHandler := handlers.NewOrganizationHandler(organizationService)
	batchHandler := handlers.NewBatchHandler(batchService)
	drugHandler := handlers.NewDrugHandler(drugService)
	transferHandler := handlers.NewTransferHandler(transferService)
	ledgerHandler := handlers.NewLedgerHandler(ledgerService)

	orgGroup := e.Group("/organizations")
	orgGroup.GET("", organizationHandler.GetOrganizations)
	orgGroup.GET("/:id", organizationHandler.GetOrganizationByID)

	batchesGroup := e.Group("/batches")
	batchesGroup.POST("", batchHandler.CreateBatch)
	batchesGroup.GET("", batchHandler.GetAllBatches)
	batchesGroup.GET("/:id/exists", batchHandler.BatchExists)
	batchesGroup.GET("/:id", batchHandler.GetBatchByID)
	batchesGroup.PATCH("/batches", batchHandler.UpdateBatch)

	ledgerGroup := e.Group("/ledger")
	ledgerGroup.POST("/init", ledgerHandler.InitLedger)

	drugsGroup := e.Group("/drugs")
	drugsGroup.POST("", drugHandler.CreateDrug)
	drugsGroup.GET("/my", drugHandler.GetMyDrugs)
	drugsGroup.GET("/my/available", drugHandler.GetMyAvailDrugs)
	drugsGroup.GET("/:drugID", drugHandler.GetDrug)
	drugsGroup.GET("/batch/:batchID", drugHandler.GetDrugByBatch)

	transferGroup := e.Group("/transfers")
	transferGroup.POST("", transferHandler.CreateTransfer)
	transferGroup.GET("/my", transferHandler.GetMyTransfers)
	transferGroup.GET("/my/outgoing", transferHandler.GetMyOutTransfer)
	transferGroup.GET("/my/incoming", transferHandler.GetMyInTransfer)
	transferGroup.POST("/accept", transferHandler.AcceptTransfer)
	transferGroup.POST("/reject", transferHandler.RejectTransfer)
	transferGroup.GET("/:id", transferHandler.GetTransfer)

	e.Logger.Fatal(e.Start(":8080"))
}
