package collector

import (
	"fmt"
	"gpu-metric-collector/pkg/api"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var GPU_METRIC_COLLECTOR_DEBUGG_LEVEL = os.Getenv("DEBUGG_LEVEL")
var MaxLanes = 16 //nvlink를 확인하기 위해 도는 레인 수

const (
	LEVEL1 = "LEVEL1"
	LEVEL2 = "LEVEL2"
	LEVEL3 = "LEVEL3"
)

type MetricCollector struct {
	HostKubeClient  *kubernetes.Clientset
	Interval        *time.Duration
	StatRequest     *http.Request
	NodeName		string
	NodeIP			string
}

type Network struct {
	NetworkRxBytes resource.Quantity
	NetworkTxBytes resource.Quantity
}

type Summary struct {
	Node NodeStats  `json:"node"`
	Pods []PodStats `json:"pods"`
}

type NodeStats struct {
	CPU     *CPUStats     `json:"cpu,omitempty"`
	Memory  *MemoryStats  `json:"memory,omitempty"`
	Network *NetworkStats `json:"network,omitempty"`
	Fs      *FsStats      `json:"fs,omitempty"`
}

type InterfaceStats struct {
	Name     string  `json:"name"`
	RxBytes  *uint64 `json:"rxBytes,omitempty"`
	RxErrors *uint64 `json:"rxErrors,omitempty"`
	TxBytes  *uint64 `json:"txBytes,omitempty"`
	TxErrors *uint64 `json:"txErrors,omitempty"`
}

type PodStats struct {
	PodRef           PodReference  		`json:"podRef"`
	Containers 		[]ContainerStats 	`json:"containers"`
	CPU              *CPUStats    		`json:"cpu,omitempty"`
	Memory           *MemoryStats  		`json:"memory,omitempty"`
	Network          *NetworkStats 		`json:"network,omitempty"`
	EphemeralStorage *FsStats     		`json:"ephemeral-storage,omitempty"`
}

type ContainerStats struct {
	Name 		string 			`json:"name"`
	CPU 		*CPUStats		`json:"cpu,omitempty"`
	Memory 		*MemoryStats 	`json:"memory,omitempty"`
}

type PodReference struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	UID       string `json:"uid"`
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

type NetworkStats struct {
	Interfaces []InterfaceStats `json:"interfaces,omitempty"`
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

func NewMetricCollector() *MetricCollector {
	hostKubeClient := api.NewClientset()
	t, _ := time.ParseDuration(os.Getenv("COLLECT_INTERVAL"))
	interval := t * time.Second

	nodeIP := os.Getenv("NODE_IP")

	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println(fmt.Sprintf("[error] NewNetricCollector() > rest.InClusterconfig error: %v", err))
	}
	token := config.BearerToken

	scheme := "https"
	url := url.URL{
		Scheme: scheme,
		Host:   net.JoinHostPort(nodeIP, strconv.Itoa(10250)),
		Path:   "/stats/summary",
	}

	statRequest, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		fmt.Println(fmt.Sprintf("[error] NewNetricCollector() > http.NewRequest() error: %v", err))
	}

	statRequest.Header.Set("Content-Type", "application/json")
	statRequest.Header.Set("Authorization", "Bearer "+token)

	return &MetricCollector{
		HostKubeClient:  hostKubeClient,
		Interval:        &interval,
		StatRequest:     statRequest,
	}
}

// func New
