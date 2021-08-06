package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	mt "github.com/brettbuddin/musictheory"
)

// C12
var blankPitch = mt.Pitch{Interval: mt.Interval{12, 0, 0}}

type Challenge struct {
	Pitches []mt.Pitch
	Bounty  int
	Tries   int
}

type GameState struct {
	LastNoteTime time.Time
	LastPitch    mt.Pitch
	Tempo        int
	Challenge    Challenge
	Score        int
}

var pitches = []string{"E2", "G2", "E3", "G3"}

func (gs *GameState) PrepareChallenge() {
	c := Challenge{Pitches: make([]mt.Pitch, 2), Bounty: 100, Tries: 0}

	c.Pitches[0] = mt.NewPitch(mt.C, mt.Natural, 3)

	pickedPitch := pitches[rand.Intn(len(pitches))]
	c.Pitches[1] = mt.MustParsePitch(pickedPitch)

	gs.Challenge = c
}

func (gs *GameState) BeginChallenge() {
	go gs.PlayChallenge()

	gs.SetLastPitch(blankPitch)
}

func (gs GameState) PlayChallenge() {
	var t = gs.Tempo * int(time.Millisecond)
	var c = gs.Challenge

	fmt.Printf("> Challenge:   |   ")
	time.Sleep(500 * time.Millisecond)
	for i, p := range c.Pitches {
		if i == 0 {
			fmt.Printf("%s  ", pName(p))
		} else {
			fmt.Printf("?   ")
		}

		synthNote(midiNoteOn, p)
		time.Sleep(time.Duration(t / 2))
		synthNote(midiNoteOff, p)
		time.Sleep(time.Duration(t / 2))
	}

	fmt.Printf("|\tBounty: %4d", gs.Challenge.Bounty)
	fmt.Println("\t\t (", pName(c.Pitches[0]), pName(c.Pitches[0]), "to replay)")
}

func (gs GameState) PlayWinning() {
	// revPitches := reverseCopy(pitches)

	time.Sleep(100 * time.Millisecond)
	for i := 0; i < 6; i++ {
		synthNote(midiNoteOn, gs.SecondPitch())
		time.Sleep(time.Duration((60 - 10*i)) * time.Millisecond)
		synthNote(midiNoteOff, gs.SecondPitch())
		time.Sleep(20 * time.Millisecond)
	}
}

func (gs *GameState) Winning() {
	gs.Score += gs.Challenge.Bounty
	fmt.Printf("\n+%d points! [SCORE: %4d]\n", int(gs.Challenge.Bounty), gs.Score)
	fmt.Println("\n. . . . . . . . . . . . . . . . . . . . ")
	fmt.Println("")
	gs.PlayWinning()
}

func (gs *GameState) Nope() {
	gs.Challenge.Tries += 1
	gs.Challenge.Bounty /= 2

	fmt.Println("No. Bounty: ", gs.Challenge.Bounty, "(", gs.Challenge.Tries, "tries )")

	// time.Sleep(200 * time.Millisecond)
	// synthNote(midiNoteOn, gs.FirstPitch())
	// time.Sleep(1 * time.Millisecond)
	// synthNote(midiNoteOff, gs.FirstPitch())
	// time.Sleep(5 * time.Millisecond)
	// synthNote(midiNoteOn, gs.FirstPitch())
	// time.Sleep(1 * time.Millisecond)
	// synthNote(midiNoteOff, gs.FirstPitch())
	// time.Sleep(500 * time.Millisecond)

	// go gs.PlayChallenge()
}

func (gs GameState) FirstPitch() mt.Pitch {
	return gs.Challenge.Pitches[0]
}

func (gs GameState) SecondPitch() mt.Pitch {
	return gs.Challenge.Pitches[1]
}

func (gs *GameState) SetLastPitch(pitch mt.Pitch) {
	gs.LastPitch = pitch
	gs.LastNoteTime = time.Now()
}

func (gs *GameState) SetTappedTempo() {
	timeDiff := time.Since(gs.LastNoteTime)
	gs.Tempo = int(math.Min(float64(timeDiff.Milliseconds()), 2000))
}

func NewGameState() GameState {
	gs := GameState{Tempo: 1000}

	gs.PrepareChallenge()

	return gs
}
