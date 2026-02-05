package screen

import "io"

func Flush(buf Accessor, w io.Writer) error {
	return nil
}

func DiffFlush(next Accessor, current Buffer, w io.Writer) error {
	return nil
}
