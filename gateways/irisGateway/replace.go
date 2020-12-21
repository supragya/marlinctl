package irisGateway

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func ReplaceCommand() *cli.Command {
	var bootstrapAddr, listenPortPeer, peerIP, peerPort, rpcPort string
	var version string

	return &cli.Command{
		Name:  "replace",
		Usage: "replace existing and create a new gateway",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "bootstrapaddr",
				Usage:       "--bootstrapaddr \"<IP1:PORT1>\"",
				Destination: &bootstrapAddr,
				Value:       "127.0.0.1:8002",
			},
			&cli.StringFlag{
				Name:        "listenportpeer",
				Usage:       "--listenportpeer \"PORT\"",
				Destination: &listenPortPeer,
				Value:       "59001",
			},
			&cli.StringFlag{
				Name:        "peerip",
				Usage:       "--peerip \"IP\"",
				Destination: &peerIP,
				Value:       "127.0.0.1",
			},
			&cli.StringFlag{
				Name:        "peerport",
				Usage:       "--peerport \"PORT\"",
				Destination: &peerPort,
				Value:       "26656",
			},
			&cli.StringFlag{
				Name:        "rpcport",
				Usage:       "--rpcport \"PORT\"",
				Destination: &rpcPort,
				Value:       "26657",
			},
			&cli.StringFlag{
				Name:        "version",
				Usage:       "--version <NUMBER>",
				Value:       "latest",
				Destination: &version,
			},
		},
		Action: func(c *cli.Context) error {
			chain := "iris"
			program := chain + "_gateway"
			out, _ := exec.Command("sudo", "supervisorctl", "status", program).Output()
			if strings.Contains(string(out), "no such process") {
				// Doesn't Already exists
				return errors.New("Doesn't already exists, use create command instead")
			}

			if err := Destroy(); err != nil {
				return err
			}
			return Create(bootstrapAddr, listenPortPeer, peerIP, peerPort, rpcPort, version)
		},
	}
}
