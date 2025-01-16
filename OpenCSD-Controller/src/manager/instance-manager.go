package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

var InstanceManager_ *InstanceManager

var (
	OPENCSD_CONTROLLER_PORT                = "40801"
	OPENCSD_INSTANCE_METRIC_COLLECTOR_PORT = "40804"
	OPENCSD_INSTANCE_METRIC_INFO_DB_PORT   = "40806"
	STORAGE_API_SERVER_PORT                = "40306"
)

var (
	INSTANCE_METRIC_MYSQL_ROOT_PASSWORD = os.Getenv("MYSQL_ROOT_PASSWORD")
	INSTANCE_METRIC_MYSQL_USER          = os.Getenv("MYSQL_USER")
	INSTANCE_METRIC_MYSQL_PASSWORD      = os.Getenv("MYSQL_PASSWORD")
)

const (
	READY    = "READY"
	NOTREADY = "NOTREADY"
	BROKEN   = "BROKEN"
	NORMAL   = "NORMAL"
) //storage status

const (
	STORAGE    = "STORAGE"
	OPERATION  = "OPERATION"
	MANAGEMENT = "MANAGEMENT"
) // layer

const (
	OPENCSD = "OPENCSD"
	MYSQL   = "MYSQL"
	GRAPHDB = "GRAPHDB"
) //instance type

const (
	SSD     = "SSD"
	CSD     = "CSD"
	CROSS   = "CROSS"
	UNKNOWN = "UNKNOWN"
) // storage type

const (
	GLUSTER  = "GLUSTER"
	LVM      = "LVM"
	HOSTPATH = "HOSTPATH"
	DEFAULT  = "DEFAULT"
) // volume type

const (
	RUNNING     = "RUNNING"
	NONE        = "NONE"
	ERROR       = "ERROR"
	AVAILABLE   = "AVAILABLE"
	UNAVAILABLE = "UNAVAILABLE"
) // engine status

type InstanceManager struct {
	ClusterName       string
	MasterNode        *Node
	ClusterConfig     ClusterConfig
	InformerFactory   informers.SharedInformerFactory
	StorageLayer      map[string]*Node             // key : node name
	OperationLayer    map[string]*Node             // key : node name
	StorageNodeInfo   map[string]*NodeStorageInfo  // key : node name
	OperationNodeInfo map[string]*NodeInstanceInfo // key : node name
	InstanceInfo      map[string]*Instance         // key : instance name
	VolumeInfo        map[string]*VolumeInfo       // key : volume name
	mu                sync.Mutex
}

func NewInstaceManager() *InstanceManager {
	return &InstanceManager{
		ClusterName:       "",
		MasterNode:        &Node{},
		StorageLayer:      make(map[string]*Node),
		OperationLayer:    make(map[string]*Node),
		StorageNodeInfo:   make(map[string]*NodeStorageInfo),
		OperationNodeInfo: make(map[string]*NodeInstanceInfo),
		InstanceInfo:      make(map[string]*Instance),
		VolumeInfo:        make(map[string]*VolumeInfo),
	}
}

type ClusterConfig struct {
	Clientset *kubernetes.Clientset
	config    *rest.Config
}

type Node struct {
	Name   string   `json:"nodeName"`
	IP     string   `json:"nodeIp"`
	Status string   `json:"status"`
	Layer  string   `json:"layer"`
	Label  []string `json:"label"`
}

