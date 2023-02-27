package game

//go:generate ../tools/gen_sprite_tags.sh ../assets/box.json box_anim.go box

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sinisterstuf/freefall/assets"
)

// BoxSize is based on the box sprite visual dimensions
const BoxSize = 5

// Box is the player character in the game
type Box struct {
	Coords image.Point
	Chute  bool
	size   int
	HitBox image.Rectangle
	State  boxAnimationTags // Current animation state
	Frame  int              // Current animation frame
	Sprite *assets.SpriteSheet
	Tick   int
}

func (b *Box) Update() error {
	b.Tick++
	b.Frame = assets.Animate(b.Frame, b.Tick, b.Sprite.Meta.FrameTags[b.State])
	if b.Frame == b.Sprite.Meta.FrameTags[b.State].To {
		switch b.State {
		case boxOpening:
			b.State = boxOpen
		case boxClosing:
			b.State = boxClosed
		}
	}
	return nil
}

func (b *Box) Draw(screen *ebiten.Image) {
	s := b.Sprite
	frame := s.Sprite[b.Frame]
	op := &ebiten.DrawImageOptions{}

	// Centre
	op.GeoM.Translate(
		float64(-frame.Position.W/2),
		float64(-frame.Position.H/2),
	)
	// Position
	op.GeoM.Translate(
		float64(b.Coords.X),
		float64(b.Coords.Y),
	)

	screen.DrawImage(
		s.Image.SubImage(image.Rect(
			frame.Position.X,
			frame.Position.Y,
			frame.Position.X+frame.Position.W,
			frame.Position.Y+frame.Position.H,
		)).(*ebiten.Image),
		op,
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
		Sprite: assets.LoadSprite("box"),
		State:  boxClosed,
	}
}

// Move moves the player upwards
func (b *Box) Pull() {
	if b.State != boxOpening && b.State != boxClosing {
		b.Chute = !b.Chute
	}
	if b.State == boxOpen {
		b.State = boxClosing
	}
	if b.State == boxClosed {
		b.State = boxOpening
	}
}
