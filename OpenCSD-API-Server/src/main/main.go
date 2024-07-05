package main

import (
	"net/http"
	"fmt"

	handler "api-server/src/handler"
)

func main() {
	fmt.Println("[OpenCSD API Server] Running..")

	//handler
	http.HandleFunc("/create/opencsd", handler.CreateOpenCSD)
	http.HandleFunc("/create/mysql", handler.CreateMySQL)
	
	http.ListenAndServe(":40800", nil)
}