type NodeStorageInfo struct {
	NodeName        string                 `json:"nodeName"`
	CsdList         []CsdEntry             `json:"csdList"`
	SsdList         []SsdEntry             `json:"ssdList"`
	NodeStorageType string                 `json:"nodeType"`
	VolumeInfo      map[string]*VolumeInfo `json:"volumeInfo"`
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

type NodeInstanceInfo struct {
	NodeName     string               `json:"nodeName"`
	InstanceInfo map[string]*Instance `json:"instanceInfo"`
}

type Instance struct {
	InstanceName        string `json:"instanceName"`
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

func (instanceManager *InstanceManager) InitInstanceManager() {
	masterNode := Node{}
	masterName := os.Getenv("NODE_NAME")
	if masterName == "" {
		fmt.Println("NODE_NAME environment variable is not set")
	}
	masterNode.Name = masterName
	instanceManager.ClusterName = "OPENCSD"

	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("InClusterConfig error:", err)
		return
	}

	instanceManager.ClusterConfig.config = config
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("NewForConfig error:", err)
		return
	}

	instanceManager.ClusterConfig.Clientset = clientset

	instanceManager.StorageLayer = make(map[string]*Node)
	instanceManager.OperationLayer = make(map[string]*Node)

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Get Node error:", err)
	} else {
		for _, node := range nodes.Items {
			newNode := &Node{}

			newNode.Name = node.Name

			for _, address := range node.Status.Addresses {
				if address.Type == "InternalIP" {
					newNode.IP = address.Address
					break
				}
			}

			for _, condition := range node.Status.Conditions {
				if condition.Type == corev1.NodeReady {
					newNode.Status = READY
				} else {
					newNode.Status = NOTREADY
				}
			}

			if node.Name == masterName {
				newNode.Layer = MANAGEMENT
				instanceManager.MasterNode = newNode
			} else if node.Labels["layer"] == "storage" {
				newNode.Layer = STORAGE
				instanceManager.StorageLayer[node.Name] = newNode
			} else if node.Labels["layer"] == "operation" {
				newNode.Layer = OPERATION
				instanceManager.OperationLayer[node.Name] = newNode
				instanceManager.OperationNodeInfo[node.Name] = &NodeInstanceInfo{
					NodeName:     node.Name,
					InstanceInfo: make(map[string]*Instance),
				}
			}
		}
	}

	instanceManager.getStorageNodeInfo()
	instanceManager.initInstanceInfo()
}

func (instanceManager *InstanceManager) getStorageNodeInfo() {
	for _, storageNode := range instanceManager.StorageLayer {
		url := "http://" + storageNode.IP + ":" + STORAGE_API_SERVER_PORT + "/node/info/storage-list"

		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		nodeStorageInfo := &NodeStorageInfo{}

		err = json.Unmarshal(body, &nodeStorageInfo)
		if err != nil {
			log.Fatal(err)
		}

		instanceManager.StorageNodeInfo[storageNode.Name] = nodeStorageInfo

		// 볼륨 정보 하드코딩
		if storageNode.Name == "storage-node1" {
			volumeName := "opencsd-csd-glusterfs"
			volumeInfo := &VolumeInfo{
				VolumeName:   volumeName,
				VolumePath:   "/mnt/gluster/client",
				NodeName:     "storage-node1",
				StorageType:  CSD,
				VolumeType:   GLUSTER,
				StorageCount: 8,
				StorageName:  "NGD-IN2510-080T4-C",
			}
			instanceManager.StorageNodeInfo[storageNode.Name].VolumeInfo[volumeName] = volumeInfo
			instanceManager.VolumeInfo[volumeName] = volumeInfo
		} else if storageNode.Name == "storage-node4" {
			volumeName := "opencsd-ssd-lvm"
			volumeInfo := &VolumeInfo{
				VolumeName:   volumeName,
				VolumePath:   "/mnt/lvm",
				NodeName:     "storage-node4",
				StorageType:  SSD,
				VolumeType:   LVM,
				StorageCount: 8,
				StorageName:  "Samsung SSD 870 QVO",
			}
			instanceManager.StorageNodeInfo[storageNode.Name].VolumeInfo[volumeName] = volumeInfo
			instanceManager.VolumeInfo[volumeName] = volumeInfo
		}
	}
}

func (instanceManager *InstanceManager) initInstanceInfo() {
	instanceManager.InformerFactory = informers.NewSharedInformerFactory(instanceManager.ClusterConfig.Clientset, 0)
	instanceManager.addEventHandler()

	excludedNamespaces := map[string]bool{
		"kube-system":           true,
		"management-controller": true,
		"storage-controller":    true,
	}

	pods, err := instanceManager.ClusterConfig.Clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error fetching pods: %v", err)
	}

	for _, pod := range pods.Items {
		if excludedNamespaces[pod.Namespace] {
			continue
		}

		if value, exists := pod.Labels["tier"]; !exists || value != "opencsd" {
			continue
		}

		instanceManager.realAddInstance(&pod)
	}
}

