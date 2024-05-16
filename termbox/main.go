package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	width  = 20
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
}

func placeFood() {
	rand.Seed(time.Now().UnixNano())
	food = Point{X: rand.Intn(width), Y: rand.Intn(height)}
}

func main() {
	initGame()
	defer termbox.Close()
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

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

	for !gameOver {
		select {
		case <-ticker.C:
			update()
			draw()
		}
	}
}

func update() {
	head := snake.Body[0]
	newHead := Point{X: head.X + snake.Dir.X, Y: head.Y + snake.Dir.Y}

	if newHead.X < 0 || newHead.X >= width || newHead.Y < 0 || newHead.Y >= height {
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
		placeFood()
	} else {
		snake.Body = snake.Body[:len(snake.Body)-1]
	}
}

func draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for _, p := range snake.Body {
		termbox.SetCell(p.X, p.Y, 'O', termbox.ColorGreen, termbox.ColorDefault)
	}
	termbox.SetCell(food.X, food.Y, 'X', termbox.ColorRed, termbox.ColorDefault)

	termbox.Flush()
}
