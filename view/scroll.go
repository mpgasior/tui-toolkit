package view

type Scroll struct {
	Index  int
	Offset int
	Margin int
}

func (s *Scroll) Update(size, total int) (start, end int) {
	if total == 0 {
		s.Index = -1
		s.Offset = 0
		return 0, 0
	}

	if s.Index < -1 {
		s.Index = -1
	} else if s.Index >= total {
		s.Index = total - 1
	}

	if s.Index < s.Offset+s.Margin {
		s.Offset = s.Index - s.Margin
	}

	if s.Index >= s.Offset+size-s.Margin {
		s.Offset = s.Index - size + s.Margin + 1
	}

	if s.Offset > total-size {
		s.Offset = total - size
	}

	if s.Offset < 0 {
		s.Offset = 0
	}

	start = max(s.Offset, 0)
	end = min(s.Offset+size, total)

	return start, end
}
