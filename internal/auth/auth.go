package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AryaJayadi/MedTrace_api/cmd/fabric"
	"github.com/AryaJayadi/MedTrace_api/internal/config"
	"github.com/AryaJayadi/MedTrace_api/internal/models/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/labstack/echo/v4"
)

// JWTCustomClaims are custom claims extending default ones.
type JWTCustomClaims struct {
	OrgID string `json:"orgId"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

const (
	// OrgContextKey is the key used to store the Fabric contract in Echo context.
	OrgContextKey = "org_contract"
	// DefaultChaincodeName is used if CHAINCODE_NAME env var is not set.
	DefaultChaincodeName = "medtrace_cc"
	// DefaultChannelName is used if CHANNEL_NAME env var is not set.
	DefaultChannelName = "medtrace"
)

func init() {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		fmt.Println("Warning: JWT_SECRET_KEY is not set. Using a default, insecure key for demonstration purposes ONLY.")
		secret = "aVeryInsecureDefaultKeyPleaseChangeThisForProduction"
	}
	jwtSecret = []byte(secret)
}

// GenerateJWT generates a new JWT for a given organization ID.
func GenerateJWT(orgID string) (string, error) {
	// Validate if the orgID is a configured and valid organization
	if _, err := config.GetOrgConfig(orgID); err != nil {
		return "", fmt.Errorf("cannot generate token for invalid organization '%s': %w", orgID, err)
	}

	claims := &JWTCustomClaims{
		OrgID: orgID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)), // Token expires in 72 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "MedTraceAPI",
			Subject:   orgID, // Subject can be the OrgID
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signedToken, nil
}

// AuthMiddleware is an Echo middleware that handles JWT authentication
// and initializes the Fabric connection context for the authenticated organization.
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing or malformed JWT")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			return echo.NewHTTPError(http.StatusUnauthorized, "Malformed Authorization header: expecting 'Bearer <token>'")
		}
		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Invalid JWT: %s", err.Error()))
		}

		if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
			orgCfg, err := config.GetOrgConfig(claims.OrgID)
			if err != nil {
				c.Logger().Errorf("AuthMiddleware: Failed to get org config for %s: %v", claims.OrgID, err)
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Cannot process request for organization %s", claims.OrgID))
			}

			orgSetup, err := fabric.Initialize(orgCfg)
			if err != nil {
				c.Logger().Errorf("AuthMiddleware: Failed to initialize Fabric for Org %s: %v", claims.OrgID, err)
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to connect to network for organization %s", claims.OrgID))
			}

			// Defer closing the gateway connection. orgSetup.Gateway is client.Gateway (a struct).
			// It's valid if fabric.Initialize succeeded.
			defer func() {
				if errClose := orgSetup.Gateway.Close(); errClose != nil {
					c.Logger().Errorf("AuthMiddleware: Error closing Fabric gateway for %s: %v", claims.OrgID, errClose)
				}
			}()

			chaincodeName := os.Getenv("CHAINCODE_NAME")
			if chaincodeName == "" {
				chaincodeName = DefaultChaincodeName
			}
			channelName := os.Getenv("CHANNEL_NAME")
			if channelName == "" {
				channelName = DefaultChannelName
			}

			network := orgSetup.Gateway.GetNetwork(channelName)
			contractInstance := network.GetContract(chaincodeName)

			c.Set(OrgContextKey, contractInstance)
			return next(c)
		}
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid JWT claims")
	}
}

// GetContractFromContext retrieves the Fabric contract from the Echo context.
// Handlers use this to get the contract instance initialized by AuthMiddleware.
func GetContractFromContext(c echo.Context) (*client.Contract, error) {
	contractVal := c.Get(OrgContextKey)
	if contractVal == nil {
		return nil, fmt.Errorf("Fabric contract not found in context, ensure AuthMiddleware is applied")
	}
	contract, ok := contractVal.(*client.Contract)
	if !ok {
		return nil, fmt.Errorf("Fabric contract in context is not of expected type *client.Contract")
	}
	return contract, nil
}

// LoginPayload defines the expected JSON structure for the login request.
type LoginPayload struct {
	Organization string `json:"organization" validate:"required"`
	Password     string `json:"password" validate:"required"`
}

func validateOrgAndPassword(org, password string) bool {
	expectedPassword := org + "asdf"
	return password == expectedPassword
}

// LoginResponseData defines the structure for successful login response value
type LoginResponseData struct {
	Token   string `json:"token"`
	OrgID   string `json:"orgId"`
	Message string `json:"message"`
}

// LogoutResponseData defines the structure for successful logout response value
type LogoutResponseData struct {
	Message string `json:"message"`
}

// LoginHandler handles the /login endpoint.
func LoginHandler(c echo.Context) error {
	payload := new(LoginPayload)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorValueResponse[LoginResponseData](http.StatusBadRequest, "Invalid request payload: %v", err))
	}
	if payload.Organization == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorValueResponse[LoginResponseData](http.StatusBadRequest, "organization is required"))
	}
	if payload.Password == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorValueResponse[LoginResponseData](http.StatusBadRequest, "Password is required"))
	}
	if !validateOrgAndPassword(payload.Organization, payload.Password) {
		return c.JSON(http.StatusUnauthorized, response.ErrorValueResponse[LoginResponseData](http.StatusUnauthorized, "Invalid organization or password"))
	}

	token, err := GenerateJWT(payload.Organization)
	if err != nil {
		c.Logger().Errorf("LoginHandler: Failed to generate JWT for OrgID '%s': %v", payload.Organization, err)
		if strings.Contains(err.Error(), "cannot generate token for invalid organization") {
			return c.JSON(http.StatusBadRequest, response.ErrorValueResponse[LoginResponseData](http.StatusBadRequest, "Invalid organization ID: %s", payload.Organization))
		}
		return c.JSON(http.StatusInternalServerError, response.ErrorValueResponse[LoginResponseData](http.StatusInternalServerError, "Login failed: could not generate authentication token."))
	}

	responseData := LoginResponseData{
		Message: "Login successful",
		Token:   token,
		OrgID:   payload.Organization,
	}
	return c.JSON(http.StatusOK, response.SuccessValueResponse(responseData))
}

// LogoutHandler handles the /logout endpoint.
func LogoutHandler(c echo.Context) error {
	responseData := LogoutResponseData{
		Message: "Logout successful. Please ensure the token is removed from client-side storage.",
	}
	return c.JSON(http.StatusOK, response.SuccessValueResponse(responseData))
}
