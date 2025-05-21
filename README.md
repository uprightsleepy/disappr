# disappr

**disappr** is a secure, ephemeral pastebin built for privacy. It encrypts your content client-side using AES-GCM, stores it in Firestore with automatic expiration via TTL, and authenticates users with Firebase. Burn-after-read support ensures your secrets don't linger longer than necessary.

> ⚠️ This project is under active development. Expect rapid iteration and improvements.

---

## ✨ Features

- 🔐 **AES-GCM Encryption** using keys stored in GCP Secret Manager
- 🧑‍💻 **Firebase JWT Authentication** for user ownership and access control
- 🗑 **Burn After Read** support — delete a paste immediately after it's viewed
- ⏳ **TTL-Based Expiration** via Firestore's built-in Time-To-Live engine
- 🌐 **Cloud Run Deployment** using Terraform and Docker
- 🧪 **Tested with Table-Driven Unit Tests** and injectable interfaces for mocking
- 🔄 **CI/CD Pipeline** with automated testing, security scanning, and deployment
- 🛡️ **Container Security Scanning** using Trivy to detect vulnerabilities
- 📊 **Code Quality Enforcement** with golangci-lint static analysis

---

## 📦 Architecture Overview

- **Backend**: Go (`net/http`)
- **Auth**: Firebase JWT (Google Secure Token service)
- **Database**: Firestore (Native mode)
- **Secrets**: GCP Secret Manager (Base64-encoded AES key)
- **Infra**: Terraform, Cloud Run, Artifact Registry
- **CI/CD**: Cloud Build with automated testing, security scanning, and deployment

---

## 🔄 CI/CD Pipeline

The project uses a comprehensive CI/CD pipeline that runs on every push to the main branch:

1. **Code Quality Checks**:
   - Static code analysis with golangci-lint
   - 80% test coverage enforcement

2. **Security Validation**:
   - Container vulnerability scanning with Trivy
   - Secret detection in codebase
   - Dependency vulnerability analysis

3. **Automated Deployment**:
   - Docker image building and publishing to Artifact Registry
   - Zero-downtime deployment to Cloud Run

This ensures that every change is thoroughly tested and securely deployed.

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
- Continuous vulnerability scanning in CI/CD pipeline
- Automated detection of secrets and sensitive information in code
- Regular dependency updates to patch security vulnerabilities

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

## 🛠️ Infrastructure as Code

The entire infrastructure is managed using Terraform, allowing for:
- Repeatable, consistent deployments
- Version-controlled infrastructure
- Easy scaling and modifications
- Automated resource provisioning

Key resources include:
- Cloud Run services for hosting the application
- Firestore database for storing encrypted pastes
- Artifact Registry for Docker image storage
- IAM permissions and service accounts

---

## 📌 TODO

- [ ] Frontend UI
- [ ] Custom domain mapping
- [ ] Enhanced observability and monitoring
- [ ] User paste management dashboard

---

## 📜 License

MIT — do whatever, just don't sell insecure pastebins.