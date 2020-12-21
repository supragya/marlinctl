package relay

import (
	"errors"
	"fmt"
	"marlinctl/util"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func LogsCommand() *cli.Command {
	var chain string
	return &cli.Command{
		Name:  "logs",
		Usage: "logs for the relay",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "chain",
				Usage:       "--chain \"<CHAIN>\"",
				Destination: &chain,
				Required:    true,
			},
		},
		Action: func(c *cli.Context) error {
			program := chain + "_relay"
			out, _ := exec.Command("sudo", "supervisorctl", "status", program).Output()
			fmt.Println(string(out))
			if strings.Contains(string(out), "no such process") {
				return errors.New("No program exists with name " + program)
			}

			return util.SupervisorTail(program)
		},
	}
}
