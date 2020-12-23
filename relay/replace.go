package relay

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func ReplaceCommand() *cli.Command {
	var chain, discovery_addrs, heartbeat_addrs, datadir string
	var discovery_port, pubsub_port uint
	var address, name string
	var version string
	var abci_version string
	var sync_mode string

	return &cli.Command{
		Name:  "replace",
		Usage: "destroy the existing and create a new relay",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "chain",
				Usage:       "--chain \"<CHAIN>\"",
				Destination: &chain,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "discovery-addrs",
				Usage:       "--discovery-addrs \"<IP1:PORT1>,<IP2:PORT2>,...\"",
				Destination: &discovery_addrs,
				Value:       "127.0.0.1:8002",
			},
			&cli.StringFlag{
				Name:        "heartbeat-addrs",
				Usage:       "--heartbeat-addrs \"<IP1:PORT1>,<IP2:PORT2>,...\"",
				Destination: &heartbeat_addrs,
				Value:       "127.0.0.1:8003",
			},
			&cli.StringFlag{
				Name:        "datadir",
				Usage:       "--datadir \"/path/to/datadir\"",
				Destination: &datadir,
				Value:       "~/.ethereum/",
			},
			&cli.UintFlag{
				Name:        "discovery-port",
				Usage:       "--discovery-port <PORT>",
				Destination: &discovery_port,
			},
			&cli.UintFlag{
				Name:        "pubsub-port",
				Usage:       "--pubsub-port <PORT>",
				Destination: &pubsub_port,
			},
			&cli.StringFlag{
				Name:        "address",
				Usage:       "--address \"0x...\"",
				Destination: &address,
			},
			&cli.StringFlag{
				Name:        "name",
				Usage:       "--name \"<NAME>\"",
				Destination: &name,
			},
			&cli.StringFlag{
				Name:        "version",
				Usage:       "--version <NUMBER>",
				Value:       "latest",
				Destination: &version,
			},
			&cli.StringFlag{
				Name:        "abci-version",
				Usage:       "--abci-version <NUMBER>",
				Value:       "latest",
				Destination: &abci_version,
			},
			&cli.StringFlag{
				Name:        "sync-mode",
				Usage:       "--sync-mode <MODE>",
				Value:       "light",
				Destination: &sync_mode,
			},
		},
		Action: func(c *cli.Context) error {
			program := chain + "_relay"
			out, _ := exec.Command("sudo", "supervisorctl", "status", program).Output()
			if strings.Contains(string(out), "no such process") {
				// Doesn't exists
				return errors.New("Doesn't already exists, use create command instead")
			}

			if err := Destroy(chain); err != nil {
				return err
			}
			return Create(chain, discovery_addrs, heartbeat_addrs, datadir, discovery_port, pubsub_port, address, name, version, abci_version, sync_mode)
		},
	}
}
