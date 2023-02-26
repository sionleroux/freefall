package game

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
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
	g.Dusts.Draw(screen)
	g.Projectiles.Draw(screen)
	g.Box.Draw(screen)
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
