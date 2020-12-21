package nearGateway

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func ReplaceCommand() *cli.Command {
	var bootstrap_addr string
	var version string

	return &cli.Command{
		Name:  "replace",
		Usage: "destroy existing and create a new gateway",
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
			chain := "near"
			program := chain + "_gateway"
			out, _ := exec.Command("sudo", "supervisorctl", "status", program).Output()
			if strings.Contains(string(out), "no such process") {
				// Doesn't Already exists
				return errors.New("Doesn't already exists, use create command instead")
			}

			if err := Destroy(); err != nil && err.Error() != "Not found" {
				return err
			}

			return Create(bootstrap_addr, version)
		},
	}
}
