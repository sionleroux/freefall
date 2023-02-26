package game

import (
	"fmt"
	"image"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sinisterstuf/freefall/nokia"
)

// Types of screens/scenes in the game
type Screen int

const (
	ScreenTitle Screen = iota // Shows when the game first starts
	ScreenGame                // The actual game itself
	ScreenMax                 // How many screens there are
)

// EOS is an End Of Screen error
// This means the screen's logic has terminated and the controlling game should
// switch to a different screen
type EOS struct {
	NextScreen Screen
}

func (e *EOS) Error() string {
	return fmt.Sprintf("end of screen, next screen: %v", e.NextScreen)
}

// A generic thing that fits into a tree of Ebitengine update-draw calls
type Entity interface {
	Update() error
	Draw(screen *ebiten.Image)
}

type TitleScreen struct {
	Background *ebiten.Image
	TouchIDs   *[]ebiten.TouchID
}

func (t *TitleScreen) Update() error {
	if IsMainActionButtonPressed(t.TouchIDs) {
		return &EOS{ScreenGame}
	}
	return nil
}

func (t *TitleScreen) Draw(screen *ebiten.Image) {
	screen.DrawImage(t.Background, &ebiten.DrawImageOptions{})
}

// GameScreen represents state for the game proper
type GameScreen struct {
	Box         *Box
	Dusts       Dusts
	Projectiles Projectiles
	Tick        int64
	TouchIDs    *[]ebiten.TouchID
}

func (g *GameScreen) Update() error {
	g.Tick++

	if g.Box.Chute {
		if g.Tick%2 == 0 {
			g.Dusts.Update()
		}
	} else {
		g.Dusts.Update()
	}

	if g.Tick%2 == 0 {
		g.Projectiles.Update()
	}

	for _, p := range g.Projectiles {
		ProjHitBox := image.Rectangle{
			p.Coords,
			p.Coords.Add(image.Pt(ProjSize, ProjSize)),
		}
		if g.Box.HitBox.Overlaps(ProjHitBox) {
			log.Printf("game over: %v hit %v", ProjHitBox, g.Box.HitBox)
			return &EOS{ScreenTitle}
		}
	}

	// Movement controls
	if IsMainActionButtonPressed(g.TouchIDs) {
		g.Box.Pull()
	}

	return nil
}

func (g *GameScreen) Draw(screen *ebiten.Image) {
	for _, d := range g.Dusts {
		ebitenutil.DrawRect(
			screen,
			float64(d.Coords.X), float64(d.Coords.Y),
			1, 1,
			nokia.PaletteOriginal.Dark(),
		)
	}

	for _, p := range g.Projectiles {
		p.Draw(screen)
	}

	ebitenutil.DrawRect(
		screen,
		float64(g.Box.HitBox.Min.X),
		float64(g.Box.HitBox.Max.Y),
		float64(g.Box.HitBox.Dx()),
		float64(g.Box.HitBox.Dy()),
		nokia.PaletteOriginal.Dark(),
	)
}

func NewGameScreen(touchIDs *[]ebiten.TouchID) *GameScreen {
	return &GameScreen{
		Box: NewBox(
			image.Pt(nokia.GameSize.X/2, nokia.GameSize.Y/2),
			BoxSize,
		),
		Dusts:       Dusts{},
		Projectiles: Projectiles{},
		TouchIDs:    touchIDs,
	}
}

// BoxSize is based on the box sprite visual dimensions
const BoxSize = 5

// Box is the player character in the game
type Box struct {
	Coords image.Point
	Chute  bool
	size   int
	HitBox image.Rectangle
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

// Dust is decorative dirt on the screen to give the illusion of motion
type Dust struct {
	Coords image.Point
}

func (d *Dust) Update() {
	// Move dusts up
	d.Coords.Y--
}

type Dusts []*Dust

func (ds *Dusts) Update() {
	const maxDusts = 5

	if len(*ds) < maxDusts {
		dsX := rand.Intn(nokia.GameSize.X)
		*ds = append(*ds, &Dust{
			image.Pt(dsX, nokia.GameSize.Y+1),
		})
	}

	for i, d := range *ds {
		d.Update()
		if d.Coords.Y < 0 {
			ds.Drop(i)
		}
	}
}

func (ds *Dusts) Drop(i int) {
	(*ds)[i] = nil
	*ds = append((*ds)[:i], (*ds)[i+1:]...)
}

// Projectile is something that flies across the screen and causes damage if it
// hits the box
type Projectile struct {
	Coords image.Point
	Tail   int
	Size   int
}

const TailMax = 10 // Maximum length of projectile tail
const TailDist = 1 // Distance between projectile and tail
const ProjSize = 2 // How big a projectile's hitbox is

func (p *Projectile) Update() {
	p.Coords.Y--
	p.Coords.X++
	if p.Tail < TailMax {
		p.Tail++
	}
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
			float64(p.Coords.X-TailDist), float64(p.Coords.Y+1),
			float64(p.Coords.X-TailDist-p.Tail), float64(p.Coords.Y+1),
			nokia.PaletteOriginal.Dark(),
		)
	}
}

type Projectiles []*Projectile

func (ps *Projectiles) Update() {
	const maxProjectiles = 2

	if len(*ps) < maxProjectiles {
		psX := rand.Intn(nokia.GameSize.X)
		*ps = append(*ps, &Projectile{
			Coords: image.Pt(psX, nokia.GameSize.Y+1),
			Size:   ProjSize,
		})
	}

	for i, p := range *ps {
		p.Update()
		if p.Coords.Y < 0 {
			ps.Drop(i)
		}
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
