package session

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var (
	WorkbenchSessionStore = make(map[string]ConnectionInfo)
	mu                    sync.Mutex
)

type ConnectionInfo struct {
	InstanceName  string `json:"instanceName"`
	UserName      string `json:"userName"`
	DbName        string `json:"dbName"`
	DbUser        string `json:"dbUser"`
	DbPassword    string `json:"dbPassword"`
	InstanceType  string `json:"instanceType"`
	OperationNode string `json:"operationNode"`
	StorageNode   string `json:"storageNode"`
}

type ConnectionInfoMin struct {
	InstanceName string `json:"instanceName"`
	UserName     string `json:"userName"`
	DbName       string `json:"dbName"`
	DbUser       string `json:"dbUser"`
	DbPassword   string `json:"dbPassword"`
}

func ConvertToConnectionInfoMin(ci ConnectionInfo) ConnectionInfoMin {
	return ConnectionInfoMin{
		InstanceName: ci.InstanceName,
		UserName:     ci.UserName,
		DbName:       ci.DbName,
		DbUser:       ci.DbUser,
		DbPassword:   ci.DbPassword,
	}
}

func CreateSessionHandler(connectionInfo ConnectionInfo, instanceType string, operationNode string, storageNode string) string {
	sessionId := uuid.New().String()

	mu.Lock()
	connectionInfo.InstanceType = instanceType
	connectionInfo.OperationNode = operationNode
	connectionInfo.StorageNode = storageNode
	WorkbenchSessionStore[sessionId] = connectionInfo
	mu.Unlock()

	return sessionId
}

func CreateDefaultSessionHandler() {
	connectionInfo := ConnectionInfo{}

	mu.Lock()
	connectionInfo.InstanceName = "keti-opencsd"
	connectionInfo.UserName = "keti-admin"
	connectionInfo.DbName = "tpch_origin"
	connectionInfo.DbUser = "root"
	connectionInfo.DbPassword = "keti"
	connectionInfo.InstanceType = "OPENCSD"
	connectionInfo.OperationNode = "operation-node1"
	connectionInfo.StorageNode = "storage-node1"
	WorkbenchSessionStore["-1"] = connectionInfo
	mu.Unlock()

	return
}

func GetConnectionInfo(sessionId string) (ConnectionInfo, bool) {
	mu.Lock()
	connectionInfo, exists := WorkbenchSessionStore[sessionId]
	mu.Unlock()

	if !exists {
		fmt.Print("no session id in WorkbenchSessionStore")
		return ConnectionInfo{}, false
	}

	return connectionInfo, true
}
