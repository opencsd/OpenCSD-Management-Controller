package collector

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	client "github.com/influxdata/influxdb/client/v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var InstanceMetricCollector_ *MetricCollector

var (
	OPENCSD_INSTANCE_METRIC_COLLECTOR_PORT = "40804"
	OPENCSD_CONTROLLER_PORT                = "40801"
)

var (
	INFLUX_CLIENT   client.HTTPClient
	INFLUX_PORT     = os.Getenv("INFLUXDB_PORT")
	INFLUX_USERNAME = os.Getenv("INFLUXDB_USER")
	INFLUX_PASSWORD = os.Getenv("INFLUXDB_PASSWORD")
	INFLUX_DB       = os.Getenv("INFLUXDB_DB")
)

const (
	READY    = "READY"
	NOTREADY = "NOTREADY"
	BROKEN   = "BROKEN"
	NORMAL   = "NORMAL"
)

type MetricCollector struct {
	NodeName           string
	NodeMetric         *NodeMetric
	InstanceMetric     map[string]*InstanceMetric
	Imutex             sync.Mutex
	config             *Config
	statSummaryRequest *http.Request
}

func NewMetricCollector() *MetricCollector {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("cannot get hostname:", err)
		hostname = ""
	}

	NodeMetric := NewNodeMetric()
	NodeMetric.InitNodeMetric()

	nodeName := os.Getenv("NODE_NAME")
	if nodeName == "" {
		fmt.Println("NODE_NAME environment variable is not set")
	}

	nodeIP := os.Getenv("NODE_IP")

	config := NewConfig()
	token := config.Config.BearerToken

	scheme := "https"
	url := url.URL{
		Scheme: scheme,
		Host:   net.JoinHostPort(nodeIP, strconv.Itoa(10250)),
		Path:   "/stats/summary",
	}

	statSummaryRequest, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		fmt.Errorf("failed to get stat summary request")
	}

	statSummaryRequest.Header.Set("Content-Type", "application/json")
	statSummaryRequest.Header.Set("Authorization", "Bearer "+token)

	return &MetricCollector{
		NodeName:           hostname,
		NodeMetric:         NodeMetric,
		InstanceMetric:     make(map[string]*InstanceMetric),
		config:             NewConfig(),
		statSummaryRequest: statSummaryRequest,
	}
}

type Config struct {
	Config    *rest.Config
	Clientset *kubernetes.Clientset
	ClusterIP string
}

func NewConfig() *Config {
	hostConfig, _ := rest.InClusterConfig()
	hostKubeClient := kubernetes.NewForConfigOrDie(hostConfig)

	return &Config{
		Config:    hostConfig,
		Clientset: hostKubeClient,
		ClusterIP: hostConfig.Host,
	}
}

type NodeMetric struct {
	mutex   sync.Mutex
	Cpu     Cpu
	Memory  Memory
	Disk    Disk
	Network Network
	Power   Power
}

func NewNodeMetric() *NodeMetric {
	return &NodeMetric{
		Cpu:     Cpu{},
		Memory:  Memory{},
		Disk:    Disk{},
		Network: Network{},
		Power:   Power{},
	}
}

type InstanceMetric struct {
	InstanceName   string `json:"instanceName"`
	CpuUsage       int64  `json:"cpuUsage"`
	MemoryUsage    int64  `json:"memoryUsage"`
	StorageUsage   int64  `json:"storageUsage"`
	NetworkRxUsage int64  `json:"networkRxUsage"`
	NetworkTxUsage int64  `json:"networkTxUsage"`
}

type Cpu struct {
	Total       int
	Used        float64
	Utilization float64
	StJiffies   StJiffies
}

type StJiffies struct {
	User   int
	Nice   int
	System int
	Idle   int
}

type Memory struct {
	Total       int64
	Used        int64
	Utilization float64
	Free        int64
	Buffers     int64
	Cached      int64
}

type Disk struct {
	Name        string
	Total       int64
	Used        int64
	Utilization float64
}

type Network struct {
	RxByte    int64
	TxByte    int64
	RxData    int64
	TxData    int64
	Bandwidth int64
}

type Power struct {
	Energy1 int64
	Energy2 int64
	Used    int64
}

type Summary struct {
	Pods []PodStats `json:"pods"`
}

type PodStats struct {
	PodRef           PodReference  `json:"podRef"`
	CPU              *CPUStats     `json:"cpu,omitempty"`
	Memory           *MemoryStats  `json:"memory,omitempty"`
	Network          *NetworkStats `json:"network,omitempty"`
	EphemeralStorage *FsStats      `json:"ephemeral-storage,omitempty"`
}

type NetworkStats struct {
	Interfaces []InterfaceStats `json:"interfaces,omitempty"`
}

type InterfaceStats struct {
	Name     string  `json:"name"`
	RxBytes  *uint64 `json:"rxBytes,omitempty"`
	RxErrors *uint64 `json:"rxErrors,omitempty"`
	TxBytes  *uint64 `json:"txBytes,omitempty"`
	TxErrors *uint64 `json:"txErrors,omitempty"`
}

