package health

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Handler Handles health endpoint
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Info("Returning value for health endpoint")
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}
