# --- Build stage ---
FROM golang:1.24-alpine AS build
WORKDIR /app
# no go.sum yet, so only copy go.mod for now:
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /echo-service ./cmd/echo-service

# --- Runtime stage ---
FROM gcr.io/distroless/static-debian12:nonroot

# Copy the statically linked binary produced in the builder stage
COPY --from=build /echo-service /echo-service

# Distroless images have no shell, so the entrypoint must be JSON‑array form
ENTRYPOINT ["/echo-service"]