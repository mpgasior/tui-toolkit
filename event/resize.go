package event

import (
	"bytes"
	"strconv"

	"github.com/mpgasior/tui-toolkit/ansi"
)

type Resize struct {
	Rows, Cols              int
	PixelHeight, PixelWidth int
}

type ResizeDecoder struct{}

func NewResizeDecoder() *ResizeDecoder {
	return &ResizeDecoder{}
}

func (d *ResizeDecoder) Decode(buf []byte) (resize Resize, ok bool) {
	if !bytes.HasPrefix(buf, ansi.InBandResizePrefixBytes) {
		return resize, false
	}

	if !ansi.InBandResizeFinalRange.Contains(buf[len(buf)-1]) {
		return resize, false
	}

	start, end := len(ansi.InBandResizePrefixBytes), len(buf)-1
	parts := bytes.Split(buf[start:end], []byte{';'})
	if len(parts) != 4 {
		return resize, false
	}

	var err error
	if resize.Rows, err = strconv.Atoi(string(parts[0])); err != nil {
		return resize, false
	}
	if resize.Cols, err = strconv.Atoi(string(parts[1])); err != nil {
		return resize, false
	}
	if resize.PixelHeight, err = strconv.Atoi(string(parts[2])); err != nil {
		return resize, false
	}
	if resize.PixelWidth, err = strconv.Atoi(string(parts[3])); err != nil {
		return resize, false
	}

	return resize, true
}
