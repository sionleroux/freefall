package game

import (
	"image"
	"math/rand"

	"github.com/sinisterstuf/freefall/nokia"
)

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
