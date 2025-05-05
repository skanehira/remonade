package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/skanehira/remonade/util"
)

func TestRunInit(t *testing.T) {
	t.Run("empty token", func(t *testing.T) {
		tmp := filepath.Join(os.TempDir(), "d")
		if !util.NotExist(tmp) {
			if err := os.RemoveAll(tmp); err != nil {
				t.Fatal(err)
			}
		}
		err := runInit(filepath.Join(tmp, "f"))

		if err != ErrEmptyToken {
			t.Fatalf("unexpected error, want=%s, got=%s", ErrEmptyToken, err)
		}
	})

	t.Run("empty token when existed file", func(t *testing.T) {
		tmp := filepath.Join(os.TempDir(), "d")
		err := runInit(filepath.Join(tmp, "f"))
		want := ErrEmptyToken
		if err != want {
			t.Fatalf("unexpected error, want=%s, got=%s", want, err)
		}
	})
}

func TestRunEdit(t *testing.T) {
	t.Run("no exists config", func(t *testing.T) {
		tmp := filepath.Join(os.TempDir(), "d")
		if !util.NotExist(tmp) {
			if err := os.RemoveAll(tmp); err != nil {
				t.Fatal(err)
			}
		}
		err := runEdit(filepath.Join(tmp, "f"))
		want := ErrNotExistConfig
		if err != want {
			t.Fatalf("unexpected error, want=%s, got=%s", want, err)
		}
	})

	t.Run("no $EDITOR", func(t *testing.T) {
		editor := os.Getenv("EDITOR")
		_ = os.Setenv("EDITOR", "")
		t.Cleanup(func() {
			_ = os.Setenv("EDITOR", editor)
		})

		tmp, err := os.CreateTemp("", "")
		if err != nil {
			t.Fatalf("unexpected error, want=nil, got=%s", err)
		}
		err = runEdit(tmp.Name())
		want := ErrEmptyEDITOR
		if err != want {
			t.Fatalf("unexpected error, want=%s, got=%s", want, err)
		}
	})
}
