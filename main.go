package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const (
	screenW = 400
	screenH = 600
	jetX    = 60
	jetSize = 40
	gravity = 0.5
	jump    = -9
	towerW  = 70
	gap     = 220
	speed   = 3

	scoreFile = "scores.txt"
)

type State int

const (
	Login State = iota
	Menu
	Gameplay
	Leaderboard
	Profile
	Shop
)

type Tower struct {
	x      float64
	h      float64
	scored bool
}
type Coin struct {
	x, y      float64
	collected bool
}

type Game struct {
	state     State
	username  string
	input     string
	menuIndex int
	jetY      float64
	vel       float64
	towers    []Tower
	coins     []Coin
	score     int
	best      int
}

var playerCoins = map[string]int{}
var ownedSkins = map[string][]string{}
var selectedSkin = map[string]string{}
var skinPrices = map[string]int{
	"Red Jet":    10,
	"Yellow Jet": 20,
}
var (
	bgImage    *ebiten.Image
	jetImage   *ebiten.Image
	towerImage *ebiten.Image
	coinImage  *ebiten.Image
)

func loadImage(path string) *ebiten.Image {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil
	}
	return ebiten.NewImageFromImage(img)
}
func loadScores() map[string]int {
	s := map[string]int{}
	f, err := os.Open(scoreFile)
	if err != nil {
		return s
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		p := strings.Split(sc.Text(), ",")
		if len(p) == 2 {
			v, _ := strconv.Atoi(p[1])
			s[p[0]] = v
		}
	}
	return s
}
func saveScores(m map[string]int) {
	f, _ := os.Create(scoreFile)
	defer f.Close()

	for k, v := range m {
		fmt.Fprintf(f, "%s,%d\n", k, v)
	}
}
func updateScore(user string, score int) {
	s := loadScores()
	if score > s[user] {
		s[user] = score
	}
	saveScores(s)
}
func leaderboard() []struct {
	name  string
	score int
} {
	type pair struct {
		name  string
		score int
	}
	var list []pair
	for k, v := range loadScores() {
		list = append(list, pair{k, v})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].score > list[j].score
	})
	var out []struct {
		name  string
		score int
	}
	for _, p := range list {
		out = append(out, struct {
			name  string
			score int
		}{p.name, p.score})
	}
	return out
}
func newTower() Tower {
	h := rand.Float64()*(screenH-gap-100) + 50
	return Tower{screenW, h, false}
}
func (g *Game) reset() {
	g.jetY = screenH / 2
	g.vel = 0
	g.towers = nil
	g.coins = nil
	g.score = 0
	g.best = loadScores()[g.username]
	if _, ok := playerCoins[g.username]; !ok {
		playerCoins[g.username] = 0
	}
	if _, ok := ownedSkins[g.username]; !ok {
		ownedSkins[g.username] = []string{"Default"}
	}
	if selectedSkin[g.username] == "" {
		selectedSkin[g.username] = "Default"
	}
	g.menuIndex = 0
}
func (g *Game) Update() error {
	switch g.state {
	case Login:
		for _, r := range ebiten.InputChars() {
			if len(g.input) < 12 {
				g.input += string(r)
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(g.input) > 0 {
			g.input = g.input[:len(g.input)-1]
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && strings.TrimSpace(g.input) != "" {
			g.username = g.input
			g.state = Menu
			g.menuIndex = 0
		}
	case Menu:
		opts := []string{"Play", "Leaderboard", "Profile", "Shop", "Quit"}
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			g.menuIndex--
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			g.menuIndex++
		}
		if g.menuIndex < 0 {
			g.menuIndex = len(opts) - 1
		}
		if g.menuIndex >= len(opts) {
			g.menuIndex = 0
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			switch opts[g.menuIndex] {
			case "Play":
				g.reset()
				g.state = Gameplay
			case "Leaderboard":
				g.state = Leaderboard
				g.menuIndex = 0
			case "Profile":
				g.state = Profile
				g.menuIndex = 0
			case "Shop":
				g.state = Shop
				g.menuIndex = 0
			case "Quit":
				os.Exit(0)
			}
		}
	case Gameplay:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.vel = jump
		}
		g.vel += gravity
		g.jetY += g.vel
		if len(g.towers) == 0 || g.towers[len(g.towers)-1].x < screenW-200 {
			g.towers = append(g.towers, newTower())
		}
		for i := range g.towers {
			g.towers[i].x -= speed
			if !g.towers[i].scored && g.towers[i].x+towerW < jetX {
				g.score++
				g.towers[i].scored = true
			}
			if jetX+jetSize > g.towers[i].x &&
				jetX < g.towers[i].x+towerW {
				if g.jetY < g.towers[i].h ||
					g.jetY+jetSize > g.towers[i].h+gap {
					updateScore(g.username, g.score)
					g.state = Menu
					g.menuIndex = 0
				}
			}
		}
		filter := g.towers[:0]
		for _, t := range g.towers {
			if t.x > -towerW {
				filter = append(filter, t)
			}
		}
		g.towers = filter
		if rand.Intn(100) < 2 && len(g.towers) > 0 {
			t := g.towers[len(g.towers)-1]
			y := rand.Float64()*(gap-40) + t.h + 20
			g.coins = append(g.coins, Coin{screenW, y, false})
		}
		for i := range g.coins {
			g.coins[i].x -= speed
			if !g.coins[i].collected &&
				g.coins[i].x < jetX+jetSize &&
				g.coins[i].x+10 > jetX &&
				g.coins[i].y > g.jetY &&
				g.coins[i].y < g.jetY+jetSize {
				playerCoins[g.username]++
				g.coins[i].collected = true
			}
		}
		cf := g.coins[:0]
		for _, c := range g.coins {
			if !c.collected && c.x > -10 {
				cf = append(cf, c)
			}
		}
		g.coins = cf
		if g.jetY < 0 || g.jetY > screenH {
			updateScore(g.username, g.score)
			g.state = Menu
			g.menuIndex = 0
		}
	case Leaderboard, Profile:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = Menu
			g.menuIndex = 0
		}
	case Shop:
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			g.menuIndex--
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			g.menuIndex++
		}
		if g.menuIndex < 0 {
			g.menuIndex = len(skinPrices) - 1
		}
		if g.menuIndex >= len(skinPrices) {
			g.menuIndex = 0
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			i := 0
			for skin, price := range skinPrices {
				if i == g.menuIndex {
					owned := false
					for _, s := range ownedSkins[g.username] {
						if s == skin {
							owned = true
							break
						}
					}
					if owned {
						selectedSkin[g.username] = skin
					} else if playerCoins[g.username] >= price {
						playerCoins[g.username] -= price
						ownedSkins[g.username] = append(ownedSkins[g.username], skin)
						selectedSkin[g.username] = skin
					}
					break
				}
				i++
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = Menu
			g.menuIndex = 0
		}
	}
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	if bgImage != nil {
		screen.DrawImage(bgImage, nil)
	} else {
		screen.Fill(color.RGBA{200, 220, 255, 255})
	}
	switch g.state {
	case Login:
		text.Draw(screen, "Enter Username", basicfont.Face7x13, 130, 260, color.Black)
		text.Draw(screen, g.input, basicfont.Face7x13, 130, 290, color.Black)
	case Menu:
		opts := []string{"Play", "Leaderboard", "Profile", "Shop", "Quit"}
		text.Draw(screen, "Welcome "+g.username, basicfont.Face7x13, 120, 120, color.Black)
		for i, o := range opts {
			c := color.Color(color.Black)
			if i == g.menuIndex {
				c = color.Color(color.RGBA{0, 200, 0, 255})
			}
			text.Draw(screen, o, basicfont.Face7x13, 160, 220+i*40, c)
		}
	case Gameplay:
		if jetImage != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(float64(jetSize)/float64(jetImage.Bounds().Dx()), float64(jetSize)/float64(jetImage.Bounds().Dy()))
			op.GeoM.Translate(jetX, g.jetY)
			screen.DrawImage(jetImage, op)
		} else {
			ebitenutil.DrawRect(screen, jetX, g.jetY, jetSize, jetSize, color.RGBA{0, 0, 255, 255})
		}
		for _, t := range g.towers {
			hTop := t.h
			hBottom := screenH - (t.h + gap)
			if towerImage != nil {
				op := &ebiten.DrawImageOptions{} //top tower
				scaleY := hTop / float64(towerImage.Bounds().Dy())
				op.GeoM.Scale(towerW/float64(towerImage.Bounds().Dx()), -scaleY)
				op.GeoM.Translate(t.x, hTop)
				screen.DrawImage(towerImage, op)
				op2 := &ebiten.DrawImageOptions{} //bottom tower
				op2.GeoM.Scale(towerW/float64(towerImage.Bounds().Dx()), hBottom/float64(towerImage.Bounds().Dy()))
				op2.GeoM.Translate(t.x, t.h+gap)
				screen.DrawImage(towerImage, op2)
			} else {
				ebitenutil.DrawRect(screen, t.x, 0, towerW, hTop, color.RGBA{0, 200, 0, 255})
				ebitenutil.DrawRect(screen, t.x, t.h+gap, towerW, hBottom, color.RGBA{0, 200, 0, 255})
			}
		}
		for _, c := range g.coins {
			if c.collected {
				continue
			}
			if coinImage != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(10/float64(coinImage.Bounds().Dx()), 10/float64(coinImage.Bounds().Dy()))
				op.GeoM.Translate(c.x, c.y)
				screen.DrawImage(coinImage, op)
			} else {
				ebitenutil.DrawRect(screen, c.x, c.y, 10, 10, color.RGBA{255, 215, 0, 255})
			}
		}
		text.Draw(screen, fmt.Sprintf("Score:%d", g.score), basicfont.Face7x13, 10, 20, color.Black)
		text.Draw(screen, fmt.Sprintf("Coins:%d", playerCoins[g.username]), basicfont.Face7x13, 10, 40, color.Black)
		text.Draw(screen, fmt.Sprintf("Best:%d", g.best), basicfont.Face7x13, 300, 20, color.Black)
	case Leaderboard:
		text.Draw(screen, "Leaderboard", basicfont.Face7x13, 150, 50, color.Black)
		leaders := leaderboard()
		for i, p := range leaders {
			if i >= 10 {
				break
			}
			text.Draw(screen, fmt.Sprintf("%d. %s - %d", i+1, p.name, p.score), basicfont.Face7x13, 100, 100+i*20, color.Black)
		}
		text.Draw(screen, "Press ESC to go back", basicfont.Face7x13, 100, 500, color.Color(color.RGBA{100, 100, 100, 255}))
	case Profile:
		text.Draw(screen, "Profile", basicfont.Face7x13, 150, 50, color.Black)
		text.Draw(screen, fmt.Sprintf("Username: %s", g.username), basicfont.Face7x13, 100, 100, color.Black)
		text.Draw(screen, fmt.Sprintf("Coins: %d", playerCoins[g.username]), basicfont.Face7x13, 100, 130, color.Black)
		text.Draw(screen, fmt.Sprintf("Owned Skins: %v", ownedSkins[g.username]), basicfont.Face7x13, 100, 160, color.Black)
		text.Draw(screen, fmt.Sprintf("Selected Skin: %s", selectedSkin[g.username]), basicfont.Face7x13, 100, 190, color.Black)
		text.Draw(screen, "Press ESC to go back", basicfont.Face7x13, 100, 500, color.Color(color.RGBA{100, 100, 100, 255}))
	case Shop:
		text.Draw(screen, "Shop", basicfont.Face7x13, 150, 50, color.Black)
		y := 100
		i := 0
		for skin, price := range skinPrices {
			status := "Buy"
			for _, s := range ownedSkins[g.username] {
				if s == skin {
					status = "Owned"
					break
				}
			}
			c := color.Color(color.Black)
			if i == g.menuIndex {
				c = color.Color(color.RGBA{0, 200, 0, 255})
			}
			text.Draw(screen, fmt.Sprintf("%s - %d coins [%s]", skin, price, status), basicfont.Face7x13, 100, y, c)
			y += 30
			i++
		}
		text.Draw(screen, "ESC to go back", basicfont.Face7x13, 50, 500, color.Color(color.RGBA{100, 100, 100, 255}))
	}
}
func (g *Game) Layout(int, int) (int, int) {
	return screenW, screenH
}
func main() {
	rand.Seed(time.Now().UnixNano())
	bgImage = loadImage("background.png")
	jetImage = loadImage("jet.png")
	towerImage = loadImage("tower.png")
	coinImage = loadImage("coin.png")
	game := &Game{state: Login}
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Ready Jet Go")

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
