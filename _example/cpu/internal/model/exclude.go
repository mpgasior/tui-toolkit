package model

type Exclude int

const (
	ExcludeNone Exclude = iota
	ExcludeExited
	ExcludeActive
	excludeSentinel
)

func (e Exclude) Next() Exclude {
	return (e + 1) % excludeSentinel
}
