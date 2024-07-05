package main

import (
	"net/http"
	"fmt"

	handler "engine-deployer/src/handler"
)

func main() {
	fmt.Println("[OpenCSD Engine Deployer] Running...")

	//handler
	http.HandleFunc("/create/query-engine", handler.CreateQueryEngine)
	http.HandleFunc("/create/storage-engine", handler.CreateStorageEngine)
	http.HandleFunc("/create/validator", handler.CreateValidator)

	http.HandleFunc("/info/instance-info", handler.Info)

	http.HandleFunc("/delete/query-engine", handler.DeleteQueryEngine)
	http.HandleFunc("/delete/storage-engine", handler.DeleteStorageEngine)
	http.HandleFunc("/delete/validator", handler.DeleteValidator)

	http.ListenAndServe(":40804", nil)	
}
