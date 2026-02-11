FROM golang:1.26-bookworm AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /app ./cmd/server

FROM gcr.io/distroless/base-debian12:nonroot
COPY --from=builder /app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
