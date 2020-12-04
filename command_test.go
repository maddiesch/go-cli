package cli_test

import (
	"context"
	"testing"

	"github.com/maddiesch/go-cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandChildren(t *testing.T) {
	var voice, input string
	var count int64

	app := cli.NewCommand(cli.CommandConfig{
		Use: "app",
		Children: cli.Children{
			"say": cli.NewCommand(cli.CommandConfig{
				Use:   "say",
				Short: "prints the given input to the standard output",
				Run: cli.CommandFunc(func(ctx context.Context, _ cli.Arguments) error {
					args := cli.ArgumentsForContext(ctx)
					flags := cli.FlagsForContext(ctx)

					voice = flags["voice"].(string)
					count = flags["count"].(int64)
					input = args["input"]

					return nil
				}),
				Flags: cli.Flags{
					"voice": cli.StringFlag{Default: "default"},
					"count": cli.IntFlag{Default: 1},
				},
				Arguments: []cli.PositionalArgument{
					{Name: "input"},
				},
			}),
		},
	})

	err := app.Execute(context.Background(), cli.Arguments{"say", "-voice", "test", `Testing the input`})

	require.NoError(t, err)

	assert.Equal(t, "test", voice)
	assert.Equal(t, "Testing the input", input)
	assert.Equal(t, int64(1), count)
}
