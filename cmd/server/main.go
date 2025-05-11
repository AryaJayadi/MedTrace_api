package main

import (
	"fmt"
	"os"
	"path/filepath" // For constructing paths robustly

	"github.com/AryaJayadi/MedTrace_api/cmd/fabric"
	"github.com/AryaJayadi/MedTrace_api/internal/handlers"
	"github.com/AryaJayadi/MedTrace_api/internal/services"

	"github.com/joho/godotenv" // Import godotenv
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// You can keep these const blocks for different orgs as a reference
// or if you ever need to switch back to hardcoded multi-org setup,
// but they won't be used directly for the .env configuration.

// --- Organization 1 Configurations (Reference) ---
const (
	_Org1Name  = "Org1" // Prefix with _ to indicate they are for reference
	_Org1mspID = "Org1MSP"
	// ... other Org1 constants
)

// --- Organization 2 Configurations (Reference) ---
const (
	_Org2Name  = "Org2"
	_Org2mspID = "Org2MSP"
	// ... other Org2 constants
)

func main() {
	// Load .env file.
	// This will load environment variables from a .env file in the current directory.
	// If .env is not found, it will not throw an error, allowing environment variables
	// to be set by the system as well.
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: Error loading .env file, will rely on system environment variables:", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Read configuration from environment variables (loaded from .env or system)
	orgName := os.Getenv("ORG_NAME")
	mspID := os.Getenv("MSP_ID")
	cryptoPath := os.Getenv("CRYPTO_PATH")
	userIdentity := os.Getenv("USER_IDENTITY") // e.g., User1@org1.example.com
	peerHostname := os.Getenv("PEER_HOSTNAME") // e.g., peer0.org1.example.com
	peerEndpoint := os.Getenv("PEER_ENDPOINT")

	if orgName == "" || mspID == "" || cryptoPath == "" || userIdentity == "" || peerHostname == "" || peerEndpoint == "" {
		// Handle missing critical environment variables
		// You might want to panic or provide default values if appropriate,
		// but for Fabric connection, these are usually essential.
		panic("Error: Missing one or more required environment variables for Fabric connection (ORG_NAME, MSP_ID, CRYPTO_PATH, USER_IDENTITY, PEER_HOSTNAME, PEER_ENDPOINT)")
	}

	// Construct full paths
	certPath := filepath.Join(cryptoPath, "users", userIdentity, "msp", "signcerts")
	keyPath := filepath.Join(cryptoPath, "users", userIdentity, "msp", "keystore")
	tlsCertPath := filepath.Join(cryptoPath, "peers", peerHostname, "tls", "ca.crt")

	orgConfig := fabric.OrgSetup{
		OrgName:      orgName,
		MSPID:        mspID,
		CertPath:     certPath,
		KeyPath:      keyPath,
		TLSCertPath:  tlsCertPath,
		PeerEndpoint: peerEndpoint,
		GatewayPeer:  peerHostname, // GatewayPeer is often the same as the peer you're connecting to
	}

	fmt.Printf("Initializing API for organization: %s (MSP: %s) using user: %s\n", orgConfig.OrgName, orgConfig.MSPID, userIdentity)
	fmt.Printf("CryptoPath: %s\n", cryptoPath)
	fmt.Printf("CertPath: %s\n", orgConfig.CertPath)
	fmt.Printf("KeyPath: %s\n", orgConfig.KeyPath)
	fmt.Printf("TLSCertPath: %s\n", orgConfig.TLSCertPath)
	fmt.Printf("PeerEndpoint: %s\n", orgConfig.PeerEndpoint)
	fmt.Printf("GatewayPeer: %s\n", orgConfig.GatewayPeer)

	orgSetup, setupErr := fabric.Initialize(orgConfig)
	if setupErr != nil {
		panic(fmt.Sprintf("Error initializing Fabric setup for %s: %v", orgConfig.OrgName, setupErr))
	}
	fmt.Printf("Successfully initialized Fabric setup for: %s\n", orgConfig.OrgName)

	chaincodeName := os.Getenv("CHAINCODE_NAME")
	if chaincodeName == "" {
		chaincodeName = "drugtrace" // Default if not set
	}

	channelName := os.Getenv("CHANNEL_NAME")
	if channelName == "" {
		channelName = "medtrace" // Default if not set
	}

	fmt.Printf("Accessing chaincode '%s' on channel '%s'\n", chaincodeName, channelName)

	network := orgSetup.Gateway.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	batchService := services.NewBatchService(contract)
	batchHandler := handlers.NewBatchHandler(batchService)

	apiGroup := e.Group("/api/v1")
	apiGroup.POST("/batches", batchHandler.CreateBatch)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Starting server on port %s\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}
