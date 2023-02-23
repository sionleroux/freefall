package nokia

import (
	"image"
	"image/color"
)

// GameSize is the screen resolution of a Nokia 3310
var GameSize image.Point = image.Point{84, 48}

const (
	ColorTransparent uint8 = iota
	ColorDark
	ColorLight
)

// Media settings based on the Nokia 3310 jam restrictions
var (
	// baseTransparent is completely transparent, used for images that aren't
	// square shaped to show the underlying colour
	baseTransparent color.Color = color.RGBA{67, 82, 61, 0}

	// PaletteOriginal is a 1-bit palette of greenish colours simulating Nokia 3310
	PaletteOriginal Palette = Palette{
		baseTransparent,
		color.RGBA{0x43, 0x52, 0x3d, 0xff}, // Dark
		color.RGBA{0xc7, 0xf0, 0xd8, 0xff}, // Light
	}

	PaletteHarsh Palette = Palette{
		baseTransparent,
		color.RGBA{0x2b, 0x3f, 0x09, 0xff}, // Dark
		color.RGBA{0x9b, 0xc7, 0x00, 0xff}, // Light
	}

	PaletteGray Palette = Palette{
		baseTransparent,
		color.RGBA{0x1a, 0x19, 0x14, 0xff}, // Dark
		color.RGBA{0x87, 0x91, 0x88, 0xff}, // Light
	}
)

// Palette wraps color.Palette with convenience methods for the Nokia greenish
// 2-bit+transparency color palettes
type Palette color.Palette

// Dark returns the darker of the colours, like 0 or OFF
func (p Palette) Dark() color.Color {
	return p[ColorDark]
}

// Light returns the lighter of the colours, like 1 or ON
func (p Palette) Light() color.Color {
	return p[ColorLight]
}
