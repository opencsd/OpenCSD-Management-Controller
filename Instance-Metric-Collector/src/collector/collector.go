package collector

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	// influxdb v1 client
	client "github.com/influxdata/influxdb/client/v2"
	"k8s.io/apimachinery/pkg/api/resource"
)

func (instanceMetricCollector *MetricCollector) RunNodeMetricRoutine() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			instanceMetricCollector.NodeMetric.mutex.Lock()

			instanceMetricCollector.updateCpu()
			instanceMetricCollector.updateMemory()
			instanceMetricCollector.updateNetwork()
			// instanceMetricCollector.updateStorage()
			instanceMetricCollector.updatePower()

			instanceMetricCollector.saveNodeMetric()

			instanceMetricCollector.NodeMetric.mutex.Unlock()
		}
	}
}

func (instanceMetricCollector *MetricCollector) updateCpu() {
	file, err := os.Open("/host/proc/stat")
	if err != nil {
		fmt.Println("cannot open file: ", err)
	} else {
		var cpuID string

		var curJiffies, diffJiffies StJiffies

		_, err = fmt.Fscanf(file, "%5s %d %d %d %d", &cpuID, &curJiffies.User, &curJiffies.Nice, &curJiffies.System, &curJiffies.Idle)
		if err != nil {
			fmt.Println("Error reading data from file:", err)
		}

		diffJiffies.User = curJiffies.User - instanceMetricCollector.NodeMetric.Cpu.StJiffies.User
		diffJiffies.Nice = curJiffies.Nice - instanceMetricCollector.NodeMetric.Cpu.StJiffies.Nice
		diffJiffies.System = curJiffies.System - instanceMetricCollector.NodeMetric.Cpu.StJiffies.System
		diffJiffies.Idle = curJiffies.Idle - instanceMetricCollector.NodeMetric.Cpu.StJiffies.Idle

		totalJiffies := diffJiffies.User + diffJiffies.Nice + diffJiffies.System + diffJiffies.Idle

		utilization := 100.0 * (1.0 - float64(diffJiffies.Idle)/float64(totalJiffies))
		instanceMetricCollector.NodeMetric.Cpu.Utilization, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", utilization), 64)
		used := float64(instanceMetricCollector.NodeMetric.Cpu.Total) * (1.0 - float64(diffJiffies.Idle)/float64(totalJiffies))
		instanceMetricCollector.NodeMetric.Cpu.Used, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", used), 64)

		instanceMetricCollector.NodeMetric.Cpu.StJiffies = curJiffies
	}
	file.Close()
}

func (instanceMetricCollector *MetricCollector) updateMemory() {
	file, err := os.Open("/host/proc/meminfo")
	if err != nil {
		fmt.Println("cannot open file: ", err)
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)

			if len(fields) < 2 {
				continue
			}

			key := fields[0]
			value, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				fmt.Println("Error parsing value:", err)
				continue
			}

			switch key {
			case "MemFree:":
				instanceMetricCollector.NodeMetric.Memory.Free = value
			case "Buffers:":
				instanceMetricCollector.NodeMetric.Memory.Buffers = value
			case "Cached:":
				instanceMetricCollector.NodeMetric.Memory.Cached = value
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
		}

		freeGB := (instanceMetricCollector.NodeMetric.Memory.Free + instanceMetricCollector.NodeMetric.Memory.Buffers + instanceMetricCollector.NodeMetric.Memory.Cached) / 1024.0 / 1024.0
		instanceMetricCollector.NodeMetric.Memory.Used = instanceMetricCollector.NodeMetric.Memory.Total - freeGB
		utilization := float64(instanceMetricCollector.NodeMetric.Memory.Used) / float64(instanceMetricCollector.NodeMetric.Memory.Total) * 100.0
		instanceMetricCollector.NodeMetric.Memory.Utilization, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", utilization), 64)
	}
	file.Close()
}

func (instanceMetricCollector *MetricCollector) updateNetwork() {
	statisticsFilePath := ""

	if _, err := os.Stat("/host/sys/class/net/eno1/statistics/"); os.IsNotExist(err) {
		statisticsFilePath = "/host/sys/class/net/enp96s0f0/statistics/"
	} else {
		statisticsFilePath = "/host/sys/class/net/eno1/statistics/"
	}

	rxBytesFieldName := statisticsFilePath + "rx_bytes"
	txBytesFieldName := statisticsFilePath + "tx_bytes"

	currentRxBytesStr, err := readStatisticsField(rxBytesFieldName)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentTxBytesStr, err := readStatisticsField(txBytesFieldName)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentRxBytes, _ := strconv.ParseInt(currentRxBytesStr, 10, 64)
	currentTxBytes, _ := strconv.ParseInt(currentTxBytesStr, 10, 64)

	instanceMetricCollector.NodeMetric.Network.RxData = currentRxBytes - instanceMetricCollector.NodeMetric.Network.RxByte
	instanceMetricCollector.NodeMetric.Network.TxData = currentTxBytes - instanceMetricCollector.NodeMetric.Network.TxByte

	instanceMetricCollector.NodeMetric.Network.Bandwidth = (instanceMetricCollector.NodeMetric.Network.RxData + instanceMetricCollector.NodeMetric.Network.TxData) / 5 * 8

	instanceMetricCollector.NodeMetric.Network.RxByte = currentRxBytes
	instanceMetricCollector.NodeMetric.Network.TxByte = currentTxBytes
}

