package main

import (
	"os"
	"path/filepath"
)

func Clean(t string, ignores []string) error {
	entries, err := os.ReadDir(t)
	if err != nil {
		return err
	}

	ignoresMap := make(map[string]bool)
	for _, item := range ignores {
		ignoresMap[item] = true
	}

	for _, entry := range entries {
		name := entry.Name()

		if ignoresMap[name] {
			continue
		}

		path := filepath.Join(t, name)

		if entry.IsDir() {
			if err := os.RemoveAll(path); err != nil {
				return err
			}
		} else {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	}

	return nil
}
