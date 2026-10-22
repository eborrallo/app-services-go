# --- Build stage ---
FROM golang:alpine AS build

# Tools needed to generate templates and build the frontend assets.
RUN apk add --update --no-cache git nodejs npm

WORKDIR /go/src/api
COPY . .

# Build the frontend assets (tailwind/postcss via parcel) into ./static.
RUN npm install
RUN npm run build

# Generate templ templates (pinned to match go.mod) and build the binary.
RUN CGO_ENABLED=0 go run github.com/a-h/templ/cmd/templ@v0.2.646 generate
RUN CGO_ENABLED=0 go build -o /go/bin/app-services-go cmd/api/main.go

# --- Runtime stage ---
FROM alpine:latest

# TLS root certificates for outbound HTTPS (blockchain RPC, SMTP, etc.).
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy only the compiled binary and the static assets it serves at /static.
COPY --from=build /go/bin/app-services-go /app/app-services-go
COPY --from=build /go/src/api/static /app/static

EXPOSE 8080

ENTRYPOINT ["/app/app-services-go"]
