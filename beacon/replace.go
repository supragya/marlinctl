package beacon

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func ReplaceCommand() *cli.Command {
	var discovery_addr string
	var heartbeat_addr string
	var beacon_addr string
	var program string
	var version string

	return &cli.Command{
		Name:  "replace",
		Usage: "destroy existing and then create a new beacon",
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
			out, _ := exec.Command("sudo", "supervisorctl", "status", program).Output()
			if strings.Contains(string(out), "no such process") {
				// Doesn't exists
				return errors.New("Doesn't already exists, use create command instead")
			}

			if err := Destroy(program); err != nil {
				return err
			}
			return Create(discovery_addr, heartbeat_addr, beacon_addr, program, version)
		},
	}
}
