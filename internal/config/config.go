package config

import (
	"fmt"
	"path/filepath"

	"github.com/AryaJayadi/MedTrace_api/cmd/fabric"
)

type OrgInfo struct {
	Name         string
	MSPID        string
	CryptoPath   string // Populated by init()
	PeerEndpoint string
	GatewayPeer  string
}

var orgConfigurations = map[string]OrgInfo{
	"Org1": {
		Name:         "Org1",
		MSPID:        "Org1MSP",
		PeerEndpoint: "dns:///localhost:7051",
		GatewayPeer:  "peer0.org1.medtrace.com",
	},
	"Org2": {
		Name:         "Org2",
		MSPID:        "Org2MSP",
		PeerEndpoint: "dns:///localhost:8051",
		GatewayPeer:  "peer0.org2.medtrace.com",
	},
	"Org3": {
		Name:         "Org3",
		MSPID:        "Org3MSP",
		PeerEndpoint: "dns:///localhost:9051",
		GatewayPeer:  "peer0.org3.medtrace.com",
	},
	"Org4": {
		Name:         "Org4",
		MSPID:        "Org4MSP",
		PeerEndpoint: "dns:///localhost:10051",
		GatewayPeer:  "peer0.org4.medtrace.com",
	},
}

// GetOrgConfig returns the fabric.OrgSetup for the specified organization.
// It dynamically constructs the necessary paths.
func GetOrgConfig(orgName string) (fabric.OrgSetup, error) {
	orgInfo, ok := orgConfigurations[orgName]
	if !ok {
		return fabric.OrgSetup{}, fmt.Errorf("organization '%s' not found in configuration", orgName)
	}

	// userAndOrgDomain constructs strings like "User1@org1.medtrace.com"
	// Assumes "User1" for all orgs for now. This could be made dynamic if needed.
	userAndOrgDomain := fmt.Sprintf("User1@%s.medtrace.com", lc(orgInfo.Name))

	return fabric.OrgSetup{
		OrgName:      orgInfo.Name,
		MSPID:        orgInfo.MSPID,
		CertPath:     filepath.Join(orgInfo.CryptoPath, "users", userAndOrgDomain, "msp", "signcerts", userAndOrgDomain+"-cert.pem"),
		KeyPath:      filepath.Join(orgInfo.CryptoPath, "users", userAndOrgDomain, "msp", "keystore"), // fabric.Initialize's PopulateWallet handles finding the key file in this directory
		TLSCertPath:  filepath.Join(orgInfo.CryptoPath, "peers", fmt.Sprintf("peer0.%s.medtrace.com", lc(orgInfo.Name)), "tls", "ca.crt"),
		PeerEndpoint: orgInfo.PeerEndpoint,
		GatewayPeer:  orgInfo.GatewayPeer,
	}, nil
}

// lc is a helper to format the organization name for path construction (e.g., "Org1" -> "org1").
func lc(s string) string {
	if len(s) > 3 && (s[:3] == "Org" || s[:3] == "org") {
		return "org" + s[3:] // Converts "Org1" to "org1", "Org12" to "org12"
	}
	// Fallback for unexpected formats. Consider more robust error handling or logging
	// if the input format can vary significantly.
	// For "orgX.medtrace.com", X should be like "org1", "org2".
	// Example: If s is "MyOrg", this would return "MyOrg". Modify if specific lowercasing is needed.
	// For current use with "Org1", "Org2", etc., this is fine.
	// If the input is guaranteed to be "OrgX", we could simplify, but this is safer.
	return s // Default to returning original string if not matching "OrgX" prefix.
}

func init() {
	// Populate CryptoPath for each organization.
	// This path is relative to where the application binary is expected to be run from
	// (e.g., MedTrace_api/cmd/server/). It navigates to the MedTrace_network directory.
	// Example: if binary is at /home/user/MedTrace_api/cmd/server/server_binary
	// then ../../../MedTrace_network will resolve to /home/user/MedTrace_network
	networkRelativePath := "../../../MedTrace_network"

	for orgKey, info := range orgConfigurations {
		// Construct the CryptoPath: e.g., "../../../MedTrace_network/organizations/peerOrganizations/org1.medtrace.com"
		// lc(info.Name) converts "Org1" to "org1", etc., for the directory name.
		info.CryptoPath = filepath.Join(networkRelativePath, "organizations", "peerOrganizations", fmt.Sprintf("%s.medtrace.com", lc(info.Name)))
		orgConfigurations[orgKey] = info // Update the map with the populated CryptoPath
	}
}
