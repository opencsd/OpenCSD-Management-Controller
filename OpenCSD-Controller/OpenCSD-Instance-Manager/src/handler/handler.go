package handler

import (
	"net/http"
	"fmt"
	"bufio"
	"os/exec"
)

func MappingInstanceStorage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[OpenCSD Instance Mananger] MappingInstanceStorage Completed\n"))
}

func UnMappingInstanceStorage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[OpenCSD Instance Mananger] UnMappingInstanceStorage Completed\n"))
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