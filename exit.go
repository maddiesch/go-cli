package cli

import (
	"fmt"
	"strconv"
)

// ExitError is an error that will exit with the given exit code
type ExitError int

func (e ExitError) Error() string {
	return fmt.Sprintf("0x%s", strconv.FormatInt(int64(e), 16))
}

var (
	// ExitErrorDead is returned when the process should exit
	ExitErrorDead ExitError = 0xdecea5ed
)

// compile time interface checks
var _ error = ExitErrorDead
