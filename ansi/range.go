package ansi

type ByteRange struct {
	Lo byte
	Hi byte
}

func (r ByteRange) Contains(b byte) bool {
	return b >= r.Lo && b <= r.Hi
}

func (r ByteRange) EndsWith(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	return r.Contains(b[len(b)-1])
}

type RangeTable []ByteRange

func (rt RangeTable) Contains(b byte) bool {
	for i := range rt {
		if b >= rt[i].Lo && b <= rt[i].Hi {
			return true
		}
	}
	return false
}
