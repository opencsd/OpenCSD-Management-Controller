package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	types "opencsd-api-server/src/type"
	"os/exec"
)

func CreateOpenCSD(w http.ResponseWriter, r *http.Request) {
	dbms_name := r.URL.Query().Get("dbname")

	CmdExec("/usr/local/bin/kubectl create ns " + dbms_name)

	CmdExec("cp /mnt/instance-deploy/opencsd/*.yaml /mnt/instance-deploy/opencsd/tmp")
	CmdExec("cp /mnt/instance-deploy/opencsd/*.sh /mnt/instance-deploy/opencsd/tmp")
	CmdExec("sed -i 's/OPENCSD_NAMESPACE/" + dbms_name + "/g' /mnt/instance-deploy/opencsd/tmp/*.yaml")
	CmdExec("sed -i 's/OPENCSD_NAMESPACE/" + dbms_name + "/g' /mnt/instance-deploy/opencsd/tmp/*.sh")

	CmdExec("bash /mnt/instance-deploy/opencsd/tmp/docker_secret.sh")
	CmdExec("/usr/local/bin/kubectl create -f /mnt/instance-deploy/opencsd/tmp/.")

	CmdExec("rm /mnt/instance-deploy/opencsd/tmp/*")

	w.Write([]byte("[OpenCSD] Completed\n"))
}

func CreateMySQL(w http.ResponseWriter, r *http.Request) {
	dbms_name := r.URL.Query().Get("dbname")

	CmdExec("/usr/local/bin/kubectl create ns " + dbms_name)

	CmdExec("cp /mnt/instance-deploy/mysql/*.yaml /mnt/instance-deploy/mysql/tmp")
	CmdExec("cp /mnt/instance-deploy/mysql/*.sh /mnt/instance-deploy/mysql/tmp")
	CmdExec("sed -i 's/MYSQL_NAMESPACE/" + dbms_name + "/g' /mnt/instance-deploy/mysql/tmp/*.yaml")
	CmdExec("sed -i 's/MYSQL_NAMESPACE/" + dbms_name + "/g' /mnt/instance-deploy/mysql/tmp/*.sh")

	CmdExec("bash /mnt/instance-deploy/mysql/tmp/docker_secret.sh")
	CmdExec("/usr/local/bin/kubectl create -f /mnt/instance-deploy/mysql/tmp/.")

	CmdExec("rm /mnt/instance-deploy/mysql/tmp/*")

	w.Write([]byte("[MYSQL] Completed\n"))
}

func ClusterNodeListHandler(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := types.ClusterStorageNodeInfo{
		StorageNodeList: make(map[string]types.WorkerNode),
	}

	response.ClusterName = types.ManagementMaster_.ClusterName

	for _, storageNode := range types.ManagementMaster_.StorageLayer {
		response.StorageNodeList[storageNode.NodeName] = storageNode
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func NodeStorageListHandler(w http.ResponseWriter, r *http.Request) {
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
				log.Fatal(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			nodeStorageInfo := types.NodeStorageInfo{}

			err = json.Unmarshal(body, &nodeStorageInfo)
			if err != nil {
				log.Fatal(err)
			}

			response[storageNode.NodeName] = nodeStorageInfo

		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func NodeStorageInfoHandler(w http.ResponseWriter, r *http.Request) {
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
				log.Fatal(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			storageInfo := types.NewStorageInfoMessage()

			err = json.Unmarshal(body, &storageInfo)
			if err != nil {
				log.Fatal(err)
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
				log.Fatal(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			nodeDiskInfo := []types.DiskMetric{}

			err = json.Unmarshal(body, &nodeDiskInfo)
			if err != nil {
				log.Fatal(err)
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
				log.Fatal(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			nodeMetricInfo := []types.NodeMetric{}

			err = json.Unmarshal(body, &nodeMetricInfo)
			if err != nil {
				log.Fatal(err)
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
				log.Fatal(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			storageMetric := types.NewStorageMetricMessage()

			err = json.Unmarshal(body, &storageMetric)
			if err != nil {
				log.Fatal(err)
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
				log.Fatal(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			cpuMetric := types.CpuMetricMessage{}

			err = json.Unmarshal(body, &cpuMetric)
			if err != nil {
				log.Fatal(err)
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
				log.Fatal(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			memoryMetric := types.MemoryMetricMessage{}

			err = json.Unmarshal(body, &memoryMetric)
			if err != nil {
				log.Fatal(err)
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
				log.Fatal(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			networkMetric := types.NetworkMetricMessage{}

			err = json.Unmarshal(body, &networkMetric)
			if err != nil {
				log.Fatal(err)
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
				log.Fatal(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			diskMetric := types.DiskMetricMessage{}

			err = json.Unmarshal(body, &diskMetric)
			if err != nil {
				log.Fatal(err)
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
				log.Fatal(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			powerMetric := types.PowerMetricMessage{}

			err = json.Unmarshal(body, &powerMetric)
			if err != nil {
				log.Fatal(err)
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