// func (instanceMetricCollector *MetricCollector) updateStorage() {
// 	cmd := exec.Command("df", "-k", "--total")
// 	output, err := cmd.Output()
// 	if err != nil {
// 		fmt.Println("Error executing command:", err)
// 		return
// 	}

// 	scanner := bufio.NewScanner(bytes.NewReader(output))
// 	scanner.Scan()

// 	for scanner.Scan() {
// 		line := scanner.Text()

// 		if strings.Contains(line, "total") {
// 			fields := strings.Fields(line)
// 			if len(fields) >= 3 {
// 				instanceMetricCollector.NodeMetric.Disk.Used, _ = strconv.ParseFloat(fields[2], 64)
// 				break
// 			}
// 		}
// 	}

// 	if instanceMetricCollector.NodeMetric.Disk.Total > 0 {
// 		utilization := (float64(instanceMetricCollector.NodeMetric.Disk.Used) / float64(instanceMetricCollector.NodeMetric.Disk.Total)) * 100
// 		instanceMetricCollector.NodeMetric.Disk.Utilization, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", utilization), 64)
// 	}
// }

func (instanceMetricCollector *MetricCollector) updatePower() {
	energyFieldName1 := "/host/sys/class/powercap/intel-rapl:0/energy_uj"
	energyFieldName2 := "/host/sys/class/powercap/intel-rapl:1/energy_uj"

	currentEnergyStr1, err := readStatisticsField(energyFieldName1)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentEnergyStr2, err := readStatisticsField(energyFieldName2)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentEnergy1, _ := strconv.ParseInt(currentEnergyStr1, 10, 64)
	currentEnergy2, _ := strconv.ParseInt(currentEnergyStr2, 10, 64)

	energyDiffJ1 := float64(currentEnergy1-instanceMetricCollector.NodeMetric.Power.Energy1) / 1e6
	energyDiffJ2 := float64(currentEnergy2-instanceMetricCollector.NodeMetric.Power.Energy2) / 1e6

	instanceMetricCollector.NodeMetric.Power.Used = int64((energyDiffJ1 + energyDiffJ2) / 1.0)
	instanceMetricCollector.NodeMetric.Power.Energy1 = currentEnergy1
	instanceMetricCollector.NodeMetric.Power.Energy2 = currentEnergy2
}

func (instanceMetricCollector *MetricCollector) saveNodeMetric() {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  INFLUX_DB,
		Precision: "s",
	})
	if err != nil {
		fmt.Println("DB NewBatchPoints Error:", err)
		return
	}

	var INFLUXDB_NODE_MEASUREMENT = "node_metric"

	fields := map[string]interface{}{
		"node_name": instanceMetricCollector.NodeName,

		"cpu_total":       instanceMetricCollector.NodeMetric.Cpu.Total,
		"cpu_usage":       instanceMetricCollector.NodeMetric.Cpu.Used,
		"cpu_utilization": instanceMetricCollector.NodeMetric.Cpu.Utilization,

		"memory_total":       instanceMetricCollector.NodeMetric.Memory.Total,
		"memory_usage":       instanceMetricCollector.NodeMetric.Memory.Used,
		"memory_utilization": instanceMetricCollector.NodeMetric.Memory.Utilization,

		// "disk_total":       instanceMetricCollector.NodeMetric.Disk.Total,
		// "disk_usage":       instanceMetricCollector.NodeMetric.Disk.Used,
		// "disk_utilization": instanceMetricCollector.NodeMetric.Disk.Utilization,

		"network_bandwidth": instanceMetricCollector.NodeMetric.Network.Bandwidth,
		"network_rx_data":   instanceMetricCollector.NodeMetric.Network.RxData,
		"network_tx_data":   instanceMetricCollector.NodeMetric.Network.TxData,

		"power_usage": instanceMetricCollector.NodeMetric.Power.Used,
	}

	pt, err := client.NewPoint(INFLUXDB_NODE_MEASUREMENT, nil, fields, time.Now())
	if err != nil {
		fmt.Println("DB NewPoint Error:", err)
		return
	}
	bp.AddPoint(pt)

	err = INFLUX_CLIENT.Write(bp)
	if err != nil {
		fmt.Println("DB Write Error:", err)
		return
	}
}

