package nearGateway

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/urfave/cli/v2"

	"marlinctl/util"
)

func CreateCommand() *cli.Command {
	var bootstrap_addr string
	var version string

	return &cli.Command{
		Name:  "create",
		Usage: "create a new gateway",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "bootstrap-addr",
				Usage:       "--bootstrap-addr \"<IP1:PORT1>\"",
				Destination: &bootstrap_addr,
				Value:       "127.0.0.1:8002",
			},
			&cli.StringFlag{
				Name:        "version",
				Usage:       "--version <NUMBER>",
				Value:       "latest",
				Destination: &version,
			},
		},
		Action: func(c *cli.Context) error {
			return Create(bootstrap_addr, version)
		},
	}
}

func Create(bootstrap_addr string, version string) error {
	chain := "near"
	fmt.Println(bootstrap_addr)
	program := chain + "_gateway"

	out, _ := exec.Command("sudo", "supervisorctl", "status", program).Output()
	if !strings.Contains(string(out), "no such process") {
		// Already exists
		return errors.New("Already exists")
	}

	// User details
	usr, err := util.GetUser()
	if err != nil {
		return err
	}

	// Version
	if version == "latest" {
		fmt.Println(program, "fetching latest binaries...")
		latestVersion, err := util.FetchLatestVersion(program)
		if err != nil {
			return err
		}
		version = latestVersion
		fmt.Println(program, "latest binary version: ", latestVersion)
	}

	// gateway executable
	err = util.Fetch("https://storage.googleapis.com/marlin-artifacts/bin/"+program+"-"+version+"-"+runtime.GOOS+"-"+runtime.GOARCH, usr.HomeDir+"/.marlin/ctl/bin/"+program+"-"+version, usr.Username, true, false)
	if err != nil {
		return err
	}

	// gateway config
	err = util.Fetch(
		"https://storage.googleapis.com/marlin-artifacts/configs/"+program+"-"+version+".conf", usr.HomeDir+"/.marlin/ctl/configs/"+program+"-"+version+".conf", usr.Username, false, false)
	if err != nil {
		return err
	}

	err = util.TemplatePlace(
		usr.HomeDir+"/.marlin/ctl/configs/"+program+"-"+version+".conf",
		"/etc/supervisor/conf.d/"+program+".conf",
		struct {
			Program, User, UserHome string
			BootstrapAddr           string
			Version                 string
		}{
			program, usr.Username, usr.HomeDir,
			bootstrap_addr,
			version,
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

	output, _ := exec.Command("supervisorctl", "status").Output()
	fmt.Print(string(output))

	return nil
}
