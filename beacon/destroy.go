package beacon

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func DestroyCommand() *cli.Command {
	var program string

	return &cli.Command{
		Name:  "destroy",
		Usage: "destroy the beacon",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "program",
				Usage:       "--program <NAME>",
				Value:       "beacon",
				Destination: &program,
			},
		},
		Action: func(c *cli.Context) error {
			return Destroy(program)
		},
	}
}

func Destroy(program string) error {
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

	fmt.Println("beacon destroyed")

	return nil
}
