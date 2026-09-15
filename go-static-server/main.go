// Package main implements a static file server that serves pre-generated
// HTML files over HTTP and HTTPS. It is a drop-in replacement for the
// nginx-based static server used in cloud-bulldozer benchmarks.
package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed html/*.html
var htmlFiles embed.FS

func main() {
	// Obtain the html/ subtree so paths resolve without the "html/" prefix.
	htmlFS, err := fs.Sub(htmlFiles, "html")
	if err != nil {
		log.Fatalf("failed to open embedded html subtree: %v", err)
	}

	fileServer := http.FileServerFS(htmlFS)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Default "/" to 128.html (the nginx index behavior).
		if r.URL.Path == "/" {
			r.URL.Path = "/128.html"
		}
		// Accept both GET and POST — serve the same static file.
		// Override the method to GET so the file server doesn't reject POST.
		r.Method = http.MethodGet

		// Serve index.html directly to prevent FileServer's automatic
		// redirect from /index.html → /.
		if r.URL.Path == "/index.html" {
			data, err := fs.ReadFile(htmlFS, "index.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
			return
		}

		fileServer.ServeHTTP(w, r)
	})

	// Listen addresses and TLS paths are configurable via environment
	// variables, with defaults matching the nginx image behavior.
	httpAddr := envOrDefault("HTTP_ADDR", "[::]:8080")
	httpsAddr := envOrDefault("HTTPS_ADDR", "[::]:8443")
	certFile := envOrDefault("TLS_CERT", "/etc/ssl/server.crt")
	keyFile := envOrDefault("TLS_KEY", "/etc/ssl/server.key")

	httpServer := &http.Server{
		Addr:        httpAddr,
		Handler:     handler,
		IdleTimeout: 600 * time.Second,
	}

	httpsServer := &http.Server{
		Addr:        httpsAddr,
		Handler:     handler,
		IdleTimeout: 600 * time.Second,
	}

	// Start HTTP in a goroutine.
	go func() {
		log.Printf("HTTP listening on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Start HTTPS on the main goroutine.
	log.Printf("HTTPS listening on %s (cert=%s, key=%s)", httpsServer.Addr, certFile, keyFile)
	if err := httpsServer.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTPS server error: %v", err)
	}
}

// envOrDefault returns the value of the named environment variable or
// fallback if it is empty / unset.
func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
