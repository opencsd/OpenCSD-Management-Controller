package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	session "opencsd-api-server/src/session"
	types "opencsd-api-server/src/type"
	"strconv"
	"strings"

	client "github.com/influxdata/influxdb/client/v2"
)

func ConnectInstance(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD API Server] ConnectInstanceHandler")

	connectionInfo := session.ConnectionInfo{}

	request, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(request, &connectionInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url := "http://" + types.OPENCSD_CONTROLLER_DNS + "/access/instance"
	reqBody, err := json.Marshal(connectionInfo)
	if err != nil {
		http.Error(w, "Failed to marshal JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		http.Error(w, "Failed to send request to controller: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read controller response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	type conn struct {
		Status        string `json:"status"`
		InstanceType  string `json:"instanceType"`
		OperationNode string `json:"operationNode"`
		StorageNode   string `json:"storageNode"`
	}

	var conn_ conn

	err = json.Unmarshal(body, &conn_)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if conn_.Status == "true" {
		sessionId := session.CreateSessionHandler(connectionInfo, conn_.InstanceType, conn_.OperationNode, conn_.StorageNode)

		response := map[string]string{
			"sessionId":    sessionId,
			"instanceType": conn_.InstanceType,
		}

		err = AddUser(connectionInfo.InstanceName, connectionInfo.UserName, conn_.OperationNode)
		if err != nil {
			fmt.Println("failed to add user to instance info db: ", err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Connection not authorized by controller", http.StatusUnauthorized)
	}
}

func MonitoringConnectionInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[OpenCSD API Server] MonitoringConnectionInfo")

	var jsonResponse []byte

	type instance_ struct {
		InstanceInfo   types.Instance            `json:"instanceInfo"`
		VolumeInfo     types.VolumeInfo          `json:"volumeInfo"`
		ConnectionInfo session.ConnectionInfoMin `json:"connectionInfo"`
	}

	sessionId := r.URL.Query().Get("session-id")
	connectionInfo, exist := session.GetConnectionInfo(sessionId)
	if !exist {
		http.Error(w, "session id error", http.StatusBadRequest)
		return
	}
	instanceName := connectionInfo.InstanceName

	response := make(map[string]instance_)

	url := "http://" + types.OPENCSD_CONTROLLER_DNS + "/info/instance"
	url += "?instance=" + instanceName

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

	instance := response[instanceName]
	instance.ConnectionInfo = session.ConvertToConnectionInfoMin(session.WorkbenchSessionStore[sessionId])

	jsonResponse, _ = json.Marshal(instance)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func NodeMetric(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := []types.NodeMetric{}

	sessionId := r.URL.Query().Get("session-id")
	count := r.URL.Query().Get("count")

	connectionInfo := session.WorkbenchSessionStore[sessionId]
	operationNode := connectionInfo.OperationNode
	operaionNodeIp := types.ManagementMaster_.OperationLayer[operationNode].NodeIP

	if count == "" {
		count = "1"
	}

	q := client.Query{
		Command:  "SELECT * FROM node_metric ORDER BY DESC LIMIT " + count + " TZ('Asia/Seoul')",
		Database: types.INSTANCE_METRIC_INFLUXDB_DB,
	}

	INFLUX_CLIENT, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://" + operaionNodeIp + ":" + types.INSTANCE_METRIC_INFLUXDB_PORT,
		Username: types.INSTANCE_METRIC_INFLUXDB_USER,
		Password: types.INSTANCE_METRIC_INFLUXDB_PASSWORD,
	})
	if err != nil {
		log.Fatal(err)
	}

	if result, err := INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
		for _, row := range result.Results[0].Series {
			for _, value := range row.Values {
				nodeMetric := types.NodeMetric{}

				nodeMetric.Time = fmt.Sprintf("%v", value[0])
				nodeMetric.CpuTotal = parseFloat(value[1])
				nodeMetric.CpuUsed = parseFloat(value[2])
				nodeMetric.CpuUtilization = parseFloat(value[3])
				nodeMetric.DiskTotal = parseFloat(value[4])
				nodeMetric.DiskUsed = parseFloat(value[5])
				nodeMetric.DiskUtilization = parseFloat(value[6])
				nodeMetric.MemoryTotal = parseFloat(value[7])
				nodeMetric.MemoryUsed = parseFloat(value[8])
				nodeMetric.MemoryUtilization = parseFloat(value[9])
				nodeMetric.NetworkBandwidth = parseFloat(value[10])
				nodeMetric.NetworkRxData = parseFloat(value[11])
				nodeMetric.NetworkTxData = parseFloat(value[12])
				nodeMetric.NodeName = fmt.Sprintf("%v", value[13])
				nodeMetric.PowerUsed = parseFloat(value[14])

				response = append(response, nodeMetric)
			}
		}
	} else {
		fmt.Println("Error executing query:", err)
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func InstanceMetric(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte
	response := []types.InstanceMetric{}

	sessionId := r.URL.Query().Get("session-id")
	count := r.URL.Query().Get("count")

	connectionInfo := session.WorkbenchSessionStore[sessionId]

	instanceName := connectionInfo.InstanceName
	operationNode := connectionInfo.OperationNode
	operationNodeIp := types.ManagementMaster_.OperationLayer[operationNode].NodeIP

	if count == "" {
		count = "1"
	}

	converted := strings.ReplaceAll(instanceName, "-", "_")
	measurementName := "instance_metric_" + converted

	q := client.Query{
		Command:  "SELECT * FROM " + measurementName + " ORDER BY DESC LIMIT " + count + " TZ('Asia/Seoul')",
		Database: types.INSTANCE_METRIC_INFLUXDB_DB,
	}

	INFLUX_CLIENT, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://" + operationNodeIp + ":" + types.INSTANCE_METRIC_INFLUXDB_PORT,
		Username: types.INSTANCE_METRIC_INFLUXDB_USER,
		Password: types.INSTANCE_METRIC_INFLUXDB_PASSWORD,
	})
	if err != nil {
		log.Fatal(err)
	}

	if result, err := INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
		for _, row := range result.Results[0].Series {
			for _, value := range row.Values {
				instanceMetric := types.InstanceMetric{}

				instanceMetric.Time = fmt.Sprintf("%v", value[0])
				instanceMetric.CpuUsage = parseFloat(value[1])
				instanceMetric.InstanceName = fmt.Sprintf("%v", value[2])
				instanceMetric.MemoryUsage = parseFloat(value[3])
				instanceMetric.NetworkRxUsage = parseFloat(value[4])
				instanceMetric.NetworkTxUsage = parseFloat(value[5])
				instanceMetric.StorageUsage = parseFloat(value[6])

				response = append(response, instanceMetric)
			}
		}
	} else {
		fmt.Println("Error executing query:", err)
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func CsdMetric(w http.ResponseWriter, r *http.Request) {
	var jsonResponse []byte

	targetCsd := r.URL.Query().Get("csd-id")
	sessionId := r.URL.Query().Get("session-id")
	count := r.URL.Query().Get("count")
	connectionInfo := session.WorkbenchSessionStore[sessionId]
	storageNode := connectionInfo.StorageNode
	storageNodeIp := types.ManagementMaster_.StorageLayer[storageNode].NodeIP
	storageNodeIp = "10.0.4.83" //temp

	url := "http://" + storageNodeIp + ":" + types.STORAGE_API_SERVER_PORT + "/storage/metric/all"
	url += "?storage=" + targetCsd
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

	jsonResponse, _ = json.Marshal(storageMetric.CsdList)
	w.Write(jsonResponse)
}

func ExecuteQuery(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query().Get("session-id")
	connectionInfo := session.WorkbenchSessionStore[sessionId]
	instanceName := connectionInfo.InstanceName

	url := "http://query-engine-instance-svc." + instanceName + ".svc.cluster.local:40100/query/run"

	queryRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(queryRequest))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	w.Write([]byte(string(body) + "\n"))
}

func GetSchemaInfo(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query().Get("session-id")
	connectionInfo := session.WorkbenchSessionStore[sessionId]
	instanceName := connectionInfo.InstanceName

	url := "http://query-engine-instance-svc." + instanceName + ".svc.cluster.local:40100/info/schema"
	url += "?db-name=" + connectionInfo.DbName

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	w.Write([]byte(string(body) + "\n"))
}

func QueryEnvInfo(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query().Get("session-id")
	connectionInfo := session.WorkbenchSessionStore[sessionId]
	instanceName := connectionInfo.InstanceName

	//get
	url := "http://query-engine-instance-svc." + instanceName + ".svc.cluster.local:40100/metadata/environment"

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	w.Write([]byte(string(body) + "\n"))
}

func QueryEnvEdit(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query().Get("session-id")
	connectionInfo := session.WorkbenchSessionStore[sessionId]
	instanceName := connectionInfo.InstanceName

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//post
	url := "http://query-engine-instance-svc." + instanceName + ".svc.cluster.local:40100/metadata/environment"

	queryEngineRequest, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
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

	w.Write([]byte("[OpenCSD API Server] Edit Environment Successfully\n"))
}

func NodeMetricMin(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query().Get("session-id")
	startTime := r.URL.Query().Get("start-time")
	endTime := r.URL.Query().Get("end-time")
	count := r.URL.Query().Get("count")
	connectionInfo := session.WorkbenchSessionStore[sessionId]
	operationNode := connectionInfo.OperationNode
	operaionNodeIp := types.ManagementMaster_.OperationLayer[operationNode].NodeIP

	var jsonResponse []byte
	response := []types.NodeMetric{}

	if count == "" {
		count = "1"
	}

	query := ""
	if startTime == "" || endTime == "" {
		query = fmt.Sprintf(
			"SELECT cpu_usage, power_usage FROM node_metric ORDER BY time DESC LIMIT %s TZ('Asia/Seoul')",
			count,
		)
	} else {
		query = fmt.Sprintf(
			"SELECT cpu_usage, power_usage FROM node_metric WHERE time > '%s' - 5s AND time < '%s' + 5s ORDER BY time DESC LIMIT %s TZ('Asia/Seoul')",
			startTime, endTime, count,
		)
	}

	q := client.Query{
		Command:  query,
		Database: types.INSTANCE_METRIC_INFLUXDB_DB,
	}

	INFLUX_CLIENT, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://" + operaionNodeIp + ":" + types.INSTANCE_METRIC_INFLUXDB_PORT,
		Username: types.INSTANCE_METRIC_INFLUXDB_USER,
		Password: types.INSTANCE_METRIC_INFLUXDB_PASSWORD,
	})
	if err != nil {
		log.Fatal(err)
	}

	if result, err := INFLUX_CLIENT.Query(q); err == nil && result.Error() == nil {
		for _, row := range result.Results[0].Series {
			for _, value := range row.Values {
				nodeMetric := types.NodeMetric{}

				nodeMetric.Time = fmt.Sprintf("%v", value[0])
				nodeMetric.CpuTotal = parseFloat(value[1])
				nodeMetric.CpuUsed = parseFloat(value[2])
				nodeMetric.CpuUtilization = parseFloat(value[3])
				nodeMetric.DiskTotal = parseFloat(value[4])
				nodeMetric.DiskUsed = parseFloat(value[5])
				nodeMetric.DiskUtilization = parseFloat(value[6])
				nodeMetric.MemoryTotal = parseFloat(value[7])
				nodeMetric.MemoryUsed = parseFloat(value[8])
				nodeMetric.MemoryUtilization = parseFloat(value[9])
				nodeMetric.NetworkBandwidth = parseFloat(value[10])
				nodeMetric.NetworkRxData = parseFloat(value[11])
				nodeMetric.NetworkTxData = parseFloat(value[12])
				nodeMetric.NodeName = fmt.Sprintf("%v", value[13])
				nodeMetric.PowerUsed = parseFloat(value[14])

				response = append(response, nodeMetric)
			}
		}
	} else {
		fmt.Println("Error executing query:", err)
	}

	jsonResponse, _ = json.Marshal(response)
	w.Write([]byte(string(jsonResponse) + "\n"))
}

func DeleteQueryLog(w http.ResponseWriter, r *http.Request) {

}

func GetQueryLog(w http.ResponseWriter, r *http.Request) {

}

func parseFloat(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return math.Round(v*100) / 100
	case int:
		v_ := float64(v)
		return math.Round(v_*100) / 100
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return math.Round(f*100) / 100
		}
	case json.Number:
		f, _ := v.Float64()
		return math.Round(f*100) / 100
	default:
		fmt.Printf("Unknown type: %T with value: %v\n", v, v)
	}
	return 0
}
