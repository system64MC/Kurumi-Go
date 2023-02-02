package Operator

import (
	"math"
	// glob "system64.net/KurumiGo/model/Globals"
	m "system64.net/KurumiGo"
)

type Operator struct {
	Tl              float32
	WaveformId      int32
	Mult            int32
	Phase           float32
	Detune          int32
	PhaseMod        bool
	PhaseRev        bool
	Feedback        float32
	Prev            float32
	UseCustomVolEnv bool
	VolEnv          []uint8
	PhaseEnv        []uint8
	CustomPhaseEnv  bool
	Interpolation   uint8
}

var WaveFuncs = []func(*Operator, float64) float64{
	func(op *Operator, x float64) float64 {
		return math.Sin(x*float64(op.Mult)*2*math.Pi) + (float64(op.Phase)*2*math.Pi + (op.getPhase() * math.Pi * 2))
	},
}

func (op *Operator) getPhase() float64 {
	myPhaseMod := 0.0
	if op.PhaseMod {
		myPhaseMod = 1.0
	}
	macro := m.Context.Macro
	macLen := m.Context.MacLen
	pRev := 1.0
	if op.PhaseRev {
		pRev = -1.0
	}
	// Put custom phase env here
	return (float64(macro) / float64(macLen-1)) * pRev * float64(op.Detune) * float64(myPhaseMod)
}

var Waveforms = []string{
	"Sine",
	"Rect. Sine",
	"Abs. Sine",
	"Quarter Sine",
	"Squished Sine",
	"Squished Abs. Sine",
	"Square",
	"Saw",
	"Rect. Saw",
	"Abs. Saw",
	"Cubed Saw",
	"Rect. Cubed Saw",
	"Abs. Cubed Saw",
	"Cubed Sine",
	"Rect. Cubed Sine",
	"Abs. Cubed Sine",
	"Quarter Cubed Sine",
	"Squished Cubed Sine",
	"Squi. Abs. Cubed Sine",
	"Triangle",
	"Rect. Triange",
	"Abs. Triangle",
	"Quarter Triangle",
	"Squished Triangle",
	"Abs. Squished Triangle",
	"Cubed Triangle",
	"Rect. Cubed Triangle",
	"Abs. Cubed Triangle",
	"Quarter Cubed Triangle",
	"Squi. Cubed Triangle",
	"Squi. Abs. Cubed Triangle",
	"Custom",
}
