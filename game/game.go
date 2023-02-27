package game

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/sinisterstuf/freefall/assets"
	"github.com/sinisterstuf/freefall/nokia"
	"github.com/tinne26/etxt"
)

const sampleRate int = 44100 // assuming "normal" sample rate
var Context *audio.Context

// Using globals vs meeting deadlines
var HighScore int = 0

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
	Background   *ebiten.Image
	TouchIDs     *[]ebiten.TouchID
	Music        *audio.Player
	SFXFall      *audio.Player
	Box          *Box
	TextRenderer *etxt.Renderer
}

func NewTitleScreen(touchIDs *[]ebiten.TouchID) *TitleScreen {
	return &TitleScreen{
		Background: assets.LoadImage("title-screen.png"),
		Music:      assets.NewMusicPlayer(assets.LoadSoundFile("freefall-maintheme.ogg", sampleRate), Context),
		SFXFall:    assets.NewSoundPlayer(assets.LoadSoundFile("sfxfall.ogg", sampleRate), Context),
		TouchIDs:   touchIDs,
		Box: NewBox(
			image.Pt(nokia.GameSize.X/2, -BoxSize),
			BoxSize,
		),
		TextRenderer: NewTextRenderer(),
	}
}

func (t *TitleScreen) Update() error {
	if !t.Music.IsPlaying() {
		t.Music.Play()
		t.Box.Coords.Y = -BoxSize * 2
	}

	if t.Box.Coords.Y < nokia.GameSize.Y+BoxSize*2 {
		t.Box.Coords.Y++
	}
	t.Box.Frame = assets.Animate(t.Box.Frame, t.Box.Tick, t.Box.Sprite.Meta.FrameTags[t.Box.State])

	if IsMainActionButtonPressed(t.TouchIDs) {
		t.Music.Pause()
		t.Music.Rewind()
		t.SFXFall.Rewind()
		t.SFXFall.Play()
		return &EOS{ScreenGame}
	}
	return nil
}

func (t *TitleScreen) Draw(screen *ebiten.Image) {
	screen.DrawImage(t.Background, &ebiten.DrawImageOptions{})
	t.Box.Draw(screen)
	if HighScore > 0 {
		txt := t.TextRenderer
		txt.SetTarget(screen)
		txt.Draw(
			fmt.Sprintf("Best: %dm", HighScore),
			screen.Bounds().Dx()/2,
			screen.Bounds().Dy()/8*7,
		)
	}
}

// GameScreen represents state for the game proper
type GameScreen struct {
	Box         *Box
	Dusts       Dusts
	Projectiles Projectiles
	Tick        int
	TouchIDs    *[]ebiten.TouchID
	SFXHit      *audio.Player
}

func (g *GameScreen) Update() error {
	g.Tick++

	g.Box.Update()

	if g.Box.Chute {
		if g.Tick%2 == 0 {
			g.Dusts.MoveUp()
			g.Projectiles.MoveUp()
		}
	} else {
		g.Dusts.MoveUp()
		g.Projectiles.MoveUp()
	}

	// Difficulty
	if g.Tick%100 == 0 {
		if maxProjectiles < 20 {
			maxProjectiles += 2
		}
	}

	g.Dusts.Update()
	g.Projectiles.Update(g.Tick)

	for _, p := range g.Projectiles {
		ProjHitBox := image.Rectangle{
			p.Coords,
			p.Coords.Add(image.Pt(ProjSize, ProjSize)),
		}
		if g.Box.HitBox.Overlaps(ProjHitBox) {
			log.Printf("game over: %v hit %v", ProjHitBox, g.Box.HitBox)
			if g.Tick > HighScore {
				HighScore = g.Tick
			}
			g.SFXHit.Rewind()
			g.SFXHit.Play()
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
			image.Pt(nokia.GameSize.X/2, nokia.GameSize.Y/6),
			BoxSize,
		),
		Dusts:       Dusts{},
		Projectiles: Projectiles{},
		TouchIDs:    touchIDs,
		SFXHit:      assets.NewSoundPlayer(assets.LoadSoundFile("sfxhit.ogg", sampleRate), Context),
	}
}

func NewTextRenderer() *etxt.Renderer {
	font := assets.LoadFont("tiny.ttf")
	r := etxt.NewStdRenderer()
	r.SetFont(font)
	r.SetAlign(etxt.YCenter, etxt.XCenter)
	r.SetColor(nokia.PaletteOriginal.Dark())
	r.SetSizePx(6)
	return r
}
