package tui

// 2b2d42-8d99ae-edf2f4-ef233c-d90429
type Theme struct {
	Primary    Color
	Secondary  Color
	Background Color
	Foreground Color
	Accent     Color
	Error      Color
	Success    Color

	Text      Style
	Border    Style
	Focus     Style
	Header    Style
	Selection Style
	Inactive  Style
}
