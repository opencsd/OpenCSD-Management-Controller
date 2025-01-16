package handler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	types "opencsd-api-server/src/type"
	"os/exec"
)

func CreateOpenCSD(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD API Server] CreateOpenCSD")
	qUrl := "http://" + types.OPENCSD_CONTROLLER_DNS + "/create/query-engine"

	instanceRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	queryEngineRequest, err := http.NewRequest("POST", qUrl, bytes.NewBuffer(instanceRequest))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queryEngineRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(queryEngineRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sUrl := "http://" + types.OPENCSD_CONTROLLER_DNS + "/create/storage-engine"

	storageEngineRequest, err := http.NewRequest("POST", sUrl, bytes.NewBuffer(instanceRequest))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	storageEngineRequest.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(storageEngineRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vUrl := "http://" + types.OPENCSD_CONTROLLER_DNS + "/create/validator"

	validatorRequest, err := http.NewRequest("POST", vUrl, bytes.NewBuffer(instanceRequest))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validatorRequest.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(validatorRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("[OpenCSD API Server] Create OpenCSD Successfully\n"))
}

func CreateMySQL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD API Server] CreateMySQL")

	qUrl := "http://" + types.OPENCSD_CONTROLLER_DNS + "/create/query-engine"

	instanceRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	queryEngineRequest, err := http.NewRequest("POST", qUrl, bytes.NewBuffer(instanceRequest))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queryEngineRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(queryEngineRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sUrl := "http://" + types.OPENCSD_CONTROLLER_DNS + "/create/mysql"

	mysqlRequest, err := http.NewRequest("POST", sUrl, bytes.NewBuffer(instanceRequest))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mysqlRequest.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(mysqlRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("[OpenCSD API Server] Create MySQL Successfully\n"))
}

func CreateGraphDB(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD API Server] CreateGraphDB")

	qUrl := "http://" + types.OPENCSD_CONTROLLER_DNS + "/create/query-engine"

	instanceRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	queryEngineRequest, err := http.NewRequest("POST", qUrl, bytes.NewBuffer(instanceRequest))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queryEngineRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(queryEngineRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sUrl := "http://" + types.OPENCSD_CONTROLLER_DNS + "/create/graphdb"

	mysqlRequest, err := http.NewRequest("POST", sUrl, bytes.NewBuffer(instanceRequest))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mysqlRequest.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(mysqlRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("[OpenCSD API Server] Create GraphDB Successfully\n"))
}

