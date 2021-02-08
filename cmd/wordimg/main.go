package main

import (
	"fmt"
	"os"
)

func main() {
	command := new(commandLine)
	if err := command.parse(os.Args); err != nil {
		os.Exit(1)
	}
	if err := command.exec(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdin, "wrote: %s\n", command.ImagePath)
}
