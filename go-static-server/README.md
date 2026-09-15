# go-static-server

A lightweight Go-based static file server, built as a drop-in replacement for
the [nginx](../nginx/) image used in cloud-bulldozer benchmarks.

## What it does

- Serves 15 pre-generated static HTML files (128 B to 2 MiB) plus an `index.html`
- Listens on **`:8080`** (HTTP) and **`:8443`** (HTTPS with self-signed TLS)
- Accepts both **GET** and **POST** on all paths (POST returns the same static file)
- Defaults `"/"` to `/128.html`
- 600 s keepalive / idle timeout
- IPv6-ready (binds `[::]`)

## Build

```bash
podman build -t go-static-server .
```

The Dockerfile uses a 3-stage build:

1. **builder** — compiles the Go binary and generates the HTML payload files
2. **certgen** — produces self-signed TLS certificates using `ssl/gencert.go`
3. **scratch** — final image containing only the static binary and certs (~9 MB)

## Run

```bash
podman run --rm -p 8080:8080 -p 8443:8443 go-static-server
```

### Custom TLS certificates

Mount your own cert/key and set the environment variables:

```bash
podman run --rm \
  -v /path/to/cert.pem:/tls/cert.pem:ro \
  -v /path/to/key.pem:/tls/key.pem:ro \
  -e TLS_CERT=/tls/cert.pem \
  -e TLS_KEY=/tls/key.pem \
  -p 8080:8080 -p 8443:8443 \
  go-static-server
```

## Generating HTML files locally

```bash
./generate_html.sh          # writes to html/
./generate_html.sh /tmp/out # custom output directory
```

## Comparison with nginx image

| Feature | nginx image | go-static-server |
|---|---|---|
| Base image | UBI 9 Minimal | `scratch` (empty) |
| Image size | ~90 MB | ~9 MB |
| Dependencies | nginx, microdnf, UBI packages | None (static binary) |
| TLS certificate generation | Pre-baked key/cert files | Generated at build time via Go |
| Configuration | nginx.conf, entrypoint.sh | Environment variables only |
| CVE surface | OS packages + nginx | Zero (no OS, no libc) |
