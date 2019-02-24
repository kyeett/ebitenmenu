package main

import (
	"fmt"
	"image/color"
	"image/draw"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/peterhellberg/gfx"

	"golang.org/x/image/colornames"

	"github.com/hajimehoshi/ebiten"
)

func (a *App) update(screen *ebiten.Image) error {

	for _, c := range a.components {
		c.Update()
	}

	switch a.components[0].CurrentState() {
	case "map":
		screen.Fill(colornames.Red)
	case "main menu":
		screen.Fill(colornames.Blueviolet)
	default:
		screen.Fill(colornames.Black)
	}

	for _, c := range a.components {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			c.Draw(screen)
		}
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %v", ebiten.CurrentTPS()))
	ebitenutil.DebugPrintAt(screen, "Press shift to open the menu", 0, screen.Bounds().Dy()-20)
	return nil
}

type Menu struct {
	pos      gfx.Vec
	disc     draw.Image
	radius   float64
	options  []string
	selected int
}

func triangleVertex(v gfx.Vec) ebiten.Vertex {
	return ebiten.Vertex{
		DstX:   float32(v.X),
		DstY:   float32(v.Y),
		SrcX:   0,
		SrcY:   0,
		ColorR: 0,
		ColorG: 0,
		ColorB: 0,
		ColorA: 1,
	}
}

func NewMenu(center gfx.Vec, radius float64) *Menu {
	img := gfx.NewImage(int(2*radius), int(2*radius), color.Transparent)
	gfx.DrawCircleFilled(img, gfx.V(radius, radius), radius, colornames.Palegreen)

	eImg, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	return &Menu{
		pos:      center,
		disc:     eImg,
		radius:   radius,
		options:  []string{"initial", "main menu", "map", "game"},
		selected: -1,
	}
}

func (m *Menu) CurrentState() string {
	if m.selected > 0 {
		return m.options[m.selected]
	}
	return ""
}

func (m *Menu) Update() {
	if !ebiten.IsKeyPressed(ebiten.KeyShift) {
		return
	}
	cX, cY := ebiten.CursorPosition()
	cursor := gfx.IV(cX, cY)
	v := m.pos.Sub(cursor)

	if v.Len() > 100 {
		m.selected = -1
		return
	}

	angle := v.Angle() + math.Pi

	// Avoid corner case
	opt := int(gfx.Clamp(float64(len(m.options))*angle/(2*math.Pi), 0, float64(len(m.options))-0.01))
	m.selected = opt
}

func (m *Menu) Draw(screen *ebiten.Image) {
	// w, h := m.radius, m.radius

	// c0, c1, c2, c3 := colornames.Red, colornames.Green, colornames.Blue, colornames.Violet

	// a := []uint8{100, 100, 100, 100}
	angle := 0.0
	if m.selected != -1 {
		angle = float64(m.selected) * 2 * math.Pi / float64(len(m.options))

		totalAngle := 2 * math.Pi / float64(len(m.options))
		extendedRadius := m.radius * 1.5
		center := gfx.V(m.radius, m.radius)
		first := gfx.Unit(angle).Scaled(extendedRadius)
		vs := []ebiten.Vertex{}
		vs = append(vs, triangleVertex(center))
		vs = append(vs, triangleVertex(center.Add(first)))
		vs = append(vs, triangleVertex(center.Add(first.Rotated(totalAngle/2))))
		vs = append(vs, triangleVertex(center.Add(first.Rotated(totalAngle))))

		tmp, _ := ebiten.NewImage(1, 1, ebiten.FilterDefault)
		tmp.Fill(color.White)
		opt := &ebiten.DrawTrianglesOptions{}
		opt.CompositeMode = ebiten.CompositeModeDestinationAtop
		tmp2, _ := ebiten.NewImage(int(2*m.radius), int(2*m.radius), ebiten.FilterDefault)
		tmp2.DrawTriangles(vs, []uint16{0, 1, 2, 0, 2, 3}, tmp, opt)

		op := &ebiten.DrawImageOptions{}
		// op.GeoM.Translate(m.pos.X-extendedRadius, m.pos.Y-extendedRadius)
		op.CompositeMode = ebiten.CompositeModeSourceIn
		tmp2.DrawImage(m.disc.(*ebiten.Image), op)
		// op.GeoM.Translate(m.pos.X-m.radius, m.pos.Y-m.radius)
		// op.ColorM.Scale(1, 1, 1, 0.5)
		op2 := &ebiten.DrawImageOptions{}
		// op3 := &ebiten.DrawImageOptions{}
		op2.GeoM.Translate(m.pos.X-m.radius, m.pos.Y-m.radius)
		screen.DrawImage(tmp2, op2)
		// screen.DrawImage(m.disc.(*ebiten.Image), op)
	}

	op2 := &ebiten.DrawImageOptions{}
	// op3 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(m.pos.X-m.radius, m.pos.Y-m.radius)
	op2.ColorM.Scale(1, 1, 1, 0.5)
	screen.DrawImage(m.disc.(*ebiten.Image), op2)
}

type UIComponent interface {
	Update()
	Draw(screen *ebiten.Image)
	CurrentState() string
}

type App struct {
	mode       string
	components []UIComponent
}

func main() {
	width, height := 400, 400
	a := App{
		mode:       "initial",
		components: []UIComponent{NewMenu(gfx.IV(width/2, height/2), 100)},
	}

	if err := ebiten.Run(a.update, width, height, 2, "menu example"); err != nil {
		log.Fatal(err)
	}
}
