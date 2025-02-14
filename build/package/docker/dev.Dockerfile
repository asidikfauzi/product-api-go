FROM golang:1.23-alpine

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

# Jalankan aplikasi dengan Air
CMD ["air", "-c", ".air.toml"]