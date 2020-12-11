package cli

import (
	"flag"
	"fmt"
	"strconv"
)

// Flags is a collection of flags
type Flags map[string]Flag

// CopyAndMerge merges the flags into a single group of flags and returns the
// new group
func (f Flags) CopyAndMerge(other Flags) Flags {
	copy := make(Flags)
	for k, v := range f {
		copy[k] = v
	}
	for k, v := range other {
		copy[k] = v
	}
	return copy
}

// Flag contains the information for
type Flag interface {
	Short() string                           // Return help text describing the flag
	DefaultValue() interface{}               // Return a default value for the flag
	CoerceValue(string) (interface{}, error) // Returns a value for the flag that has been converted into the specified value for the flag.
}

// IntFlag represents a flag that contains an int value
type IntFlag struct {
	Default     int64
	Description string
}

// Short provides the implementation for Flag
func (f IntFlag) Short() string {
	return f.Description
}

// DefaultValue provides the implementation for Flag
func (f IntFlag) DefaultValue() interface{} {
	return f.Default
}

// CoerceValue provides the implementation for Flag
func (f IntFlag) CoerceValue(in string) (interface{}, error) {
	if in == "" {
		return f.DefaultValue(), nil
	}
	return strconv.ParseInt(in, 10, 64)
}

// StringFlag represents a flag that contains a string value
type StringFlag struct {
	Default     string
	Description string
}

// Short provides the implementation for Flag
func (f StringFlag) Short() string {
	return f.Description
}

// DefaultValue provides the implementation for Flag
func (f StringFlag) DefaultValue() interface{} {
	return f.Default
}

// CoerceValue provides the implementation for Flag
func (f StringFlag) CoerceValue(in string) (interface{}, error) {
	return in, nil
}

// EnumFlag represents a flag that has a finite number of possible values
type EnumFlag struct {
	Description    string
	Default        string
	PossibleValues []string
}

// Short provides the implementation for Flag
func (f EnumFlag) Short() string {
	return f.Description
}

// DefaultValue provides the implementation for Flag
func (f EnumFlag) DefaultValue() interface{} {
	return f.Default
}

// CoerceValue provides the implementation for Flag
func (f EnumFlag) CoerceValue(in string) (interface{}, error) {
	if in == "" {
		return f.Default, nil
	}

	for _, v := range f.PossibleValues {
		if v == in {
			return in, nil
		}
	}

	return nil, fmt.Errorf("Invalid value '%s'", in)
}

var _ Flag = new(IntFlag)
var _ Flag = new(EnumFlag)
var _ Flag = new(StringFlag)

type parseOutput struct {
	Values    map[string]interface{}
	Arguments Arguments
}

func (f Flags) parse(name string, args Arguments) (*parseOutput, error) {
	if len(f) == 0 {
		return &parseOutput{Arguments: args}, nil
	}

	fSet := flag.NewFlagSet(name, flag.ContinueOnError)

	entries := make(map[string]*string)

	for name, flag := range f {
		str := new(string)

		fSet.StringVar(str, name, "", flag.Short())

		entries[name] = str
	}

	if err := fSet.Parse(args); err != nil {
		return nil, err
	}

	out := &parseOutput{
		Arguments: fSet.Args(),
		Values:    make(map[string]interface{}),
	}

	for name, flag := range f {
		value := entries[name]
		if value == nil {
			out.Values[name] = flag.DefaultValue()
		} else {
			val, err := flag.CoerceValue(*value)
			if err != nil {
				return nil, err
			}
			out.Values[name] = val
		}
	}

	return out, nil
}

// FlagValues provides helpful accessors
type FlagValues map[string]interface{}

// StringValue returns the string value for the specified key
func (v FlagValues) StringValue(name string) string {
	if raw, ok := v[name]; ok {
		if str, ok := raw.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", raw)
	}
	return ""
}

// IntValue returns the int value for the specified key
func (v FlagValues) IntValue(name string) int64 {
	if raw, ok := v[name]; ok {
		if i, ok := raw.(int64); ok {
			return i
		}
	}
	return 0
}
