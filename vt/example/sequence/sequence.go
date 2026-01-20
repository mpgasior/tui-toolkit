package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nimelo/tui-go/bufiox"
	"github.com/nimelo/tui-go/termx"
	"github.com/nimelo/tui-go/vt"
)

func main() {
	terminal, _ := termx.NewTerminal(os.Stdin, os.Stdout)
	defer terminal.Close()

	restoreInput, _ := terminal.MakeRaw()
	defer restoreInput()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	scanner := bufiox.NewContextScanner(terminal)
	scanner.Split(vt.ScanUtf8)

	for scanner.Scan(ctx) {
		bytes := scanner.Bytes()
		str := string(bytes)

		fmt.Fprintf(terminal, "%s\r\n", str)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(terminal, "Error: %v\r\n", err)
	}
}
