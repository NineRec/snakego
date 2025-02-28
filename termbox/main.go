package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	width  = 40
	height = 20
)

type Point struct {
	X int
	Y int
}

type Snake struct {
	Body []Point
	Dir  Point
}

var (
	snake    Snake
	food     Point
	score    int
	gameOver bool
)

func initGame() {
	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	snake = Snake{
		Body: []Point{{X: width / 2, Y: height / 2}},
		Dir:  Point{X: 1, Y: 0},
	}
	placeFood()
	score = 0
}

func placeFood() {
	rand.Seed(time.Now().UnixNano())
	food = Point{X: rand.Intn(width-2) + 1, Y: rand.Intn(height-2) + 1}
}

func main() {
	initGame()
	defer termbox.Close()

	go func() {
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyArrowUp:
					if snake.Dir.Y == 0 {
						snake.Dir = Point{X: 0, Y: -1}
					}
				case termbox.KeyArrowDown:
					if snake.Dir.Y == 0 {
						snake.Dir = Point{X: 0, Y: 1}
					}
				case termbox.KeyArrowLeft:
					if snake.Dir.X == 0 {
						snake.Dir = Point{X: -1, Y: 0}
					}
				case termbox.KeyArrowRight:
					if snake.Dir.X == 0 {
						snake.Dir = Point{X: 1, Y: 0}
					}
				case termbox.KeyEsc:
					gameOver = true
				}
			}
		}
	}()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for !gameOver {
		<-ticker.C
		update()
		draw()
	}

	// Display final score
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	message := fmt.Sprintf("Game Over! Final Score: %d", score)
	drawText(width/2-len(message)/2, height/2, message)
	termbox.Flush()
	time.Sleep(3 * time.Second)
}

func update() {
	head := snake.Body[0]
	newHead := Point{X: head.X + snake.Dir.X, Y: head.Y + snake.Dir.Y}

	if newHead.X <= 0 || newHead.X >= width-1 || newHead.Y <= 0 || newHead.Y >= height-1 {
		gameOver = true
		return
	}

	for _, p := range snake.Body {
		if p == newHead {
			gameOver = true
			return
		}
	}

	snake.Body = append([]Point{newHead}, snake.Body...)
	if newHead == food {
		score++
		placeFood()
	} else {
		snake.Body = snake.Body[:len(snake.Body)-1]
	}
}

func draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Draw borders
	for x := 0; x < width; x++ {
		termbox.SetCell(x, 0, '_', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x, height-1, '_', termbox.ColorWhite, termbox.ColorDefault)
	}
	for y := 0; y < height; y++ {
		termbox.SetCell(0, y, '|', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(width-1, y, '|', termbox.ColorWhite, termbox.ColorDefault)
	}

	// Draw snake
	for _, p := range snake.Body {
		termbox.SetCell(p.X, p.Y, 'O', termbox.ColorGreen|termbox.AttrBold, termbox.ColorDefault)
	}
	termbox.SetCell(food.X, food.Y, 'x', termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault)

	// Draw score
	scoreText := fmt.Sprintf("Score: %d", score)
	drawText(2, height, scoreText)

	termbox.Flush()
}

func drawText(x, y int, text string) {
	for i, c := range text {
		termbox.SetCell(x+i, y, c, termbox.ColorYellow, termbox.ColorDefault)
	}
}
