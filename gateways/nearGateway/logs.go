package nearGateway

import (
	"errors"
	"fmt"
	"marlinctl/util"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func LogsCommand() *cli.Command {
	return &cli.Command{
		Name:  "logs",
		Usage: "logs for the near gateway/bridge",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) error {
			chain := "near"
			program := chain + "_gateway"

			out, _ := exec.Command("sudo", "supervisorctl", "status", program).Output()
			fmt.Println(string(out))
			if strings.Contains(string(out), "no such process") {
				return errors.New("No program exists with name " + program)
			}

			return util.SupervisorTail(program)
		},
	}
}
