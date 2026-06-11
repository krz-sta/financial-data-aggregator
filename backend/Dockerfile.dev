FROM golang:1.25.9-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/api/main.go 

FROM alpine:latest 

WORKDIR /app

COPY --from=builder /server /server

CMD ["/server"]