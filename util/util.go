package util

import (
	"fmt"
	"os"
)

func ExitError(msg interface{}) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
