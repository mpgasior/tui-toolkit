package draw

func String(sf SurfaceMutator, str string, opts ...TextWriterOption) {
	w := NewTextWriter(sf, opts...)
	for _, r := range str {
		w.WriteRune(r)
	}
}
