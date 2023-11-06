package main

import (
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"ocp.performance.io/perfapp/internal/perf"
	"ocp.performance.io/perfapp/pkg/euler"
	"ocp.performance.io/perfapp/pkg/health"
	"ocp.performance.io/perfapp/pkg/ready"
	"ocp.performance.io/perfapp/pkg/utils"
)

var tables []map[string]string

func main() {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	go handleInterrupt(c)
	if os.Getenv("POSTGRESQL_HOSTNAME") != "" {
		go func() {
			if os.Getenv("POSTGRESQL_RETRY_INTERVAL") != "" {
				retryInt, err := strconv.Atoi(os.Getenv("POSTGRESQL_RETRY_INTERVAL"))
				if err != nil {
					utils.ErrorHandler(err)
				}
				perf.DB.RetryInt = retryInt
			}
			perf.Connect2Db()
			tables = append(tables, euler.Tables, ready.Tables)
			if err := perf.CreateTables(tables); err != nil {
				utils.ErrorHandler(err)
			}
		}()
		http.HandleFunc("/euler", euler.Handler)
		http.HandleFunc("/ready", ready.Handler)
	}
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", health.Handler)
	log.Printf("Listening at 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		utils.ErrorHandler(err)
	}
}

func handleInterrupt(c <-chan os.Signal) {
	<-c
	log.Println("Interrupt signal received")
	os.Exit(0)
}
