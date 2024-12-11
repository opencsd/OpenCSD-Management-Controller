package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"instance-metric-collector/src/collector"
	"instance-metric-collector/src/handler"

	"github.com/influxdata/influxdb/client/v2"
)

func main() {
	for {
		var err error
		collector.INFLUX_CLIENT, err = client.NewHTTPClient(client.HTTPConfig{
			Addr:     "http://localhost:" + collector.INFLUX_PORT,
			Username: collector.INFLUX_USERNAME,
			Password: collector.INFLUX_PASSWORD,
		})
		if err != nil {
			log.Fatal(err)
		}
		_, _, pingErr := collector.INFLUX_CLIENT.Ping(5 * time.Second)
		if pingErr == nil {
			fmt.Println("[OpenCSD Instance Metric Collector] Connected to InfluxDB!")
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	defer collector.INFLUX_CLIENT.Close()

	collector.InstanceMetricCollector_ = collector.NewMetricCollector()

	go collector.InstanceMetricCollector_.RunNodeMetricRoutine()
	go collector.InstanceMetricCollector_.RunInstanceMetricRoutine()

	fmt.Println("[OpenCSD Instance Metric Collector] run on 0.0.0.0:", collector.OPENCSD_INSTANCE_METRIC_COLLECTOR_PORT)

	http.HandleFunc("/create/instance", handler.CreateInstance)
	http.HandleFunc("/delete/instance", handler.DeleteInstance)

	http.ListenAndServe(":"+collector.OPENCSD_INSTANCE_METRIC_COLLECTOR_PORT, nil)
}
