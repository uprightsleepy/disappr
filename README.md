
# disappr

**disappr** is a secure, ephemeral pastebin built for privacy. It encrypts your content client-side using AES-GCM, stores it in Firestore with automatic expiration via TTL, and authenticates users with Firebase. Burn-after-read support ensures your secrets don’t linger longer than necessary.

> ⚠️ This project is under active development. Expect rapid iteration and improvements.

---

## ✨ Features

- 🔐 **AES-GCM Encryption** using keys stored in GCP Secret Manager
- 🧑‍💻 **Firebase JWT Authentication** for user ownership and access control
- 🗑 **Burn After Read** support — delete a paste immediately after it's viewed
- ⏳ **TTL-Based Expiration** via Firestore’s built-in Time-To-Live engine
- 🌐 **Cloud Run Deployment** using Terraform and Docker
- 🧪 **Tested with Table-Driven Unit Tests** and injectable interfaces for mocking

---

## 📦 Architecture Overview

- **Backend**: Go (`net/http`)
- **Auth**: Firebase JWT (Google Secure Token service)
- **Database**: Firestore (Native mode)
- **Secrets**: GCP Secret Manager (Base64-encoded AES key)
- **Infra**: Terraform, Cloud Run, Artifact Registry
- **CI/CD**: Manual Docker build + Terraform apply (currently)

---

## 🚀 Getting Started

### 🔧 Prerequisites

- [Go 1.21+](https://golang.org/)
- [gcloud CLI](https://cloud.google.com/sdk)
- [Docker](https://www.docker.com/)
- A Firebase project with auth enabled

### 🛠 Local Development

```bash
export FIREBASE_PROJECT_ID=your-project-id
export GCP_PROJECT=your-gcp-project-id

go run ./main.go
```

Use a bearer token from Firebase Auth in your requests.

---

## 🧪 Running Tests

```bash
go test ./... -v -cover
```

All major logic paths are covered by table-driven tests, including:
- JWT validation
- AES encryption/decryption
- Secret Manager mocking
- Auth middleware

---

## 🔐 Security Model

- All pastes are encrypted using AES-GCM before being stored
- Keys are stored securely in GCP Secret Manager
- Pastes are scoped by `OwnerID` (Firebase `sub`)
- TTL ensures pastes disappear after expiration (or immediate if burn-after-read)

---

## 📄 API Overview

### `POST /api/v1/paste`

Create a paste.

**Headers:**
```
Authorization: Bearer <JWT TOKEN>
```

**Body:**
```json
{
  "content": "secret text",
  "expires_in_minutes": 60,
  "burn_after_read": true
}
```

---

### `GET /api/v1/view?id=<paste_id>`

View a paste (burns if configured).

---

## 📌 TODO

- [ ] Frontend UI
- [ ] Custom domain mapping

---

## 📜 License

MIT — do whatever, just don’t sell insecure pastebins.
