package Synth

import "system64.net/KurumiGo/model/Operator"

type Synth struct {
	WaveLen   int
	WaveHei   int
	MacLen    int
	Macro     int
	Operators []Operator.Operator
	ModMatrix [][]bool
	SmoothWin int
	Gain      float64
}

func ConstructSynth() *Synth {
	context := &Synth{WaveLen: 32, WaveHei: 32, MacLen: 64, Macro: 0, SmoothWin: 1, Gain: 1.0}
	context.ModMatrix = [][]bool{
		{false, false, false, false},
		{true, false, false, false},
		{false, true, false, false},
		{false, false, true, false},
	}
	for i := 0; i < 4; i++ {
		if i == 3 {
			context.Operators = append(context.Operators,
				Operator.Operator{
					Tl:              1,
					WaveformId:      0,
					Mult:            1,
					Phase:           0,
					Detune:          1,
					PhaseMod:        false,
					Feedback:        0,
					Prev:            0,
					UseCustomVolEnv: false,
					Interpolation:   0})
		} else {
			context.Operators = append(context.Operators,
				Operator.Operator{
					Tl:              0,
					WaveformId:      0,
					Mult:            1,
					Phase:           0,
					Detune:          1,
					PhaseMod:        false,
					Feedback:        0,
					Prev:            0,
					UseCustomVolEnv: false,
					Interpolation:   0})
		}
	}
	return context
}
