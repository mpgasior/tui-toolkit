package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

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

	ctxReader, err := termio.NewReader(fd)
	if err != nil {
		log.Fatal(err)
	}
	defer ctxReader.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	r := termio.NewEscReader(ctx, ctxReader.ReadContext, 2*time.Second)

	ctrlC := [8]byte{3}
	for {
		if ctx.Err() != nil {
			break
		}

		var buf [8]byte
		n, err := r.Read(buf[:])

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
