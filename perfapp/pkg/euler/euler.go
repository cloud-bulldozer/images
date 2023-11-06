package euler

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"ocp.performance.io/perfapp/internal/perf"
	log "github.com/sirupsen/logrus"
)

// Tables Euler workload required tables
var Tables = map[string]string{"euler": "CREATE TABLE IF NOT EXISTS euler (date TIMESTAMP, elapsed FLOAT(24))"}

// Handler Handle requests to compute euler number aproximation
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Computing euler approximation")
	now := time.Now()
	calcEuler()
	insert := fmt.Sprintf("INSERT INTO euler VALUES ('%s', '%f')", now.Format(time.RFC3339), time.Since(now).Seconds())
	if err := perf.QueryDB(insert); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, err.Error())
	} else {
		fmt.Fprintln(w, "Ok")
		log.Printf("Euler approximation computed in %f seconds", time.Since(now).Seconds())
		perf.HTTPRequestDuration.Observe(time.Since(now).Seconds())
	}
}

func calcEuler() {
	var n float64
	var x float64
	for math.E > x {
		x = math.Pow((1 + 1/n), n)
		n++
	}
}
