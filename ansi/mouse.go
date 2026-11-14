package ansi

const (
	MouseNormalModeOn  = CSI + "?1000h"
	MouseNormalModeOff = CSI + "?1000l"

	MouseButtonModeOn  = CSI + "?1002h"
	MouseButtonModeOff = CSI + "?1002l"

	MouseAnyModeOn  = CSI + "?1003h"
	MouseAnyModeOff = CSI + "?1003l"

	MouseSGRModeOn  = CSI + "?1006h"
	MouseSGRModeOff = CSI + "?1006l"
)

var (
	MouseNormalModeOnBytes  = []byte(MouseNormalModeOn)
	MouseNormalModeOffBytes = []byte(MouseNormalModeOff)

	MouseButtonModeOnBytes  = []byte(MouseButtonModeOn)
	MouseButtonModeOffBytes = []byte(MouseButtonModeOff)

	MouseAnyModeOnBytes  = []byte(MouseAnyModeOn)
	MouseAnyModeOffBytes = []byte(MouseAnyModeOff)

	MouseSGRModeOnBytes  = []byte(MouseSGRModeOn)
	MouseSGRModeOffBytes = []byte(MouseSGRModeOff)
)

const MouseLegacyPrefix = CSI + "M"

var (
	MouseLegacyPrefixBytes      = []byte(MouseLegacyPrefix)
	MouseLegacyLength           = len(MouseLegacyPrefixBytes) + 3
	MouseLegacyCtrlIndex        = len(MouseLegacyPrefixBytes)
	MouseLegacyColIndex         = len(MouseLegacyPrefixBytes) + 1
	MouseLegacyRowIndex         = MouseLegacyColIndex + 1
	MouseLegacyOffset      byte = 32
)

const MouseSGRPrefix = CSI + "<"

var (
	MouseSGRPrefixBytes = []byte(MouseSGRPrefix)
	MouseSGRFinalRange  = RangeTable{
		{Lo: 'm', Hi: 'm'},
		{Lo: 'M', Hi: 'M'},
	}
)
