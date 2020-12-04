package cli

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/go-playground/validator/v10"
)

// FatalError provies a standard "chokepoint" for throwing all fatal errors
var FatalError = func(ctx context.Context, err error) {
	if invalid, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range invalid {
			ErrLog.Printf("Validation Error: %s", fieldErr.Error())
		}
	}

	au := color(ctx)

	runtime.KeepAlive(au)

	fmt.Fprintf(os.Stderr, "⚠️  %s %v\n", au.Red("FATAL ERROR:").String(), err)

	os.Exit(1)
}
