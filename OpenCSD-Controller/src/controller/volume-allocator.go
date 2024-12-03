package controller

import (
	"fmt"
	"net/http"
	manager "opencsd-controller/src/manager"
)

type CreateVolumeInfo struct {
	NodeName    string `json:"nodeName"`
	StorageType string `json:"storageType"`
	VolumeName  string `json:"volumeName"`
}

func AllocateVolume(volumeInfo CreateVolumeInfo) error {
	return nil
}

func DeAllocateVolume(volumeName string) error {
	return nil
}

func CreateVolumeDirectory(volumeName string, instanceName string) error {
	storageNodeName := manager.InstanceManager_.VolumeInfo[volumeName].NodeName
	storageNodeIp := manager.InstanceManager_.StorageLayer[storageNodeName].IP
	rootPath := manager.InstanceManager_.VolumeInfo[volumeName].VolumePath
	volumePath := rootPath + "/" + instanceName

	url := "http://" + storageNodeIp + ":" + manager.STORAGE_API_SERVER_PORT + "/directory/create"
	url += "?path=" + volumePath

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to create directory to %s: %v", storageNodeIp, err)
	}
	defer resp.Body.Close()

	fmt.Printf("Folder '%s' created successfully on %s\n", volumePath, storageNodeIp)
	return nil
}
