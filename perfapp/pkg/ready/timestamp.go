package ready

import (
	"fmt"
	"net/http"
	"time"

	"ocp.performance.io/perfapp/internal/perf"
	log "github.com/sirupsen/logrus"
)

// Tables Euler workload required tables
var Tables = map[string]string{"ts": "CREATE TABLE IF NOT EXISTS ts (date TIMESTAMP)"}

// Handler Handle timestamp requests
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Info("Inserting timestamp record in table")
	now := time.Now()
	insert := fmt.Sprintf("INSERT INTO ts VALUES ('%s')", now.Format(time.RFC3339))
	if err := perf.QueryDB(insert); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Fprintln(w, "Ok")
		log.Printf("Timestamp inserted in %v ns", time.Since(now).Nanoseconds())
		perf.HTTPRequestDuration.Observe(time.Since(now).Seconds())
	}
}
