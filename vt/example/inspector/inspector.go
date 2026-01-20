package main

import (
	"context"
	"os"
	"time"

	"github.com/nimelo/tui-go/termx"
)

func main() {
	terminal, _ := termx.NewTerminal(os.Stdin, os.Stdout)
	defer terminal.Close()

	restoreInput, _ := terminal.MakeRaw()
	defer restoreInput()

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}
