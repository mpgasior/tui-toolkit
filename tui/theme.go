package tui

type Theme struct {
	Surface   Color
	OnSurface Color

	Primary   Color
	Secondary Color
	Tertiary  Color

	Accent Color
	Danger Color

	Outline        Color
	OutlineVariant Color
}

func NewTheme() Theme {
	return Theme{
		Surface:   ColorHex(0x2B2D42),
		OnSurface: ColorHex(0xEDF2F4),

		Outline:        ColorHex(0x8D99AE),
		OutlineVariant: ColorHex(0x8D99AE),

		Primary:   ColorHex(0xEF233C),
		Secondary: ColorHex(0xD90429),
		Tertiary:  ColorHex(0xD90429),

		Accent: ColorHex(0xEF233C),
		Danger: ColorHex(0xD90429),
	}
}
