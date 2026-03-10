package main

import (
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenwidth  = 400
	screenheight = 600
	jetsize      = 20
	towerwidth   = 50
	towergap     = 150
	jetspeed     = 5
	gravity      = 0.5
	jetstrength  = -10
	fps          = 60
)

type Game struct {
	jetY            float64
	velocity        float64
	towers          []Tower
	score           int
	jetImage        *ebiten.Image
	towerImage      *ebiten.Image
	backgroundImage *ebiten.Image
}
type Tower struct {
	x      float64
	height float64
	scored bool
}

func loadImage(path string) *ebiten.Image {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	return ebiten.NewImageFromImage(img)
}
func createTower() Tower {

	height := rand.Float64()*300 + 50

	return Tower{
		x:      screenwidth,
		height: height,
		scored: false,
	}
}

//update loop

func (g *Game) Update() error {

	// Jump
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.velocity = jetstrength
	}
	g.velocity += gravity
	g.jetY += g.velocity
	if len(g.towers) == 0 || g.towers[len(g.towers)-1].x < screenwidth-200 {
		g.towers = append(g.towers, createTower())
	}

	for i := range g.towers {

		g.towers[i].x -= fps

		// Score system
		if g.towers[i].x+towerwidth < 50 && !g.towers[i].scored {
			g.score++
			g.towers[i].scored = true
		}
	}

	// Remove old towers
	if len(g.towers) > 0 && g.towers[0].x < -towerwidth {
		g.towers = g.towers[1:]
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// fcking images
	screen.DrawImage(g.backgroundImage, nil)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(50, g.jetY)
	screen.DrawImage(g.jetImage, op)
	for _, tower := range g.towers {

		// towers
		top := &ebiten.DrawImageOptions{}
		top.GeoM.Scale(1, tower.height/float64(g.towerImage.Bounds().Dy()))
		top.GeoM.Translate(tower.x, 0)
		screen.DrawImage(g.towerImage, top)
		bottom := &ebiten.DrawImageOptions{}
		bottomHeight := screenheight - (tower.height + towergap)
		bottom.GeoM.Scale(1, bottomHeight/float64(g.towerImage.Bounds().Dy()))
		bottom.GeoM.Translate(tower.x, tower.height+towergap)
		screen.DrawImage(g.towerImage, bottom)
	}

	// score
	ebitenutil.DebugPrint(screen, "Score: "+strconv.Itoa(g.score))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenwidth, screenheight
}

func main() {

	rand.Seed(time.Now().UnixNano())

	game := &Game{
		jetY: screenheight / 2,
	}
	game.jetImage = loadImage("jet.png")
	game.towerImage = loadImage("tower.png")
	game.backgroundImage = loadImage("background.png")
	ebiten.SetWindowSize(screenwidth, screenheight)
	ebiten.SetWindowTitle("Ready... jetniggers Set Go!")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
