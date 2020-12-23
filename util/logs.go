package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func SupervisorTail(program string) error {
	cmd := exec.Command("sudo", "supervisorctl", "tail", "-f", program)

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for logs", err)
		return err
	}
	go readPipe(cmdReader)

	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StderrPipe for logs", err)
		return err
	}
	go readPipe(cmdErrReader)

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting logs", err)
		return err
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for logs", err)
		return err
	}
	return nil
}

func readPipe(reader io.ReadCloser) {
	scannerErr := bufio.NewScanner(reader)
	for scannerErr.Scan() {
		fmt.Printf("%s\n", scannerErr.Text())
	}
}
