package auth_test

import (
	"os"
	"testing"
	"time"

	"disappr.io/auth"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testSigningKey = []byte("test-secret")

func createToken(t *testing.T, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(testSigningKey)
	require.NoError(t, err)
	return signed
}

func cloneClaims(original jwt.MapClaims, overrides map[string]interface{}) jwt.MapClaims {
	newClaims := jwt.MapClaims{}
	for k, v := range original {
		newClaims[k] = v
	}
	for k, v := range overrides {
		newClaims[k] = v
	}
	return newClaims
}

func TestVerifyFirebaseJWT(t *testing.T) {
	_ = os.Setenv("FIREBASE_PROJECT_ID", "disappr-io")

	origKeyfunc := auth.JwtTestKeyfunc()
	auth.SetJwtTestKeyfunc(func(token *jwt.Token) (interface{}, error) {
		return testSigningKey, nil
	})
	defer auth.SetJwtTestKeyfunc(origKeyfunc)

	validClaims := jwt.MapClaims{
		"aud": "disappr-io",
		"iss": "https://securetoken.google.com/disappr-io",
		"sub": "user123",
		"exp": time.Now().Add(5 * time.Minute).Unix(),
	}

	tests := []struct {
		name        string
		claims      jwt.MapClaims
		modifyToken func(string) string
		expectError bool
	}{
		{
			name:   "valid token",
			claims: validClaims,
		},
		{
			name:        "invalid audience",
			claims:      cloneClaims(validClaims, map[string]interface{}{"aud": "wrong-aud"}),
			expectError: true,
		},
		{
			name:        "invalid issuer",
			claims:      cloneClaims(validClaims, map[string]interface{}{"iss": "bad-issuer"}),
			expectError: true,
		},
		{
			name:        "malformed token",
			modifyToken: func(_ string) string { return "not.a.jwt" },
			expectError: true,
		},
		{
			name:        "expired token",
			claims:      cloneClaims(validClaims, map[string]interface{}{"exp": time.Now().Add(-10 * time.Minute).Unix()}),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tokenStr string
			if tt.modifyToken != nil {
				tokenStr = tt.modifyToken("")
			} else {
				tokenStr = createToken(t, tt.claims)
			}

			claims, err := auth.VerifyFirebaseJWT(tokenStr)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.claims["sub"], (*claims)["sub"])
			}
		})
	}
}
