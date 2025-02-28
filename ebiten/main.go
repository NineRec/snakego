package main

import (
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 640
	screenHeight = 480
	gridSize     = 20

	updateDelay = 10 // 更新延迟，增加这个值可以减慢速度
)

type Game struct {
	snake     []Position
	direction Position
	food      Position
	gameOver  bool
	score     int
	tickCount int
}

type Position struct {
	X, Y int
}

func (g *Game) Update() error {
	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.initGame()
		}
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && g.direction != (Position{0, 1}) {
		g.direction = Position{0, -1}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && g.direction != (Position{0, -1}) {
		g.direction = Position{0, 1}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && g.direction != (Position{1, 0}) {
		g.direction = Position{-1, 0}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && g.direction != (Position{-1, 0}) {
		g.direction = Position{1, 0}
	}

	g.tickCount++
	if g.tickCount < updateDelay {
		return nil
	}
	g.tickCount = 0

	head := Position{g.snake[0].X + g.direction.X, g.snake[0].Y + g.direction.Y}

	if head.X < 0 || head.X >= screenWidth/gridSize || head.Y < 0 || head.Y >= screenHeight/gridSize {
		g.gameOver = true
		return nil
	}

	for _, s := range g.snake[1:] {
		if head == s {
			g.gameOver = true
			return nil
		}
	}

	g.snake = append([]Position{head}, g.snake...)
	if head == g.food {
		g.score++
		g.food = g.generateFood()
	} else {
		g.snake = g.snake[:len(g.snake)-1]
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	for _, s := range g.snake {
		vector.DrawFilledRect(screen, float32(s.X*gridSize), float32(s.Y*gridSize), gridSize, gridSize, color.RGBA{0, 255, 0, 255}, false)
	}

	vector.DrawFilledRect(screen, float32(g.food.X*gridSize), float32(g.food.Y*gridSize), gridSize, gridSize, color.RGBA{255, 0, 0, 255}, false)

	if g.gameOver {
		msg := "Game Over! Press R to Restart"

		text.Draw(screen, msg, basicfont.Face7x13, 10, 10, color.White)
	} else {
		msg := "Score: " + strconv.Itoa(g.score)
		text.Draw(screen, msg, basicfont.Face7x13, 10, 10, color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) generateFood() Position {
	rand.Seed(time.Now().UnixNano())
	return Position{
		X: rand.Intn(screenWidth / gridSize),
		Y: rand.Intn(screenHeight / gridSize),
	}
}

func (g *Game) initGame() {
	g.snake = []Position{{X: screenWidth / (2 * gridSize), Y: screenHeight / (2 * gridSize)}}
	g.direction = Position{1, 0}
	g.food = g.generateFood()
	g.gameOver = false
	g.score = 0
}

func main() {
	game := &Game{}
	game.initGame()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake Game")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
