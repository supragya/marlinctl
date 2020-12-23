package relay

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"marlinctl/util"
)

type EthAbci struct{}

func (abci *EthAbci) Create(datadir string, version string, syncmode string) error {
	// User details
	usr, err := util.GetUser()
	if err != nil {
		return err
	}

	program := "geth"

	if version == "latest" {
		fmt.Println(program, "fetching latest binaries...")
		latestVersion, err := util.FetchLatestVersion(program)
		if err != nil {
			return err
		}
		version = latestVersion
		fmt.Println(program, "latest binary version: ", latestVersion)
	}
	// geth executable
	err = util.Fetch("https://storage.googleapis.com/marlin-artifacts/bin/"+program+"-"+version+"-"+runtime.GOOS+"-"+runtime.GOARCH, usr.HomeDir+"/.marlin/ctl/bin/"+program+"-"+version, usr.Username, true, false)
	if err != nil {
		return err
	}

	// geth config
	err = util.Fetch("https://storage.googleapis.com/marlin-artifacts/configs/"+program+"-"+version+".conf", usr.HomeDir+"/.marlin/ctl/configs/"+program+"-"+version+".conf", usr.Username, false, false)
	if err != nil {
		return err
	}

	err = util.TemplatePlace(
		usr.HomeDir+"/.marlin/ctl/configs/"+program+"-"+version+".conf",
		"/etc/supervisor/conf.d/"+program+".conf",
		struct {
			Program, User, UserHome, DataDir, Version, SyncMode string
		}{
			program, usr.Username, usr.HomeDir, datadir, version, syncmode,
		},
	)
	if err != nil {
		return err
	}

	_, err = exec.Command("supervisorctl", "reread").Output()
	if err != nil {
		return err
	}

	_, err = exec.Command("supervisorctl", "add", program).Output()
	if err != nil {
		return err
	}

	return nil
}

func (abci *EthAbci) Destroy() error {
	program := "geth"

	out, _ := exec.Command("sudo", "supervisorctl", "status", program).Output()
	if strings.Contains(string(out), "no such process") {
		return errors.New("Not found")
	}

	_, err := exec.Command("sudo", "supervisorctl", "stop", program).Output()
	if err != nil {
		return err
	}

	_, err = exec.Command("sudo", "supervisorctl", "remove", program).Output()
	if err != nil {
		return err
	}

	_, err = exec.Command("sudo", "rm", "/etc/supervisor/conf.d/"+program+".conf").Output()
	if err != nil {
		return err
	}

	_, err = exec.Command("supervisorctl", "reread").Output()
	if err != nil {
		return err
	}

	fmt.Println(program + " destroyed")

	return nil
}
