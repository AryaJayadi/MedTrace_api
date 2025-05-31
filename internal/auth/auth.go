package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AryaJayadi/MedTrace_api/cmd/fabric"
	"github.com/AryaJayadi/MedTrace_api/internal/config"
	"github.com/AryaJayadi/MedTrace_api/internal/models/dto/auth"
	"github.com/AryaJayadi/MedTrace_api/internal/models/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/labstack/echo/v4"
)

// JWTCustomClaims are custom claims extending default ones.
type JWTCustomClaims struct {
	OrgID     string `json:"orgId"`
	TokenType string `json:"tokenType"` // e.g., "access", "refresh"
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

	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"

	AccessTokenDuration  = time.Hour * 1      // Access token valid for 1 hour
	RefreshTokenDuration = time.Hour * 24 * 7 // Refresh token valid for 7 days
)

func init() {
	secret := os.Getenv("JWT_SECRET_KEY")
	jwtSecret = []byte(secret)
}

// GenerateAccessToken generates a new short-lived access JWT for a given organization ID.
func GenerateAccessToken(orgID string) (string, error) {
	if _, err := config.GetOrgConfig(orgID); err != nil {
		return "", fmt.Errorf("cannot generate access token for invalid organization '%s': %w", orgID, err)
	}

	claims := &JWTCustomClaims{
		OrgID:     orgID,
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "MedTraceAPI",
			Subject:   orgID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}
	return signedToken, nil
}

// GenerateRefreshToken generates a new long-lived refresh JWT for a given organization ID.
func GenerateRefreshToken(orgID string) (string, error) {
	if _, err := config.GetOrgConfig(orgID); err != nil {
		return "", fmt.Errorf("cannot generate refresh token for invalid organization '%s': %w", orgID, err)
	}

	claims := &JWTCustomClaims{
		OrgID:     orgID,
		TokenType: TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "MedTraceAPI",
			Subject:   orgID, // Subject can be the OrgID, used to identify whom the token refers to
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}
	return signedToken, nil
}

// AuthMiddleware is an Echo middleware that handles JWT authentication
// and initializes the Fabric connection context for the authenticated organization.
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, response.ErrorValueResponse[interface{}](http.StatusUnauthorized, "Missing or malformed JWT"))
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			return c.JSON(http.StatusUnauthorized, response.ErrorValueResponse[interface{}](http.StatusUnauthorized, "Malformed Authorization header: expecting 'Bearer <token>'"))
		}
		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil {
			// Handle specific JWT errors like expiry
			if err == jwt.ErrTokenExpired {
				return c.JSON(http.StatusUnauthorized, response.ErrorValueResponse[interface{}](http.StatusUnauthorized, "Access token has expired"))
			}
			return c.JSON(http.StatusUnauthorized, response.ErrorValueResponse[interface{}](http.StatusUnauthorized, fmt.Sprintf("Invalid access token: %s", err.Error())))
		}

		if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
			if claims.TokenType != TokenTypeAccess {
				return c.JSON(http.StatusForbidden, response.ErrorValueResponse[interface{}](http.StatusForbidden, "Invalid token type: an access token is required"))
			}

			orgCfg, err := config.GetOrgConfig(claims.OrgID)
			if err != nil {
				c.Logger().Errorf("AuthMiddleware: Failed to get org config for %s: %v", claims.OrgID, err)
				return c.JSON(http.StatusInternalServerError, response.ErrorValueResponse[interface{}](http.StatusInternalServerError, fmt.Sprintf("Cannot process request for organization %s", claims.OrgID)))
			}

			orgSetup, err := fabric.Initialize(orgCfg)
			if err != nil {
				c.Logger().Errorf("AuthMiddleware: Failed to initialize Fabric for Org %s: %v", claims.OrgID, err)
				return c.JSON(http.StatusInternalServerError, response.ErrorValueResponse[interface{}](http.StatusInternalServerError, fmt.Sprintf("Failed to connect to network for organization %s", claims.OrgID)))
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
		return c.JSON(http.StatusUnauthorized, response.ErrorValueResponse[interface{}](http.StatusUnauthorized, "Invalid access token claims"))
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

func validateOrgAndPassword(org, password string) bool {
	expectedPassword := org + "asdf"
	return password == expectedPassword
}

// LoginHandler handles the /login endpoint.
func LoginHandler(c echo.Context) error {
	payload := new(auth.PayloadLogin)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorValueResponse[auth.LoginResponseData](http.StatusBadRequest, "Invalid request payload: %v", err))
	}
	if payload.Organization == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorValueResponse[auth.LoginResponseData](http.StatusBadRequest, "organization is required"))
	}
	if payload.Password == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorValueResponse[auth.LoginResponseData](http.StatusBadRequest, "Password is required"))
	}
	if !validateOrgAndPassword(payload.Organization, payload.Password) {
		return c.JSON(http.StatusUnauthorized, response.ErrorValueResponse[auth.LoginResponseData](http.StatusUnauthorized, "Invalid organization or password"))
	}

	accessToken, err := GenerateAccessToken(payload.Organization)
	if err != nil {
		c.Logger().Errorf("LoginHandler: Failed to generate access token for OrgID '%s': %v", payload.Organization, err)
		if strings.Contains(err.Error(), "cannot generate access token for invalid organization") {
			return c.JSON(http.StatusBadRequest, response.ErrorValueResponse[auth.LoginResponseData](http.StatusBadRequest, "Invalid organization ID: %s", payload.Organization))
		}
		return c.JSON(http.StatusInternalServerError, response.ErrorValueResponse[auth.LoginResponseData](http.StatusInternalServerError, "Login failed: could not generate access token."))
	}

	refreshToken, err := GenerateRefreshToken(payload.Organization)
	if err != nil {
		c.Logger().Errorf("LoginHandler: Failed to generate refresh token for OrgID '%s': %v", payload.Organization, err)
		return c.JSON(http.StatusInternalServerError, response.ErrorValueResponse[auth.LoginResponseData](http.StatusInternalServerError, "Login failed: could not generate refresh token."))
	}

	responseData := auth.LoginResponseData{ // Use the updated LoginResponseData struct
		Message:      "Login successful",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		OrgID:        payload.Organization,
	}
	return c.JSON(http.StatusOK, response.SuccessValueResponse(responseData))
}

