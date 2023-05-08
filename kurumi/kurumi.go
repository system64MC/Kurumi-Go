package kurumi

// typedef unsigned char Uint8;
// typedef unsigned short Uint16;
// typedef short Int16;
// void Wavetable(void *userdata, Uint8 *stream, int len);
import "C"
import (
	"math"
	"strconv"
	"strings"
)

type Synth struct {
	WaveLen    int32
	WaveHei    int32
	MacLen     int32
	Macro      int32
	Operators  []Operator
	OpOutputs  []float32
	ModMatrix  [][]bool
	SmoothWin  int32
	Gain       float32
	Oversample int32

	FilterEnabled bool
	Cutoff        float32
	Pitch         int32
	Resonance     float32
	FilterType    int32

	FilterAdsrEnabled bool
	FilterStart       float32
	FilterAttack      int32
	FilterDecay       int32
	FilterSustain     float32

	SongPlaying          bool
	Normalize            bool
	NewNormalizeBehavior bool
}

type Adsr struct {
	Start   float32
	Attack  int32
	Decay   int32
	Sustain float32
}

var SynthContext *Synth

func ConstructSynth() *Synth {
	context := &Synth{WaveLen: 32, WaveHei: 31, MacLen: 64, Macro: 0, SmoothWin: 0, Gain: 1.0, Oversample: 1}
	context.ModMatrix = [][]bool{
		{false, false, false, false},
		{true, false, false, false},
		{false, true, false, false},
		{false, false, true, false},
	}
	context.OpOutputs = []float32{0, 0, 0, 1}
	for i := 0; i < 4; i++ {
		if i == 3 {
			context.Operators = append(context.Operators,
				Operator{
					Tl:              1,
					Adsr:            Adsr{0, 0, 0, 1},
					VolEnv:          []uint8{255},
					PhaseEnv:        []uint8{0},
					Wavetable:       []uint8{16, 25, 30, 31, 30, 29, 26, 25, 25, 28, 31, 28, 18, 11, 10, 13, 17, 20, 22, 20, 15, 6, 0, 2, 6, 5, 3, 1, 0, 0, 1, 4},
					MorphWave:       []uint8{16, 20, 15, 11, 11, 24, 30, 31, 28, 20, 10, 2, 0, 3, 5, 0, 16, 31, 26, 28, 31, 29, 21, 11, 3, 0, 1, 7, 20, 20, 16, 11},
					WaveformId:      0,
					Mult:            1,
					Phase:           0,
					Detune:          0,
					PhaseMod:        false,
					Feedback:        0,
					Prev:            0,
					UseCustomVolEnv: false,
					MorphTime:       64,
					DutyCycle:       0.5,
					PwmAdsr:         Adsr{0.5, 0, 0, 0.5},
					Interpolation:   0})
		} else {
			context.Operators = append(context.Operators,
				Operator{
					Tl:        0,
					Adsr:      Adsr{0, 0, 0, 0},
					VolEnv:    []uint8{255},
					PhaseEnv:  []uint8{0},
					Wavetable: []uint8{16, 25, 30, 31, 30, 29, 26, 25, 25, 28, 31, 28, 18, 11, 10, 13, 17, 20, 22, 20, 15, 6, 0, 2, 6, 5, 3, 1, 0, 0, 1, 4},
					MorphWave: []uint8{16, 20, 15, 11, 11, 24, 30, 31, 28, 20, 10, 2, 0, 3, 5, 0, 16, 31, 26, 28, 31, 29, 21, 11, 3, 0, 1, 7, 20, 20, 16, 11},

					WaveformId:      0,
					Mult:            1,
					Phase:           0,
					Detune:          0,
					PhaseMod:        false,
					Feedback:        0,
					Prev:            0,
					UseCustomVolEnv: false,
					MorphTime:       64,
					DutyCycle:       0.5,
					PwmAdsr:         Adsr{0.5, 0, 0, 0.5},
					Interpolation:   0})
		}
	}
	return context
}

/*------------------------------------------------*/

