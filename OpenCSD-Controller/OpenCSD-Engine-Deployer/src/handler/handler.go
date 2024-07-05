package handler

import (
	"net/http"
	"fmt"
	"bufio"
	"os/exec"
)

func CreateQueryEngine(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[OpenCSD Engine Deployer] CreateQueryEngine Completed\n"))
}

func CreateStorageEngine(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[OpenCSD Engine Deployer] CreateStorageEngine Completed\n"))
}

func CreateValidator(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[OpenCSD Engine Deployer] CreateValidator Completed\n"))
}

func Info(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[OpenCSD Engine Deployer] Info Completed\n"))
}

func DeleteQueryEngine(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[OpenCSD Engine Deployer] DeleteQueryEngine Completed\n"))
}

func DeleteStorageEngine(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[OpenCSD Engine Deployer] DeleteStorageEngine Completed\n"))
}

func DeleteValidator(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[OpenCSD Engine Deployer] DeleteValidator Completed\n"))
}

func CmdExec(cmdStr string) error{
	cmd := exec.Command("bash", "-c", cmdStr)
	stdoutReader, _ := cmd.StdoutPipe()
	stdoutScanner := bufio.NewScanner(stdoutReader)
	go func() {
		for stdoutScanner.Scan() {
			fmt.Println(stdoutScanner.Text())
		}
	}()
	stderrReader, _ := cmd.StderrPipe()
	stderrScanner := bufio.NewScanner(stderrReader)
	go func() {
		for stderrScanner.Scan() {
			fmt.Println(stderrScanner.Text())
		}
	}()
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error : %v \n", err)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Error: %v \n", err)
	}

	return nil
}