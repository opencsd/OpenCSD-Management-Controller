package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"instance-metric-collector/src/collector"

	"github.com/influxdata/influxdb/client/v2"
)

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
