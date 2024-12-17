package main

import (
	"fmt"
	"net/http"

	handler "opencsd-controller/src/handler"
	manager "opencsd-controller/src/manager"
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
	manager.InstanceManager_ = manager.NewInstaceManager()
	manager.InstanceManager_.InitInstanceManager()

	quitChan := make(chan struct{})

	go func() {
		manager.InstanceManager_.InformerFactory.Start(quitChan)
	}()

	fmt.Println("[OpenCSD Controller] run on 0.0.0.0:", manager.OPENCSD_CONTROLLER_PORT)

	mux := http.NewServeMux()

	//handler
	mux.HandleFunc("/create/query-engine", handler.CreateQueryEngineHandler)     //post, CreateInstanceInfo
	mux.HandleFunc("/create/storage-engine", handler.CreateStorageEngineHandler) //post, CreateInstanceInfo
	mux.HandleFunc("/create/validator", handler.CreateValidatorHandler)          //post, CreateInstanceInfo
	mux.HandleFunc("/create/mysql", handler.CreateMysqlHandler)                  //post, CreateInstanceInfo
	mux.HandleFunc("/create/graphdb", handler.CreateGraphdbHandler)              //post, CreateInstanceInfo

	mux.HandleFunc("/info/instance", handler.InstanceInfoHandler) // ?node=&instance=
	mux.HandleFunc("/info/volume", handler.VolumeInfoHandler)     // ?node=

	mux.HandleFunc("/access/instance", handler.AccessInstanceHandler) //post, ConnectionInfo

	mux.HandleFunc("/volume/allocate", handler.VolumeAllocateHandler)
	mux.HandleFunc("/volume/deallocate", handler.VolumeDeAllocateHandler) // ?volume=

	mux.HandleFunc("/delete/query-engine", handler.DeleteQueryEngineHandler)     // ?instance=
	mux.HandleFunc("/delete/storage-engine", handler.DeleteStorageEngineHandler) // ?instance=
	mux.HandleFunc("/delete/validator", handler.DeleteValidatorHandler)          // ?instance=
	mux.HandleFunc("/delete/mysql", handler.DeleteMysqlHandler)                  // ?instance=
	mux.HandleFunc("/delete/graphdb", handler.DeleteGraphdbHandler)              // ?instance=
	mux.HandleFunc("/delete/namespace", handler.DeleteNamespace)                 // ?instance=

	http.ListenAndServe(":"+manager.OPENCSD_CONTROLLER_PORT, enableCORS(mux))

	close(quitChan)
}
