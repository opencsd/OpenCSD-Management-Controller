package handler

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"instance-metric-collector/src/collector"

	_ "github.com/go-sql-driver/mysql"
	"github.com/influxdata/influxdb/client/v2"
)

func CreateInstance(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Instance Metric Collector] CreateInstance")

	dbName := r.URL.Query().Get("instance")
	convertedDbName := strings.ReplaceAll(dbName, "-", "_")

	dsn := fmt.Sprintf("%s:%s@tcp(localhost:%s)/", collector.MYSQL_USERNAME, collector.MYSQL_ROOT_PASSWORD, collector.MYSQL_PORT)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Failed to connect to MySQL: %v", err)
		http.Error(w, "Failed to connect to database"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Printf("Failed to ping MySQL: %v", err)
		http.Error(w, "Database connection error"+err.Error(), http.StatusInternalServerError)
		return
	}

	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", convertedDbName)
	_, err = db.Exec(createDBQuery)
	if err != nil {
		fmt.Printf("Failed to create database: %v", err)
		http.Error(w, "Failed to create database"+err.Error(), http.StatusInternalServerError)
		return
	}
	db.Close()

	dsn = fmt.Sprintf("%s:%s@tcp(localhost:%s)/%s", collector.MYSQL_USERNAME, collector.MYSQL_ROOT_PASSWORD, collector.MYSQL_PORT, convertedDbName)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Failed to reconnect to MySQL database: %v", err)
		http.Error(w, "Failed to reconnect to database"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	sqlBytes, err := ioutil.ReadFile("/etc/config/base.sql")
	if err != nil {
		fmt.Printf("Failed to read SQL file: %v", err)
		http.Error(w, "Failed to read SQL file"+err.Error(), http.StatusInternalServerError)
		return
	}

	createTableQueries := string(sqlBytes)

	queries := strings.Split(createTableQueries, ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		_, err := db.Exec(query)
		if err != nil {
			fmt.Printf("Failed to execute query: %v\nQuery: %s\n", err, query)
			http.Error(w, "Failed to execute query: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Database and table created successfully\n"))
}

func DeleteInstance(w http.ResponseWriter, r *http.Request) {
	// InstanceMetricCollector,influx에서 instance measurement 삭제 : instance_metric_{instance namespace}

	fmt.Println("[OpenCSD Instance Metric Collector] DeleteInstance")

	instance := r.URL.Query().Get("instance")

	if _, exists := collector.InstanceMetricCollector_.InstanceMetric[instance]; exists {
		delete(collector.InstanceMetricCollector_.InstanceMetric, instance)
	}

	parsedInstanceName := strings.ReplaceAll(instance, "-", "_")
	var INFLUXDB_INSTANCE_MEASUREMENT = "instance_metric_" + parsedInstanceName

	query := fmt.Sprintf(`DROP MEASUREMENT "%s"`, INFLUXDB_INSTANCE_MEASUREMENT)

	q := client.NewQuery(query, collector.INFLUX_DB, "")
	if response, err := collector.INFLUX_CLIENT.Query(q); err == nil && response.Error() == nil {
		fmt.Printf("Measurement '%s' has been dropped successfully\n", INFLUXDB_INSTANCE_MEASUREMENT)
	} else {
		if err != nil {
			log.Fatalf("Error executing query: %v", err)
		} else {
			log.Fatalf("Query error: %v", response.Error())
		}
	}

	w.Write([]byte("[OpenCSD Instance Metric Collector] Delete Instance Successfully\n"))
}
