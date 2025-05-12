package main

import (
	"fmt"
	"os"

	"github.com/dihedron/overlay/metadata"
	"github.com/jessevdk/go-flags"
)

func main() {

	if len(os.Args) == 2 && (os.Args[1] == "version" || os.Args[1] == "--version") {
		metadata.Print(os.Stdout)
		os.Exit(0)
	} else if len(os.Args) == 3 && os.Args[1] == "version" && (os.Args[2] == "--verbose" || os.Args[2] == "-v") {
		metadata.PrintFull(os.Stdout)
		os.Exit(0)
	}

	var command Command

	var parser = flags.NewParser(&command, flags.Default)

	if args, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	} else {
		if err := command.Execute(args); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}
}
