package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	err := SyncVite()
	if errors.Is(err, ErrNoTemplatesFound) {
		fmt.Fprintln(os.Stderr, "Vite: ", err)
		os.Exit(1)
	} else if err != nil {
		fmt.Println(err)
	}
}
