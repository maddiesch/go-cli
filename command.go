package cli

import (
	"context"
	"errors"
	"fmt"
)

// Arguments provies the type for the Command Arguments
type Arguments []string

// Command provides the interface for executing commands
type Command interface {
	Name() string

	Execute(context.Context, Arguments) error
}

// Children provides the type for child commands
type Children []Command

// CommandConfig provides the required values for creating a new Command
type CommandConfig struct {
	Use       string // Provide a use statement for auto-generated help
	Short     string // Provide a short description for auto-generated help
	Long      string // Provide a long description for auto-generated help
	Flags     Flags
	Arguments []PositionalArgument
	Children  Children
	Run       Command
}

// NewCommand returns a new Command with Flags & Arguments
func NewCommand(config CommandConfig) Command {
	if len(config.Flags) == 0 {
		config.Flags = make(Flags)
	}

	if len(config.Arguments) == 0 {
		config.Arguments = make([]PositionalArgument, 0)
	}

	if config.Run != nil && len(config.Children) > 0 {
		panic(fmt.Errorf("unable to create a new command (%s) with both Run and Children", config.Use))
	}

	return &ContainerCommand{
		Use:       config.Use,
		Short:     config.Short,
		Long:      config.Long,
		Flags:     config.Flags,
		Arguments: config.Arguments,
		Children:  config.Children,
		Run:       config.Run,
	}
}

// ContainerCommand wraps a SubCommand with flags, descriptions, and other details.
type ContainerCommand struct {
	Use       string
	Short     string
	Long      string
	Flags     Flags
	Arguments []PositionalArgument
	Run       Command
	Children  Children
}

// Name provides the needed implementation for the Command interface
func (c *ContainerCommand) Name() string {
	return c.Use
}

// Execute provides the needed implementation for the Command interface
func (c *ContainerCommand) Execute(ctx context.Context, args Arguments) error {
	out, err := c.Flags.parse(c.Use, args)
	if err != nil {
		return err
	}

	if c.Run == nil && len(c.Children) == 0 {
		return fmt.Errorf("Command (%s) not properly configured! Missing either Run or Children", c.Use)
	}

	ctx = context.WithValue(ctx, flagValueContextKey, FlagValues(out.Values))
	ctx = context.WithValue(ctx, parentCommandValueContextKey, c)

	if c.Run != nil {
		return c.run(ctx, out.Arguments)
	}

	return c.child(ctx, out.Arguments)
}

func (c *ContainerCommand) run(ctx context.Context, args Arguments) error {
	posArgs := make(ArgumentValues)

	for _, pos := range c.Arguments {
		if len(args) == 0 {
			if pos.DefaultValue != nil {
				posArgs[pos.Name] = *(pos.DefaultValue)
			} else {
				return fmt.Errorf("missing required argument %s", pos.Name)
			}
		} else {
			posArgs[pos.Name] = args[0]
			args = args[1:]
		}
	}

	ctx = context.WithValue(ctx, argsValueContextKey, posArgs)

	return c.Run.Execute(ctx, args)
}

func (c *ContainerCommand) child(ctx context.Context, args Arguments) error {
	if len(args) == 0 {
		return errors.New("must specify a sub-command, try 'help'")
	}

	name := args[0]
	args = args[1:]

	for _, child := range c.Children {
		if child.Name() == name {
			return child.Execute(ctx, args)
		}
	}

	return errors.New("invalid sub-command name")
}

// CommandFunc provides a command that will execute the given function
func CommandFunc(fn func(context.Context, Arguments) error) Command {
	return &commandFunc{
		handler: fn,
	}
}

type commandFunc struct {
	handler func(context.Context, Arguments) error
}

func (c *commandFunc) Name() string {
	return "<unnamed>"
}

func (c *commandFunc) Execute(ctx context.Context, args Arguments) error {
	return c.handler(ctx, args)
}

// compile time type validation
var _ Command = new(commandFunc)
var _ Command = new(ContainerCommand)

// PositionalArgument contains information about a positional argument
type PositionalArgument struct {
	Name         string `validate:"required"`
	Description  string
	DefaultValue *string
}
