package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"disappr.io/auth"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func mockVerify(fn func(token string) (*jwt.MapClaims, error)) func() {
	orig := auth.VerifyFirebaseJWT
	auth.VerifyFirebaseJWT = fn
	return func() {
		auth.VerifyFirebaseJWT = orig
	}
}

func TestRequireAuth(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		mockVerify func(token string) (*jwt.MapClaims, error)
		expectCode int
		expectUser string
	}{
		{
			name:       "valid token",
			authHeader: "Bearer validtoken",
			mockVerify: func(token string) (*jwt.MapClaims, error) {
				claims := jwt.MapClaims{"sub": "user123"}
				return &claims, nil
			},
			expectCode: http.StatusOK,
			expectUser: "user123",
		},
		{
			name:       "missing Authorization header",
			authHeader: "",
			expectCode: http.StatusUnauthorized,
		},
		{
			name:       "invalid token",
			authHeader: "Bearer badtoken",
			mockVerify: func(token string) (*jwt.MapClaims, error) {
				return nil, http.ErrAbortHandler
			},
			expectCode: http.StatusUnauthorized,
		},
		{
			name:       "missing sub claim",
			authHeader: "Bearer nosub",
			mockVerify: func(token string) (*jwt.MapClaims, error) {
				claims := jwt.MapClaims{}
				return &claims, nil
			},
			expectCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockVerify != nil {
				defer mockVerify(tt.mockVerify)()
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userID := r.Context().Value(auth.UserIDKey)
				assert.Equal(t, tt.expectUser, userID)
				w.WriteHeader(http.StatusOK)
			})

			handler := auth.RequireAuth(next)
			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectCode, rec.Code)
		})
	}
}
