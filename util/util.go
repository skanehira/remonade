package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func NotExist(path string) bool {
	_, err := os.Stat(path)
	return errors.Is(err, os.ErrNotExist)
}

func ExitError(msg interface{}) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func Create(path string) error {
	base := filepath.Dir(path)

	if NotExist(base) {
		if err := os.Mkdir(base, os.ModePerm); err != nil {
			return err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}
