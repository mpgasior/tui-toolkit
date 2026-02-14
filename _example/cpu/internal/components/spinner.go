package components

type Spinner struct {
	ID     string
	Frames []rune
	index  int
}

func NewSpinner(id string) *Spinner {
	return &Spinner{
		ID: id,
	}
}

func (s *Spinner) Next() {
	if len(s.Frames) == 0 {
		return
	}
	s.index = (s.index + 1) % len(s.Frames)
}

func (s *Spinner) Frame() rune {
	if len(s.Frames) == 0 {
		return ' '
	}

	return s.Frames[s.index]
}
