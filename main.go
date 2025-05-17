// ghostbin-backend/main.go

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"disappr.io/crypto"
	"disappr.io/secrets"
)

type Paste struct {
	ID             string    `json:"id"`
	EncryptedData  string    `json:"encrypted_data"`
	BurnAfterRead  bool      `json:"burn_after_read"`
	ExpiresAt      time.Time `json:"expires_at"`
	Viewed         bool      `json:"-"`
}

var firestoreClient *firestore.Client

func main() {
	ctx := context.Background()
	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		log.Fatal("GCP_PROJECT environment variable must be set")
	}

	var err error
	firestoreClient, err = firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer firestoreClient.Close()

	http.HandleFunc("/api/v1/paste", createPasteHandler)
	http.HandleFunc("/api/v1/view", viewPasteHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func createPasteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB limit

	var req struct {
		Content         string `json:"content"`
		ExpiresInMinutes int    `json:"expires_in_minutes"`
		BurnAfterRead    bool   `json:"burn_after_read"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	key, err := secrets.GetEncryptionKey(r.Context())
	if err != nil {
		http.Error(w, "Failed to load encryption key", http.StatusInternalServerError)
		return
	}
	encrypted, err := crypto.Encrypt(req.Content, key)
	if err != nil {
		http.Error(w, "Encryption failed", http.StatusInternalServerError)
		return
	}

	pasteID := uuid.NewString()
	paste := Paste{
		ID:            pasteID,
		EncryptedData: encrypted,
		BurnAfterRead: req.BurnAfterRead,
		ExpiresAt:     time.Now().Add(time.Duration(req.ExpiresInMinutes) * time.Minute),
	}

	_, err = firestoreClient.Collection("pastes").Doc(pasteID).Set(r.Context(), paste)
	if err != nil {
		http.Error(w, "Failed to store paste", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"url":        fmt.Sprintf("/api/v1/view?id=%s", pasteID),
		"expires_at": paste.ExpiresAt,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func viewPasteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	doc, err := firestoreClient.Collection("pastes").Doc(id).Get(r.Context())
	if err != nil {
		http.Error(w, "Paste not found", http.StatusNotFound)
		return
	}

	var paste Paste
	if err := doc.DataTo(&paste); err != nil {
		http.Error(w, "Failed to parse paste", http.StatusInternalServerError)
		return
	}

	if paste.ExpiresAt.Before(time.Now()) || paste.Viewed {
		http.Error(w, "Paste expired or already viewed", http.StatusGone)
		return
	}

	key, err := secrets.GetEncryptionKey(r.Context())
	if err != nil {
		http.Error(w, "Failed to load encryption key", http.StatusInternalServerError)
		return
	}

	decrypted, err := crypto.Decrypt(paste.EncryptedData, key)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusInternalServerError)
		return
	}

	if paste.BurnAfterRead {
		_, _ = firestoreClient.Collection("pastes").Doc(id).Update(r.Context(), []firestore.Update{
			{Path: "Viewed", Value: true},
		})
	}

	resp := map[string]string{"content": decrypted}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
