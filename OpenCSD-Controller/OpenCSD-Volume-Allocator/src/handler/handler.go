package handler

import (
	"net/http"
	"fmt"
	"bufio"
	"os/exec"
)

func StorageNodeInfo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[OpenCSD Volume Allocator] StorageNodeInfo Completed\n"))
}

func StorageVolumeInfo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[OpenCSD Volume Allocator] StorageVolumeInfo Completed\n"))
}

func AllocateVolume(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[OpenCSD Volume Allocator] AllocateVolume Completed\n"))
}

func DeallocateVolume(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[OpenCSD Volume Allocator] DeallocateVolume Completed\n"))
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