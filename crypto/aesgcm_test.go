package crypto_test

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"disappr.io/crypto"
)

func generateRandomKey(t *testing.T) []byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)
	return key
}

func TestEncryptDecrypt(t *testing.T) {
	key := generateRandomKey(t)

	tests := []struct {
		name        string
		input       string
		modifyCT    func(string) string
		key         []byte
		expectError bool
	}{
		{
			name:  "valid encryption and decryption",
			input: "hello world",
			key:   key,
		},
		{
			name:        "invalid base64",
			input:       "baddata",
			key:         key,
			modifyCT:    func(_ string) string { return "!!!notbase64===" },
			expectError: true,
		},
		{
			name:        "wrong key size",
			input:       "hello world",
			key:         key[:10],
			expectError: true,
		},
		{
			name:        "tampered ciphertext",
			input:       "secret",
			key:         key,
			modifyCT:    func(s string) string { return s[:len(s)-2] + "zz" },
			expectError: true,
		},
		{
			name:        "short ciphertext (less than nonce)",
			input:       "hi",
			key:         key,
			modifyCT:    func(_ string) string { return base64.StdEncoding.EncodeToString([]byte("123")) },
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct, err := crypto.Encrypt(tt.input, key)
			require.NoError(t, err)

			if tt.modifyCT != nil {
				ct = tt.modifyCT(ct)
			}

			plaintext, err := crypto.Decrypt(ct, tt.key)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.input, plaintext)
			}
		})
	}
}
