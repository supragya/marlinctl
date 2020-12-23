package nearGateway

import (
	"github.com/urfave/cli/v2"
)

var NearGateway = cli.Command{
	Name:  "near",
	Usage: "create, start or stop near gateway",
	Subcommands: []*cli.Command{
		CreateCommand(),
		DestroyCommand(),
		ReplaceCommand(),
		LogsCommand(),
	},
}
