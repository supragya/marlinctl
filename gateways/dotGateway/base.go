package dotGateway

import (
	"github.com/urfave/cli/v2"
)

var DotGateway = cli.Command{
	Name:  "dot",
	Usage: "create, start or stop dot gateway",
	Subcommands: []*cli.Command{
		CreateCommand(),
		DestroyCommand(),
		ReplaceCommand(),
		LogsCommand(),
	},
}
