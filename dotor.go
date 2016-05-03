package main

import (
	"fmt"
	"os"

	"github.com/januswel/dotor/core"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, r)
		}
	}()

	if len(os.Args) != 3 {
		panic(fmt.Errorf("Usage: %s <setting file name> <source path>", os.Args[0]))
	}
	settingsFileName := os.Args[1]
	sourcePath := os.Args[2]

	if err := core.Execute(settingsFileName, sourcePath); err != nil {
		panic(err)
	}
}