func (instanceManager *InstanceManager) getNodePort(namespace string) string {
	serviceName := "query-engine-instance-svc"

	service, err := instanceManager.ClusterConfig.Clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Error fetching service: %v", err)
	}

	for _, port := range service.Spec.Ports {
		if port.NodePort != 0 {
			nodeport := strconv.FormatInt(int64(port.NodePort), 10)
			return nodeport
		}
	}

	return ""
}

func (instanceManager *InstanceManager) getStorageEngineInfo(namespace string, pod *v1.Pod) (string, string) {
	instanceType, volumeName := "", ""

	for _, volume := range pod.Spec.Volumes {
		if strings.HasPrefix(volume.Name, "opencsd") {
			volumeName = volume.Name
			break
		}
	}

	if strings.HasPrefix(pod.Name, "storage-engine-instance") {
		instanceType = OPENCSD
	} else if strings.HasPrefix(pod.Name, "mysql") {
		instanceType = MYSQL
	} else if strings.HasPrefix(pod.Name, "graphdb") {
		instanceType = GRAPHDB
	}

	return instanceType, volumeName
}

func (instanceManager *InstanceManager) addEventHandler() {
	instanceManager.InformerFactory.Core().V1().Pods().Informer().AddEventHandler(
		cache.FilteringResourceEventHandler{
			FilterFunc: func(obj interface{}) bool {
				switch t := obj.(type) {
				case *v1.Pod:
					return instanceManager.isInstance(t)
				case cache.DeletedFinalStateUnknown:
					if _, ok := t.Obj.(*v1.Pod); ok {
						return true
					}
					return false
				default:
					return false
				}
			},
			Handler: cache.ResourceEventHandlerFuncs{
				AddFunc:    instanceManager.addInstance,
				UpdateFunc: instanceManager.updateInstance,
				DeleteFunc: instanceManager.deleteInstance,
			},
		},
	)
}

func (instanceManager *InstanceManager) isInstance(obj interface{}) bool {
	pod, _ := obj.(*v1.Pod)

	if value, exists := pod.Labels["tier"]; exists && value == "opencsd" && pod.Spec.NodeName != "" {
		return true
	}

	return false
}

func (instanceManager *InstanceManager) addInstance(obj interface{}) {
	pod, _ := obj.(*v1.Pod)
	instanceManager.realAddInstance(pod)
}

func (instanceManager *InstanceManager) realAddInstance(pod *v1.Pod) {
	instanceManager.mu.Lock()
	defer instanceManager.mu.Unlock()

	instanceName := pod.Namespace
	nodeName := pod.Spec.NodeName

	if instance, exists := instanceManager.InstanceInfo[instanceName]; exists {
		if strings.HasPrefix(pod.Name, "query-engine-instance") {
			instance.AccessPort = instanceManager.getNodePort(instanceName)
			instance.QueryEngineStatus = parsingStatus(string(pod.Status.Phase))
		} else { // "storage-engine-instance" or "mysql" or "graphdb"
			instanceType, volumeName := instanceManager.getStorageEngineInfo(instanceName, pod)

			instance.InstanceType = instanceType
			instance.StorageNode = instanceManager.VolumeInfo[volumeName].NodeName
			instance.VolumeName = volumeName
			instance.StorageEngineStatus = parsingStatus(string(pod.Status.Phase))
		}
		instanceManager.updateInstanceStatus(nodeName, instanceName)
	} else {
		if strings.HasPrefix(pod.Name, "query-engine-instance") {
			newInstance := &Instance{
				InstanceName:        instanceName,
				AccessPort:          instanceManager.getNodePort(instanceName),
				OperationNode:       pod.Spec.NodeName,
				QueryEngineStatus:   parsingStatus(string(pod.Status.Phase)),
				StorageEngineStatus: NONE,
				InstanceStatus:      UNAVAILABLE,
			}
			instanceManager.InstanceInfo[instanceName] = newInstance
			instanceManager.OperationNodeInfo[nodeName].InstanceInfo[instanceName] = newInstance
		} else { // "storage-engine-instance" or "mysql" or "graphdb"
			instanceType, volumeName := instanceManager.getStorageEngineInfo(instanceName, pod)

			newInstance := &Instance{
				InstanceName:        instanceName,
				InstanceType:        instanceType,
				OperationNode:       pod.Spec.NodeName,
				StorageNode:         instanceManager.VolumeInfo[volumeName].NodeName,
				VolumeName:          volumeName,
				QueryEngineStatus:   NONE,
				StorageEngineStatus: parsingStatus(string(pod.Status.Phase)),
				InstanceStatus:      UNAVAILABLE,
			}
			instanceManager.InstanceInfo[instanceName] = newInstance
			instanceManager.OperationNodeInfo[nodeName].InstanceInfo[instanceName] = newInstance
		}
	}
}

