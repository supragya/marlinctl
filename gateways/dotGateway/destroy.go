package dotGateway

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func DestroyCommand() *cli.Command {
	return &cli.Command{
		Name:  "destroy",
		Usage: "destroy the gateway",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) error {
			return Destroy()
		},
	}
}

func Destroy() error {
	chain := "dot"
	{
		program := chain + "_gateway"

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

		fmt.Println(program + " destroyed")
	}

	{
		program := chain + "_bridge"

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

		fmt.Println(program + " destroyed")
	}
	_, err := exec.Command("supervisorctl", "reread").Output()
	if err != nil {
		return err
	}

	return nil
}
