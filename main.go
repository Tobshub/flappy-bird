package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SCREEN_WIDTH  = 800
	SCREEN_HEIGHT = 650
)

const (
	PIPE_WIDTH       = 100
	PIPE_SPEED       = 4
	PIPE_GAP_SIZE    = 150
	PIPE_SPAWN_SPEED = 3

	G = .1

	PLAYER_SIZE = 25
	PLAYER_JUMP = 35
)

type Pipe struct {
	X int32
	Y int32

	Passed bool
}

type Player struct {
	Position   rl.Vector2
	Fall_Speed float32
}

var (
	PLAYER Player
	Pipes  = []Pipe{}
)

var (
	has_lost bool = false
	score         = 0
)

func InitGame() {
	has_lost = false
	score = 0

	PLAYER = Player{
		Position:   rl.NewVector2(SCREEN_WIDTH/2, SCREEN_HEIGHT/2),
		Fall_Speed: 0,
	}

	Pipes = []Pipe{}
}

func DrawGame() {
	rl.DrawCircle(int32(PLAYER.Position.X), int32(PLAYER.Position.Y), PLAYER_SIZE, rl.Red)

	for _, pipe := range Pipes {
		rl.DrawRectangle(pipe.X, 0, PIPE_WIDTH, pipe.Y, rl.DarkGray)                                                // top
		rl.DrawRectangle(pipe.X, pipe.Y+PIPE_GAP_SIZE, PIPE_WIDTH, SCREEN_HEIGHT-pipe.Y-PIPE_GAP_SIZE, rl.DarkGray) // bottom
	}

	if has_lost {
		game_over_text, font_size := "GAME OVER", int32(32)
		rl.DrawText(game_over_text, SCREEN_WIDTH/2-rl.MeasureText(game_over_text, font_size)/2, SCREEN_HEIGHT/2-16, font_size, rl.Red)
		game_over_instructions, font_size := "Press Enter to restart", int32(24)
		rl.DrawText(game_over_instructions, SCREEN_WIDTH/2-rl.MeasureText(game_over_instructions, font_size)/2, SCREEN_HEIGHT/2+16, font_size, rl.Gray)
	}

	rl.DrawText(fmt.Sprintf("Score: %d", score), 10, 10, 20, rl.DarkGreen)
}

var spawn_frame_counter = -1

func UpdateGame() {
	if !has_lost {
		spawn_frame_counter++

		if spawn_frame_counter%(1000/int(PIPE_SPAWN_SPEED*float32(PIPE_SPEED))) == 0 {
			new_pipe := Pipe{
				X:      SCREEN_WIDTH,
				Y:      rl.GetRandomValue(PIPE_GAP_SIZE+PLAYER_SIZE, SCREEN_HEIGHT-PIPE_GAP_SIZE-PLAYER_SIZE),
				Passed: false,
			}
			Pipes = append(Pipes, new_pipe)

			spawn_frame_counter = 0
		}

		if rl.IsKeyPressed(rl.KeySpace) {
			PLAYER.Position.Y -= PLAYER_JUMP
			PLAYER.Fall_Speed = 0
		} else {
			PLAYER.Position.Y += PLAYER.Fall_Speed
			PLAYER.Fall_Speed += G
		}

		for i := 0; i < len(Pipes); i++ {
			pipe := &Pipes[i]
			pipe.X -= PIPE_SPEED

			r := float32(PLAYER_SIZE - 4) // allowance

			if pipe.X+PIPE_WIDTH < 0 {
				Pipes = append(Pipes[:i], Pipes[i+1:]...)
			} else if PLAYER.Position.X+r > float32(pipe.X) && PLAYER.Position.X-r < float32(pipe.X+PIPE_WIDTH) {
				if PLAYER.Position.Y-r < float32(pipe.Y) || PLAYER.Position.Y+r > float32(pipe.Y+PIPE_GAP_SIZE) {
					has_lost = true
					break
				} else if !pipe.Passed {
					score++
					pipe.Passed = true
				}
			}
		}
	} else {
		if rl.IsKeyPressed(rl.KeyEnter) {
			InitGame()
		}
	}
}

func main() {
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Flappy Bird")
	rl.SetTargetFPS(60)

	InitGame()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		DrawGame()
		UpdateGame()

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
