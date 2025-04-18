# --- Build stage ---
FROM golang:1.24-alpine AS build
WORKDIR /app
# no go.sum yet, so only copy go.mod for now:
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /echo-service ./cmd/echo-service

# --- Runtime stage ---
FROM alpine:3.21
COPY --from=build /echo-service /usr/local/bin/echo-service
EXPOSE 8080
ENV PORT=8080
ENTRYPOINT ["/usr/local/bin/echo-service"]