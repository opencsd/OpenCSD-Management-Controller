package main

import (
	"fmt"
	"net/http"

	handler "opencsd-api-server/src/handler"
	session "opencsd-api-server/src/session"
	types "opencsd-api-server/src/type"
)

// CORS 미들웨어
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // 모든 도메인 허용
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// OPTIONS 요청에 대해 허용된 CORS 설정만 반환
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	types.ManagementMaster_ = &types.MasterNode{}
	types.ManagementMaster_.InitCluster()

	session.CreateDefaultSessionHandler()

	fmt.Println("[OpenCSD API Server] run on 0.0.0.0:", types.OPENCSD_API_SERVER_PORT)

	mux := http.NewServeMux()

	//instance
	mux.HandleFunc("/instance/create/opencsd", handler.CreateOpenCSD)   //post, instanceCreationInfo
	mux.HandleFunc("/instance/create/mysql", handler.CreateMySQL)       //post, instanceCreationInfo
	mux.HandleFunc("/instance/create/graphdb", handler.CreateGraphDB)   //post, instanceCreationInfo
	mux.HandleFunc("/instance/delete/opencsd", handler.DeleteOpenCSD)   //?instance=
	mux.HandleFunc("/instance/delete/mysql", handler.DeleteMySQL)       //?instance=
	mux.HandleFunc("/instance/delete/graphdb", handler.DeleteGraphDB)   //?instance=
	mux.HandleFunc("/instance/delete/instance", handler.DeleteInstance) //?instance=
	mux.HandleFunc("instance/update/session", handler.UpdateSessionInfo)

	//volume
	mux.HandleFunc("/volume/allocate", handler.VolumeAllocateHandler)     //post, volumeCreationInfo
	mux.HandleFunc("/volume/deallocate", handler.VolumeDeAllocateHandler) //?volume=

	//cluster
	mux.HandleFunc("/cluster/node-list", handler.ClusterNodeListHandler)  // ?layer=
	mux.HandleFunc("/cluster/node/volume", handler.VolumeInfoHandler)     // ?node=
	mux.HandleFunc("/cluster/node/instance", handler.InstanceInfoHandler) // ?node=&instance=
	mux.HandleFunc("/cluster/storage-node-list", handler.ClusterStorageNodeListHandler)

	//node
	mux.HandleFunc("/node/info/storage-list", handler.NodeStorageListHandler) // ?node=&count=
	mux.HandleFunc("/node/info/storage", handler.NodeStorageInfoHandler)      // ?node=&count=
	mux.HandleFunc("/node/metric/disk", handler.NodeDiskInfoHandler)          // ?node=&count=
	mux.HandleFunc("/node/metric/all", handler.NodeMetricInfoHandler)         // ?node=&count=

	//storage
	mux.HandleFunc("/storage/info", handler.NodeStorageInfoHandler)                // ?node=&storage=&count=
	mux.HandleFunc("/storage/metric/all", handler.StorageMetricInfoHandler)        // ?node=&storage=&count=
	mux.HandleFunc("/storage/metric/cpu", handler.StorageMetricCpuHandler)         // ?node=&storage=&count=
	mux.HandleFunc("/storage/metric/memory", handler.StorageMetricMemoryHandler)   // ?node=&storage=&count=
	mux.HandleFunc("/storage/metric/network", handler.StorageMetricNetworkHandler) // ?node=&storage=&count=
	mux.HandleFunc("/storage/metric/disk", handler.StorageMetricDiskHandler)       // ?node=&storage=&count=

	//workbench
	mux.HandleFunc("/workbench/main/connection", handler.ConnectInstance)                     //post, ConnectionInfo
	mux.HandleFunc("/workbench/monitoring/connection-info", handler.MonitoringConnectionInfo) //get, ?session-id=
	mux.HandleFunc("/workbench/monitoring/metric/node", handler.NodeMetric)                   //get, ?session-id=&count=
	mux.HandleFunc("/workbench/monitoring/metric/instance", handler.InstanceMetric)           //get, ?session-id=&count=
	mux.HandleFunc("/workbench/monitoring/metric/csd", handler.CsdMetric)
	mux.HandleFunc("/workbench/query/metric", handler.NodeMetricMin) //get, ?session-id=&count=&start-time=&end-time=
	mux.HandleFunc("/workbench/query/run", handler.RunQuery)
	mux.HandleFunc("/workbench/query/terminate", handler.TerminateQuery)
	mux.HandleFunc("/workbench/query/log/get", handler.GetQueryLog)       //get, ?session-id=?log-id=
	mux.HandleFunc("/workbench/query/log/delete", handler.DeleteQueryLog) //get, ?session-id=
	mux.HandleFunc("/workbench/query-ssd/run", handler.RunQuerySsd)
	// mux.HandleFunc("/workbench/query/schema-info", handler.GetSchemaInfo)                     //get, ?session-id=
	// mux.HandleFunc("/workbench/query/environment", handler.QueryEnvInfo)                      //get, post, ?session-id=

	//gluesys
	mux.HandleFunc("/dashboard/cluster/storage-node-list", handler.ClusterStorageNodeListHandler) //1. ClusterNodeList
	mux.HandleFunc("/dashboard/node/info/storage-list", handler.NodeStorageListHandler)           //2. NodeStorageList
	mux.HandleFunc("/dashboard/node/info/storage", handler.NodeStorageInfoHandler)                //4. NodeStorageInfo
	mux.HandleFunc("/dashboard/node/metric/disk", handler.NodeDiskInfoHandler)                    //3. NodeDiskInfo
	mux.HandleFunc("/dashboard/node/metric/all", handler.NodeMetricInfoHandler)                   //5. NodeMetricInfo
	mux.HandleFunc("/dashboard/storage/info", handler.NodeStorageInfoHandler)                     //6. StorageInfo
	mux.HandleFunc("/dashboard/storage/metric/all", handler.StorageMetricInfoHandler)             //7. CSDMetricInfo
	mux.HandleFunc("/dashboard/storage/metric/cpu", handler.StorageMetricCpuHandler)              //8. CSDCpuInfo
	mux.HandleFunc("/dashboard/storage/metric/memory", handler.StorageMetricMemoryHandler)        //9. CSDMemInfo
	mux.HandleFunc("/dashboard/storage/metric/network", handler.StorageMetricNetworkHandler)      //10. CSDNetInfo
	mux.HandleFunc("/dashboard/storage/metric/disk", handler.StorageMetricDiskHandler)            //11. CSDDiskInfo

	// CORS 미들웨어와 함께 서버 실행
	http.ListenAndServe(":"+types.OPENCSD_API_SERVER_PORT, enableCORS(mux))
}
