package handler

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"

	controller "opencsd-controller/src/controller"
	manager "opencsd-controller/src/manager"
)

func createInstanceInfoDatabase(instanceName string, operationNode string) error {
	operaionNodeIp := manager.InstanceManager_.OperationLayer[operationNode].IP

	url := "http://" + operaionNodeIp + ":" + manager.OPENCSD_INSTANCE_METRIC_COLLECTOR_PORT + "/create/instance"
	url += "?instance=" + instanceName

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Request failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func CreateQueryEngineHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Controller] CreateQueryEngineHandler")

	instanceInfo := controller.CreateInstanceInfo{}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(string(requestBody))

	err = json.Unmarshal(requestBody, &instanceInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.CreateNamespace(instanceInfo.InstanceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.CreateQueryEngineDeployment(instanceInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if instanceInfo.InstanceType == manager.OPENCSD {
		err = createInstanceInfoDatabase(instanceInfo.InstanceName, instanceInfo.OperationNode)
		if err != nil {
			w.Write([]byte("create instance info database failed\n"))
			return
		}
	}

	w.Write([]byte("[OpenCSD Controller] Create QueryEngine Successfully\n"))
}

func CreateStorageEngineHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Controller] CreateStorageEngineHandler")

	instanceInfo := controller.CreateInstanceInfo{}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(string(requestBody))

	err = json.Unmarshal(requestBody, &instanceInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.CreateNamespace(instanceInfo.InstanceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.CreateVolumeDirectory(instanceInfo.VolumeName, instanceInfo.InstanceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.CreateStorageEngineDeployment(instanceInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Write([]byte("[OpenCSD Controller] Create StorageEngine Successfully\n"))
}

func CreateValidatorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Controller] CreateStorageEngineHandler")

	instanceInfo := controller.CreateInstanceInfo{}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(string(requestBody))

	err = json.Unmarshal(requestBody, &instanceInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.CreateNamespace(instanceInfo.InstanceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.CreateValidatorDeployment(instanceInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Write([]byte("[OpenCSD Controller] Create Validator Successfully\n"))
}

func CreateMysqlHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Controller] CreateMysqlHandler")

	instanceInfo := controller.CreateInstanceInfo{}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(string(requestBody))

	err = json.Unmarshal(requestBody, &instanceInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.CreateNamespace(instanceInfo.InstanceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.CreateVolumeDirectory(instanceInfo.VolumeName, instanceInfo.InstanceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.CreateMySQLDeployment(instanceInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Write([]byte("[OpenCSD Controller] Create MySQL Successfully\n"))
}

func CreateGraphdbHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Controller] CreateGraphdbHandler")

	instanceInfo := controller.CreateInstanceInfo{}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(string(requestBody))

	err = json.Unmarshal(requestBody, &instanceInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.CreateNamespace(instanceInfo.InstanceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.CreateVolumeDirectory(instanceInfo.VolumeName, instanceInfo.InstanceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.CreateGraphDBDeployment(instanceInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Write([]byte("[OpenCSD Controller] Create GraphDB Successfully\n"))
}

func DeleteInstance(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Controller] DeleteInstance")

	targetInstance := r.URL.Query().Get("instance")
	fmt.Println("Delete Instance; ", targetInstance)

	err := controller.DeleteInstance(targetInstance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("[OpenCSD Controller] Delete Instance Successfully\n"))
}

func DeleteQueryEngineHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Controller] DeleteQueryEngineHandler")

	targetInstance := r.URL.Query().Get("instance")
	fmt.Println("Delete Query Engine; ", targetInstance)

	err := controller.DeleteQueryEngineDeployment(targetInstance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("[OpenCSD Controller] Delete QueryEngine Successfully\n"))
}

func DeleteStorageEngineHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Controller] DeleteStorageEngineHandler")

	targetInstance := r.URL.Query().Get("instance")
	fmt.Println("Delete Storage Engine; ", targetInstance)

	err := controller.DeleteStorageEngineDeployment(targetInstance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("[OpenCSD Controller] Delete StorageEngine Successfully\n"))
}

func DeleteValidatorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Controller] DeleteValidatorHandler")

	targetInstance := r.URL.Query().Get("instance")
	fmt.Println("Delete Validator Engine; ", targetInstance)

	err := controller.DeleteValidatorDeployment(targetInstance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("[OpenCSD Controller] Delete Validator Successfully\n"))
}

func DeleteMysqlHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Controller] DeleteMysqlHandler")

	targetInstance := r.URL.Query().Get("instance")
	fmt.Println("Delete MySQL; ", targetInstance)

	err := controller.DeleteMySQLDeployment(targetInstance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("[OpenCSD Controller] Delete MySQL Successfully\n"))
}

func DeleteGraphdbHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Controller] DeleteGraphdbHandler")

	targetInstance := r.URL.Query().Get("instance")
	fmt.Println("Delete Graphdb; ", targetInstance)

	err := controller.DeleteGraphDBDeployment(targetInstance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("[OpenCSD Controller] Delete GraphDB Successfully\n"))
}

func VolumeAllocateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Controller] VolumeAllocateHandler")

	volumeInfo := controller.CreateVolumeInfo{}
	responseBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(string(responseBody))

	err = json.Unmarshal(responseBody, &volumeInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.AllocateVolume(volumeInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("[OpenCSD Controller] Allocate Volume Successfully\n"))
}

func VolumeDeAllocateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Controller] VolumeDeAllocateHandler")

	targetVolume := r.URL.Query().Get("volume")
	fmt.Println("Delete Volume; ", targetVolume)

	err := controller.DeAllocateVolume(targetVolume)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("[OpenCSD Controller] Deallocate Volume Successfully\n"))
}

func InstanceInfoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Controller] InstanceInfoHandler")
	var jsonResponse []byte

	type instance_ struct {
		InstanceInfo manager.Instance   `json:"instanceInfo"`
		VolumeInfo   manager.VolumeInfo `json:"volumeInfo"`
	}
	response := map[string]instance_{}

	instanceName := r.URL.Query().Get("instance")
	nodeName := r.URL.Query().Get("node")

	for key1, operationNodeInfo := range manager.InstanceManager_.OperationNodeInfo {
		if nodeName == "" || nodeName == key1 {
			for key2, instance := range operationNodeInfo.InstanceInfo {
				if instanceName == "" || instanceName == key2 {
					var myInstance instance_
					myInstance.InstanceInfo = *instance
					if instance.StorageEngineStatus == manager.RUNNING {
						myInstance.VolumeInfo = *manager.InstanceManager_.VolumeInfo[instance.VolumeName]
					}
					response[key2] = myInstance
				}
			}
		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func VolumeInfoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Controller] VolumeInfoHandler")
	var jsonResponse []byte

	type VolumeMessage struct {
		Name       string                         `json:"nodeName"`
		IP         string                         `json:"nodeIp"`
		Status     string                         `json:"status"`
		VolumeInfo map[string]*manager.VolumeInfo `json:"volumeInfo"`
	}

	response := map[string]VolumeMessage{}

	targetNode := r.URL.Query().Get("node")

	for _, storageNode := range manager.InstanceManager_.StorageLayer {
		var volumeMessage VolumeMessage
		if targetNode == "" || targetNode == storageNode.Name {
			volumeMessage.Name = storageNode.Name
			volumeMessage.IP = storageNode.IP
			volumeMessage.Status = storageNode.Status
			volumeMessage.VolumeInfo = manager.InstanceManager_.StorageNodeInfo[storageNode.Name].VolumeInfo
			response[storageNode.Name] = volumeMessage
		}
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func AccessInstanceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD Controller] AccessInstanceHandler")

	type ConnectionInfo struct {
		InstanceName string `json:"instanceName"`
		UserName     string `json:"userName"`
		DbName       string `json:"dbName"`
		DbUser       string `json:"dbUser"`
		DbPassword   string `json:"dbPassword"`
	}

	connectionInfo := ConnectionInfo{}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(requestBody, &connectionInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(string(requestBody))

	instanceType := manager.InstanceManager_.InstanceInfo[connectionInfo.InstanceName].InstanceType
	operationNode := manager.InstanceManager_.InstanceInfo[connectionInfo.InstanceName].OperationNode
	storageNode := manager.InstanceManager_.InstanceInfo[connectionInfo.InstanceName].StorageNode
	response, dbHost := "", ""

	if instanceType == manager.GRAPHDB {
		dbHost = "graphdb-dbms-svc." + connectionInfo.InstanceName + ".svc.cluster.local"
		response = fmt.Sprintf(`{"status": "false", "instanceType": "%s"}`, instanceType)
	} else {
		if instanceType == manager.MYSQL {
			dbHost = "mysql-dbms-svc." + connectionInfo.InstanceName + ".svc.cluster.local"
		} else { // OPENCSD
			dbHost = "storage-engine-dbms-svc." + connectionInfo.InstanceName + ".svc.cluster.local"
		}
		dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
			"root",                    // 사용자 이름
			connectionInfo.DbPassword, // 비밀번호
			dbHost,                    // 데이터베이스 호스트
			connectionInfo.DbName,     // 데이터베이스 이름
		)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			fmt.Printf("Error opening database connection: %v\n", err)
			response = fmt.Sprintf(`{"status": "false", "instanceType": "%s"}`, instanceType)
		}
		defer db.Close()

		err = db.Ping()
		if err != nil {
			fmt.Printf("Error connecting to database: %v\n", err)
			response = fmt.Sprintf(`{"status": "false", "instanceType": "%s"}`, instanceType)
		} else {
			response = fmt.Sprintf(`{"status": "true", "instanceType": "%s", "operationNode": "%s", "storageNode": "%s"}`, instanceType, operationNode, storageNode)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
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
