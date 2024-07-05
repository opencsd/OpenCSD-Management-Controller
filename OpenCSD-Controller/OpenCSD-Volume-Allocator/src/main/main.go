package main

import (
	"net/http"
	"fmt"

	handler "volume-allocator/src/handler"
)

func main() {
	fmt.Println("[OpenCSD Volume Allocator] Running...")

	//handler
	http.HandleFunc("/info/StorageNodeInfo", handler.StorageNodeInfo)
	http.HandleFunc("/info/StorageVolumeInfo", handler.StorageVolumeInfo)

	http.HandleFunc("/allocate/AllocateVolume", handler.AllocateVolume)
	http.HandleFunc("/allocate/DeallocateVolume", handler.DeallocateVolume)

	http.ListenAndServe(":40806", nil)
}
