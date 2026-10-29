package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/mpgasior/tui-toolkit/termio"
)

func main() {
	mods := []termio.SessionModifier{
		termio.WithAltScreen,
	}
	t, err := termio.OpenSession(os.Stdin, os.Stdout, mods...)
	if err != nil {
		log.Fatal(err)
	}
	defer t.Close()

	ctrlC := [8]byte{3}
	e := [8]byte{'e'}
	for {
		var buf [8]byte
		n, err := t.ReadContext(context.Background(), buf[:])

		if err != nil {
			fmt.Fprintf(t, "%v\r\n", err)
			continue
		}

		if n == 0 {
			continue
		}

		fmt.Fprintf(t, "%v\r\n", buf)

		if buf == ctrlC {
			break
		}

		if buf == e {
			t.Suspend(func() {
				fmt.Fprint(t, "editing\n")
				editor := os.Getenv("EDITOR")
				cmd := exec.Command(editor)
				cmd.Env = os.Environ()
				cmd.Stdin = t.In()
				cmd.Stdout = t.Out()
				cmd.Stderr = t.Out()

				if err := cmd.Run(); err != nil {
					fmt.Fprintf(t, "%v\r\n", err)
				}
			})
		}
	}
}
