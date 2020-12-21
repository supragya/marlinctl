package relay

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func DestroyCommand() *cli.Command {
	var chain string

	return &cli.Command{
		Name:  "destroy",
		Usage: "destroy the relay",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "chain",
				Usage:       "--chain \"<CHAIN>\"",
				Destination: &chain,
				Required:    true,
			},
		},
		Action: func(c *cli.Context) error {
			return Destroy(chain)
		},
	}
}

func Destroy(chain string) error {
	program := chain + "_relay"

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

	// Destroy abci
	if abci, found := abciMap[chain]; found {
		err := abci.Destroy()
		if err != nil {
			return err
		}
	} else {
		return errors.New("Unrecognized chain")
	}

	fmt.Println(program + " destroyed")

	return nil
}