func (instanceMetricCollector *MetricCollector) RunInstanceMetricRoutine() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			instanceMetricCollector.Imutex.Lock()

			instanceMetricCollector.updateInstanceMetric()
			instanceMetricCollector.saveInstanceMetric()

			instanceMetricCollector.Imutex.Unlock()
		}
	}
}

func (instanceMetricCollector *MetricCollector) updateInstanceMetric() {
	if summary, err := instanceMetricCollector.getStats(); err == nil {
		for _, pod := range summary.Pods {
			podName := pod.PodRef.Name

			if strings.HasPrefix(podName, "storage-engine") {
				instanceName := pod.PodRef.Namespace

				instanceMetric := NewInstanceMetric(instanceName)

				for _, container := range pod.Containers {
					if container.CPU.UsageNanoCores != nil {
						instanceMetric.CpuUsage += int64(math.Ceil(float64(*container.CPU.UsageNanoCores) / 1000000))
					}

					if container.Memory.WorkingSetBytes != nil {
						instanceMetric.MemoryUsage += int64(float64(*container.Memory.WorkingSetBytes) / float64(1024*1024))
					}

					if pod.Network != nil {
						var RX_Usage uint64 = 0
						var TX_Usage uint64 = 0

						for _, Interface := range pod.Network.Interfaces {
							RX_Usage = RX_Usage + *Interface.RxBytes
							TX_Usage = TX_Usage + *Interface.TxBytes
						}

						var networkRXBytes resource.Quantity
						var networkTXBytes resource.Quantity

						networkRXBytes = *uint64Quantity(RX_Usage, 0)
						networkRXBytes.Format = resource.BinarySI

						networkTXBytes = *uint64Quantity(TX_Usage, 0)
						networkTXBytes.Format = resource.BinarySI

						instanceMetric.NetworkRxUsage, _ = networkRXBytes.AsInt64()
						instanceMetric.NetworkTxUsage, _ = networkTXBytes.AsInt64()
					}
				}

				if pod.EphemeralStorage.UsedBytes != nil {
					instanceMetric.StorageUsage = int64(*pod.EphemeralStorage.UsedBytes)
				}

				instanceMetricCollector.InstanceMetric[instanceName] = instanceMetric
			}
		}
	}
}

func (instanceMetricCollector *MetricCollector) getStats() (*Summary, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: transport,
	}

	response, err := client.Do(instanceMetricCollector.statSummaryRequest)
	if err != nil {
		return nil, fmt.Errorf("[error] get node network stats error: %v", err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("[error] get node network stats error: %v", err)
	}

	summary := &Summary{}

	if err := json.Unmarshal(body, &summary); err != nil {
		return nil, fmt.Errorf("[error] get node network stats error: %v", err)
	}

	return summary, nil
}

func uint64Quantity(val uint64, scale resource.Scale) *resource.Quantity {
	if val <= math.MaxInt64 {
		return resource.NewScaledQuantity(int64(val), scale)
	}

	return resource.NewScaledQuantity(int64(val/10), resource.Scale(1)+scale)
}

func (instanceMetricCollector *MetricCollector) saveInstanceMetric() {
	for instanceName, instanceMetric := range instanceMetricCollector.InstanceMetric {
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  INFLUX_DB,
			Precision: "s",
		})
		if err != nil {
			fmt.Println("DB NewBatchPoints Error:", err)
			return
		}

		parsedInstanceName := strings.ReplaceAll(instanceName, "-", "_")
		var INFLUXDB_INSTANCE_MEASUREMENT = "instance_metric_" + parsedInstanceName

		fields := map[string]interface{}{
			"instance_name":   instanceMetric.InstanceName,
			"cpu_usage":       instanceMetric.CpuUsage,
			"memory_usage":    instanceMetric.MemoryUsage,
			"storage_usage":   instanceMetric.StorageUsage,
			"network_rx_data": instanceMetric.NetworkRxUsage,
			"network_tx_data": instanceMetric.NetworkTxUsage,
		}

		pt, err := client.NewPoint(INFLUXDB_INSTANCE_MEASUREMENT, nil, fields, time.Now())
		if err != nil {
			fmt.Println("DB NewPoint Error:", err)
			return
		}
		bp.AddPoint(pt)

		err = INFLUX_CLIENT.Write(bp)
		if err != nil {
			fmt.Println("DB Write Error:", err)
			return
		}
	}
}
