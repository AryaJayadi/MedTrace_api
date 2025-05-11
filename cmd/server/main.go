package main

import (
	"fmt"
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
		fmt.Println("Error initializing setup for Org1: ", err)
	}

	chaincodeName := "drugtrace"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "medtrace"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := orgSetup.Gateway.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	batchService := services.NewBatchService(contract)

	batchHandler := handlers.NewBatchHandler(batchService)

	e.POST("/batches", batchHandler.CreateBatch)

	e.Logger.Fatal(e.Start(":8080"))
}