var volROM = [...]float64{0.0, 0.00390625, 0.0078125, 0.01171875, 0.015625, 0.01953125, 0.0234375, 0.02734375, 0.03125, 0.03515625, 0.0390625, 0.04296875, 0.046875, 0.05078125, 0.0546875, 0.05859375, 0.0625, 0.06640625, 0.0703125, 0.07421875, 0.078125, 0.08203125, 0.0859375, 0.08984375, 0.09375, 0.09765625, 0.1015625, 0.10546875, 0.109375, 0.11328125, 0.1171875, 0.12109375, 0.125, 0.12890625, 0.1328125, 0.13671875, 0.140625, 0.14453125, 0.1484375, 0.15234375, 0.15625, 0.16015625, 0.1640625, 0.16796875, 0.171875, 0.17578125, 0.1796875, 0.18359375, 0.1875, 0.19140625, 0.1953125, 0.19921875, 0.203125, 0.20703125, 0.2109375, 0.21484375, 0.21875, 0.22265625, 0.2265625, 0.23046875, 0.234375, 0.23828125, 0.2421875, 0.24609375, 0.25, 0.25390625, 0.2578125, 0.26171875, 0.265625, 0.26953125, 0.2734375, 0.27734375, 0.28125, 0.28515625, 0.2890625, 0.29296875, 0.296875, 0.30078125, 0.3046875, 0.30859375, 0.3125, 0.31640625, 0.3203125, 0.32421875, 0.328125, 0.33203125, 0.3359375, 0.33984375, 0.34375, 0.34765625, 0.3515625, 0.35546875, 0.359375, 0.36328125, 0.3671875, 0.37109375, 0.375, 0.37890625, 0.3828125, 0.38671875, 0.390625, 0.39453125, 0.3984375, 0.40234375, 0.40625, 0.41015625, 0.4140625, 0.41796875, 0.421875, 0.42578125, 0.4296875, 0.43359375, 0.4375, 0.44140625, 0.4453125, 0.44921875, 0.453125, 0.45703125, 0.4609375, 0.46484375, 0.46875, 0.47265625, 0.4765625, 0.48046875, 0.484375, 0.48828125, 0.4921875, 0.49609375, 0.5, 0.50390625, 0.5078125, 0.51171875, 0.515625, 0.51953125, 0.5234375, 0.52734375, 0.53125, 0.53515625, 0.5390625, 0.54296875, 0.546875, 0.55078125, 0.5546875, 0.55859375, 0.5625, 0.56640625, 0.5703125, 0.57421875, 0.578125, 0.58203125, 0.5859375, 0.58984375, 0.59375, 0.59765625, 0.6015625, 0.60546875, 0.609375, 0.61328125, 0.6171875, 0.62109375, 0.625, 0.62890625, 0.6328125, 0.63671875, 0.640625, 0.64453125, 0.6484375, 0.65234375, 0.65625, 0.66015625, 0.6640625, 0.66796875, 0.671875, 0.67578125, 0.6796875, 0.68359375, 0.6875, 0.69140625, 0.6953125, 0.69921875, 0.703125, 0.70703125, 0.7109375, 0.71484375, 0.71875, 0.72265625, 0.7265625, 0.73046875, 0.734375, 0.73828125, 0.7421875, 0.74609375, 0.75, 0.75390625, 0.7578125, 0.76171875, 0.765625, 0.76953125, 0.7734375, 0.77734375, 0.78125, 0.78515625, 0.7890625, 0.79296875, 0.796875, 0.80078125, 0.8046875, 0.80859375, 0.8125, 0.81640625, 0.8203125, 0.82421875, 0.828125, 0.83203125, 0.8359375, 0.83984375, 0.84375, 0.84765625, 0.8515625, 0.85546875, 0.859375, 0.86328125, 0.8671875, 0.87109375, 0.875, 0.87890625, 0.8828125, 0.88671875, 0.890625, 0.89453125, 0.8984375, 0.90234375, 0.90625, 0.91015625, 0.9140625, 0.91796875, 0.921875, 0.92578125, 0.9296875, 0.93359375, 0.9375, 0.94140625, 0.9453125, 0.94921875, 0.953125, 0.95703125, 0.9609375, 0.96484375, 0.96875, 0.97265625, 0.9765625, 0.98046875, 0.984375, 0.98828125, 0.9921875, 1}

type Dest = int

const DestWave Dest = 0
const DestMorph Dest = 1
const DestVolEnv Dest = 2
const DestPhaseEnv Dest = 3

func ApplyStringToWaveform(opId int, str string, destination Dest) {
	strArr := strings.Split(str, " ")
	bArr := make([]uint8, 0)
	for _, v := range strArr {
		p, err := strconv.ParseUint(v, 10, 8)
		if err == nil {
			bArr = append(bArr, uint8(p))
		} else {
			println(err)
		}
	}

	if len(bArr) == 0 {
		bArr = append(bArr, 0)
	}

	switch destination {
	case 0:
		SynthContext.Operators[opId].Wavetable = bArr
	case 1:
		SynthContext.Operators[opId].MorphWave = bArr
	case 2:
		SynthContext.Operators[opId].VolEnv = bArr
	case 3:
		SynthContext.Operators[opId].PhaseEnv = bArr
	}
}

