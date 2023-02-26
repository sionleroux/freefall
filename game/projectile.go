package game

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sinisterstuf/freefall/nokia"
)

// Projectile is something that flies across the screen and causes damage if it
// hits the box
type Projectile struct {
	Coords   image.Point
	Tail     int
	Size     int
	Velocity int
}

const TailMax = 10 // Maximum length of projectile tail
const TailDist = 1 // Distance between projectile and tail
const ProjSize = 2 // How big a projectile's hitbox is

func (p *Projectile) Update() {
	p.Coords.X = p.Coords.X + p.Velocity
	if p.Tail < TailMax {
		p.Tail++
	}
}

func (p *Projectile) MoveUp() {
	p.Coords.Y--
}

func (p *Projectile) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(
		screen,
		float64(p.Coords.X), float64(p.Coords.Y),
		ProjSize, ProjSize,
		nokia.PaletteOriginal.Dark(),
	)
	if p.Tail > 0 {
		ebitenutil.DrawLine(
			screen,
			float64(p.Coords.X-(ProjSize+TailDist)*p.Velocity), float64(p.Coords.Y+1),
			float64(p.Coords.X-(ProjSize+TailDist+p.Tail)*p.Velocity), float64(p.Coords.Y+1),
			nokia.PaletteOriginal.Dark(),
		)
	}
}

type Projectiles []*Projectile

func (ps *Projectiles) Draw(screen *ebiten.Image) {
	for _, p := range *ps {
		p.Draw(screen)
	}
}

func (ps *Projectiles) Update() {
	const maxProjectiles = 2

	if len(*ps) < maxProjectiles {
		spawnSide := rand.Intn(2) * nokia.GameSize.X // left or right of screen
		var velocity int
		if spawnSide == 0 {
			velocity = 1
		} else {
			velocity = -1
		}
		*ps = append(*ps, &Projectile{
			Coords:   image.Pt(spawnSide, nokia.GameSize.Y+1),
			Size:     ProjSize,
			Velocity: velocity,
		})
	}

	for i, p := range *ps {
		p.Update()
		if p.Coords.Y < 0 {
			ps.Drop(i)
		}
	}
}

func (ps *Projectiles) MoveUp() {
	for _, p := range *ps {
		p.MoveUp()
	}
}

func (ps *Projectiles) Drop(i int) {
	(*ps)[i] = nil
	*ps = append((*ps)[:i], (*ps)[i+1:]...)
}

// Main action button is 5, like in the middle of a Nokia 3310
// Fallbacks for people without a numpad:
//   - Spacebar on desktop
//   - Tap the screen on mobile
func IsMainActionButtonPressed(TouchIDs *[]ebiten.TouchID) bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyNumpad5) ||
		inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		len(*TouchIDs) > 0 {
		return true
	}
	return false
}
