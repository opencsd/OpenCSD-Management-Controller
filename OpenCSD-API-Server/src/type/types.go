package types

import (
	"context"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var ManagementMaster_ *MasterNode

var (
	OPENCSD_API_SERVER_PORT                = "40800"
	STORAGE_API_SERVER_PORT                = "40306"
	OPENCSD_INSTANCE_METRIC_COLLECTOR_PORT = "40804"
	OPENCSD_CONTROLLER_DNS                 = "opencsd-controller-svc.management-controller.svc.cluster.local:40801"
)

var (
	INSTANCE_METRIC_INFLUXDB_DB       = os.Getenv("INFLUXDB_DB")
	INSTANCE_METRIC_INFLUXDB_PASSWORD = os.Getenv("INFLUXDB_PASSWORD")
	INSTANCE_METRIC_INFLUXDB_USER     = os.Getenv("INFLUXDB_USER")
	INSTANCE_METRIC_INFLUXDB_PORT     = os.Getenv("INFLUXDB_PORT")
)

var (
	INSTANCE_METRIC_MYSQL_USER          = "root"
	INSTANCE_METRIC_MYSQL_ROOT_PASSWORD = os.Getenv("MYSQL_ROOT_PASSWORD")
	INSTANCE_METRIC_MYSQL_PORT          = os.Getenv("MYSQL_PORT")
)

const (
	READY    = "READY"
	NOTREADY = "NOTREADY"
)

const (
	STORAGE   = "STORAGE"
	OPERATION = "OPERATION"
)

type MasterNode struct {
	ClusterName    string
	NodeName       string
	MasterIP       string
	MatserStatus   string
	clusterConfig  ClusterConfig
	StorageLayer   map[string]WorkerNode
	OperationLayer map[string]WorkerNode
}

type ClusterConfig struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
}

type WorkerNode struct {
	NodeName string `json:"nodeName"`
	NodeIP   string `json:"nodeIp"`
	Status   string `json:"status"`
	Layer    string `json:"layer"`
}

func (masterNode *MasterNode) InitCluster() {
	masterName := os.Getenv("NODE_NAME")
	if masterName == "" {
		fmt.Println("NODE_NAME environment variable is not set")
	}
	masterNode.NodeName = masterName
	masterNode.ClusterName = "OPENCSD"

	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("InClusterConfig error:", err)
		return
	}

	masterNode.clusterConfig.config = config
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("NewForConfig error:", err)
		return
	}

	masterNode.clusterConfig.clientset = clientset

	masterNode.StorageLayer = make(map[string]WorkerNode)
	masterNode.OperationLayer = make(map[string]WorkerNode)

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Get Node error:", err)
	} else {
		for _, node := range nodes.Items {
			if node.Name == masterName {
				for _, address := range node.Status.Addresses {
					if address.Type == "InternalIP" {
						masterNode.MasterIP = address.Address
						break
					}
				}

				for _, condition := range node.Status.Conditions {
					if condition.Type == corev1.NodeReady {
						masterNode.MatserStatus = READY
					} else {
						masterNode.MatserStatus = NOTREADY
					}
				}
			} else {
				var workerNode WorkerNode
				workerNode.NodeName = node.Name

				for _, condition := range node.Status.Conditions {
					if condition.Type == corev1.NodeReady {
						workerNode.Status = READY
					} else {
						workerNode.Status = NOTREADY
					}
				}

				for _, address := range node.Status.Addresses {
					if address.Type == "InternalIP" {
						workerNode.NodeIP = address.Address
						break
					}
				}

				if node.Labels["layer"] == "storage" {
					workerNode.Layer = STORAGE
					masterNode.StorageLayer[node.Name] = workerNode
				} else if node.Labels["layer"] == "operation" {
					workerNode.Layer = OPERATION
					masterNode.OperationLayer[node.Name] = workerNode
				}
			}
		}
	}
}

