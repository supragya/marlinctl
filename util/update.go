package util

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/inconshreveable/go-update"
)

var currentVersion = "1.3.0"

func CheckAndUpdate() bool {
	latestVersion, err := FetchLatestVersion("marlinctl")
	if err != nil || latestVersion == currentVersion {
		return false
	}

	fmt.Println("New version found! Updating...")
	err = Update("https://storage.googleapis.com/marlin-artifacts/bin/marlinctl-" + runtime.GOOS + "-" + runtime.GOARCH)

	if err != nil {
		fmt.Println("Error while updating. Falling back... ", err)
		return false
	}
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Updated to version", latestVersion, "Reloading...")
	err = cmd.Start()
	if err != nil {
		fmt.Println("Unable to reload updated cli. Falling back... ", err)
		return false
	}
	err = cmd.Wait() // waiting for the child process to complete
	if err != nil {
		fmt.Println("error occured: ", err)
	}
	return true
}

func Update(url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Fetch error")
	}
	return update.Apply(resp.Body, update.Options{})
}
