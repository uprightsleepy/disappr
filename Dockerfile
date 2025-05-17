# ---- Builder stage ----
FROM golang:1.23.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o disappr .

# ---- Final stage ----
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app/disappr /disappr

USER nonroot:nonroot
ENTRYPOINT ["/disappr"]
