package kurumi

type Operator struct {
	Tl              float32
	Reverse         bool
	Adsr            Adsr
	WaveformId      int32
	Mult            int32
	Phase           float32
	Detune          int32
	PhaseMod        bool
	PhaseRev        bool
	Feedback        float32
	Prev            float32
	Curr            float32
	UseCustomVolEnv bool
	VolEnv          []uint8
	PhaseEnv        []uint8
	Wavetable       []uint8
	MorphWave       []uint8
	CustomPhaseEnv  bool
	Interpolation   int32
	Morphing        bool
	MorphTime       int32
	ModMode         int32
	DutyCycle       float32
	PwmAdsr         Adsr
	PwmAdsrEnabled  bool

	IsEnvelopeEnabled bool
}

var CopiedOp *Operator = nil

func CopyOp(op *Operator) {
	CopiedOp = &Operator{}
	*CopiedOp = *op
	volEnv := make([]uint8, len(op.VolEnv))
	phaseEnv := make([]uint8, len(op.PhaseEnv))
	wav := make([]uint8, len(op.Wavetable))
	morph := make([]uint8, len(op.MorphWave))

	copy(volEnv, op.VolEnv)
	copy(phaseEnv, op.PhaseEnv)
	copy(wav, op.Wavetable)
	copy(morph, op.MorphWave)
	CopiedOp.VolEnv = volEnv
	CopiedOp.PhaseEnv = phaseEnv
	CopiedOp.Wavetable = wav
	CopiedOp.MorphWave = morph
}

func PasteOp(op *Operator) {
	*op = *CopiedOp
	volEnv := make([]uint8, len(CopiedOp.VolEnv))
	phaseEnv := make([]uint8, len(CopiedOp.PhaseEnv))
	wav := make([]uint8, len(CopiedOp.Wavetable))
	morph := make([]uint8, len(CopiedOp.MorphWave))

	copy(volEnv, CopiedOp.VolEnv)
	copy(phaseEnv, CopiedOp.PhaseEnv)
	copy(wav, CopiedOp.Wavetable)
	copy(morph, CopiedOp.MorphWave)
	op.VolEnv = volEnv
	op.PhaseEnv = phaseEnv
	op.Wavetable = wav
	op.MorphWave = morph
	Synthesize()
}

func (op *Operator) oscillate(x float64) float64 {
	if op.Reverse {

		return -WaveFuncs[int(op.WaveformId)%len(WaveFuncs)](op, x)
	}
	return WaveFuncs[int(op.WaveformId)%len(WaveFuncs)](op, x)

}

func (op *Operator) getMult() float64 {
	if op.Mult != 0 {
		return float64(op.Mult)
	}
	return 0.5
}

func (op *Operator) getPhase() float64 {
	macro := SynthContext.Macro
	macLen := SynthContext.MacLen

	if !op.CustomPhaseEnv {
		// Anti divide-by-0
		if macLen == 1 {
			macLen = 2
		}
		return (float64(macro) / float64(macLen-1)) * float64(op.Detune)
	}

	if len(op.PhaseEnv) < 1 {
		return 0
	}
	return (float64(op.PhaseEnv[int(Clamp(0, int(macro), len(op.PhaseEnv)-1))]) / 255.0) * float64(op.Detune)
}

func (op *Operator) getFB() float64 {
	return float64(op.Feedback) * (float64(op.Prev) / float64(6*op.getMult()))
}

func (op *Operator) getFB3() float64 {
	return float64(op.Feedback) * (float64(op.Prev) / (op.getMult() / float64(SynthContext.WaveLen*SynthContext.Oversample)))
}

func (op *Operator) getVolume() float64 {
	if !op.IsEnvelopeEnabled {
		return float64(op.Tl)
	}
	if op.UseCustomVolEnv {
		return op.customEnv()
	}
	return op.adsr()
}

func (op *Operator) customEnv() float64 {
	index := int(ClampF64(0.0, float64(SynthContext.Macro), float64(len(op.VolEnv))))
	return (float64(op.Tl) * (volROM[op.VolEnv[Clamp(0, index, len(op.VolEnv)-1)]&0b11111111]))
}

func (op *Operator) adsr() float64 {
	macro := SynthContext.Macro
	// Attack
	if macro <= op.Adsr.Attack {
		if op.Adsr.Attack <= 0 {
			return float64(op.Tl)
		}
		return LinearInterpolation(0, 0, float64(op.Adsr.Attack), float64(op.Tl), float64(macro))
	}

	// Decay and Sustain
	if macro > op.Adsr.Attack && macro < (op.Adsr.Attack+op.Adsr.Decay) {
		if op.Adsr.Decay <= 0 {
			return float64(op.Adsr.Sustain)
		}
		return LinearInterpolation(float64(op.Adsr.Attack), float64(op.Tl), float64(op.Adsr.Attack+op.Adsr.Decay), float64(op.Adsr.Sustain), float64(macro))
	}
	return float64(op.Adsr.Sustain)
}

func (op *Operator) pwmAdsr() float64 {
	macro := SynthContext.Macro
	// Attack
	if macro <= op.PwmAdsr.Attack {
		if op.PwmAdsr.Attack <= 0 {
			return float64(op.DutyCycle)
		}
		return LinearInterpolation(0, float64(op.PwmAdsr.Start), float64(op.PwmAdsr.Attack), float64(op.DutyCycle), float64(macro))
	}

	// Decay and Sustain
	if macro > op.PwmAdsr.Attack && macro < (op.PwmAdsr.Attack+op.PwmAdsr.Decay) {
		if op.PwmAdsr.Decay <= 0 {
			return float64(op.PwmAdsr.Sustain)
		}
		return LinearInterpolation(float64(op.PwmAdsr.Attack), float64(op.DutyCycle), float64(op.PwmAdsr.Attack+op.PwmAdsr.Decay), float64(op.PwmAdsr.Sustain), float64(macro))
	}
	return float64(op.PwmAdsr.Sustain)
}

func (op *Operator) getDutyCycle() float64 {
	if op.PwmAdsrEnabled {
		return op.pwmAdsr()
	}
	return float64(op.DutyCycle)
}
