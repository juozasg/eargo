package main

import (
	"fmt"
	"math"
	"time"

	mt "github.com/brettbuddin/musictheory"
)

// C12
var blankPitch = mt.Pitch{Interval: mt.Interval{12, 0, 0}}

type Challenge struct {
	Pitches []mt.Pitch
	Bounty  int
}

type GameState struct {
	LastNoteTime time.Time
	LastPitch    mt.Pitch
	Tempo        int
	Challenge    Challenge
}

func (gs *GameState) PrepareChallenge() {
	c := Challenge{Pitches: make([]mt.Pitch, 2), Bounty: 100}

	c.Pitches[0] = mt.NewPitch(mt.C, mt.Natural, 3)
	c.Pitches[1] = mt.NewPitch(mt.E, mt.Natural, 3)

	gs.Challenge = c
}

func (gs GameState) BeginChallenge() {
	fmt.Printf("* Your challenge:\n")
	gs.PlayChallenge()

	gs.SetLastPitch(blankPitch)
}

func (gs GameState) PlayChallenge() {
	var t = gs.Tempo * int(time.Millisecond)
	var c = gs.Challenge

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
