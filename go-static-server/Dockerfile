# --- Stage 1: Build the Go static server ---
FROM docker.io/library/golang:1.22 AS builder

WORKDIR /build

# Copy module file first for layer caching.
COPY go.mod ./

# Copy source and generate HTML files.
COPY main.go ./
COPY generate_html.sh ./
RUN chmod +x generate_html.sh && ./generate_html.sh html

# Build a statically-linked binary.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /go-static-server .

# --- Stage 2: Generate self-signed TLS certificates ---
FROM docker.io/library/golang:1.22 AS certgen

WORKDIR /certgen
COPY ssl/gencert.go .

RUN mkdir -p /certs && \
    go run gencert.go -cert /certs/server.crt -key /certs/server.key

# --- Stage 3: Minimal scratch image ---
FROM scratch

# Copy the static binary.
COPY --from=builder /go-static-server /go-static-server

# Copy TLS certificates.
COPY --from=certgen /certs/server.crt /etc/ssl/server.crt
COPY --from=certgen /certs/server.key /etc/ssl/server.key

# Expose HTTP and HTTPS ports.
EXPOSE 8080 8443

# Environment variables for TLS cert/key paths (can be overridden).
ENV TLS_CERT=/etc/ssl/server.crt
ENV TLS_KEY=/etc/ssl/server.key

ENTRYPOINT ["/go-static-server"]
