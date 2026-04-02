FROM golang:1.25-bookworm AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o doc-metadata ./service/document-metadata/main.go

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/doc-metadata .
ENV HTTP_PORT=8080
ENV GRPC_PORT=50051
EXPOSE 8080 50051
CMD ["./doc-metadata"]
