package cli

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Output writes a single line to the standard output
func Output(value string) {
	fmt.Fprintln(os.Stderr, value)
}