func (instanceManager *InstanceManager) updateInstance(oldObj interface{}, newObj interface{}) {
	instanceManager.mu.Lock()
	defer instanceManager.mu.Unlock()

	oldPod, _ := oldObj.(*v1.Pod)
	newPod, _ := newObj.(*v1.Pod)

	instanceName := newPod.Namespace
	nodeName := newPod.Spec.NodeName

	if oldPod.Spec.NodeName != newPod.Spec.NodeName {
		instanceManager.realAddInstance(newPod)
		return
	}

	if oldPod.Status.Phase != newPod.Status.Phase {
		status := parsingStatus(string(newPod.Status.Phase))

		if strings.HasPrefix(newPod.Name, "query-engine-instance") {
			instanceManager.InstanceInfo[instanceName].QueryEngineStatus = status
			instanceManager.OperationNodeInfo[nodeName].InstanceInfo[instanceName].QueryEngineStatus = status
		} else {
			instanceManager.InstanceInfo[instanceName].StorageEngineStatus = status
			instanceManager.OperationNodeInfo[nodeName].InstanceInfo[instanceName].StorageEngineStatus = status
		}

		instanceManager.updateInstanceStatus(nodeName, instanceName)
	}
}

func (instanceManager *InstanceManager) deleteInstance(obj interface{}) {
	instanceManager.mu.Lock()
	defer instanceManager.mu.Unlock()

	pod, _ := obj.(*v1.Pod)
	fmt.Println("deleteInstance ", pod.Name)

	instanceName := pod.Namespace
	nodeName := pod.Spec.NodeName

	if strings.HasPrefix(pod.Name, "query-engine-instance") {
		if instanceManager.OperationNodeInfo[nodeName].InstanceInfo[instanceName].StorageEngineStatus == NONE {
			delete(instanceManager.OperationNodeInfo[nodeName].InstanceInfo, instanceName)
			delete(instanceManager.InstanceInfo, instanceName)
		} else {
			instanceManager.OperationNodeInfo[nodeName].InstanceInfo[instanceName].QueryEngineStatus = NONE
			instanceManager.OperationNodeInfo[nodeName].InstanceInfo[instanceName].AccessPort = ""
			instanceManager.updateInstanceStatus(nodeName, instanceName)
		}
	} else {
		if instanceManager.OperationNodeInfo[nodeName].InstanceInfo[instanceName].QueryEngineStatus == NONE {
			delete(instanceManager.OperationNodeInfo[nodeName].InstanceInfo, instanceName)
			delete(instanceManager.InstanceInfo, instanceName)
		} else {
			instanceManager.OperationNodeInfo[nodeName].InstanceInfo[instanceName].StorageEngineStatus = NONE
			instanceManager.OperationNodeInfo[nodeName].InstanceInfo[instanceName].StorageNode = ""
			instanceManager.OperationNodeInfo[nodeName].InstanceInfo[instanceName].VolumeName = ""
			instanceManager.updateInstanceStatus(nodeName, instanceName)
		}
	}
}

func (instanceManager *InstanceManager) updateInstanceStatus(nodeName string, instanceName string) {
	if instanceManager.InstanceInfo[instanceName].QueryEngineStatus == RUNNING && instanceManager.InstanceInfo[instanceName].StorageEngineStatus == RUNNING {
		instanceManager.InstanceInfo[instanceName].InstanceStatus = AVAILABLE
		instanceManager.OperationNodeInfo[nodeName].InstanceInfo[instanceName].InstanceStatus = AVAILABLE
	} else {
		instanceManager.InstanceInfo[instanceName].InstanceStatus = UNAVAILABLE
		instanceManager.OperationNodeInfo[nodeName].InstanceInfo[instanceName].InstanceStatus = UNAVAILABLE
	}
}

func parsingStatus(status string) string {
	if status == string(corev1.PodRunning) {
		return RUNNING
	} else {
		return ERROR
	}
}
