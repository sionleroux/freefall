// Copyright 2021 Siôn le Roux.  All rights reserved.
// Use of this source code is subject to an MIT-style
// licence which can be found in the LICENSE file.

package main

import (
	"errors"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sinisterstuf/freefall/assets"
	"github.com/sinisterstuf/freefall/game"
	"github.com/sinisterstuf/freefall/nokia"
)

func main() {
	windowScale := 10
	ebiten.SetWindowSize(nokia.GameSize.X*windowScale, nokia.GameSize.Y*windowScale)
	ebiten.SetWindowTitle("Freefall")
	ebiten.SetTPS(15)

	TouchIDs := []ebiten.TouchID{}

	game := &Game{
		Size: nokia.GameSize,
		Screens: []game.Entity{
			&game.TitleScreen{
				Background: assets.LoadImage("title-screen.png"),
				TouchIDs:   &TouchIDs,
			},
			game.NewGameScreen(&TouchIDs),
		},
		TouchIDs: &TouchIDs,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// Game represents the main Ebitengine Game state that coordinates screens
type Game struct {
	Size     image.Point       // Physical game dimensions
	TouchIDs *[]ebiten.TouchID // Re-usable touch ID list
	Screens  []game.Entity     // A slice of all possible game screens
	Screen   game.Screen       // The current screen
}

// Layout is hardcoded for now, may be made dynamic in future
func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.Size.X, g.Size.Y
}

// Update calculates game logic
func (g *Game) Update() error {

	// This should only be needed once per update
	*g.TouchIDs = inpututil.AppendJustPressedTouchIDs((*g.TouchIDs)[:0])

	// Pressing Q any time quits immediately
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return errors.New("game quit by player")
	}

	// Pressing F toggles full-screen
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		if ebiten.IsFullscreen() {
			ebiten.SetFullscreen(false)
		} else {
			ebiten.SetFullscreen(true)
		}
	}

	err := g.Screens[g.Screen].Update()
	var EOS *game.EOS
	if errors.As(err, &EOS) {
		next := EOS.NextScreen
		log.Println("resetting game screen to:", next)
		// Should be better logic here but right now it's just really important
		// to reset this
		g.Screens[game.ScreenGame] = game.NewGameScreen(g.TouchIDs)
		g.Screen = next
		return nil
	}

	return err
}

// Draw draws the game screen by one frame
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(nokia.PaletteOriginal.Light())
	g.Screens[g.Screen].Draw(screen)
}
