package main

import (
	"context"
	"fmt"
	"os"
	"time"

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

	scanner := vt.NewSequenceScanner(terminal, vt.TODO)

	for scanner.ScanContext(ctx) {
		seq := scanner.Sequence()

		fmt.Fprintf(terminal, "%s: [% X] \r\n", seq.Type.String(), seq.Data)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(terminal, "Error: %v\r\n", err)
	}
}
