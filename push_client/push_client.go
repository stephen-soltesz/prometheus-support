package main

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	completionTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_backup_last_completion_timestamp_seconds",
		Help: "The timestamp of the last completion of a DB backup, successful or not.",
	})
	successTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_backup_last_success_timestamp_seconds",
		Help: "The timestamp of the last successful completion of a DB backup.",
	})
	duration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_backup_duration_seconds",
		Help: "The duration of the last DB backup in seconds.",
	})
	records = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_backup_records_processed",
		Help: "The number of records processed in the last DB backup.",
	})
)

func performBackup() (int, error) {
	// Perform the backup and return the number of backed up records and any
	// applicable error.
	// ...
	return 42, nil
}

func main() {
	// registry := prometheus.NewRegistry()
	prometheus.MustRegister(completionTime, duration, records)
	// Note that successTime is not registered at this time.

	start := time.Now()

	n, err := performBackup()

	records.Set(float64(n))
	duration.Set(time.Since(start).Seconds())
	completionTime.SetToCurrentTime()
	if err != nil {
		fmt.Println("DB backup failed:", err)
	} else {
		// Only now register successTime.
		prometheus.MustRegister(successTime)
		successTime.SetToCurrentTime()
	}
	// AddFromGatherer is used here rather than FromGatherer to not delete a
	// previously pushed success timestamp in case of a failure of this
	// backup.
	if err := push.FromGatherer(
		"db_backup", push.HostnameGroupingKey(),
		"http://35.184.247.140:9091",
		prometheus.DefaultGatherer,
	); err != nil {
		fmt.Println("Could not push to Pushgateway:", err)
	}
	fmt.Println("vim-go")
}
