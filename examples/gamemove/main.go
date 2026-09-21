// Command gamemove picks the next move in a toy Snake game, following the
// same pattern as the dozens of Jev-plays-a-game demos in the community
// (jev-snake, jev-tetris, Jev Pac-Man, and others): the game engine owns
// the rules and state, and one Choice question per tick picks the move.
//
// See https://github.com/hellogumbo/awesome-jev for the real games this
// pattern is drawn from.
package main

import (
	"context"
	"fmt"
	"log"

	typesafe "github.com/atharvamhaske/typesafe-sdk-go"
)

func main() {
	client, err := typesafe.NewClient() // reads TYPESAFE_API_KEY
	if err != nil {
		log.Fatal(err)
	}

	state := `Snake board, 10x10 grid. Head at (4,5), moving right.
Body: [(3,5), (2,5)].
Food at (7,5).
Wall or body directly right of head: no.`

	resp, err := client.SystemOne(context.Background(), state,
		map[string]typesafe.Question{
			"move": typesafe.Choice{
				Instructions: "Pick the next move that avoids collisions and moves toward the food.",
				Criteria: map[string]string{
					"up":    "Move up",
					"down":  "Move down",
					"left":  "Move left",
					"right": "Move right",
				},
			},
		})
	if err != nil {
		log.Fatal(err)
	}

	move := resp.Choices()["move"]
	fmt.Printf("move: %s (confidence %.2f)\n", move.Choice, move.Confidence)
}