type PodReference struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type CPUStats struct {
	Time                 v1.Time `json:"time"`
	UsageNanoCores       *uint64 `json:"usageNanoCores,omitempty"`
	UsageCoreNanoSeconds *uint64 `json:"usageCoreNanoSeconds,omitempty"`
}

type MemoryStats struct {
	Time            v1.Time `json:"time"`
	AvailableBytes  *uint64 `json:"availableBytes,omitempty"`
	UsageBytes      *uint64 `json:"usageBytes,omitempty"`
	WorkingSetBytes *uint64 `json:"workingSetBytes,omitempty"`
	RSSBytes        *uint64 `json:"rssBytes,omitempty"`
	PageFaults      *uint64 `json:"pageFaults,omitempty"`
	MajorPageFaults *uint64 `json:"majorPageFaults,omitempty"`
}

type FsStats struct {
	Time           v1.Time `json:"time"`
	AvailableBytes *uint64 `json:"availableBytes,omitempty"`
	CapacityBytes  *uint64 `json:"capacityBytes,omitempty"`
	UsedBytes      *uint64 `json:"usedBytes,omitempty"`
	InodesFree     *uint64 `json:"inodesFree,omitempty"`
	Inodes         *uint64 `json:"inodes,omitempty"`
	InodesUsed     *uint64 `json:"inodesUsed,omitempty"`
}

func (nodeMetric *NodeMetric) InitNodeMetric() {
	nodeMetric.mutex.Lock()
	defer nodeMetric.mutex.Unlock()

	{
		cmd := exec.Command("grep", "-c", "processor", "/proc/cpuinfo")
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error: Command execution failed:", err)
		} else {
			coreCountStr := strings.TrimSpace(string(output))
			coreCount, err := strconv.Atoi(coreCountStr)
			if err != nil {
				fmt.Println("Error: Failed to convert core count to integer:", err)
			} else {
				nodeMetric.Cpu.Total = coreCount
			}
		}
	}

	{
		file, err := os.Open("/proc/stat")
		if err != nil {
			fmt.Println("cannot open file: ", err)
		} else {
			var cpuID string
			_, err = fmt.Fscanf(file, "%5s %d %d %d %d", &cpuID, &nodeMetric.Cpu.StJiffies.User, &nodeMetric.Cpu.StJiffies.Nice, &nodeMetric.Cpu.StJiffies.System, &nodeMetric.Cpu.StJiffies.Idle)
			if err != nil {
				fmt.Println("Error reading data from file:", err)
			}
		}
		file.Close()
	}

	{
		file, err := os.Open("/proc/meminfo")
		if err != nil {
			fmt.Println("cannot open file: ", err)
		} else {
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()

				if strings.HasPrefix(line, "MemTotal:") {
					fields := strings.Fields(line)
					if len(fields) >= 2 {
						nodeMetric.Memory.Total, err = strconv.ParseInt(fields[1], 10, 64)
						if err != nil {
							fmt.Println("Error parsing memory value:", err)
						}
					}
					break
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading file:", err)
			}
		}
		file.Close()
	}

	{
		statisticsFilePath := ""

		if _, err := os.Stat("/sys/class/net/eno1/statistics/"); os.IsNotExist(err) {
			statisticsFilePath = "/sys/class/net/enp96s0f0/statistics/"
		} else {
			statisticsFilePath = "/sys/class/net/eno1/statistics/"
		}

		rxBytesFieldName := statisticsFilePath + "rx_bytes"
		txBytesFieldName := statisticsFilePath + "tx_bytes"

		rxBytes, err := readStatisticsField(rxBytesFieldName)
		if err != nil {
			fmt.Println(err)
			return
		}

		txBytes, err := readStatisticsField(txBytesFieldName)
		if err != nil {
			fmt.Println(err)
			return
		}

		nodeMetric.Network.RxByte, _ = strconv.ParseInt(rxBytes, 10, 64)
		nodeMetric.Network.TxByte, _ = strconv.ParseInt(txBytes, 10, 64)
	}

	{
		cmd := exec.Command("df", "-k", "--total")
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error executing command:", err)
			return
		}

		scanner := bufio.NewScanner(bytes.NewReader(output))
		scanner.Scan()

		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, "total") {
				fields := strings.Fields(line)
				if len(fields) >= 3 {
					nodeMetric.Disk.Total, _ = strconv.ParseInt(fields[1], 10, 64)
					break
				}
			}
		}
	}
}

func readStatisticsField(fieldName string) (string, error) {
	data, err := os.ReadFile(fieldName)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %v", fieldName, err)
	}

	value := strings.TrimSpace(string(data))
	return value, nil
}

func extractCSDId(addr string) string {
	parts := strings.Split(addr, ".")
	if len(parts) > 0 {
		id := parts[2]
		return id
	}
	return ""
}

func convertSizeToMB(sizeStr string) int64 {
	unit := sizeStr[len(sizeStr)-1:]
	sizeValue, _ := strconv.ParseFloat(sizeStr[:len(sizeStr)-1], 64)
	switch unit {
	case "T":
		return int64(sizeValue * 1024 * 1024)
	case "G":
		return int64(sizeValue * 1024)
	case "M":
		return int64(sizeValue)
	default:
		return 0
	}
}
