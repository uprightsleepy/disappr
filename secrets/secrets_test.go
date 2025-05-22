package secrets

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"

	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/googleapis/gax-go/v2"
	"github.com/stretchr/testify/assert"
)

// mockSecretManagerClient is a mock implementation of the secretManagerClient interface.
type mockSecretManagerClient struct {
	accessSecretVersionFunc func(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error)
	closeFunc               func() error
}

func (m *mockSecretManagerClient) AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
	if m.accessSecretVersionFunc != nil {
		return m.accessSecretVersionFunc(ctx, req, opts...)
	}
	// Default behavior: Return an empty response and no error, or an error if appropriate
	return &secretmanagerpb.AccessSecretVersionResponse{}, nil
}

func (m *mockSecretManagerClient) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil // Default behavior if not set
}

func TestGetEncryptionKey_Success(t *testing.T) {
	// Define the expected key (decoded) and the base64 encoded version
	expectedKey := []byte("this-is-a-test-key-32-bytes-long") // 32 bytes
	encodedKey := base64.StdEncoding.EncodeToString(expectedKey)

	// Create a mock SecretManager client instance
	mockClientInstance := &mockSecretManagerClient{
		accessSecretVersionFunc: func(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
			return &secretmanagerpb.AccessSecretVersionResponse{
				Payload: &secretmanagerpb.SecretPayload{
					Data: []byte(encodedKey),
				},
			}, nil
		},
		closeFunc: func() error {
			return nil
		},
	}

	// Define the newClientFunc to return the mock client instance
	mockNCF := func(ctx context.Context) (SecretManagerClient, error) {
		return mockClientInstance, nil
	}

	// Call GetEncryptionKey with the mock newClientFunc
	ctx := context.Background()
	// Set GCP_PROJECT for the test environment
	t.Setenv("GCP_PROJECT", "test-project")
	key, err := GetEncryptionKey(ctx, mockNCF)

	// Assertions
	assert.NoError(t, err, "GetEncryptionKey should not return an error")
	assert.Equal(t, expectedKey, key, "Returned key does not match expected key")
}

func TestGetEncryptionKey_ClientError(t *testing.T) {
	expectedErrorMsg := "mock NewClient error"

	// Define a newClientFunc that returns an error
	mockNCF := func(ctx context.Context) (SecretManagerClient, error) {
		return nil, errors.New(expectedErrorMsg)
	}

	// Call GetEncryptionKey with the mock newClientFunc
	ctx := context.Background()
	// Set GCP_PROJECT for the test environment, though it won't be used if client creation fails
	t.Setenv("GCP_PROJECT", "test-project")
	key, err := GetEncryptionKey(ctx, mockNCF)

	// Assertions
	assert.Error(t, err, "GetEncryptionKey should return an error")
	assert.Nil(t, key, "Returned key should be nil on client error")
	assert.Contains(t, err.Error(), "failed to create secret manager client", "Error message should indicate client creation failure")
	assert.Contains(t, err.Error(), expectedErrorMsg, "Error message should contain the original error")
}

func TestGetEncryptionKey_AccessSecretVersionError(t *testing.T) {
	expectedErrorMsg := "mock AccessSecretVersion error"

	// Create a mock SecretManager client instance
	mockClientInstance := &mockSecretManagerClient{
		accessSecretVersionFunc: func(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
			// Return the predefined error for AccessSecretVersion
			return nil, errors.New(expectedErrorMsg)
		},
		closeFunc: func() error {
			return nil
		},
	}

	// Define the newClientFunc to return the mock client instance (simulating successful client creation)
	mockNCF := func(ctx context.Context) (SecretManagerClient, error) {
		return mockClientInstance, nil
	}

	// Call GetEncryptionKey with the mock newClientFunc
	ctx := context.Background()
	// Set GCP_PROJECT for the test environment
	t.Setenv("GCP_PROJECT", "test-project")
	key, err := GetEncryptionKey(ctx, mockNCF)

	// Assertions
	assert.Error(t, err, "GetEncryptionKey should return an error due to AccessSecretVersion failure")
	assert.Nil(t, key, "Returned key should be nil on AccessSecretVersion error")
	assert.Contains(t, err.Error(), "failed to access secret version", "Error message should indicate AccessSecretVersion failure")
	assert.Contains(t, err.Error(), expectedErrorMsg, "Error message should contain the original AccessSecretVersion error")
}

func TestGetEncryptionKey_DecodeError(t *testing.T) {
	invalidBase64String := "this-is-not-base64"

	// Create a mock SecretManager client instance
	mockClientInstance := &mockSecretManagerClient{
		accessSecretVersionFunc: func(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
			return &secretmanagerpb.AccessSecretVersionResponse{
				Payload: &secretmanagerpb.SecretPayload{
					Data: []byte(invalidBase64String),
				},
			}, nil
		},
		closeFunc: func() error {
			return nil
		},
	}

	// Define the newClientFunc to return the mock client instance
	mockNCF := func(ctx context.Context) (SecretManagerClient, error) {
		return mockClientInstance, nil
	}

	// Call GetEncryptionKey with the mock newClientFunc
	ctx := context.Background()
	// Set GCP_PROJECT for the test environment
	t.Setenv("GCP_PROJECT", "test-project")
	key, err := GetEncryptionKey(ctx, mockNCF)

	// Assertions
	assert.Error(t, err, "GetEncryptionKey should return an error due to base64 decoding failure")
	assert.Nil(t, key, "Returned key should be nil on decoding error")
	assert.Contains(t, err.Error(), "failed to decode secret data", "Error message should indicate decoding failure")
	// Check for the specific Go base64 error text
	assert.Contains(t, err.Error(), "illegal base64 data", "Error message should contain specific base64 error")
}
