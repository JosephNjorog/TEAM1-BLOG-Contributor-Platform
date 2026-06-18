# Migrations are compiled into the binary via go:embed, so the final
# image only needs the binary itself plus CA certs (for outbound HTTPS to
# Neon, Resend, Cloudinary, and the Avalanche RPC endpoint).

FROM golang:1.26-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/team1-api ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /out/team1-api ./team1-api

EXPOSE 8080
ENTRYPOINT ["./team1-api"]
