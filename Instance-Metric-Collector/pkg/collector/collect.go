package collector

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sync"
	"time"
	"strings"
	"strconv"
	"os"
	"log"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/wait"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"

	client "github.com/influxdata/influxdb/client/v2"

	//linuxproc "github.com/c9s/goprocinfo/linux"

	// "database/sql"
	// _ "fmt"
	// "log"
	// "strconv"
	// _ "context"
	// _ "encoding/json"
	// "types"
	// "client"

	//mysql 
	// "github.com/go-sql-driver/mysql"
	// _ "github.com/docker/docker/api/types"
	// "github.com/docker/docker/client"
	// "github.com/tecbot/gorocksdb"

	// influxdb v1
	// client "github.com/influxdata/influxdb/client/v2"
)

const (
	Core1_Energy_File = "/mnt/power/intel-rapl:0/energy_uj"
	Core2_Energy_File = "/mnt/power/intel-rapl:1/energy_uj"
)

var (
	OpenCSD_METRIC_COLLECTOR_IP = os.Getenv("OpenCSD_METRIC_COLLECTOR_IP")
    OpenCSD_METRIC_COLLECTOR_PORT = os.Getenv("OpenCSD_METRIC_COLLECTOR_PORT")
	NODE_METRIC_DB_NAME = "node_metric"
	INSTANCE_METRIC_DB_NAME = "keti_opencsd"
	power_before = 0.0
	cpu_before = 0
)

var curJiffies stJiffies
var prevJiffies stJiffies
var diffJiffies stJiffies
var cpu_usage_tick int64

type NodeMetric struct {
	Name				string 		`json:"name"`
	IP           	 	string     	`json:"ip"`
	Status				string 		`json:"status"`
	CpuUsageCore		int64 		`json:"cpuusagecore"`
	CpuUsageTick		int64 		`json:"cpuusagetick"`
	CpuTotal	  		int64 		`json:"cputotal"`
	MemUsage		  	int64  		`json:"memusage"`
	MemTotal	  		int64 		`json:"memtotal"`
	DiskUsage		  	int64  	   	`json:"diskusage"`
	DiskTotal  			int64 		`json:"disktotal"`
	NetworkRx	  		int64     	`json:"networkrx"`
	NetworkTx	  		int64     	`json:"networktx"`
	PowerUsage		  	int64  	   	`json:"powerusage"`
}

type PodMetric struct {
	Name				string 		`json:"name"`
	Namespace			string 		`json:"namespace"`
	CpuUsage		  	int64 		`json:"cpuusage"`
	MemUsage		  	int64  		`json:"memusage"`
	DiskUsage		  	int64  	   	`json:"diskusage"`
	NetworkRx	  		int64     	`json:"networkrx"`
	NetworkTx	  		int64     	`json:"networktx"`
}

type stJiffies struct {
    user 		int64
    nice		int64
    system		int64
    idle		int64
}

func Power_Consume(file_path string) float64 {
	energy_file, err := ioutil.ReadFile(file_path)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to open energy file")
		return 0.0
	}

	energy := strings.TrimRight(string(energy_file), "\n")
	energy_float, err := strconv.ParseFloat(energy, 64)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to parse energy value")
		return 0.0
	}

	energy_float /= 1000000.0
	return energy_float
}

