/* BRIEF DESIGN:

 _______________________________________
|  | | | |  |  | | | | | |  |  | | | |  |
|  | | | |  |  | | | | | |  |  | | | |  |
|  | | | |  |  | | | | | |  |  | | | |  |
|  |_| |_|  |  |_| |_| |_|  |  |_| |_|  |
|   |   |   |   |   |   |   |   |   |   |
|   |   |   |   |   |   |   |   |   |   |
|___|___|___|___|___|___|___|___|___|___|
ASCII Art: Alexander Craxton
https://asciiart.website/index.php?art=music/pianos


The Eargo Game
	1. adjust difficulty. print bounty
	2. print and play root note
	3. play rest of the sequence


	info statistics for 5, 10, 50 recent attemps
		average tries to succeed
		average bounty


	difficulty/score multipliers:
		phraseLengthScores := int[]{0, 1, 2, 3, 5, 8, 10, 20, 30, 50, 80, 100, 200}
		octaves range
			order C3, C2, C4, C5, C1, C6, C7
			score 1   2   4   6   8   10  12
		scales
			n of total 12 major/minor scales. bounty = max(1, (3 * n-1))
			major scale for now
		intervals
		 	up + down: 10x multiplier
			ranked multiplier sums:  O(1), P5(2), 3rd(2), 2nd(5), 6th(8), 7th(10), b3rd(20), b2nd(30)


	total score:
	 	current_score += attempt value
		attempt values = 1, 0.5, 0.1, 0.1,...

	progressive difficulty modes
		a) adjust it yourself
		b) the game does it for you and tracks how you go

	controls:
		root + down octave = next one please. TODO: and make it easier
		root + up octave = next one and TODO: make it harder
		top root twice = playback at this tempo

Next steps todo:
	1. generate sequence at fixed difficulty
	2. calculate bounty
	3. game loop with controls, scoring and output
*/

package main

import (
	"fmt"

	mt "github.com/brettbuddin/musictheory"
	// "github.com/brettbuddin/musictheory/intervals"
)

type Challenge struct {
	Pitches []mt.Pitch
	Bounty  int
}

func playChallenge(c Challenge) {
	fmt.Println("Your challenge is (beepty, boopty): ", c)
}

func newChallenge() Challenge {
	c := Challenge{Pitches: make([]mt.Pitch, 2), Bounty: 100}

	c.Pitches[0] = mt.NewPitch(mt.C, mt.Natural, 4)
	c.Pitches[1] = mt.NewPitch(mt.E, mt.Natural, 4)

	playChallenge(c)
	return c
}

func gameLoop() {
	fmt.Println("** Eago game is ready! **")
	// fmt.Println("** Press ESC to quit **")
	fmt.Println("")
	challenge := newChallenge()

	for {
		select {
		case i := <-keyboardInput:
			fmt.Printf("Key: %d\n", i)
		case e := <-noteEvents:
			semitones := int(e.Data1) - 24
			pitch := mt.Pitch{Interval: mt.Semitones(semitones)}
			// fmt.Printf("Note %#x: %d\n", e.Status, e.Data1)
			go synthNote(e.Status, pitch)
			fmt.Println("almost there:", challenge)
		}
	}
}
