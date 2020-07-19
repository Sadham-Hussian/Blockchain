package main

import (
	"os"

	"github.com/Sadham-Hussian/Blockchain/cli"
)

func main() {
	defer os.Exit(0)
	cli := cli.CommandLine{}
	cli.Run()
}
