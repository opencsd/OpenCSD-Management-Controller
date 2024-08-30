package main

import (
	"net/http"
	"fmt"

	handler "volume-allocator/src/handler"
)

func main() {
	fmt.Println("[OpenCSD Volume Allocator] Running...")

	//handler
	http.HandleFunc("/info/storage-node-info", handler.StorageNodeInfo)
	http.HandleFunc("/info/storage-volume-info", handler.StorageVolumeInfo)

	http.HandleFunc("/allocate/allocate-volume", handler.AllocateVolume)
	http.HandleFunc("/allocate/deallocate-volume", handler.DeallocateVolume)

	http.ListenAndServe(":40806", nil)
}
