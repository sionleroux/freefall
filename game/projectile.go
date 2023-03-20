package game

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sinisterstuf/freefall/nokia"
)

// Projectile is something that flies across the screen and causes damage if it
// hits the box
type Projectile struct {
	Coords   Point
	Tail     int
	Size     int
	Velocity float64 // Direction and speed
	Spacing  int     // How far away to place the next one
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
			p.Coords.X-(ProjSize+TailDist)*p.Velocity, p.Coords.Y+1,
			p.Coords.X-(ProjSize+TailDist+float64(p.Tail))*p.Velocity, p.Coords.Y+1,
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

var maxProjectiles = 2

func (ps *Projectiles) Update(tick int) {
	if len(*ps) == 0 {
		ps.Spawn(tick)
	}

	if len(*ps) < maxProjectiles && tick > (*ps)[len(*ps)-1].Spacing {
		ps.Spawn(tick)
	}

	for i, p := range *ps {
		if tick%2 == 0 {
			p.Update()
		}
		if p.Coords.Y < 0 {
			ps.Drop(i)
		}
	}
}

const maxSpacing = 15

func (ps *Projectiles) Spawn(tick int) {
	spawnSide := rand.Intn(2) * nokia.GameSize.X // left or right of screen
	speedMin, speedMax := 0.8, 3.0
	speed := speedMin + rand.Float64()*(speedMax-speedMin)
	var velocity float64
	if spawnSide == 0 {
		velocity = speed
	} else {
		velocity = -speed
	}
	*ps = append(*ps, &Projectile{
		Coords:   Point{float64(spawnSide), float64(nokia.GameSize.Y + 1)},
		Size:     ProjSize,
		Velocity: velocity,
		Spacing:  tick + rand.Intn(maxSpacing),
	})
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