// RefreshTokenHandler handles the /auth/refresh endpoint.
func RefreshTokenHandler(c echo.Context) error {
	reqPayload := new(auth.PayloadRefreshToken) // Use the defined RefreshTokenRequest DTO
	if err := c.Bind(reqPayload); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorValueResponse[auth.RefreshTokenResponseData](http.StatusBadRequest, "Invalid request payload: %v", err))
	}

	if reqPayload.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorValueResponse[auth.RefreshTokenResponseData](http.StatusBadRequest, "Refresh token is required"))
	}

	token, err := jwt.ParseWithClaims(reqPayload.RefreshToken, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		// Handle specific JWT errors like expiry
		if err == jwt.ErrTokenExpired {
			return c.JSON(http.StatusUnauthorized, response.ErrorValueResponse[auth.RefreshTokenResponseData](http.StatusUnauthorized, "Refresh token has expired"))
		}
		return c.JSON(http.StatusUnauthorized, response.ErrorValueResponse[auth.RefreshTokenResponseData](http.StatusUnauthorized, fmt.Sprintf("Invalid refresh token: %s", err.Error())))
	}

	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		if claims.TokenType != TokenTypeRefresh {
			return c.JSON(http.StatusForbidden, response.ErrorValueResponse[auth.RefreshTokenResponseData](http.StatusForbidden, "Invalid token type: not a refresh token"))
		}
		if claims.OrgID == "" {
			return c.JSON(http.StatusUnauthorized, response.ErrorValueResponse[auth.RefreshTokenResponseData](http.StatusUnauthorized, "Refresh token missing organization ID"))
		}

		// Validate if the organization from the token still exists/is valid
		if _, orgErr := config.GetOrgConfig(claims.OrgID); orgErr != nil {
			c.Logger().Warnf("RefreshTokenHandler: Organization '%s' from refresh token no longer valid: %v", claims.OrgID, orgErr)
			return c.JSON(http.StatusUnauthorized, response.ErrorValueResponse[auth.RefreshTokenResponseData](http.StatusUnauthorized, "Organization from refresh token is no longer valid"))
		}

		newAccessToken, err := GenerateAccessToken(claims.OrgID)
		if err != nil {
			c.Logger().Errorf("RefreshTokenHandler: Failed to generate new access token for OrgID '%s': %v", claims.OrgID, err)
			return c.JSON(http.StatusInternalServerError, response.ErrorValueResponse[auth.RefreshTokenResponseData](http.StatusInternalServerError, "Failed to generate new access token"))
		}

		// OPTIONAL: Implement Refresh Token Rotation
		// For enhanced security, you can issue a new refresh token here as well.
		// This invalidates the used refresh token.
		// newRefreshToken, err := GenerateRefreshToken(claims.OrgID)
		// if err != nil {
		//    c.Logger().Errorf("RefreshTokenHandler: Failed to generate new refresh token for OrgID '%s': %v", claims.OrgID, err)
		//    // Decide if this failure should prevent issuing the new access token
		// }

		responseData := auth.RefreshTokenResponseData{
			AccessToken: newAccessToken,
			// NewRefreshToken: newRefreshToken, // If implementing rotation
		}
		return c.JSON(http.StatusOK, response.SuccessValueResponse(responseData))
	}

	return c.JSON(http.StatusUnauthorized, response.ErrorValueResponse[auth.RefreshTokenResponseData](http.StatusUnauthorized, "Invalid refresh token claims"))
}

// LogoutHandler handles the /logout endpoint.
func LogoutHandler(c echo.Context) error {
	responseData := auth.LogoutResponseData{
		Message: "Logout successful. Please ensure the token is removed from client-side storage.",
	}
	return c.JSON(http.StatusOK, response.SuccessValueResponse(responseData))
}
