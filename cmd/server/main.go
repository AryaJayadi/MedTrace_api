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

// --- Organization 1 Configurations ---
const (
	Org1Name         = "Org1"
	Org1mspID        = "Org1MSP"
	Org1cryptoPath   = "../../test-network/organizations/peerOrganizations/org1.example.com" // Adjust if your path is different
	Org1certPath     = Org1cryptoPath + "/users/User1@org1.example.com/msp/signcerts"
	Org1keyPath      = Org1cryptoPath + "/users/User1@org1.example.com/msp/keystore"
	Org1tlsCertPath  = Org1cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	Org1peerEndpoint = "dns:///localhost:7051" // Default for Fabric test-network Org1
	Org1gatewayPeer  = "peer0.org1.example.com"
)

// --- Organization 2 Configurations ---
const (
	Org2Name         = "Org2"                                                                // EXAMPLE: Replace with actual Org2 name
	Org2mspID        = "Org2MSP"                                                             // EXAMPLE: Replace with actual Org2 MSP ID
	Org2cryptoPath   = "../../test-network/organizations/peerOrganizations/org2.example.com" // EXAMPLE: Adjust path for Org2
	Org2certPath     = Org2cryptoPath + "/users/User1@org2.example.com/msp/signcerts"        // EXAMPLE: Adjust user and path for Org2
	Org2keyPath      = Org2cryptoPath + "/users/User1@org2.example.com/msp/keystore"         // EXAMPLE: Adjust user and path for Org2
	Org2tlsCertPath  = Org2cryptoPath + "/peers/peer0.org2.example.com/tls/ca.crt"           // EXAMPLE: Adjust peer and path for Org2
	Org2peerEndpoint = "dns:///localhost:9051"                                               // EXAMPLE: Adjust port for Org2 peer (default for Fabric test-network Org2 is 9051)
	Org2gatewayPeer  = "peer0.org2.example.com"                                              // EXAMPLE: Adjust gateway peer for Org2
)

// --- Organization 3 Configurations (Example - Add as many as you need) ---
const (
	Org3Name         = "Org3"                                                                // EXAMPLE: Replace with actual Org3 name
	Org3mspID        = "Org3MSP"                                                             // EXAMPLE: Replace with actual Org3 MSP ID
	Org3cryptoPath   = "../../test-network/organizations/peerOrganizations/org3.example.com" // EXAMPLE: Adjust path for Org3
	Org3certPath     = Org3cryptoPath + "/users/User1@org3.example.com/msp/signcerts"        // EXAMPLE: Adjust user and path for Org3
	Org3keyPath      = Org3cryptoPath + "/users/User1@org3.example.com/msp/keystore"         // EXAMPLE: Adjust user and path for Org3
	Org3tlsCertPath  = Org3cryptoPath + "/peers/peer0.org3.example.com/tls/ca.crt"           // EXAMPLE: Adjust peer and path for Org3
	Org3peerEndpoint = "dns:///localhost:11051"                                              // EXAMPLE: Adjust port for Org3 peer (default for Fabric test-network Org3 if using a similar pattern might be 11051)
	Org3gatewayPeer  = "peer0.org3.example.com"                                              // EXAMPLE: Adjust gateway peer for Org3
)

// Choose which organization this instance of the API will represent
// This could also be determined by an environment variable or a config file
const (
	ActiveOrgName      = Org1Name         // Or Org2Name, Org3Name
	ActiveMspID        = Org1mspID        // Or Org2mspID, Org3mspID
	ActiveCryptoPath   = Org1cryptoPath   // Or Org2cryptoPath, Org3cryptoPath
	ActiveCertPath     = Org1certPath     // Or Org2certPath, Org3certPath
	ActiveKeyPath      = Org1keyPath      // Or Org2keyPath, Org3keyPath
	ActiveTlsCertPath  = Org1tlsCertPath  // Or Org2tlsCertPath, Org3tlsCertPath
	ActivePeerEndpoint = Org1peerEndpoint // Or Org2peerEndpoint, Org3peerEndpoint
	ActiveGatewayPeer  = Org1gatewayPeer  // Or Org2gatewayPeer, Org3gatewayPeer
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Use the "Active" organization's configuration for initialization
	orgConfig := fabric.OrgSetup{
		OrgName:      ActiveOrgName,
		MSPID:        ActiveMspID,
		CertPath:     ActiveCertPath,
		KeyPath:      ActiveKeyPath,
		TLSCertPath:  ActiveTlsCertPath,
		PeerEndpoint: ActivePeerEndpoint,
		GatewayPeer:  ActiveGatewayPeer,
	}

	fmt.Printf("Initializing API for organization: %s (MSP: %s)\n", orgConfig.OrgName, orgConfig.MSPID)

	orgSetup, err := fabric.Initialize(orgConfig)
	if err != nil {
		// Consider more robust error handling or logging to a file in production
		fmt.Printf("Error initializing setup for %s: %v\n", orgConfig.OrgName, err)
		// os.Exit(1) or panic("Failed to initialize Fabric setup") might be appropriate
		// depending on whether the application can run without Fabric.
		panic(fmt.Sprintf("Error initializing setup for %s: %v", orgConfig.OrgName, err))
	}
	fmt.Printf("Successfully initialized Fabric setup for: %s\n", orgConfig.OrgName)

	chaincodeName := "drugtrace"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "medtrace"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	fmt.Printf("Accessing chaincode '%s' on channel '%s'\n", chaincodeName, channelName)

	network := orgSetup.Gateway.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	batchService := services.NewBatchService(contract)
	batchHandler := handlers.NewBatchHandler(batchService)

	// API Grouping (Optional but good practice)
	apiGroup := e.Group("/api/v1") // Example API versioning

	// Routes
	apiGroup.POST("/batches", batchHandler.CreateBatch)
	// Add more routes here as needed, e.g.:
	// apiGroup.GET("/batches/:id", batchHandler.GetBatchByID)
	// apiGroup.PUT("/batches/:id", batchHandler.UpdateBatch)
	// apiGroup.GET("/batches", batchHandler.GetAllBatches)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	fmt.Printf("Starting server on port %s\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}
