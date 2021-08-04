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
	"math"
	"time"

	mt "github.com/brettbuddin/musictheory"
)

type Challenge struct {
	Pitches []mt.Pitch
	Bounty  int
}

func pName(p mt.Pitch) string {
	return p.Name(mt.AscNames)
}

func playChallenge(c Challenge) {
	var t = tempo * int(time.Millisecond)
	for i, p := range c.Pitches {
		if i == 0 {
			fmt.Println(pName(p))
		} else {
			fmt.Println("?")
		}

		synthNote(midiNoteOn, p)
		time.Sleep(time.Duration(t / 2))
		synthNote(midiNoteOff, p)
		time.Sleep(time.Duration(t / 2))
	}

	fmt.Println("* Tap", pName(c.Pitches[0]), "twice to replay the sequence (with tempo)")
}

func newChallenge() Challenge {
	c := Challenge{Pitches: make([]mt.Pitch, 2), Bounty: 100}

	c.Pitches[0] = mt.NewPitch(mt.C, mt.Natural, 3)
	c.Pitches[1] = mt.NewPitch(mt.E, mt.Natural, 3)

	return c
}

var noPitch mt.Pitch
var challenge Challenge
var lastNoteTime time.Time
var lastPitch mt.Pitch
var tempo int = 1000

func prepareChallenge() {
	challenge = newChallenge()

	fmt.Println("Your challenge: ")
	time.Sleep(300 * time.Millisecond)
	playChallenge(challenge)

	lastPitch = noPitch
	lastNoteTime = time.Now()
}

func gameLoop() {
	noPitch, _ = mt.ParsePitch("C12")

	fmt.Println("** Eargo game is ready! **")
	fmt.Println("")

	prepareChallenge()

	for {
		select {
		case i := <-keyboardInput:
			fmt.Printf("Key: %d\n", i)
		case e := <-noteEvents:
			semitones := int(e.Data1) - 24
			pitch := mt.Pitch{Interval: mt.Semitones(semitones)}
			// fmt.Printf("Note %#x: %d\n", e.Status, e.Data1)
			go synthNote(e.Status, pitch)
			if e.Status == midiNoteOff {
				// fmt.Println("challenge:", challenge)
				// fmt.Println("last pitch: ", pName(lastPitch))
				// fmt.Println("this pitch: ", pName(pitch))

				if lastPitch == challenge.Pitches[0] {
					if pitch == challenge.Pitches[0] {
						timeDiff := time.Since(lastNoteTime)
						tempo = int(math.Min(float64(timeDiff.Milliseconds()), 2000))
						time.Sleep(500 * time.Millisecond)

						fmt.Printf("* Your challenge:\n")
						playChallenge(challenge)

						lastPitch = noPitch
						lastNoteTime = time.Now()

						continue
					} else if pitch == challenge.Pitches[1] {
						fmt.Println("* Yes! You win ", challenge.Bounty, "points!")
						time.Sleep(400 * time.Millisecond)
						prepareChallenge()

						continue
					} else {
						fmt.Println("Nope")
					}
				}

				lastPitch = pitch
				lastNoteTime = time.Now()
			}

		}
	}
}
