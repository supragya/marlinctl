package irisGateway

import (
	"github.com/urfave/cli/v2"
)

var IrisGateway = cli.Command{
	Name:  "iris",
	Usage: "create, start or stop iris gateway",
	Subcommands: []*cli.Command{
		CreateCommand(),
		DestroyCommand(),
		ReplaceCommand(),
		LogsCommand(),
	},
}
