package apperrors

import (
	"fmt"
	"os"
)

func Exit(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
