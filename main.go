package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	SCREEN_WIDTH  = 800
	SCREEN_HEIGHT = 650
)

const (
	G = 1.25

	PIPE_WIDTH       = 100
	PIPE_SPEED       = 2
	PIPE_GAP_SIZE    = 150
	PIPE_SPAWN_SPEED = 5

	PLAYER_SPEED = 1
	PLAYER_SIZE  = 25
	PLAYER_JUMP  = 30
)

type Pipe struct {
	X int32
	Y int32
}

type Player struct {
	Position rl.Vector2
	Angle    int32
}

var (
	PLAYER Player
	Pipes  = []Pipe{}
)

var has_lost bool = false

func InitGame() {
	has_lost = false

	PLAYER = Player{
		Position: rl.NewVector2(SCREEN_WIDTH/2, SCREEN_HEIGHT/2),
		Angle:    0,
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
	}
}

var spawn_frame_counter = -1

func UpdateGame() {
	if !has_lost {
		spawn_frame_counter++

		if spawn_frame_counter%(1000/PIPE_SPAWN_SPEED) == 0 {
			new_pipe := Pipe{
				X: SCREEN_WIDTH,
				Y: rl.GetRandomValue(PIPE_GAP_SIZE+PLAYER_SIZE, SCREEN_HEIGHT-PIPE_GAP_SIZE-PLAYER_SIZE),
			}
			Pipes = append(Pipes, new_pipe)

			spawn_frame_counter = 0
		}

		if rl.IsKeyPressed(rl.KeySpace) {
			PLAYER.Position.Y -= PLAYER_JUMP
		} else {
			PLAYER.Position.Y += G
		}

		for i := 0; i < len(Pipes); i++ {
			pipe := &Pipes[i]
			pipe.X -= PIPE_SPEED

			r := float32(PLAYER_SIZE / 2)

			if pipe.X+PIPE_WIDTH < 0 {
				Pipes = append(Pipes[:i], Pipes[i+1:]...)
			} else if PLAYER.Position.X+r > float32(pipe.X) && PLAYER.Position.X-r < float32(pipe.X+PIPE_WIDTH) {
				if PLAYER.Position.Y-r < float32(pipe.Y) || PLAYER.Position.Y+r > float32(pipe.Y+PIPE_GAP_SIZE) {
					has_lost = true
					break
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
