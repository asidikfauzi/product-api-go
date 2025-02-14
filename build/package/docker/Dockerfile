FROM golang:1.23-alpine AS builder

# Install make
RUN apk add --no-cache make

# Set working directory to the application root (product-api-go)
WORKDIR /go/src/product-api-go

# Salin seluruh proyek dari root `product-api-go`, bukan dari `docker/`
COPY ../../../ ./

# Jalankan `make migrate`
RUN make migrate

# Jalankan `make seed`
RUN make seed

# Jalankan `make all`
RUN make all

# Final stage: minimal Alpine image
FROM alpine:3.18.4
RUN apk add --no-cache ca-certificates

# Copy binary hasil build ke dalam container
COPY --from=builder /go/src/product-api-go/bin/product ./
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /

# Set timezone dan environment
ENV TZ=Asia/Jakarta
ENV ZONEINFO=/zoneinfo.zip

# Expose port 8080
EXPOSE 8080

# Gunakan path absolut untuk entrypoint
ENTRYPOINT ["./product"]
