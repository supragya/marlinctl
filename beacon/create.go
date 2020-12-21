package beacon

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
	var discovery_addr string
	var heartbeat_addr string
	var beacon_addr string
	var program string
	var version string

	return &cli.Command{
		Name:  "create",
		Usage: "create a new beacon",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "program",
				Usage:       "--program <NAME>",
				Value:       "beacon",
				Destination: &program,
			},
			&cli.StringFlag{
				Name:        "discovery-addr",
				Usage:       "--discovery-addr <IP:PORT>",
				DefaultText: "127.0.0.1:8002",
				Destination: &discovery_addr,
			},
			&cli.StringFlag{
				Name:        "heartbeat-addr",
				Usage:       "--heartbeat-addr <IP:PORT>",
				DefaultText: "127.0.0.1:8003",
				Destination: &heartbeat_addr,
			},
			&cli.StringFlag{
				Name:        "bootstrap-addr",
				Usage:       "--bootstrap-addr <IP:PORT>",
				Aliases:     []string{"beacon-addr"},
				Destination: &beacon_addr,
			},
			&cli.StringFlag{
				Name:        "version",
				Usage:       "--version <NUMBER>",
				Value:       "latest",
				Destination: &version,
			},
		},
		Action: func(c *cli.Context) error {
			return Create(discovery_addr, heartbeat_addr, beacon_addr, program, version)
		},
	}
}

func Create(discovery_addr string,
	heartbeat_addr string,
	beacon_addr string,
	program string,
	version string) error {

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

	if version == "latest" {
		fmt.Println("beacon fetching latest binaries...")
		latestVersion, err := util.FetchLatestVersion("beacon")
		if err != nil {
			return err
		}
		version = latestVersion
		fmt.Println("beacon latest binary version: ", latestVersion)
	}
	// Beacon executable
	err = util.Fetch("https://storage.googleapis.com/marlin-artifacts/bin/beacon-"+version+"-"+runtime.GOOS+"-"+runtime.GOARCH, usr.HomeDir+"/.marlin/ctl/bin/beacon-"+version, usr.Username, true, false)
	if err != nil {
		return err
	}

	// Beacon config
	err = util.Fetch("https://storage.googleapis.com/marlin-artifacts/configs/beacon-"+version+".conf", usr.HomeDir+"/.marlin/ctl/configs/beacon-"+version+".conf", usr.Username, false, false)
	if err != nil {
		return err
	}

	err = util.TemplatePlace(
		usr.HomeDir+"/.marlin/ctl/configs/beacon-"+version+".conf",
		"/etc/supervisor/conf.d/"+program+".conf",
		struct {
			Program, User, UserHome                  string
			DiscoveryAddr, HeartbeatAddr, BeaconAddr string
			Version                                  string
		}{
			program, usr.Username, usr.HomeDir, discovery_addr, heartbeat_addr, beacon_addr, version,
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

	output, err := exec.Command("supervisorctl", "status").Output()
	fmt.Print(string(output))

	return nil
}
