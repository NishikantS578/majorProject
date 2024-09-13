package main

import (
	"fmt"
	"io"
	"majorProject/compiler/repl"
	"os"
)

func main() {
	var args = os.Args

	var ioSrc io.Reader
	var ioDest io.Writer
	var err error

	if len(args) == 1 {
		ioSrc = os.Stdin
		ioDest = os.Stdout
	} else {
		ioSrc, err = os.Open(args[1])

		if err != nil {
			fmt.Println("Could not read file: " + args[1])
			os.Exit(0)
		}

		ioDest, err = os.Create("app.asm")

		if err != nil {
			fmt.Println("Could not create file: app.asm")
			os.Exit(0)
		}
	}

	repl.Run(ioSrc, ioDest)
}