type ClusterNodeInfo struct {
	ClusterName string                `json:"clusterName"`
	NodeList    map[string]WorkerNode `json:"nodeList"`
}

type CsdEntry struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type SsdEntry struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type NodeStorageInfo struct {
	NodeName string     `json:"nodeName"`
	CsdList  []CsdEntry `json:"csdList"`
	SsdList  []SsdEntry `json:"ssdList"`
	NodeType string     `json:"nodeType"`
}

type NodeMetric struct {
	Time              string  `json:"timestamp"`
	NodeName          string  `json:"name"`
	CpuTick           int     `json:"cpuTick"`
	CpuTotal          float64 `json:"cpuTotal"`
	CpuUsed           float64 `json:"cpuUsed"`
	CpuUtilization    float64 `json:"cpuUtilization"`
	MemoryTotal       float64 `json:"memoryTotal"`
	MemoryUsed        float64 `json:"memoryUsed"`
	MemoryUtilization float64 `json:"memoryUtilization"`
	DiskTotal         float64 `json:"diskTotal"`
	DiskUsed          float64 `json:"diskUsed"`
	DiskUtilization   float64 `json:"diskUtilization"`
	NetworkRxData     float64 `json:"networkRxData"`
	NetworkTxData     float64 `json:"networkTxData"`
	NetworkBandwidth  float64 `json:"networkBandwidth"`
	PowerUsed         float64 `json:"powerUsed"`
}

type CpuMetric struct {
	Time           string  `json:"timestamp"`
	Name           string  `json:"name"`
	CpuTotal       float64 `json:"cpuTotal"`
	CpuUsed        float64 `json:"cpuUsed"`
	CpuUtilization float64 `json:"cpuUtilization"`
}

type CpuMetricMessage map[string][]CpuMetric

type PowerMetric struct {
	Time      string  `json:"timestamp"`
	Name      string  `json:"name"`
	PowerUsed float64 `json:"powerUsed"`
}

type PowerMetricMessage map[string][]PowerMetric

type MemoryMetric struct {
	Time              string  `json:"timestamp"`
	Name              string  `json:"name"`
	MemoryTotal       float64 `json:"memoryTotal"`
	MemoryUsed        float64 `json:"memoryUsed"`
	MemoryUtilization float64 `json:"memoryUtilization"`
}

type MemoryMetricMessage map[string][]MemoryMetric

type NetworkMetric struct {
	Time             string  `json:"timestamp"`
	Name             string  `json:"name"`
	NetworkRxData    float64 `json:"networkRxData"`
	NetworkTxData    float64 `json:"networkTxData"`
	NetworkBandwidth float64 `json:"networkBandwidth"`
}

type NetworkMetricMessage map[string][]NetworkMetric

type DiskMetric struct {
	Time            string  `json:"timestamp"`
	Name            string  `json:"name"`
	DiskTotal       float64 `json:"diskTotal"`
	DiskUsed        float64 `json:"diskUsed"`
	DiskUtilization float64 `json:"diskUtilization"`
}

type DiskMetricMessage map[string][]DiskMetric

type SsdMetric struct {
	Time            string  `json:"timestamp"`
	Id              string  `json:"id"`
	Name            string  `json:"name"`
	DiskTotal       float64 `json:"diskTotal"`
	DiskUsed        float64 `json:"diskUsed"`
	DiskUtilization float64 `json:"diskUtilization"`
	Status          string  `json:"status"`
}

