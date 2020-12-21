package irisGateway

import (
	"errors"
	"fmt"
	"marlinctl/util"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func LogsCommand() *cli.Command {
	var bridge bool
	return &cli.Command{
		Name:  "logs",
		Usage: "logs for the iris gateway/bridge",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "bridge",
				Usage:       "--bridge=true",
				Destination: &bridge,
				Value:       false,
			},
		},
		Action: func(c *cli.Context) error {
			chain := "iris"
			var program string
			if bridge {
				program = chain + "_bridge"
			} else {
				program = chain + "_gateway"
			}

			out, _ := exec.Command("sudo", "supervisorctl", "status", program).Output()
			fmt.Println(string(out))
			if strings.Contains(string(out), "no such process") {
				return errors.New("No program exists with name " + program)
			}

			return util.SupervisorTail(program)
		},
	}
}
