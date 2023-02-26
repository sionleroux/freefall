package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sinisterstuf/freefall/nokia"
)

// BoxSize is based on the box sprite visual dimensions
const BoxSize = 5

// Box is the player character in the game
type Box struct {
	Coords image.Point
	Chute  bool
	size   int
	HitBox image.Rectangle
}

func (b *Box) Update() error {
	panic("not implemented") // TODO: Implement
}

func (b *Box) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(
		screen,
		float64(b.HitBox.Min.X),
		float64(b.HitBox.Max.Y),
		float64(b.HitBox.Dx()),
		float64(b.HitBox.Dy()),
		nokia.PaletteOriginal.Dark(),
	)
}

func NewBox(coords image.Point, size int) *Box {
	boxOffset := image.Pt(size/2, size/2)
	return &Box{
		Coords: coords,
		size:   size,
		HitBox: image.Rectangle{
			coords.Sub(boxOffset),
			coords.Add(boxOffset),
		},
	}
}

// Move moves the player upwards
func (b *Box) Pull() {
	b.Chute = !b.Chute
}
