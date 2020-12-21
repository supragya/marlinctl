package relay

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/urfave/cli/v2"

	"marlinctl/util"
)

func CreateCommand() *cli.Command {
	var chain, discovery_addrs, heartbeat_addrs, datadir string
	var discovery_port, pubsub_port uint
	var address, name string
	var version string
	var abci_version string
	var sync_mode string

	return &cli.Command{
		Name:  "create",
		Usage: "create a new relay",
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
			return Create(chain, discovery_addrs, heartbeat_addrs, datadir, discovery_port, pubsub_port, address, name, version, abci_version, sync_mode)
		},
	}
}

func Create(chain, discovery_addrs, heartbeat_addrs, datadir string,
	discovery_port, pubsub_port uint,
	address, name string,
	version string,
	abci_version string,
	sync_mode string) error {

	program := chain + "_relay"

	out, _ := exec.Command("sudo", "supervisorctl", "status", program).Output()
	if !strings.Contains(string(out), "no such process") {
		// Already exists
		return errors.New("Already exists")
	}

	// Set up abci first
	if abci, found := abciMap[chain]; found {
		err := abci.Create(datadir, abci_version, sync_mode)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Unrecognized chain")
	}

	// User details
	usr, err := util.GetUser()
	if err != nil {
		return err
	}

	// Tilde expansion
	if datadir == "~" {
		datadir = usr.HomeDir
	} else if strings.HasPrefix(datadir, "~/") {
		datadir = filepath.Join(usr.HomeDir, datadir[2:])
	}

	// Version
	if version == "latest" {
		fmt.Println(program, "fetching latest binaries...")
		latestVersion, err := util.FetchLatestVersion(program)
		if err != nil {
			return err
		}
		version = latestVersion
		fmt.Println(program, "latest binary version: ", latestVersion)
	}

	// relay executable
	err = util.Fetch("https://storage.googleapis.com/marlin-artifacts/bin/"+program+"-"+version+"-"+runtime.GOOS+"-"+runtime.GOARCH, usr.HomeDir+"/.marlin/ctl/bin/"+program+"-"+version, usr.Username, true, false)
	if err != nil {
		return err
	}

	// relay config
	err = util.Fetch("https://storage.googleapis.com/marlin-artifacts/configs/"+program+"-"+version+".conf", usr.HomeDir+"/.marlin/ctl/configs/"+program+"-"+version+".conf", usr.Username, false, false)
	if err != nil {
		return err
	}

	err = util.TemplatePlace(
		usr.HomeDir+"/.marlin/ctl/configs/"+program+"-"+version+".conf",
		"/etc/supervisor/conf.d/"+program+".conf",
		struct {
			Program, User, UserHome                 string
			DiscoveryAddrs, HeartbeatAddrs, Datadir string
			DiscoveryPort, PubsubPort               uint
			Address, Name                           string
			Version                                 string
		}{
			program, usr.Username, usr.HomeDir,
			discovery_addrs, heartbeat_addrs, datadir,
			discovery_port, pubsub_port,
			address, name,
			version,
		},
	)
	if err != nil {
		return err
	}

	_, err = exec.Command("supervisorctl", "reread").Output()
	if err != nil {
		return err
	}

	_, err = exec.Command("supervisorctl", "add", program).Output()
	if err != nil {
		return err
	}

	output, _ := exec.Command("supervisorctl", "status").Output()
	fmt.Print(string(output))

	return nil
}
