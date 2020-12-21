package beacon

import (
	"errors"
	"fmt"
	"marlinctl/util"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func LogsCommand() *cli.Command {
	var program string
	return &cli.Command{
		Name:  "logs",
		Usage: "logs for the beacon",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "program",
				Usage:       "--program <NAME>",
				Value:       "beacon",
				Destination: &program,
			},
		},
		Action: func(c *cli.Context) error {
			out, _ := exec.Command("sudo", "supervisorctl", "status", program).Output()
			fmt.Println(string(out))
			if strings.Contains(string(out), "no such process") {
				return errors.New("No program exists with name " + program)
			}

			return util.SupervisorTail(program)
		},
	}
}
