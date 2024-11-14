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
	OPENCSD_API_SERVER_PORT       = os.Getenv("OPENCSD_API_SERVER_PORT")
	OPENCSD_ENGINE_DEPLOYER_PORT  = os.Getenv("OPENCSD_ENGINE_DEPLOYER_PORT")
	OPENCSD_INSTANCE_MANAGER_PORT = os.Getenv("OPENCSD_INSTANCE_MANAGER_PORT")
	OPENCSD_VOLUME_ALLOCATOR_PORT = os.Getenv("OPENCSD_VOLUME_ALLOCATOR_PORT")
	STORAGE_API_SERVER_PORT       = os.Getenv("STORAGE_API_SERVER_PORT")
)

const (
	READY    = "READY"
	NOTREADY = "NOTREADY"
)

const (
	STORAGE   = "STORAGE"
	OPERATION = "OPERATOR"
)

type MasterNode struct {
	ClusterName   string
	NodeName      string
	MasterIP      string
	clusterConfig ClusterConfig
	StorageLayer  []WorkerNode
	OperatorLayer []WorkerNode
}

type ClusterConfig struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
}

type WorkerNode struct {
	NodeName string `json:"node_name"`
	NodeIP   string `json:"node_ip"`
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

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Get Node error:", err)
	} else {
		for _, node := range nodes.Items {
			var workerNode WorkerNode
			for _, condition := range node.Status.Conditions {
				if condition.Type == corev1.NodeReady {
					workerNode.Status = READY
					if node.Name == masterName {
						for _, address := range node.Status.Addresses {
							if address.Type == "InternalIP" {
								masterNode.MasterIP = address.Address
								break
							}
						}
					} else {
						for _, address := range node.Status.Addresses {
							if address.Type == "InternalIP" {
								workerNode.NodeIP = address.Address
								workerNode.NodeName = node.Name

								if node.Labels["layer"] == "storage" {
									workerNode.Layer = STORAGE
									masterNode.StorageLayer = append(masterNode.StorageLayer, workerNode)
								} else if node.Labels["layer"] == "operator" {
									workerNode.Layer = OPERATION
									masterNode.OperatorLayer = append(masterNode.OperatorLayer, workerNode)
								}

								break
							}
						}
					}
				} else {
					workerNode.NodeName = node.Name
					workerNode.Status = NOTREADY
				}
			}
		}
	}
}

type ClusterStorageNodeInfo struct {
	ClusterName     string                `json:"cluster_name"`
	StorageNodeList map[string]WorkerNode `json:"storage_node_list"`
}

type CsdEntry struct {
	CsdId  string `json:"csd_id"`
	Status string `json:"status"`
}

type NodeStorageInfo struct {
	NodeName string     `json:"node_name"`
	CsdList  []CsdEntry `json:"csd_list"`
	SsdList  []string   `json:"ssd_list"`
	NodeType string     `json:"node_type"`
}

type NodeMetric struct {
	Time              string  `json:"timestamp"`
	NodeName          string  `json:"name"`
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
	SsdList map[string][]SsdMetric `json:"ssd_list"`
	CsdList map[string][]CsdMetric `json:"csd_list"`
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
	SsdList map[string][]SsdMetric    `json:"ssd_list"`
	CsdList map[string][]CsdMetricMin `json:"csd_list"`
}

func NewStorageInfoMessage() StorageInfoMessage {
	return StorageInfoMessage{
		SsdList: make(map[string][]SsdMetric),
		CsdList: make(map[string][]CsdMetricMin),
	}
}
