package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mpgasior/tui-toolkit/termio"
	"golang.org/x/term"
)

func main() {
	fd := int(os.Stdin.Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		log.Fatal(err)
	}
	defer term.Restore(fd, state)

	r, err := termio.NewReader(fd)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	ctrlC := [8]byte{3}
	for {
		var buf [8]byte
		n, err := r.ReadContext(context.Background(), buf[:])

		if err != nil {
			fmt.Printf("%v\r\n", err)
			continue
		}

		if n == 0 {
			continue
		}

		fmt.Printf("%v\r\n", buf)

		if buf == ctrlC {
			break
		}
	}
}
