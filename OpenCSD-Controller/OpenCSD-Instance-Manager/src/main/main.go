package main

import (
	"net/http"
	"fmt"

	handler "instance-manager/src/handler"
)

func main() {
	fmt.Println("[OpenCSD Instance Manager] Running...")

	//handler
	http.HandleFunc("/map", handler.MappingInstanceStorage)
	http.HandleFunc("/unmap", handler.UnMappingInstanceStorage)

	http.ListenAndServe(":40805", nil)
}