type CsdMetric struct {
	Time                 string  `json:"timestamp"`
	Id                   string  `json:"id"`
	Name                 string  `json:"name"`
	Ip                   string  `json:"ip"`
	CpuTotal             float64 `json:"cpuTotal"`
	CpuUsed              float64 `json:"cpuUsed"`
	CpuUtilization       float64 `json:"cpuUtilization"`
	MemoryTotal          float64 `json:"memoryTotal"`
	MemoryUsed           float64 `json:"memoryUsed"`
	MemoryUtilization    float64 `json:"memoryUtilization"`
	DiskTotal            float64 `json:"diskTotal"`
	DiskUsed             float64 `json:"diskUsed"`
	DiskUtilization      float64 `json:"diskUtilization"`
	NetworkRxData        float64 `json:"networkRxData"`
	NetworkTxData        float64 `json:"networkTxData"`
	NetworkBandwidth     float64 `json:"networkBandwidth"`
	CsdMetricScore       float64 `json:"csdMetricScore"`
	CsdWorkingBlockCount float64 `json:"csdWorkingBlockCount"`
	Status               string  `json:"status"`
}

type StorageMetricMessage struct {
	SsdList map[string][]SsdMetric `json:"ssdList"`
	CsdList map[string][]CsdMetric `json:"csdList"`
}

func NewStorageMetricMessage() StorageMetricMessage {
	return StorageMetricMessage{
		SsdList: make(map[string][]SsdMetric),
		CsdList: make(map[string][]CsdMetric),
	}
}

type CsdMetricMin struct {
	Time            string  `json:"timestamp"`
	Id              string  `json:"id"`
	Name            string  `json:"name"`
	Ip              string  `json:"ip"`
	DiskTotal       float64 `json:"diskTotal"`
	DiskUsed        float64 `json:"diskUsed"`
	DiskUtilization float64 `json:"diskUtilization"`
	CsdMetricScore  float64 `json:"csdMetricScore"`
	Status          string  `json:"status"`
}

type StorageInfoMessage struct {
	SsdList map[string][]SsdMetric    `json:"ssdList"`
	CsdList map[string][]CsdMetricMin `json:"csdList"`
}

func NewStorageInfoMessage() StorageInfoMessage {
	return StorageInfoMessage{
		SsdList: make(map[string][]SsdMetric),
		CsdList: make(map[string][]CsdMetricMin),
	}
}

type Instance struct {
	InstanceName        string `json:"instanceName"`
	StorageEngineUid    string `json:"storageEngineUid"`
	QueryEngineUid      string `json:"queryEngineUid"`
	AccessPort          string `json:"accessPort"`
	InstanceType        string `json:"instanceType"`
	OperationNode       string `json:"operationNode"`
	StorageNode         string `json:"storageNode"`
	VolumeName          string `json:"volumeName"`
	QueryEngineStatus   string `json:"queryEngineStatus"`
	StorageEngineStatus string `json:"storageEngineStatus"`
	InstanceStatus      string `json:"instanceStatus"`
}

type VolumeInfo struct {
	VolumeName string `json:"volumeName"`
	VolumePath string `json:"volumePath"`
	NodeName   string `json:"nodeName"`
	// SizeTotal     float64 `json:"sizeTotal"`
	// SizeUsed      float64 `json:"sizeUsed"`
	// SizeAvailable float64 `json:"sizeAvailable"`
	// Utilization   float64 `json:"instanceUtilization"`
	StorageType  string `json:"storageType"`
	VolumeType   string `json:"volumeType"`
	StorageCount int    `json:"storageCount"`
	StorageName  string `json:"storageName"`
}

type InstanceMetric struct {
	Time           string  `json:"timestamp"`
	InstanceName   string  `json:"instanceName"`
	CpuUsage       float64 `json:"cpuUsage"`
	MemoryUsage    float64 `json:"memoryUsage"`
	StorageUsage   float64 `json:"storageUsage"`
	NetworkRxUsage float64 `json:"networkRxUsage"`
	NetworkTxUsage float64 `json:"networkTxUsage"`
}

type EnvironmentInfo struct {
	DbName     string `json:"dbName"`
	DbType     string `json:"dbType"`
	Algorithm  string `json:"algorithm"`
	BlockCount string `json:"blockCount"`
	DbSize     string `json:"dbSize"`
}
