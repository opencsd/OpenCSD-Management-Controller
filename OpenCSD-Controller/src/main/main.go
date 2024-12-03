package main

import (
	"fmt"
	"net/http"

	handler "opencsd-controller/src/handler"
	manager "opencsd-controller/src/manager"
)

func main() {
	manager.InstanceManager_ = manager.NewInstaceManager()
	manager.InstanceManager_.InitInstanceManager()

	quitChan := make(chan struct{})

	go func() {
		manager.InstanceManager_.InformerFactory.Start(quitChan)
	}()

	fmt.Println("[OpenCSD Controller] run on 0.0.0.0:", manager.OPENCSD_CONTROLLER_PORT)

	//handler
	http.HandleFunc("/create/query-engine", handler.CreateQueryEngineHandler)     //post, CreateInstanceInfo
	http.HandleFunc("/create/storage-engine", handler.CreateStorageEngineHandler) //post, CreateInstanceInfo
	http.HandleFunc("/create/validator", handler.CreateValidatorHandler)          //post, CreateInstanceInfo
	http.HandleFunc("/create/mysql", handler.CreateMysqlHandler)                  //post, CreateInstanceInfo
	http.HandleFunc("/create/graphdb", handler.CreateGraphdbHandler)              //post, CreateInstanceInfo

	http.HandleFunc("/info/instance", handler.InstanceInfoHandler) // ?node=&instance=
	http.HandleFunc("/info/volume", handler.VolumeInfoHandler)     // ?node=

	http.HandleFunc("/access/instance", handler.AccessInstanceHandler) //post, ConnectionInfo

	http.HandleFunc("/volume/allocate", handler.VolumeAllocateHandler)
	http.HandleFunc("/volume/deallocate", handler.VolumeDeAllocateHandler) // ?volume=

	http.HandleFunc("/delete/query-engine", handler.DeleteQueryEngineHandler)     // ?instance=
	http.HandleFunc("/delete/storage-engine", handler.DeleteStorageEngineHandler) // ?instance=
	http.HandleFunc("/delete/validator", handler.DeleteValidatorHandler)          // ?instance=
	http.HandleFunc("/delete/mysql", handler.DeleteMysqlHandler)                  // ?instance=
	http.HandleFunc("/delete/graphdb", handler.DeleteGraphdbHandler)              // ?instance=
	http.HandleFunc("/delete/namespace", handler.DeleteNamespace)                 // ?instance=

	http.ListenAndServe(":"+manager.OPENCSD_CONTROLLER_PORT, nil)

	close(quitChan)
}