func DeleteOpenCSD(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD API Server] DeleteOpenCSD")

	instance := r.URL.Query().Get("instance")

	qUrl := "http://" + types.OPENCSD_CONTROLLER_DNS + "/delete/query-engine"
	qUrl += "?instance=" + instance

	resp, err := http.Get(qUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	sUrl := "http://" + types.OPENCSD_CONTROLLER_DNS + "/delete/storage-engine"
	sUrl += "?instance=" + instance

	resp, err = http.Get(sUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	w.Write([]byte("[OpenCSD API Server] Delete OpenCSD Successfully\n"))
}

func DeleteMySQL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD API Server] DeleteMySQL")

	instance := r.URL.Query().Get("instance")

	qUrl := "http://" + types.OPENCSD_CONTROLLER_DNS + "/delete/query-engine"
	qUrl += "?instance=" + instance

	resp, err := http.Get(qUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	sUrl := "http://" + types.OPENCSD_CONTROLLER_DNS + "/delete/mysql"
	sUrl += "?instance=" + instance

	resp, err = http.Get(sUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	w.Write([]byte("[OpenCSD API Server] Delete OpenCSD Instance Successfully\n"))
}

func DeleteGraphDB(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD API Server] DeleteGraphDB")

	instance := r.URL.Query().Get("instance")

	qUrl := "http://" + types.OPENCSD_CONTROLLER_DNS + "/delete/query-engine"
	qUrl += "?instance=" + instance

	resp, err := http.Get(qUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(string(body))

	sUrl := "http://" + types.OPENCSD_CONTROLLER_DNS + "/delete/graphdb"
	sUrl += "?instance=" + instance

	resp, err = http.Get(sUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	w.Write([]byte("[OpenCSD API Server] Delete OpenCSD Instance Successfully\n"))
}

func DeleteInstance(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD API Server] DeleteInstance")

	instance := r.URL.Query().Get("instance")

	qUrl := "http://" + types.OPENCSD_CONTROLLER_DNS + "/delete/instance"
	qUrl += "?instance=" + instance

	resp, err := http.Get(qUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	w.Write([]byte("[OpenCSD API Server] Delete Instance Successfully\n"))
}

func VolumeAllocateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD API Server] VolumeAllocateHandler")

	var jsonResponse []byte

	type VolumeInfo struct {
		VolumeName   string  `json:"volumeName"`
		VolumePath   string  `json:"volumePath"`
		SizeTotal    float64 `json:"sizeTotal"`
		SizeUsed     float64 `json:"sizeUsed"`
		Utilization  float64 `json:"instanceUtilization"`
		StorageType  string  `json:"storageType"`
		VolumeType   string  `json:"volumeType"`
		StorageCount int     `json:"storageCount"`
		StorageName  string  `json:"storageName"`
	}

	response := VolumeInfo{}

	sUrl := "http://" + types.OPENCSD_CONTROLLER_DNS + "/volume/allocate"

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	volumeRequest, err := http.NewRequest("POST", sUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	volumeRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(volumeRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func VolumeDeAllocateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD API Server] VolumeDeAllocateHandler")

	volume := r.URL.Query().Get("volume")

	qUrl := "http://" + types.OPENCSD_CONTROLLER_DNS + "/volume/deallocate"
	qUrl += "?volume=" + volume

	resp, err := http.Get(qUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	w.Write([]byte("[OpenCSD API Server] DeAllocate Volume Successfully\n"))
}

func ClusterStorageNodeListHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD API Server] ClusterStorageNodeListHandler")

	var jsonResponse []byte
	response := types.ClusterNodeInfo{
		NodeList: make(map[string]types.WorkerNode),
	}

	response.ClusterName = types.ManagementMaster_.ClusterName

	for _, storageNode := range types.ManagementMaster_.StorageLayer {
		response.NodeList[storageNode.NodeName] = storageNode
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func ClusterNodeListHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD API Server] ClusterNodeListHandler")

	var jsonResponse []byte
	response := types.ClusterNodeInfo{
		NodeList: make(map[string]types.WorkerNode),
	}

	response.ClusterName = types.ManagementMaster_.ClusterName

	layer := r.URL.Query().Get("layer")
	if layer == "operation" {
		response.NodeList = types.ManagementMaster_.OperationLayer
	} else if layer == "storage" {
		response.NodeList = types.ManagementMaster_.StorageLayer
	} else {
		for key, value := range types.ManagementMaster_.OperationLayer {
			response.NodeList[key] = value
		}
		for key, value := range types.ManagementMaster_.StorageLayer {
			response.NodeList[key] = value
		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func VolumeInfoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD API Server] VolumeInfoHandler")

	var jsonResponse []byte

	type VolumeInfo struct {
		VolumeName string `json:"volumeName"`
		VolumePath string `json:"volumePath"`
		// SizeTotal   float64 `json:"sizeTotal"`
		// SizeUsed    float64 `json:"sizeUsed"`
		// Utilization float64 `json:"instanceUtilization"`
		StorageType  string `json:"storageType"`
		VolumeType   string `json:"volumeType"`
		StorageCount int    `json:"storageCount"`
		StorageName  string `json:"storageName"`
	}

	type VolumeMessage struct {
		Name       string                `json:"nodeName"`
		IP         string                `json:"nodeIp"`
		Status     string                `json:"status"`
		VolumeInfo map[string]VolumeInfo `json:"volumeInfo"`
	}

	response := map[string]VolumeMessage{}

	targetNode := r.URL.Query().Get("node")

	url := "http://" + types.OPENCSD_CONTROLLER_DNS + "/info/volume"
	url += "?node=" + targetNode

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func InstanceInfoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD API Server] InstanceInfoHandler")

	type instance_ struct {
		InstanceInfo types.Instance   `json:"instanceInfo"`
		VolumeInfo   types.VolumeInfo `json:"volumeInfo"`
	}

	var jsonResponse []byte
	response := map[string]instance_{}

	targetInstance := r.URL.Query().Get("instance")
	targetNode := r.URL.Query().Get("node")

	url := "http://" + types.OPENCSD_CONTROLLER_DNS + "/info/instance"
	url += "?instance=" + targetInstance
	url += "&node=" + targetNode

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func NodeStorageListHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD API Server] NodeStorageListHandler")

	var jsonResponse []byte
	response := map[string]types.NodeStorageInfo{}

	targetNode := r.URL.Query().Get("node")
	count := r.URL.Query().Get("count")

	for _, storageNode := range types.ManagementMaster_.StorageLayer {
		if targetNode == "" || targetNode == storageNode.NodeName {
			url := "http://" + storageNode.NodeIP + ":" + types.STORAGE_API_SERVER_PORT + "/node/info/storage-list"
			url += "?count=" + count

			resp, err := http.Get(url)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			nodeStorageInfo := types.NodeStorageInfo{}

			err = json.Unmarshal(body, &nodeStorageInfo)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			response[storageNode.NodeName] = nodeStorageInfo

		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func NodeStorageInfoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD API Server] NodeStorageInfoHandler")

	var jsonResponse []byte
	response := map[string]types.StorageInfoMessage{}

	targetNode := r.URL.Query().Get("node")
	targetStorage := r.URL.Query().Get("storage")
	count := r.URL.Query().Get("count")

	for _, storageNode := range types.ManagementMaster_.StorageLayer {
		if targetNode == "" || targetNode == storageNode.NodeName {

			url := "http://" + storageNode.NodeIP + ":" + types.STORAGE_API_SERVER_PORT + "/node/info/storage"
			url += "?storage=" + targetStorage
			url += "&count=" + count

			resp, err := http.Get(url)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			storageInfo := types.NewStorageInfoMessage()

			err = json.Unmarshal(body, &storageInfo)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			response[storageNode.NodeName] = storageInfo

		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func NodeDiskInfoHandler(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := map[string][]types.DiskMetric{}

	targetNode := r.URL.Query().Get("node")
	count := r.URL.Query().Get("count")

	for _, storageNode := range types.ManagementMaster_.StorageLayer {
		if targetNode == "" || targetNode == storageNode.NodeName {
			url := "http://" + storageNode.NodeIP + ":" + types.STORAGE_API_SERVER_PORT + "/node/metric/disk"
			url += "?count=" + count

			resp, err := http.Get(url)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			nodeDiskInfo := []types.DiskMetric{}

			err = json.Unmarshal(body, &nodeDiskInfo)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			response[storageNode.NodeName] = nodeDiskInfo

		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func NodeMetricInfoHandler(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := map[string][]types.NodeMetric{}

	targetNode := r.URL.Query().Get("node")
	count := r.URL.Query().Get("count")

	for _, storageNode := range types.ManagementMaster_.StorageLayer {
		if targetNode == "" || targetNode == storageNode.NodeName {
			url := "http://" + storageNode.NodeIP + ":" + types.STORAGE_API_SERVER_PORT + "/node/metric/all"
			url += "?count=" + count

			resp, err := http.Get(url)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			nodeMetricInfo := []types.NodeMetric{}

			err = json.Unmarshal(body, &nodeMetricInfo)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			response[storageNode.NodeName] = nodeMetricInfo

		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func StorageMetricInfoHandler(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := map[string]types.StorageMetricMessage{}

	targetNode := r.URL.Query().Get("node")
	targetStorage := r.URL.Query().Get("storage")
	count := r.URL.Query().Get("count")

	for _, storageNode := range types.ManagementMaster_.StorageLayer {
		if targetNode == "" || targetNode == storageNode.NodeName {
			url := "http://" + storageNode.NodeIP + ":" + types.STORAGE_API_SERVER_PORT + "/storage/metric/all"
			url += "?storage=" + targetStorage
			url += "&count=" + count

			resp, err := http.Get(url)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			storageMetric := types.NewStorageMetricMessage()

			err = json.Unmarshal(body, &storageMetric)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			response[storageNode.NodeName] = storageMetric

		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func StorageMetricCpuHandler(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := map[string]types.CpuMetricMessage{}

	targetNode := r.URL.Query().Get("node")
	targetStorage := r.URL.Query().Get("storage")
	count := r.URL.Query().Get("count")

	for _, storageNode := range types.ManagementMaster_.StorageLayer {
		if targetNode == "" || targetNode == storageNode.NodeName {
			url := "http://" + storageNode.NodeIP + ":" + types.STORAGE_API_SERVER_PORT + "/storage/metric/cpu"
			url += "?storage=" + targetStorage
			url += "&count=" + count

			resp, err := http.Get(url)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			cpuMetric := types.CpuMetricMessage{}

			err = json.Unmarshal(body, &cpuMetric)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			response[storageNode.NodeName] = cpuMetric

		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func StorageMetricMemoryHandler(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := map[string]types.MemoryMetricMessage{}

	targetNode := r.URL.Query().Get("node")
	targetStorage := r.URL.Query().Get("storage")
	count := r.URL.Query().Get("count")

	for _, storageNode := range types.ManagementMaster_.StorageLayer {
		if targetNode == "" || targetNode == storageNode.NodeName {
			url := "http://" + storageNode.NodeIP + ":" + types.STORAGE_API_SERVER_PORT + "/storage/metric/memory"
			url += "?storage=" + targetStorage
			url += "&count=" + count

			resp, err := http.Get(url)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			memoryMetric := types.MemoryMetricMessage{}

			err = json.Unmarshal(body, &memoryMetric)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			response[storageNode.NodeName] = memoryMetric

		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func StorageMetricNetworkHandler(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := map[string]types.NetworkMetricMessage{}

	targetNode := r.URL.Query().Get("node")
	targetStorage := r.URL.Query().Get("storage")
	count := r.URL.Query().Get("count")

	for _, storageNode := range types.ManagementMaster_.StorageLayer {
		if targetNode == "" || targetNode == storageNode.NodeName {
			url := "http://" + storageNode.NodeIP + ":" + types.STORAGE_API_SERVER_PORT + "/storage/metric/network"
			url += "?storage=" + targetStorage
			url += "&count=" + count

			resp, err := http.Get(url)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			networkMetric := types.NetworkMetricMessage{}

			err = json.Unmarshal(body, &networkMetric)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			response[storageNode.NodeName] = networkMetric

		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func StorageMetricDiskHandler(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := map[string]types.DiskMetricMessage{}

	targetNode := r.URL.Query().Get("node")
	targetStorage := r.URL.Query().Get("storage")
	count := r.URL.Query().Get("count")

	for _, storageNode := range types.ManagementMaster_.StorageLayer {
		if targetNode == "" || targetNode == storageNode.NodeName {
			url := "http://" + storageNode.NodeIP + ":" + types.STORAGE_API_SERVER_PORT + "/storage/metric/disk"
			url += "?storage=" + targetStorage
			url += "&count=" + count

			resp, err := http.Get(url)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			diskMetric := types.DiskMetricMessage{}

			err = json.Unmarshal(body, &diskMetric)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			response[storageNode.NodeName] = diskMetric

		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func StorageMetricPowerHandler(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := map[string]types.PowerMetricMessage{}

	targetNode := r.URL.Query().Get("node")
	targetStorage := r.URL.Query().Get("storage")
	count := r.URL.Query().Get("count")

	for _, storageNode := range types.ManagementMaster_.StorageLayer {
		if targetNode == "" || targetNode == storageNode.NodeName {
			url := "http://" + storageNode.NodeIP + ":" + types.STORAGE_API_SERVER_PORT + "/storage/metric/power"
			url += "?storage=" + targetStorage
			url += "&count=" + count

			resp, err := http.Get(url)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			powerMetric := types.PowerMetricMessage{}

			err = json.Unmarshal(body, &powerMetric)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			response[storageNode.NodeName] = powerMetric

		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func CmdExec(cmdStr string) error {
	cmd := exec.Command("bash", "-c", cmdStr)
	stdoutReader, _ := cmd.StdoutPipe()
	stdoutScanner := bufio.NewScanner(stdoutReader)
	go func() {
		for stdoutScanner.Scan() {
			fmt.Println(stdoutScanner.Text())
		}
	}()
	stderrReader, _ := cmd.StderrPipe()
	stderrScanner := bufio.NewScanner(stderrReader)
	go func() {
		for stderrScanner.Scan() {
			fmt.Println(stderrScanner.Text())
		}
	}()
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error : %v \n", err)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Error: %v \n", err)
	}

	return nil
}
