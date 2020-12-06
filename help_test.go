package cli_test

import (
	"context"
	"runtime"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/maddiesch/go-cli"
)

func TestHelp(t *testing.T) {
	cmd := cli.NewCommand(cli.CommandConfig{
		Use:   "cli",
		Short: "the root cli app",
		Long:  "cli provides a simple interface for creating robust, well-documented, and fast CLI applications.",
		Children: cli.Children{
			cli.NewCommand(cli.CommandConfig{
				Use:   "foo",
				Short: "called when you need to foo",
				Run: cli.CommandFunc(func(context.Context, cli.Arguments) error {
					return nil
				}),
			}),
		},
	})

	help := cli.Help(context.Background(), cmd)

	spew.Dump(help)

	runtime.KeepAlive(help)
}
