package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	snakeLegth   = 256
	squareSize   = 31
	screenWidth  = 800
	screenHeight = 450
)

var (
	darkBlue  = color.RGBA{0, 82, 172, 255}
	blue      = color.RGBA{0, 121, 241, 255}
	skyBlue   = color.RGBA{102, 191, 255, 255}
	rayWhite  = color.RGBA{245, 245, 245, 255}
	lightGray = color.RGBA{200, 200, 200, 255}
)

type Vector2 struct {
	X float32
	Y float32
}

type Snake struct {
	position Vector2
	size     Vector2
	speed    Vector2
	color    color.RGBA
}

type Food struct {
	position Vector2
	size     Vector2
	active   bool
	color    color.RGBA
}

type Game struct {
	framesCounter int
	gameOver      bool
	pause         bool
	fruit         Food
	snake         [snakeLegth]Snake
	snakePosition [snakeLegth]Vector2
	allowMove     bool
	offset        Vector2
	counterTail   int
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("classic game: snake")
	game := &Game{}
	game.initGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) initGame() {
	g.framesCounter = 0
	g.gameOver = false
	g.pause = false
	g.fruit = Food{}
	g.snake = [snakeLegth]Snake{}
	g.snakePosition = [snakeLegth]Vector2{}
	g.allowMove = false
	g.offset = Vector2{}
	g.counterTail = 1

	g.offset.X = screenHeight % squareSize
	g.offset.Y = screenHeight % squareSize

	for i := 0; i < snakeLegth; i++ {
		g.snake[i].position = Vector2{X: g.offset.X / 2, Y: g.offset.Y / 2}
		g.snake[i].size = Vector2{X: squareSize, Y: squareSize}
		g.snake[i].speed = Vector2{X: squareSize}

		if i == 0 {
			g.snake[i].color = darkBlue
		} else {
			g.snake[i].color = blue
		}
	}

	for i := 0; i < snakeLegth; i++ {
		g.snakePosition[i] = Vector2{}
	}

	g.fruit.size = Vector2{X: squareSize, Y: squareSize}
	g.fruit.color = skyBlue
	g.fruit.active = false
}

func (g *Game) Update() error {
	if !g.gameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			g.pause = !g.pause
		}

		if !g.pause {
			// Player control
			if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && (g.snake[0].speed.X == 0) && g.allowMove {
				g.snake[0].speed = Vector2{X: squareSize}
				g.allowMove = false
			}
			if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && (g.snake[0].speed.X == 0) && g.allowMove {
				g.snake[0].speed = Vector2{X: -squareSize}
				g.allowMove = false
			}
			if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && (g.snake[0].speed.Y == 0) && g.allowMove {
				g.snake[0].speed = Vector2{Y: -squareSize}
				g.allowMove = false
			}
			if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && (g.snake[0].speed.Y == 0) && g.allowMove {
				g.snake[0].speed = Vector2{Y: squareSize}
				g.allowMove = false
			}

			// Snake movement
			for i := 0; i < g.counterTail; i++ {
				g.snakePosition[i] = g.snake[i].position
			}

			// The snake moves 1 squareSize every 5 frames
			if (g.framesCounter % 5) == 0 {
				// Snake movement logic
				for i := 0; i < g.counterTail; i++ {
					if i == 0 { // Head of the snake
						g.snake[0].position.X += g.snake[0].speed.X
						g.snake[0].position.Y += g.snake[0].speed.Y
						g.allowMove = true
					} else { // Body of the snake
						g.snake[i].position = g.snakePosition[i-1]
					}
				}
			}

			// Wall behaviour
			if ((g.snake[0].position.X) > (screenWidth - g.offset.X)) || ((g.snake[0].position.Y) > (screenHeight - g.offset.Y)) || (g.snake[0].position.X < 0) || (g.snake[0].position.Y < 0) {
				g.gameOver = true
				// go from the opposite side
				// if g.snake[0].position.X < 0 { // left
				// 	g.snake[0].position.X = screenWidth - g.offset.X // right
				// } else if g.snake[0].position.X > screenWidth-g.offset.X { // right
				// 	g.snake[0].position.X = 0
				// } else if g.snake[0].position.Y < 0 { // top
				// 	g.snake[0].position.Y = screenHeight - g.offset.Y // bottom
				// } else if g.snake[0].position.Y > screenHeight-g.offset.Y { // bottom
				// 	g.snake[0].position.Y = 0 // top
				// }
			}

			// Collision with yourself
			for i := 1; i < g.counterTail; i++ {
				if (g.snake[0].position.X == g.snake[i].position.X) && (g.snake[0].position.Y == g.snake[i].position.Y) {
					g.gameOver = true
				}
			}

			// Fruit position calculation
			if !g.fruit.active {
				g.fruit.active = true
				g.fruit.position = Vector2{
					X: float32(rand.Intn((screenWidth/squareSize)-1))*squareSize + (g.offset.X / 2),
					Y: float32(rand.Intn((screenHeight/squareSize)-1))*squareSize + (g.offset.Y / 2),
				}

				for i := 0; i < g.counterTail; i++ {
					for (g.fruit.position.X == g.snake[i].position.X) && (g.fruit.position.Y == g.snake[i].position.Y) {
						g.fruit.position = Vector2{
							X: float32(rand.Intn((screenWidth/squareSize)-1))*squareSize + (g.offset.X / 2),
							Y: float32(rand.Intn((screenHeight/squareSize)-1))*squareSize + (g.offset.Y / 2),
						}
						i = 0
					}
				}
			}

			// Collision
			if (g.snake[0].position.X < (g.fruit.position.X+g.fruit.size.X) && (g.snake[0].position.X+g.snake[0].size.X) > g.fruit.position.X) && (g.snake[0].position.Y < (g.fruit.position.Y+g.fruit.size.Y) && (g.snake[0].position.Y+g.snake[0].size.Y) > g.fruit.position.Y) {
				g.snake[g.counterTail].position = g.snakePosition[g.counterTail-1]
				g.counterTail++
				g.fruit.active = false
			}

			g.framesCounter++
		}
	} else {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.initGame()
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(rayWhite)

	if !g.gameOver {
		// Draw grid lines
		for i := 0; i < screenWidth/squareSize+1; i++ {
			vector.StrokeLine(screen,
				float32(squareSize*i)+g.offset.X/2, g.offset.Y/2,
				float32(squareSize*i)+g.offset.X/2, screenHeight-g.offset.Y/2,
				1, lightGray, false)
		}

		for i := 0; i < screenHeight/squareSize+1; i++ {
			vector.StrokeLine(screen,
				g.offset.X/2, float32(squareSize*i)+g.offset.Y/2,
				screenWidth-g.offset.X/2, float32(squareSize*i)+g.offset.Y/2,
				1, lightGray, false)
		}

		// Draw snake
		for i := 0; i < g.counterTail; i++ {
			vector.DrawFilledRect(screen,
				g.snake[i].position.X, g.snake[i].position.Y,
				g.snake[i].size.X, g.snake[i].size.Y,
				g.snake[i].color, false)
		}

		// Draw fruit to pick
		vector.DrawFilledRect(screen,
			g.fruit.position.X, g.fruit.position.Y,
			g.fruit.size.X, g.fruit.size.Y,
			g.fruit.color, false)

		if g.pause {
			ebitenutil.DebugPrint(screen, "GAME PAUSED")
		}

	} else {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("GAME OVER\n\nSCORE: %d\n\nPRESS [ENTER] TO PLAY AGAIN", g.counterTail-1))
	}
}
