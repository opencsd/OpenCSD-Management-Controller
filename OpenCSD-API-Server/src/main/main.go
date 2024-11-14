package main

import (
	"fmt"
	"net/http"

	handler "opencsd-api-server/src/handler"
	types "opencsd-api-server/src/type"
)

func main() {
	types.ManagementMaster_ = &types.MasterNode{}
	types.ManagementMaster_.InitCluster()

	fmt.Println("[OpenCSD API Server] run on 0.0.0.0:", types.OPENCSD_API_SERVER_PORT)

	http.HandleFunc("/instance/create/opencsd", handler.CreateOpenCSD)
	http.HandleFunc("/instance/create/mysql", handler.CreateMySQL)

	//1. ClusterNodeList
	// http.HandleFunc("/dashboard/cluster/storage-node-list", handler.ClusterNodeListHandler)
	http.HandleFunc("/dashboard/cluster/storage-node-list", handler.ClusterNodeListHandler)

	//2. NodeStorageList
	// http.HandleFunc("/dashboard/node/storagelist", handler.NodeStorageListHandler)
	// http.HandleFunc("/storagepage/storage/storagelist", handler.NodeStorageListHandler)
	// http.HandleFunc("/diskpage/disk/storagelist", handler.NodeStorageListHandler)
	http.HandleFunc("/dashboard/node/info/storage-list", handler.NodeStorageListHandler)

	//4. NodeStorageInfo
	// http.HandleFunc("/dashboard/storage/storageinfo", handler.NodeStorageInfoHandler)
	http.HandleFunc("/dashboard/node/info/storage", handler.NodeStorageInfoHandler)

	//3. NodeDiskInfo
	// http.HandleFunc("/dashboard/node/diskinfo", handler.NodeDiskInfoHandler)
	http.HandleFunc("/dashboard/node/metric/disk", handler.NodeDiskInfoHandler)

	//5. NodeMetricInfo
	// http.HandleFunc("/dashboard/node/metricinfo", handler.NodeMetricInfoHandler)
	http.HandleFunc("/dashboard/node/metric/all", handler.NodeMetricInfoHandler)

	//6. StorageInfo
	// http.HandleFunc("/storagepage/storage/storageinfo", handler.StorageInfoHandler)
	http.HandleFunc("/dashboard/storage/info", handler.NodeStorageInfoHandler)

	//7. CSDMetricInfo
	// http.HandleFunc("/storagepage/storage/csdmetricinfo", handler.CSDMetricInfoHandler)
	http.HandleFunc("/dashboard/storage/metric/all", handler.StorageMetricInfoHandler)

	//8. CSDCpuInfo
	// http.HandleFunc("/storagepage/storage/csdcpuinfo", handler.CPUInfoHandler)
	http.HandleFunc("/dashboard/storage/metric/cpu", handler.StorageMetricCpuHandler)

	//9. CSDMemInfo
	// http.HandleFunc("/storagepage/storage/csdmeminfo", handler.MemInfoHandler)
	http.HandleFunc("/dashboard/storage/metric/memory", handler.StorageMetricMemoryHandler)

	//10. CSDNetInfo
	// http.HandleFunc("/storagepage/storage/csdnetinfo", handler.NetInfoHandler)
	http.HandleFunc("/dashboard/storage/metric/network", handler.StorageMetricNetworkHandler)

	//11. CSDDiskInfo
	// http.HandleFunc("/storagepage/storage/csddiskinfo", handler.DiskInfoHandler)
	http.HandleFunc("/dashboard/storage/metric/disk", handler.StorageMetricDiskHandler)

	// http.HandleFunc("/dashboard/storage/metric/power", handler.StorageMetricPowerHandler)

	http.ListenAndServe(":"+types.OPENCSD_API_SERVER_PORT, nil)
}
