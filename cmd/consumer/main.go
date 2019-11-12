package main

import (
	"fmt"
	"github.com/fusion-app/fusion-app/cmd/consumer/commands"
	"os"
)

func main() {
	if err := commands.NewCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

