package secrets

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/googleapis/gax-go/v2"
	// Ensure option is imported if secretmanager.NewClient has options in its signature
	// "google.golang.org/api/option"
)

// secretManagerClient defines the interface for the Secret Manager client.
// This allows for mocking in tests.
type secretManagerClient interface {
	AccessSecretVersion(context.Context, *secretmanagerpb.AccessSecretVersionRequest, ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error)
	Close() error
}

// newClientFunc defines a function type for creating a new secretManagerClient.
// This is used to allow mocking of client creation.
type newClientFunc func(ctx context.Context) (secretManagerClient, error)

func GetEncryptionKey(ctx context.Context, ncf newClientFunc) ([]byte, error) {
	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		return nil, fmt.Errorf("GCP_PROJECT environment variable not set")
	}
	name := fmt.Sprintf("projects/%s/secrets/disappr-aes-key/versions/latest", projectID)

	client, err := ncf(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret manager client: %w", err)
	}
	defer client.Close()

	accessReq := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := client.AccessSecretVersion(ctx, accessReq)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version: %w", err)
	}

	key, err := base64.StdEncoding.DecodeString(string(result.Payload.Data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret data: %w", err)
	}

	return key, nil
}
