package cli

import (
	"context"
	"os"
)

var (
	flagValueContextKey          = struct{ string }{"flags"}
	argsValueContextKey          = struct{ string }{"args"}
	parentCommandValueContextKey = struct{ string }{"parent"}
)

// FlagValues provides helpful accessors
type FlagValues map[string]interface{}

// ArgumentValues provides helpful accessors
type ArgumentValues map[string]string

// BaseContext is used to setup default values for the command
func BaseContext(ctx context.Context) context.Context {
	au := newAU(os.Stderr)

	return context.WithValue(ctx, colorContextKey, au)
}

// ArgumentsForContext returns the arguments for the given context
func ArgumentsForContext(ctx context.Context) ArgumentValues {
	if raw := ctx.Value(argsValueContextKey); raw != nil {
		if args, ok := raw.(ArgumentValues); ok {
			return args
		}
	}
	return make(ArgumentValues)
}

// FlagsForContext returns the flags for the given context
func FlagsForContext(ctx context.Context) FlagValues {
	if raw := ctx.Value(flagValueContextKey); raw != nil {
		if flags, ok := raw.(FlagValues); ok {
			return flags
		}
	}
	return make(FlagValues)
}

func parentCommandForContext(ctx context.Context) Command {
	if raw := ctx.Value(parentCommandValueContextKey); raw != nil {
		if cmd, ok := raw.(Command); ok {
			return cmd
		}
	}
	return nil
}
