package cli

import (
	"log"
	"os"
)

// ErrLog provides all logging needed
var ErrLog = log.New(os.Stderr, "", 0)