func nodeMetricInsert(metricData *NodeMetric){
	INFLUX_IP := os.Getenv("INFLUX_IP") 
	INFLUX_PORT := os.Getenv("INFLUX_PORT")
    INFLUX_USERNAME := os.Getenv("INFLUX_USERNAME")
    INFLUX_PASSWORD := os.Getenv("INFLUX_PASSWORD")

	c, err := client.NewHTTPClient(client.HTTPConfig{ // InfluxDB 연결
		Addr: "http://" + INFLUX_IP + ":" + INFLUX_PORT, 
		Username: INFLUX_USERNAME,
		Password: INFLUX_PASSWORD,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{ // 배치 포인트 생성
		Database:  NODE_METRIC_DB_NAME, // DB 이름
 		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	//NODE
	var INFLUXDB_NODE_MEASUREMENT = "node_monitoring" 

	fields := map[string]interface{}{
		"cpu_usage_nanocore": metricData.CpuUsageCore,
		"cpu_usage_tick": metricData.CpuUsageTick,
		"memory_usage": metricData.MemUsage,
		"disk_usage": metricData.DiskUsage,
		"network_rx_byte": metricData.NetworkRx,
		"network_tx_byte": metricData.NetworkTx,
		"power_usage": metricData.PowerUsage,
	}

	pt, err := client.NewPoint(INFLUXDB_NODE_MEASUREMENT, nil, fields, time.Now()) 
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt) // 생성해둔 배치 포인트에 새로운 데이터 추가 

	err = c.Write(bp) //influxdb에 write
	if err != nil {
		log.Fatal(err)
	}
}

func podMetricInsert(podmetricData *PodMetric){
	INFLUX_IP := os.Getenv("INFLUX_IP") 
	INFLUX_PORT := os.Getenv("INFLUX_PORT")
    INFLUX_USERNAME := os.Getenv("INFLUX_USERNAME")
    INFLUX_PASSWORD := os.Getenv("INFLUX_PASSWORD")

	c, err := client.NewHTTPClient(client.HTTPConfig{ // InfluxDB 연결
		Addr: "http://" + INFLUX_IP + ":" + INFLUX_PORT, 
		Username: INFLUX_USERNAME,
		Password: INFLUX_PASSWORD,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{ // 배치 포인트 생성
		Database:  INSTANCE_METRIC_DB_NAME, // DB 이름
 		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	//POD
	var INFLUXDB_POD_MEASUREMENT = "pod_monitoring" 

	pod_fields := map[string]interface{}{
		"cpu_usage": podmetricData.CpuUsage,
		"memory_usage": podmetricData.MemUsage,
		"disk_usage": podmetricData.DiskUsage,
		"network_rx_usage": podmetricData.NetworkRx,
		"network_tx_usage": podmetricData.NetworkTx,
	}

	pt, err := client.NewPoint(INFLUXDB_POD_MEASUREMENT, nil, pod_fields, time.Now()) 
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt) // 생성해둔 배치 포인트에 새로운 데이터 추가 

	err = c.Write(bp) //influxdb에 write
	if err != nil {
		log.Fatal(err)
	}
}

func (m *MetricCollector) RunMetricCollector(ctx context.Context, wg *sync.WaitGroup) {
	go wait.UntilWithContext(ctx, m.MetricCollectingCycle, 5 * time.Second)
}

func (m *MetricCollector) MetricCollectingCycle(ctx context.Context) {
	
	// // SSD monitoring
	// var dbInfo DbInfo
	// dbInfo.MYSQL_USERNAME = "keti"
	// dbInfo.MYSQL_PASSWORD = "ketilinux"
	// dbInfo.MYSQL_IP = "10.0.4.80"
	// dbInfo.MYSQL_PORT = "3306"
	// dbInfo.MYSQL_DB = "information_schema"

	// db, err := sql.Open("mysql", makeDbUrl(dbInfo))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// // SSDMetricCollector 인스턴스 생성
	// ssdCollector := SSDMetricCollector{} //초기값으로 구조체?
	// var ssdMetricData SSDMetricData
	// ssdCollector.DDLInit(db, &ssdMetricData)

	// ssdCollector.getConnClientCnt(db, &ssdMetricData)
	// ssdCollector.getDDLCnt(db,&ssdMetricData)
	// ssdCollector.getCacheHitRatio(db, &ssdMetricData)
	// ssdCollector.getCacheUsageRatio(db, &ssdMetricData)
	// ssdCollector.getNetworkRatio(db, &ssdMetricData)
	// ssdCollector.getDiskUsageRatio(db, &ssdMetricData)
	// ssdCollector.getDDLCnt(db, &ssdMetricData)
	// SSDMetricInsert(&ssdMetricData)
	// fmt.Println(ssdMetricData)
	

	if summary, err := m.getStats(); err == nil {
		//NODE

		//CPU Tick
		stat, err := os.ReadFile("/mnt/cpu/stat")
		if err != nil {
			fmt.Println("err : cpu stat read fail")
		}
		lines := strings.Split(string(stat), "\n")

		cpu_core_num := 0.0
		cpu_usage_percent := 0.0
		cpu_usage_tick = 0
		var totalJiffies int64

		for _, line := range lines {
			fields := strings.Fields(line)

			if len(fields) > 0 {
				if fields[0] == "cpu"{
					if prevJiffies.user != 0 {
						cpu_before = 1
					}
					curJiffies.user, _ = strconv.ParseInt(fields[1], 10, 64)
					curJiffies.nice, _ = strconv.ParseInt(fields[2], 10, 64)
					curJiffies.system, _ = strconv.ParseInt(fields[3], 10, 64)
					curJiffies.idle, _ = strconv.ParseInt(fields[4], 10, 64)
			
					fmt.Println("curJiffies : ", curJiffies)
					fmt.Println("prevJiffies : ", prevJiffies)
					fmt.Println("diffJiffies : ", diffJiffies)
			
					diffJiffies.user =  curJiffies.user - prevJiffies.user
					diffJiffies.nice =  curJiffies.nice - prevJiffies.nice
					diffJiffies.system =  curJiffies.system - prevJiffies.system
					diffJiffies.idle =  curJiffies.idle - prevJiffies.idle
			
					totalJiffies = diffJiffies.user + diffJiffies.nice + diffJiffies.system + diffJiffies.idle
			
					tmp := float64(diffJiffies.idle) / float64(totalJiffies)
					cpu_usage_percent =  1 - tmp

					fmt.Println("cpu_usage_percent : ", cpu_usage_percent)

					prevJiffies = curJiffies

				}else if strings.HasPrefix(fields[0],"cpu") {
					cpu_core_num++
				}
			}
		}

		if cpu_before == 1 {
			cpu_usage := cpu_core_num * cpu_usage_percent
			cpu_usage_tick = totalJiffies - diffJiffies.idle
			fmt.Println("cpu_core_num : ", cpu_core_num)
			fmt.Println("cpu_usage_core : ", cpu_usage)
			fmt.Println("cpu_usage_tick : ", cpu_usage_tick)
		}

		//NETWORK
		var RX_Usage uint64 = 0
		var TX_Usage uint64 = 0

		for _, Interface := range summary.Node.Network.Interfaces {
			RX_Usage = RX_Usage + *Interface.RxBytes
			TX_Usage = TX_Usage + *Interface.TxBytes
		}

		var networkRXBytes resource.Quantity
		var networkTXBytes resource.Quantity

		networkRXBytes = *uint64Quantity(RX_Usage, 0)
		networkRXBytes.Format = resource.BinarySI

		networkTXBytes = *uint64Quantity(TX_Usage, 0)
		networkTXBytes.Format = resource.BinarySI

		//POWER
		var power_current float64
		var node_power float64
		energy1 := Power_Consume(Core1_Energy_File)
		energy2 := Power_Consume(Core2_Energy_File)
		power_current = energy1 + energy2
		if power_before == 0 {
			power_before = power_current
			node_power = 0
		}else {
			if power_current >= power_before {
				node_power = (power_current - power_before)/5
				power_before = power_current
			} else {
				power_before = power_current
				node_power = 0
			}
		}

		//CPU MEM DISK
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		nodes, _ := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})

		//INSERT NODE METRIC DATA
		var node_metric NodeMetric
		nodename := os.Getenv("NODE_NAME")

		for _, node := range nodes.Items {
			if node.Name == nodename {
				cpu_quan := node.Status.Capacity["cpu"]
				cpu_quan_string := fmt.Sprintf("%+v", cpu_quan)
				cpu_s_trim := strings.TrimLeft(cpu_quan_string,"{")
				cpu_s_split_1 := strings.Split(cpu_s_trim, " ")
				cpu_total := strings.Split(cpu_s_split_1[0], ":")
				cpu_total_int, _ := strconv.ParseInt(cpu_total[len(cpu_total)-1], 10, 64)

				mem_quan := node.Status.Capacity["memory"]
				mem_quan_string := fmt.Sprintf("%+v", mem_quan)
				mem_s_trim := strings.TrimLeft(mem_quan_string,"{")
				mem_s_split_1 := strings.Split(mem_s_trim, " ")
				mem_total := strings.Split(mem_s_split_1[0], ":")
				mem_total_int, _ := strconv.ParseInt(mem_total[len(mem_total)-1], 10, 64)

				fs_quan := node.Status.Capacity["ephemeral-storage"]
				fs_quan_string := fmt.Sprintf("%+v", fs_quan)
				fs_s_trim := strings.TrimLeft(fs_quan_string,"{")
				fs_s_split_1 := strings.Split(fs_s_trim, " ")
				fs_total := strings.Split(fs_s_split_1[0], ":")
				fs_total_int, _ := strconv.ParseInt(fs_total[len(fs_total)-1], 10, 64)

				//node_metric.ID = 1
				node_metric.Name = node.Name
				node_metric.IP = node.Status.Addresses[0].Address
				node_metric.Status = "Ready"
				node_metric.CpuUsageCore =  int64(*summary.Node.CPU.UsageNanoCores)
				node_metric.CpuUsageTick =  cpu_usage_tick
				node_metric.CpuTotal = cpu_total_int
				node_metric.MemUsage = int64(*summary.Node.Memory.UsageBytes)
				node_metric.MemTotal = mem_total_int
				node_metric.DiskUsage = int64(*summary.Node.Fs.UsedBytes)
				node_metric.DiskTotal = fs_total_int
				node_metric.NetworkRx, _ = networkRXBytes.AsInt64()
				node_metric.NetworkTx, _ = networkTXBytes.AsInt64()
				node_metric.PowerUsage = int64(node_power)

				fmt.Println("**nodemetric: ", &node_metric)
				nodeMetricInsert(&node_metric)

				break;
			}
		}

		//POD
		for _, pod := range summary.Pods {
			//pod_namespace := strings.Replace(INSTANCE_METRIC_DB_NAME, "_", "-", -1)
			pod_namespace := "keti-opencsd"
			pod_name_TF := strings.HasPrefix(pod.PodRef.Name, "storage-engine")
			if pod.PodRef.Namespace == pod_namespace && pod_name_TF {
				//podMetric := &metric.PodMetric{}
				var pod_metric PodMetric

				pod_metric.Name = pod.PodRef.Name
				pod_metric.Namespace = pod.PodRef.Namespace

//				if pod.CPU.UsageNanoCores != nil {
//					pod_metric.CpuUsage = int64(*pod.CPU.UsageNanoCores)
//				}

//				if pod.Memory.UsageBytes != nil {
//					pod_metric.MemUsage = int64(*pod.Memory.UsageBytes)
//				}

				var pods_cpu int64 = 0
				var pods_mem int64 = 0

				for _, container := range pod.Containers {
					
					if container.CPU.UsageNanoCores != nil {
						pods_cpu = pods_cpu + int64(*container.CPU.UsageNanoCores)
					}

					if container.Memory.WorkingSetBytes != nil {
						pods_mem = pods_mem + int64(*container.Memory.WorkingSetBytes)
					}
				}

				pod_metric.CpuUsage = pods_cpu
				pod_metric.MemUsage = pods_mem
				
				if pod.EphemeralStorage.UsedBytes != nil {
					pod_metric.DiskUsage = int64(*pod.EphemeralStorage.UsedBytes)
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

					pod_metric.NetworkRx, _ = networkRXBytes.AsInt64()
					pod_metric.NetworkTx, _ = networkTXBytes.AsInt64()
				}

				fmt.Println("**podmetric: ", &pod_metric)
				podMetricInsert(&pod_metric)
			
			}
		}

	} else {
		fmt.Println(err)
	}

}

func (m *MetricCollector) getStats() (*Summary, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: transport,
	}

	response, err := client.Do(m.StatRequest)
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