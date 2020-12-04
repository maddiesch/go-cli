package cli

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/logrusorgru/aurora/v3"
	"github.com/mattn/go-isatty"
)

// Helper provides the interface for generating a Help
type Helper interface {
	Help(context.Context, int, *strings.Builder)
}

var _ Helper = new(ContainerCommand)

// Help generates the help text for the given interface
func Help(ctx context.Context, helper interface{}) string {
	if helper, ok := helper.(Helper); ok {
		var builder strings.Builder

		helper.Help(ctx, 0, &builder)

		return builder.String()
	}

	return fmt.Sprintf("%v", helper)
}

// Help provides the implementation for cli.Helper
func (c *ContainerCommand) Help(ctx context.Context, depth int, buf *strings.Builder) {
	switch depth {
	case 0:
		c.w0(ctx, buf)
	case 1:
		c.w1(ctx, buf)
	case 2:
		c.w2(ctx, buf)
	}
}

func (c *ContainerCommand) w0(ctx context.Context, buf *strings.Builder) {
	au := color(ctx)

	buf.WriteString(au.Bold(c.Use).String())

	if len(c.Short) > 0 {
		if len(c.Use) > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(au.Faint(c.Short).String())
	}

	buf.WriteRune('\n')

	if len(c.Long) > 0 {
		buf.WriteString(c.Long)
		buf.WriteRune('\n')
	}

	if len(c.Children) > 0 {
		buf.WriteRune('\n')
		buf.WriteString(au.Underline(au.Bold(au.Magenta("commands"))).String())
		buf.WriteRune('\n')

		for name, child := range c.Children {
			if helper, ok := child.(Helper); ok {
				if container, ok := child.(*ContainerCommand); ok {
					if container.Use == "" {
						container.Use = name
					}
				}
				helper.Help(ctx, 1, buf)
				buf.WriteRune('\n')
			}
		}
	}

	buf.WriteRune('\n')
}

func (c *ContainerCommand) w1(ctx context.Context, buf *strings.Builder) {
	au := color(ctx)

	buf.WriteString("")
	buf.WriteString(c.Use)
	buf.WriteString(" -- ")
	buf.WriteString(c.Short)
	if d := c.Long; len(d) > 0 {
		buf.WriteString("\n\n")
		buf.WriteString(d)
		buf.WriteString("\n")
	}
	if len(c.Flags) > 0 {
		buf.WriteString(au.Gray(8, "\n  Flags:\n").String())

		for _, name := range sortedFlagNames(c.Flags) {
			buf.WriteString(flagHelp(4, au, name, c.Flags[name]) + "\n")
		}
	}
}

func (c *ContainerCommand) w2(ctx context.Context, buf *strings.Builder) {
}

func flagHelp(depth int, au aurora.Aurora, name string, f Flag) string {
	var s strings.Builder
	for i := 0; i < depth; i++ {
		s.WriteRune(' ')
	}

	tName := strings.Replace(strings.Replace(reflect.TypeOf(f).String(), "cli.", "", 1), "Flag", "", 1)

	s.WriteString(au.Cyan("-" + name).String())
	s.WriteRune(' ')
	s.WriteString(au.Gray(18, fmt.Sprintf("[%s]", tName)).String())
	s.WriteString(" - ")
	s.WriteString(f.Short())

	if d := f.DefaultValue(); d != nil {
		s.WriteString(aurora.Gray(14, fmt.Sprintf(" (default '%v')", d)).String())
	}

	if e, ok := f.(EnumFlag); ok {
		s.WriteRune('\n')
		for i := 0; i < depth+len(name)+2; i++ {
			s.WriteRune(' ')
		}
		s.WriteString(au.Gray(12, "Possible Values: ").String())

		s.WriteString(strings.Join(e.PossibleValues, ", "))
	}

	return s.String()
}

func color(ctx context.Context) aurora.Aurora {
	if untyped := ctx.Value(colorContextKey); untyped != nil {
		if au, ok := untyped.(aurora.Aurora); ok {
			return au
		}
	}

	return aurora.NewAurora(false)
}

var colorContextKey = &struct{ string }{"_color"}

// SelfHelp returns a new command that will print help for the current container
func SelfHelp() Command {
	return &helpCommand{}
}

type helpCommand struct {
}

func (h *helpCommand) Execute(ctx context.Context, args Arguments) error {
	out := os.Stdout

	au := newAU(out)

	parent := parentCommandForContext(ctx)

	ctx = context.WithValue(ctx, colorContextKey, au)

	fmt.Fprintf(out, "%s\n", Help(ctx, parent))

	return nil
}

var _ Command = new(helpCommand)

func sortedFlagNames(f Flags) []string {
	out := make([]string, 0, len(f))
	for k := range f {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func newAU(out *os.File) aurora.Aurora {
	return aurora.NewAurora(isatty.IsTerminal(out.Fd()))
}