var WaveOutput = make([]int, 0)

func Normalize(wavetable []float64) []float64 {
	waveMin, waveMax := minandmaxFloat(wavetable)
	mult := math.Max(waveMin, waveMax)
	for i := 0; i < len(wavetable); i++ {
		wavetable[i] = wavetable[i] * 1 / mult
	}
	return wavetable
}

func Synthesize() {
	if SynthContext.NewNormalizeBehavior {
		SynthesizeNew()
		return
	}
	SynthesizeOld()
}

func SynthesizeNew() {
	ResetFB()
	WaveOutput = make([]int, 0)
	myLen := int(SynthContext.WaveLen)
	oversample := int(SynthContext.Oversample)
	myTmp := make([]float64, myLen*oversample)
	// Preheat
	for x := 0; x < int(myLen*oversample); x++ {
		myTmp[x] = (fm(float64(x)))
	}
	for x := 0; x < int(myLen*oversample); x++ {
		myTmp[x] = (fm(float64(x)))
	}
	if SynthContext.FilterEnabled {
		myTmp = filter(myTmp)
	}

	myTmp = smooth(myTmp)

	myOutFloat := make([]float64, SynthContext.WaveLen)
	myOut := make([]int, SynthContext.WaveLen)
	tmpLen := len(myTmp)

	for c := 0; c < tmpLen; c += oversample {
		res := 0.0
		for i := 0; i < oversample; i++ {
			res += myTmp[c+i]
		}
		res = res / float64(oversample)
		myOutFloat[c/oversample] = res
	}

	if SynthContext.Normalize {
		myOutFloat = Normalize(myOutFloat)
		for c := 0; c < len(myOutFloat); c++ {
			myOut[c] = int(math.Round((myOutFloat[c] + 1) * (float64(SynthContext.WaveHei) / 2.0)))
		}
	} else {
		for c := 0; c < len(myOutFloat); c++ {
			myOut[c] = int(ClampF64(0, math.Round((myOutFloat[c]+1)*(float64(SynthContext.WaveHei)/2.0)), float64(SynthContext.WaveHei)))
		}
	}

	WaveOutput = myOut
}

func SynthesizeOld() {
	ResetFB()
	WaveOutput = make([]int, 0)
	myLen := int(SynthContext.WaveLen)
	oversample := int(SynthContext.Oversample)
	myTmp := make([]float64, myLen*oversample)
	// Preheat
	for x := 0; x < int(myLen*oversample); x++ {
		myTmp[x] = (fm(float64(x)))
	}
	for x := 0; x < int(myLen*oversample); x++ {
		myTmp[x] = (fm(float64(x)))
	}
	if SynthContext.FilterEnabled {
		myTmp = filter(myTmp)
	}

	myTmp = smooth(myTmp)

	if SynthContext.Normalize {
		myTmp = Normalize(myTmp)
	} else {
		for x := 0; x < len(myTmp); x++ {
			myTmp[x] = ClampF64(-1, myTmp[x]*float64(SynthContext.Gain), 1)
		}
	}

	myOut := make([]int, SynthContext.WaveLen)
	tmpLen := len(myTmp)

	for c := 0; c < tmpLen; c += oversample {
		res := 0.0
		for i := 0; i < oversample; i++ {
			res += myTmp[c+i]
		}
		res = res / float64(oversample)
		tmp := int(math.Round((res + 1) * (float64(SynthContext.WaveHei) / 2.0)))
		myOut[c/oversample] = tmp
	}
	WaveOutput = myOut
}

var WaveStr = ""

func GenerateWaveStr() {
	str := ""
	for _, n := range WaveOutput {
		str += strconv.Itoa(n) + " "
	}
	str += ";"
	WaveStr = str
}

var WaveSeqStr = ""

func GenerateWaveSeqStr() {
	str := ""
	tmpMac := SynthContext.Macro
	for i := 0; i < int(SynthContext.MacLen); i++ {
		SynthContext.Macro = int32(i)
		Synthesize()
		GenerateWaveStr()

		str += WaveStr + "\n"
	}
	SynthContext.Macro = tmpMac
	Synthesize()
	WaveSeqStr = str
}

const (
	toneHz   = 440
	sampleHz = 48000
)
