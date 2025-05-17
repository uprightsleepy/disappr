package secrets

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

func GetEncryptionKey(ctx context.Context) ([]byte, error) {
	projectID := os.Getenv("GCP_PROJECT")
	name := fmt.Sprintf("projects/%s/secrets/disappr-aes-key/versions/latest", projectID)

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	accessReq := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := client.AccessSecretVersion(ctx, accessReq)
	if err != nil {
		return nil, err
	}

	key, err := base64.StdEncoding.DecodeString(string(result.Payload.Data))
	if err != nil {
		return nil, err
	}

	return key, nil
}
