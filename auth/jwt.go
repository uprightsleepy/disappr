package auth

import (
	"fmt"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"os"
)

var keySet *keyfunc.JWKS

var jwtKeyfunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
	if keySet == nil {
		return nil, fmt.Errorf("keySet not initialized")
	}
	return keySet.Keyfunc(token)
}

func InitFirebaseVerifier() error {
	jwksURL := "https://www.googleapis.com/service_accounts/v1/jwk/securetoken@system.gserviceaccount.com"

	options := keyfunc.Options{
		RefreshErrorHandler: func(err error) {
			fmt.Printf("JWKS refresh error: %v\n", err)
		},
		RefreshInterval:    time.Hour,
		RefreshTimeout:     10 * time.Second,
		RefreshUnknownKID:  true,
	}

	var err error
	keySet, err = keyfunc.Get(jwksURL, options)
	return err
}

var VerifyFirebaseJWT = verifyFirebaseJWT

func verifyFirebaseJWT(tokenString string) (*jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, jwtKeyfunc)
    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    projectID := os.Getenv("FIREBASE_PROJECT_ID")

    if claims["aud"] != projectID {
        return nil, fmt.Errorf("invalid aud")
    }
    if claims["iss"] != fmt.Sprintf("https://securetoken.google.com/%s", projectID) {
        return nil, fmt.Errorf("invalid iss")
    }

    return &claims, nil
}

// test-only: overrideable JWT keyfunc
func SetJwtTestKeyfunc(fn jwt.Keyfunc) {
	jwtKeyfunc = fn
}

// test-only: retrieve current keyfunc
func JwtTestKeyfunc() jwt.Keyfunc {
	return jwtKeyfunc
}
