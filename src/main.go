/*
 * Tetris
 */

package main

import (
	"math/rand"
	"time"

	"github.com/mattn/go-tty"
)

type point struct {
	x, y int
}

type block struct {
	points []point
}

type game struct {
	current    block
	h, w       int
	lines      int
	level      int
	matrix     [][]int
	over       bool
	pos        point
	score      int
	tetrominos []block
}

func (g *game) input() {
	tty, _ := tty.Open()

	input, _ := tty.ReadRune()
	switch input {
	case 'w', 'j':
		g.current.rotate()
	case 'a', 'h':
		g.pos.x--
		if g.collides() {
			g.pos.x++
		}
	case 's', 'k':
		g.pos.y++
		if g.collides() {
			g.pos.y--
		}
	case 'd', 'l':
		g.pos.x++
		if g.collides() {
			g.pos.x--
		}
	case 'q':
		g.over = true
	}
	tty.Close()
}

func (b *block) rotate() {
	// Calculate the center point
	// of the block.
	cx, cy := 0, 0
	for _, p := range b.points {
		cx += p.x
		cy += p.y
	}
	cx /= len(b.points)
	cy /= len(b.points)

	// Rotate the block clockwise
	// around the center point.
	for i, p := range b.points {
		b.points[i] = point{cy - p.y + cx, p.x - cx + cy}
	}
}

func (g *game) update() {
	if g.over {
		return
	}

	// Increase the level every 10 lines.
	if g.lines/10 > g.level {
		g.level++
	}

	// Calculate the speed based on the level.
	speed := 1000 - g.level*100

	// Move the current piece down
	g.pos.y++

	// Check for collisions
	if g.collides() {
		// Move the piece back up and
		// add it to the game board.
		g.pos.y--
		g.add()

		// Choose a new piece to start with.
		g.current = g.tetrominos[rand.Intn(len(g.tetrominos))]

		// Reset the position, placing the
		// new piece at the top of the board
		// and centered horizontally.
		g.pos = point{len(g.matrix[0]) / 2, 0}

		// Check for game over by calling
		// collides() again. If the new piece
		// collides with the game board, the
		// game is over.
		if g.collides() {
			g.over = true
		}
	}

	time.Sleep(time.Nanosecond * time.Duration(speed))

	// Clear any full rows
	g.clear()
}

func (g *game) draw() {
	// Clear the screen.
	print("\033[H\033[2J")

	sb := ""

	// Draw the background, and the
	// pieces that have already fallen.
	for _, row := range g.matrix {
		for _, col := range row {
			if col == 0 {
				sb += "."
			} else {
				sb += "#"
			}
		}
		sb += "\n"
	}

	// Draw the current shape.
	for _, p := range g.current.points {
		x := g.pos.x + p.x
		y := g.pos.y + p.y
		sb = sb[:y*(len(g.matrix[0])+1)+x] + "#" + sb[y*(len(g.matrix[0])+1)+x+1:]
	}

	// Print the game.
	println(sb)

	// Print the score.
	println("Score:", g.score)

	// Print the level.
	println("Level:", g.level)

	// Print game controls.
	println("asd to move, w to rotate, q to quit")
}

func (g *game) collides() bool {
	for _, p := range g.current.points {
		x := g.pos.x + p.x
		y := g.pos.y + p.y
		// Check for collisions with the game board and other pieces.
		if y >= len(g.matrix) || x < 0 || x >= len(g.matrix[0]) || (y >= 0 && g.matrix[y][x] != 0) {
			return true
		}
	}
	return false
}

func (g *game) add() {
	for _, p := range g.current.points {
		x := g.pos.x + p.x
		y := g.pos.y + p.y
		if y >= 0 {
			g.matrix[y][x] = 1
		}
	}
}

func (g *game) clear() {
	// Count the number of full rows.
	cleared := 0

	// Check for full rows.
	for y := 0; y < len(g.matrix); y++ {
		full := true
		for x := 0; x < len(g.matrix[y]); x++ {
			if g.matrix[y][x] == 0 {
				full = false
				break
			}
		}
		// If the row is full, remove it and
		// add a new empty row at the top.
		if full {
			copy(g.matrix[1:], g.matrix[:y])
			g.matrix[0] = make([]int, len(g.matrix[0]))
			cleared++
		}
	}

	// Update the score and level.
	if cleared > 0 {
		g.lines += cleared
		g.score += 100 << cleared
	}
}

func (g *game) init(h, w int) {
	// Initialize the game board to
	// the height and width provided.
	g.matrix = make([][]int, h)
	for i := range g.matrix {
		g.matrix[i] = make([]int, w)
	}

	// Choose a random piece to start with.
	g.current = g.tetrominos[rand.Intn(len(g.tetrominos))]
	// Reset the position, placing the
	// new piece at the top of the board.
	g.pos = point{w / 2, 0}
}

func main() {
	game := new(game)
	game.h = 20
	game.w = 10

	game.tetrominos = []block{
		{[]point{{0, 0}, {0, 1}, {0, 2}, {0, 3}}},
		{[]point{{0, 0}, {0, 1}, {1, 0}, {1, 1}}},
		{[]point{{0, 0}, {0, 1}, {0, 2}, {1, 2}}},
		{[]point{{0, 0}, {0, 1}, {0, 2}, {-1, 2}}},
		{[]point{{0, 0}, {0, 1}, {1, 1}, {1, 2}}},
		{[]point{{0, 0}, {0, 1}, {-1, 1}, {-1, 2}}},
		{[]point{{0, 0}, {0, 1}, {0, 2}, {-1, 1}}},
	}

	// Initialize the game.
	game.init(game.h, game.w)

	// Run the game loop.
	for !game.over {
		game.draw()
		game.input()
		game.update()
	}

	println("Game over!")
}
