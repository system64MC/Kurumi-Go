package kurumi

// typedef unsigned char Uint8;
// typedef unsigned short Uint16;
// typedef short Int16;
// void Wavetable(void *userdata, Uint8 *stream, int len);
import "C"
import (
	"encoding/json"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
	"fmt"

	"github.com/ncruces/zenity"

	"github.com/veandco/go-sdl2/sdl"
)

type Adsr struct {
	Start   float32
	Attack  int32
	Decay   int32
	Sustain float32
}

type Operator struct {
	Tl              float32
	Reverse bool
	Adsr            Adsr
	WaveformId      int32
	Mult            int32
	Phase           float32
	Detune          int32
	PhaseMod        bool
	PhaseRev        bool
	Feedback        float32
	Prev            float32
	Curr float32
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
	DutyCycle 		float32
	PwmAdsr			Adsr
	PwmAdsrEnabled  bool

	IsEnvelopeEnabled bool
}

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
	//filter
	FilterEnabled bool
	Cutoff  float32
	Pitch int32
	Resonance  float32
	FilterType int32

	FilterAdsrEnabled bool
	FilterStart float32
	FilterAttack int32
	FilterDecay int32
	FilterSustain float32

	SongPlaying bool
	Normalize bool
}

func EncodeJson() []byte {
	// synthJson, err := json.Marshal(SynthContext)
	// if err != nil {
    //     fmt.Println("Error:", err)
    // }

	data := map[string]interface{}{
		"Format": "vampire",
		"Synth":  SynthContext,
	}

	synthJson, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return synthJson
}

func SaveJson() error {

	path, errZen := zenity.SelectFileSave(
		zenity.ConfirmOverwrite(),
		zenity.Filename("output"),
		zenity.FileFilters{
			{"Kurumi Vampire Patch files", []string{"*.kvp"}, false},
		})
	if(errZen == zenity.ErrCanceled) {
		return errZen
	}
	if !strings.HasSuffix(path, ".kvp") {
        path += ".kvp"
    }

	data := EncodeJson()
	file, err := os.Create(path)
	if(err != nil) {
		return err
	}
	
	_, err2 := file.Write(data)
	if(err2 != nil) {
		return err2
	}
	return nil
}

func LoadJson() error {
	path, errZen := zenity.SelectFile(
		zenity.FileFilters{
			{"Kurumi Vampire Patch files", []string{"*.kvp"}, false},
		})

	if(errZen == zenity.ErrCanceled) {
			return errZen
	}

	jsonData, err := os.Open(path)
	defer jsonData.Close()

	decoder := json.NewDecoder(jsonData)

	// Create a map to decode the JSON into
	var data map[string]interface{}

	// Decode the JSON
	err = decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	// Extract the Synth object from the map
	synth := &Synth{}
	format := data["Format"].(string)

	if format != "vampire" {
		panic("Invalid format")
	}

	synthMap, ok := data["Synth"].(map[string]interface{})
	if ok {
		// Decode the Synth map into a Synth struct
		bytes, err := json.Marshal(synthMap)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(bytes, synth)
		if err != nil {
			panic(err)
		}
	}

	SynthContext.ModMatrix = synth.ModMatrix
	
	SynthContext.Cutoff = synth.Cutoff
	SynthContext.FilterAdsrEnabled = synth.FilterAdsrEnabled
	SynthContext.FilterStart = synth.FilterStart
	SynthContext.FilterAttack = synth.FilterAttack
	SynthContext.FilterDecay = synth.FilterDecay
	SynthContext.FilterSustain = synth.FilterSustain
	SynthContext.FilterType = synth.FilterType
	SynthContext.FilterEnabled = synth.FilterEnabled
	SynthContext.Pitch = synth.Pitch
	SynthContext.Resonance = synth.Resonance
	
	SynthContext.Operators = synth.Operators
	SynthContext.OpOutputs = synth.OpOutputs
	
	SynthContext.Normalize = synth.Normalize
	
	SynthContext.WaveLen = synth.WaveLen
	SynthContext.WaveHei = synth.WaveHei
	
	SynthContext.MacLen = synth.MacLen
	SynthContext.Macro = synth.Macro
	
	SynthContext.SmoothWin = synth.SmoothWin
	
	SynthContext.Gain = synth.Gain
	
	SynthContext.Oversample = synth.Oversample
	
	Synthesize()

	
	fmt.Printf("Format: %s\n", data["Format"])
	fmt.Printf("Synth: %+v\n", synth)

	return nil
}

var SynthContext *Synth

func ConstructSynth() *Synth {
	context := &Synth{WaveLen: 32, WaveHei: 31, MacLen: 64, Macro: 0, SmoothWin: 0, Gain: 1.0, Oversample: 1 }
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
					DutyCycle: 		 0.5,
					PwmAdsr: Adsr{0.5, 0, 0, 0.5},
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
					DutyCycle: 		 0.5,
					PwmAdsr: 		 Adsr{0.5, 0, 0, 0.5},
					Interpolation:   0})
		}
	}
	return context
}

func moduloF64(a float64, b float64) float64 {
	tmp := math.Mod(a, b) + b
	return math.Mod(tmp, b)
}

func moduloInt(a int, b int) int {
	return ((a % b) + b) % b
}

func minandmax(values []uint8) (uint8, uint8) {
	min := values[0]   //assign the first element equal to min
	max := values[0]  //assign the first element equal to max
	for _, number := range values {
		if number < min {
			min = number
		}
		if number > max {
			max = number
		}
	}
	return min, max
}

func minandmaxFloat(values []float64) (float64, float64) {
	min := values[0]   //assign the first element equal to min
	max := values[0]  //assign the first element equal to max
	for _, number := range values {
		if number < min {
			min = number
		}
		if number > max {
			max = number
		}
	}
	min = math.Abs(min)
	return min, max
}

func lerp(x, y, a float64) float64 {
	return x * (1 - a) + y * a
}

func getWTSample(op *Operator, x float64) float64 {
	if(op.Morphing) {
		a := interpolate(x * float64(len(op.Wavetable)) * float64(op.Mult) + (float64(op.Phase) * float64(len(op.Wavetable)) + op.getPhase() * float64(len(op.Wavetable))), op, op.Wavetable)
		b := interpolate(x * float64(len(op.MorphWave)) * float64(op.Mult) + (float64(op.Phase) * float64(len(op.MorphWave)) + op.getPhase() * float64(len(op.MorphWave))), op, op.MorphWave)
		c := float64(SynthContext.Macro) / float64(op.MorphTime)
		return lerp(a, b, c)
	}
	return interpolate(x * float64(len(op.Wavetable)) * float64(op.Mult) + (float64(op.Phase) * float64(len(op.Wavetable)) + op.getPhase() * float64(len(op.Wavetable))), op, op.Wavetable)
}

func noInterpolation(x float64, wt []uint8) float64 {
	t := x
	idx := math.Floor(t)
	len := len(wt)
	_, max := minandmax(wt)
	myMax := float64(max) / 2.0
	s0 := float64(wt[moduloInt(int(idx), len)]) / myMax - 1
	return s0
}

func linearInterpolation(x float64, wt []uint8) float64 {
	t := x
	idx := int(math.Floor(t))
	len := len(wt)
	_, m := minandmax(wt)
	max := float64(m) / 2.0
	mu := (t - float64(idx))
	s0 := float64(wt[moduloInt(idx, len)]) / max - 1.0
	s1 := float64(wt[moduloInt(idx + 1, len)]) / max - 1.0
	return s0 + mu * s1 - (mu * s0)
}

func cosineInterpolation(x float64, wt []uint8) float64 {
	t := x
	idx := int(math.Floor(t))
	len := len(wt)
	_, m := minandmax(wt)
	max := float64(m) / 2.0
	mu := (t - float64(idx))
	muCos := (1 - math.Cos(mu * math.Pi) / 2)
	s0 := float64(wt[moduloInt(idx, len)]) / max - 1.0
	s1 := float64(wt[moduloInt(idx + 1, len)]) / max - 1.0
	return s0 + muCos * s1 - (muCos * s0)
}

func cubicInterpolation(x float64, wt []uint8) float64 {
	t := x
	idx := int(math.Floor(t))
	len := len(wt)
	_, m := minandmax(wt)
	max := float64(m) / 2.0
	s0 := float64(wt[moduloInt(idx - 1, len)]) / max - 1.0
	s1 := float64(wt[moduloInt(idx, len)]) / max - 1.0
	s2 := float64(wt[moduloInt(idx + 1, len)]) / max - 1.0
	s3 := float64(wt[moduloInt(idx + 2, len)]) / max - 1.0
	mu := (t - float64(idx))
	mu2 := mu * mu
	a0 := -0.5 * s0 + 1.5 * s1 - 1.5 * s2 + 0.5 * s3
    a1 := s0 - 2.5 * s1 + 2 * s2 - 0.5 * s3
    a2 := -0.5 * s0 + 0.5 * s2
    a3 := s1
	return (a0 * mu * mu2 + a1 * mu2 + a2 * mu + a3)
}

func interpolate(x float64, op *Operator, wt []uint8) float64 {
	// wt := op.Wavetable
	switch op.Interpolation {
	case 0:
		return noInterpolation(x, wt)
	case 1:
		return linearInterpolation(x, wt)
	case 2:
		return cosineInterpolation(x, wt)
	case 3:
		return cubicInterpolation(x, wt)
	}
	return noInterpolation(x, wt)
}

var CopiedOp *Operator = nil

func CopyOp(op *Operator) {
	CopiedOp = &Operator {}
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
/*------------------------------------------------*/

func sine(op *Operator, x float64) float64 {
	return math.Sin((x * float64(op.Mult) * 2 * math.Pi) + (float64(op.Phase)*2*math.Pi + (op.getPhase() * math.Pi * 2)))
}
func rectSine(op *Operator, x float64) float64 {
	return math.Max(sine(op, x), 0)
}
func absSine(op *Operator, x float64) float64 {
	return math.Abs(sine(op, x))
}
func quarterSine(op *Operator, x float64) float64 {
	if math.Mod((x*float64(op.Mult)+float64(op.Phase)+op.getPhase()), 0.5) <= 0.25 {
		return absSine(op, x)
	}
	return 0
}
func squishedSine(op *Operator, x float64) float64 {
	if sine(op, x) > 0 {
		return math.Sin((x * float64(op.Mult) * 4 * math.Pi) + (float64(op.Phase)*4*math.Pi + (op.getPhase() * math.Pi * 4)))
	}
	return 0
}
func squishedRectSine(op *Operator, x float64) float64 {
	return math.Max(squishedSine(op, x), 0)
}
func squishedAbsSine(op *Operator, x float64) float64 {
	return math.Abs(squishedSine(op, x))
}

func square(op *Operator, x float64) float64 {
	width := op.getDutyCycle()
	a := moduloF64(x * math.Pi * 2 * float64(op.Mult) + (float64(op.Phase) * math.Pi * 2 + op.getPhase() * math.Pi * 2), math.Pi * 2)
	if a >= (math.Pi * width * 2) {
		return -1
	}
	return 1
}

func rectSquare(op *Operator, x float64) float64 {
	a := square(op, x)
	if(a < 0) {
		a = 0
	}
	return a
}

func saw(op *Operator, x float64) float64 {
	return math.Atan(math.Tan(x*math.Pi*float64(op.Mult)+(float64(op.Phase)*math.Pi+(op.getPhase()*math.Pi)))) / (math.Pi * 0.5)
}
func rectSaw(op *Operator, x float64) float64 {
	return math.Max(saw(op, x), 0)
}
func absSaw(op *Operator, x float64) float64 {
	a := saw(op, x)
	if a < 0 {
		return a + 1
	}
	return a
}

func cubSaw(op *Operator, x float64) float64 {
	a := saw(op, x)
	return math.Pow(a, 3)
}
func rectCubSaw(op *Operator, x float64) float64 {
	return math.Max(cubSaw(op, x), 0)
}
func absCubSaw(op *Operator, x float64) float64 {
	a := absSaw(op, x)
	return math.Pow(a, 3)
}

func cubedSine(op *Operator, x float64) float64 {
	a := sine(op, x)
	return math.Pow(a, 3)
}
func cubedRectSine(op *Operator, x float64) float64 {
	a := rectSine(op, x)
	return math.Pow(a, 3)
}
func cubedAbsSine(op *Operator, x float64) float64 {
	a := absSine(op, x)
	return math.Pow(a, 3)
}
func cubedQuarterSine(op *Operator, x float64) float64 {
	a := quarterSine(op, x)
	return math.Pow(a, 3)
}
func cubedSquishedSine(op *Operator, x float64) float64 {
	a := squishedSine(op, x)
	return math.Pow(a, 3)
}
func cubedRectSquiSine(op *Operator, x float64) float64 {
	a := squishedRectSine(op, x)
	return math.Pow(a, 3)
}
func cubedAbsSquiSine(op *Operator, x float64) float64 {
	a := squishedAbsSine(op, x)
	return math.Pow(a, 3)
}

func triangle(op *Operator, x float64) float64 {
	return math.Asin(sine(op, x)) / (math.Pi * 0.5)
}
func rectTriangle(op *Operator, x float64) float64 {
	return math.Max(triangle(op, x), 0)
}
func absTriangle(op *Operator, x float64) float64 {
	return math.Abs(triangle(op, x))
}
func quarterTriangle(op *Operator, x float64) float64 {
	if math.Mod((x*float64(op.Mult)+(float64(op.Phase)+op.getPhase())), 0.5) <= 0.25 {
		return triangle(op, x)
	}
	return 0
}
func squishedTriangle(op *Operator, x float64) float64 {
	if sine(op, x) > 0 {
		return math.Asin(math.Sin((x * float64(op.Mult) * 4 * math.Pi)   + (float64(op.Phase) * math.Pi * 4 + (op.getPhase() * math.Pi * 4)) ))/ (math.Pi/2)
	}
	return 0
}
func rectSquiTriangle(op *Operator, x float64) float64 {
	return math.Max(squishedTriangle(op, x), 0)
}
func absSquiTriangle(op *Operator, x float64) float64 {
	return math.Abs(squishedTriangle(op, x))
}
func cubedTriangle(op *Operator, x float64) float64 {
	a := triangle(op, x)
	return math.Pow(a, 3)
}
func cubedRectTri(op *Operator, x float64) float64 {
	a := rectTriangle(op, x)
	return math.Pow(a, 3)
}
func cubedAbsTri(op *Operator, x float64) float64 {
	a := absTriangle(op, x)
	return math.Pow(a, 3)
}
func cubedQuartTri(op *Operator, x float64) float64 {
	a := quarterTriangle(op, x)
	return math.Pow(a, 3)
}
func cubedSquiTri(op *Operator, x float64) float64 {
	a := squishedTriangle(op, x)
	return math.Pow(a, 3)
}
func cubedRectSquiTri(op *Operator, x float64) float64 {
	a := rectSquiTriangle(op, x)
	return math.Pow(a, 3)
}
func cubedAbsSquiTri(op *Operator, x float64) float64 {
	a := absSquiTriangle(op, x)
	return math.Pow(a, 3)
}

func noise(op *Operator, x float64) float64 {
	return (rand.Float64() * 2) - 1
}

func custom(op *Operator, x float64) float64 {
	return getWTSample(op, x)
}
func rectCustom(op *Operator, x float64) float64 {
	return math.Max(custom(op, x), 0)
}
func absCustom(op *Operator, x float64) float64 {
	return math.Abs(custom(op, x))
}
func cubedCustom(op *Operator, x float64) float64 {
	a := custom(op, x)
	return math.Pow(a, 3)
}

var WaveFuncs = []func(*Operator, float64) float64{
	sine,
	rectSine,
	absSine,
	quarterSine,
	squishedSine,
	squishedRectSine,
	squishedAbsSine,
	square,
	rectSquare,
	saw,
	rectSaw,
	absSaw,
	cubSaw,
	rectCubSaw,
	absCubSaw,
	cubedSine,
	cubedRectSine,
	cubedAbsSine,
	cubedQuarterSine,
	cubedSquishedSine,
	cubedRectSquiSine,
	cubedAbsSquiSine,
	triangle,
	rectTriangle,
	absTriangle,
	quarterTriangle,
	squishedTriangle,
	rectSquiTriangle,
	absSquiTriangle,
	cubedTriangle,
	cubedRectTri,
	cubedAbsTri,
	cubedQuartTri,
	cubedSquiTri,
	cubedRectSquiTri,
	cubedAbsSquiTri,
	noise,
	custom,
	rectCustom,
	absCustom,
	cubedCustom,
}

var Waveforms = []string{
	"Sine",
	"Rect. Sine",
	"Abs. Sine",
	"Quarter Sine",
	"Squished Sine",
	"Squished Rect. Sine",
	"Squished Abs. Sine",
	"Pulse",
	"Rectified Pulse",
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
	"Squi. Rect. Cubed Sine",
	"Squi. Abs. Cubed Sine",
	"Triangle",
	"Rect. Triange",
	"Abs. Triangle",
	"Quarter Triangle",
	"Squished Triangle",
	"Rect. Squished Triangle",
	"Abs. Squished Triangle",
	"Cubed Triangle",
	"Rect. Cubed Triangle",
	"Abs. Cubed Triangle",
	"Quarter Cubed Triangle",
	"Squi. Cubed Triangle",
	"Squi. Rect. Cubed Triangle",
	"Squi. Abs. Cubed Triangle",
	"Noise",
	"Custom",
	"Rect. Custom",
	"Abs. Custom",
	"Cubed Custom",
}

var ModModes = []string{
	"FM",
	"OR",
	"XOR",
	"AND",
	"NAND",
	"ADD",
	"SUB",
	"MUL",
}

var Interpolations = []string{
	"None",
	"Linear",
	"Cosine",
	"Cubic",
}

var volROM = [...]float64 {0.0, 0.00390625, 0.0078125, 0.01171875, 0.015625, 0.01953125, 0.0234375, 0.02734375, 0.03125, 0.03515625, 0.0390625, 0.04296875, 0.046875, 0.05078125, 0.0546875, 0.05859375, 0.0625, 0.06640625, 0.0703125, 0.07421875, 0.078125, 0.08203125, 0.0859375, 0.08984375, 0.09375, 0.09765625, 0.1015625, 0.10546875, 0.109375, 0.11328125, 0.1171875, 0.12109375, 0.125, 0.12890625, 0.1328125, 0.13671875, 0.140625, 0.14453125, 0.1484375, 0.15234375, 0.15625, 0.16015625, 0.1640625, 0.16796875, 0.171875, 0.17578125, 0.1796875, 0.18359375, 0.1875, 0.19140625, 0.1953125, 0.19921875, 0.203125, 0.20703125, 0.2109375, 0.21484375, 0.21875, 0.22265625, 0.2265625, 0.23046875, 0.234375, 0.23828125, 0.2421875, 0.24609375, 0.25, 0.25390625, 0.2578125, 0.26171875, 0.265625, 0.26953125, 0.2734375, 0.27734375, 0.28125, 0.28515625, 0.2890625, 0.29296875, 0.296875, 0.30078125, 0.3046875, 0.30859375, 0.3125, 0.31640625, 0.3203125, 0.32421875, 0.328125, 0.33203125, 0.3359375, 0.33984375, 0.34375, 0.34765625, 0.3515625, 0.35546875, 0.359375, 0.36328125, 0.3671875, 0.37109375, 0.375, 0.37890625, 0.3828125, 0.38671875, 0.390625, 0.39453125, 0.3984375, 0.40234375, 0.40625, 0.41015625, 0.4140625, 0.41796875, 0.421875, 0.42578125, 0.4296875, 0.43359375, 0.4375, 0.44140625, 0.4453125, 0.44921875, 0.453125, 0.45703125, 0.4609375, 0.46484375, 0.46875, 0.47265625, 0.4765625, 0.48046875, 0.484375, 0.48828125, 0.4921875, 0.49609375, 0.5, 0.50390625, 0.5078125, 0.51171875, 0.515625, 0.51953125, 0.5234375, 0.52734375, 0.53125, 0.53515625, 0.5390625, 0.54296875, 0.546875, 0.55078125, 0.5546875, 0.55859375, 0.5625, 0.56640625, 0.5703125, 0.57421875, 0.578125, 0.58203125, 0.5859375, 0.58984375, 0.59375, 0.59765625, 0.6015625, 0.60546875, 0.609375, 0.61328125, 0.6171875, 0.62109375, 0.625, 0.62890625, 0.6328125, 0.63671875, 0.640625, 0.64453125, 0.6484375, 0.65234375, 0.65625, 0.66015625, 0.6640625, 0.66796875, 0.671875, 0.67578125, 0.6796875, 0.68359375, 0.6875, 0.69140625, 0.6953125, 0.69921875, 0.703125, 0.70703125, 0.7109375, 0.71484375, 0.71875, 0.72265625, 0.7265625, 0.73046875, 0.734375, 0.73828125, 0.7421875, 0.74609375, 0.75, 0.75390625, 0.7578125, 0.76171875, 0.765625, 0.76953125, 0.7734375, 0.77734375, 0.78125, 0.78515625, 0.7890625, 0.79296875, 0.796875, 0.80078125, 0.8046875, 0.80859375, 0.8125, 0.81640625, 0.8203125, 0.82421875, 0.828125, 0.83203125, 0.8359375, 0.83984375, 0.84375, 0.84765625, 0.8515625, 0.85546875, 0.859375, 0.86328125, 0.8671875, 0.87109375, 0.875, 0.87890625, 0.8828125, 0.88671875, 0.890625, 0.89453125, 0.8984375, 0.90234375, 0.90625, 0.91015625, 0.9140625, 0.91796875, 0.921875, 0.92578125, 0.9296875, 0.93359375, 0.9375, 0.94140625, 0.9453125, 0.94921875, 0.953125, 0.95703125, 0.9609375, 0.96484375, 0.96875, 0.97265625, 0.9765625, 0.98046875, 0.984375, 0.98828125, 0.9921875, 1}


func Clamp(low, val, high int) int {
	if val < low {
		return low
	}
	if val > high {
		return high
	}
	return val
}
func (op *Operator) getPhase() float64 {
	macro := SynthContext.Macro
	macLen := SynthContext.MacLen

	if !op.CustomPhaseEnv {
	// Anti divide-by-0
		if(macLen == 1) {
			macLen = 2
		}
		return (float64(macro) / float64(macLen-1)) * float64(op.Detune)
	}

	if(len(op.PhaseEnv) < 1) {
		return 0
	}
	return (float64(op.PhaseEnv[int(Clamp(0, int(macro), len(op.PhaseEnv) - 1))]) / 255.0) * float64(op.Detune)
}



type Dest = int

const DestWave Dest = 0
const DestMorph Dest = 1
const DestVolEnv Dest = 2
const DestPhaseEnv Dest = 3

func ApplyStringToWaveform(opId int,str string, destination Dest) {
	strArr := strings.Split(str, " ")
	// println(*str)
	bArr := make([]uint8, 0)
	for _, v := range strArr {
		p, err := strconv.ParseUint(v, 10, 8)
		if(err == nil) {
			bArr = append(bArr, uint8(p))
			println(p)
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

func ClampF64(low, val, high float64) float64 {
	if val < low {
		return low
	}
	if val > high {
		return high
	}
	return val
}

func moduloFix(a, b int) int {
	return ((a % b) + b) % b
}

func smooth(arr []float64) []float64 {
	out := make([]float64, 0)
	for i := 0; i < len(arr); i++ {
		smp := 0.0
		for j := -SynthContext.SmoothWin; j <= SynthContext.SmoothWin; j++ {
			smp += arr[moduloFix(i+int(j), len(arr))]
		}
		avg := smp / (float64(SynthContext.SmoothWin*2) + 1)
		out = append(out, avg)
	}
	return out
}

var WaveOutput = make([]int, 0)

var samples = make([]float64, 4)

func logicMod(x float64, modValue float64, opId int) float64 {
	op := &SynthContext.Operators[opId]
	switch op.ModMode {
	case 0: // FM
		return op.oscillate(x+modValue+op.getFB()) * op.getVolume()
	case 1: // OR
		a := int(math.Round((modValue + 1) * 32767.5))
		b := int(math.Round(((op.oscillate(x) * op.getVolume()) + (1 * op.getVolume())) * 32767.5))
		return (float64(a|b)/32767.5) - (1 * op.getVolume())
	case 2: // XOR
		a := int(math.Round((modValue + 1) * 32767.5))
		b := int(math.Round(((op.oscillate(x) * op.getVolume()) + (1 * op.getVolume())) * 32767.5))
		return (float64(a^b)/32767.5) - (1 * op.getVolume())
	case 3: // AND
		a := int(math.Round((modValue + 1) * 32767.5))
		b := int(math.Round(((op.oscillate(x) * op.getVolume()) + (1 * op.getVolume())) * 32767.5))
		return (float64(a&b)/32767.5) - (1 * op.getVolume())
	case 4: // NAND
		a := int(math.Round((modValue + 1) * 32767.5))
		b := int(math.Round(((op.oscillate(x) * op.getVolume()) + (1 * op.getVolume())) * 32767.5))
		return float64(^int(a&b))/32767.5 - (1 * op.getVolume())
	case 5: // ADD
		return modValue + (op.oscillate(x) * op.getVolume())
	case 6: // SUB
		return op.oscillate(x)*op.getVolume() - modValue
	case 7: // MUL
		return modValue * (op.oscillate(x) * op.getVolume())
	}
	return op.oscillate(x+modValue+op.getFB()) * op.getVolume()
}

func fm(x float64) float64 {
	x = x / float64(SynthContext.WaveLen*SynthContext.Oversample)
	matrix := SynthContext.ModMatrix
	for op := 0; op < 4; op++ {
		sum := 0.0
		for mod := 0; mod < 4; mod++ {
			if matrix[op][mod] {
				sum += samples[mod]
			}
		}
		samples[op] = logicMod(x, sum, op)
		SynthContext.Operators[op].Prev = float32(samples[op])
	}
	output := 0.0
	for o := 0; o < 4; o++ {
		output += samples[o] * float64(SynthContext.OpOutputs[o])
	}
	return output
}

func fm2(x float64) float64 {
	x = x / float64(65536)
	matrix := SynthContext.ModMatrix
	for op := 0; op < 4; op++ {
		sum := 0.0
		for mod := 0; mod < 4; mod++ {
			if matrix[op][mod] {
				sum += samples[mod]
			}
		}
		samples[op] = logicMod(x, sum, op)
		SynthContext.Operators[op].Prev = float32(samples[op])
	}
	output := 0.0
	for o := 0; o < 4; o++ {
		output += samples[o] * float64(SynthContext.OpOutputs[o])
	}
	return output
}

func ResetFB() {
	for i := 0; i < 4; i++ {
		SynthContext.Operators[i].Prev = 0
		SynthContext.Operators[i].Curr = 0
		samples[i] = 0
	}
}

func Normalize(wavetable []float64) []float64 {
	waveMin, waveMax := minandmaxFloat(wavetable)
	mult := math.Max(waveMin, waveMax)
	for i := 0; i < len(wavetable); i++ {
		wavetable[i] = wavetable[i] * 1/mult
	}
	return wavetable
}

func Synthesize() {
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
	if(SynthContext.FilterEnabled) {
		myTmp = filter(myTmp)
	}

	myTmp = smooth(myTmp)

	if(SynthContext.Normalize) {
		myTmp = Normalize(myTmp)
	} else {
		for x := 0; x < len(myTmp); x++ {
			myTmp[x] = ClampF64(-1, myTmp[x] * float64(SynthContext.Gain), 1)
		}
	}

	myOut := make([]int, 0)
	tmpLen := len(myTmp)
	for c := 0; c < tmpLen; c += oversample {
		res := 0.0
		for i := 0; i < oversample; i++ {
			res += myTmp[c+i]
		}
		res = res / float64(oversample)
		tmp := int(math.Round((res + 1) * (float64(SynthContext.WaveHei) / 2.0)))
		myOut = append(myOut, tmp)
	}
	WaveOutput = myOut
}


type LowPassFilter struct {
	alpha, lastY float64
}

// Create a new low-pass filter
func NewLowPassFilter(cutoff, sampleRate float64) *LowPassFilter {
	dt := 1.0 / sampleRate
	rc := 1.0 / (2.0 * math.Pi * cutoff)
	alpha := dt / (rc + dt)
	return &LowPassFilter{alpha: alpha}
}

// Filter a sample
func (f *LowPassFilter) Filter(x float64) float64 {
	y := f.lastY + f.alpha*(x-f.lastY)
	f.lastY = y
	return y
}

// Resample a waveform to a new length
func Resample(input []float64, outputLength int) []float64 {
	// Calculate the original length of the input waveform
	originalLength := len(input)

	// Calculate the ratio of the output length to the original length
	ratio := float64(outputLength) / float64(originalLength)

	// Calculate the target sampling rate
	targetSampleRate := float64(outputLength) / float64(originalLength)

	// Create a new slice to hold the output waveform
	output := make([]float64, outputLength)

	// Create a low-pass filter with a cutoff frequency of half the target sampling rate
	cutoff := targetSampleRate / 2.0
	filter := NewLowPassFilter(cutoff, float64(originalLength))

	// Initialize the filter state
	filter.lastY = input[0]

	// Resample the waveform
	for i := 0; i < outputLength; i++ {
		// Calculate the index of the closest input sample
		j := int(math.Floor(float64(i) / ratio))

		// Filter the input sample
		y := filter.Filter(input[j])

		// Add the filtered sample to the output waveform
		output[i] = y
		println(y)
	}

	return output
}



func Synthesize22() {
	WaveOutput = make([]int, 0)
	myTmp := make([]float64, 65536)
	for x := 0; x < int(65536); x++ {
		myTmp[x] = fm(float64(x))
	}

	for x := 0; x < len(myTmp); x++ {
		myTmp[x] = ClampF64(-1, myTmp[x] * float64(SynthContext.Gain), 1)
	}

	downsampled := downsampleWaveTableLinear(myTmp, int(SynthContext.WaveLen))
	myOut := make([]int, 0)

	for c := 0; c < len(downsampled); c++ {
		res := downsampled[c]
		// println(res)
		tmp := int(math.Round((res + 1) * (float64(SynthContext.WaveHei) / 2.0)))
		myOut = append(myOut, tmp)
	}
	WaveOutput = myOut
}

func downsampleWaveTableLinear2(wavetable []float64, newSampleRate int) []float64 {
    // Compute the downsampling factor
    downsampleFactor := float64(len(wavetable)) / float64(newSampleRate)
    
    // Create a low-pass filter with a cutoff frequency of 1/downsampleFactor
    filterSize := 2 * int(downsampleFactor) + 1
    filter := make([]float64, filterSize)
    cutoff := 1 / downsampleFactor
    for i := 0; i < filterSize; i++ {
        t := float64(i - filterSize/2) / float64(newSampleRate)
        if t == 0 {
            filter[i] = 2 * cutoff
        } else {
            filter[i] = math.Sin(2 * math.Pi * cutoff * t) / (math.Pi * t)
        }
    }
    
    // Normalize the filter
    sum := 0.0
    for i := range filter {
        sum += filter[i]
    }
    for i := range filter {
        filter[i] /= sum
    }
    
    // Convolve the filter with the waveform to obtain the downsampled signal
    downsampled := make([]float64, newSampleRate)
    for i := 0; i < newSampleRate; i++ {
        start := int(float64(i) * downsampleFactor)
        for j := 0; j < filterSize; j++ {
            if start+j < len(wavetable) {
                downsampled[i] += wavetable[start+j] * filter[j]
            }
        }
    }
    
    return downsampled
}

func downsampleWaveTableLinear(input []float64, outputLen int) []float64 {
	// Compute the input and output sample rates
	inputLen := len(input)
	sampleRateIn := float64(inputLen)
	sampleRateOut := float64(outputLen)

	// Compute the ratio between the input and output sample rates
	ratio := sampleRateIn / sampleRateOut

	// Allocate memory for the output signal
	output := make([]float64, outputLen)

	// Compute the index of the first sample in the input signal
	// indexIn := 0

	// Interpolate each sample in the output signal
	for i := 0; i < outputLen; i++ {
		// Compute the index of the current sample in the input signal
		index := int(float64(i) * ratio)

		// Compute the fractional distance between the current sample index and the next sample index
		frac := float64(i)*ratio - float64(index)

		// Interpolate the current sample using linear interpolation
		if index < inputLen-1 {
			output[i] = (1-frac)*input[index] + frac*input[index+1]
		} else {
			// If we've reached the end of the input signal, repeat the last sample
			output[i] = input[inputLen-1]
		}
	}

	return output
}


func (op *Operator) oscillate(x float64) float64 {
	if op.Reverse {

		return -WaveFuncs[int(op.WaveformId)%len(WaveFuncs)](op, x)
	}
	return WaveFuncs[int(op.WaveformId)%len(WaveFuncs)](op, x)

}

func (op *Operator) getFB() float64 {
	return float64(op.Feedback) * (float64(op.Prev) / float64(6*op.Mult))
}

func (op *Operator) getFB3() float64 {
	return float64(op.Feedback) * (float64(op.Prev) / float64(op.Mult / (SynthContext.WaveLen * SynthContext.Oversample)))
}

func (op *Operator) getVolume() float64 {
	if(!op.IsEnvelopeEnabled) {
		return float64(op.Tl)
	}
	if op.UseCustomVolEnv {
		return op.customEnv()
	}
	return op.adsr()
}

func (op *Operator) customEnv() float64 {
	index := int(ClampF64(0.0, float64(SynthContext.Macro), float64(len(op.VolEnv))))
	return (float64(op.Tl) * (volROM[op.VolEnv[index] & 0b11111111]))
}

func (op *Operator) adsr() float64 {
	macro := SynthContext.Macro
	// Attack
	if macro <= op.Adsr.Attack {
		if op.Adsr.Attack <= 0 {
			return float64(op.Tl)
		}
		return LinearInterpolation(0, 0, float64(op.Adsr.Attack), float64(op.Tl),float64(macro))
	}

	// Decay and Sustain
	if macro > op.Adsr.Attack && macro < (op.Adsr.Attack + op.Adsr.Decay) {
		if op.Adsr.Decay <= 0 {
			return float64(op.Adsr.Sustain)
		}
		return LinearInterpolation(float64(op.Adsr.Attack), float64(op.Tl), float64(op.Adsr.Attack + op.Adsr.Decay), float64(op.Adsr.Sustain), float64(macro))
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
		return LinearInterpolation(0, float64(op.PwmAdsr.Start), float64(op.PwmAdsr.Attack), float64(op.DutyCycle),float64(macro))
	}

	// Decay and Sustain
	if macro > op.PwmAdsr.Attack && macro < (op.PwmAdsr.Attack + op.PwmAdsr.Decay) {
		if op.PwmAdsr.Decay <= 0 {
			return float64(op.PwmAdsr.Sustain)
		}
		return LinearInterpolation(float64(op.PwmAdsr.Attack), float64(op.DutyCycle), float64(op.PwmAdsr.Attack + op.PwmAdsr.Decay), float64(op.PwmAdsr.Sustain), float64(macro))
	}
	return float64(op.PwmAdsr.Sustain)
}

func (op *Operator) getDutyCycle() float64 {
	if(op.PwmAdsrEnabled) {
		return op.pwmAdsr()
	}
	return float64(op.DutyCycle)
}


func LinearInterpolation(x1 float64, y1 float64, x2 float64, y2 float64, x float64) float64 {
	slope := (y2 - y1) / (x2 - x1)
	return y1 + (slope * (x - x1))
}

func adsrFilter() float64 {
	macro := SynthContext.Macro
	// Attack
	if macro <= SynthContext.FilterAttack {
		if SynthContext.FilterAttack <= 0 {
			return float64(SynthContext.Cutoff)
		}
		return LinearInterpolation(0, float64(SynthContext.FilterStart), float64(SynthContext.FilterAttack), float64(SynthContext.Cutoff),float64(macro))
	}

	// Decay and Sustain
	if macro > SynthContext.FilterAttack && macro < (SynthContext.FilterAttack + SynthContext.FilterDecay) {
		if SynthContext.FilterDecay <= 0 {
			return float64(SynthContext.FilterSustain)
		}
		return LinearInterpolation(float64(SynthContext.FilterAttack), float64(SynthContext.Cutoff), float64(SynthContext.FilterAttack + SynthContext.FilterDecay), float64(SynthContext.FilterSustain), float64(macro))
	}
	return float64(SynthContext.FilterSustain)
}

func ApplyAlg(alg int) {
	switch alg {
	case 0:
		SynthContext.ModMatrix = [][]bool{
			{false, false, false, false},
			{true, false, false, false},
			{false, true, false, false},
			{false, false, true, false},
		}
		SynthContext.OpOutputs = []float32{0, 0, 0, 1}
	case 1:
		SynthContext.ModMatrix = [][]bool{
			{false, false, false, false},
			{false, false, false, false},
			{true, true, false, false},
			{false, false, true, false},
		}
		SynthContext.OpOutputs = []float32{0, 0, 0, 1}
	case 2:
		SynthContext.ModMatrix = [][]bool{
			{false, false, false, false},
			{false, false, false, false},
			{false, true, false, false},
			{true, false, true, false},
		}
		SynthContext.OpOutputs = []float32{0, 0, 0, 1}
	case 3:
		SynthContext.ModMatrix = [][]bool{
			{false, false, false, false},
			{true, false, false, false},
			{false, false, false, false},
			{false, false, true, false},
		}
		SynthContext.OpOutputs = []float32{0, 1, 0, 1}
	case 4:
		SynthContext.ModMatrix = [][]bool{
			{false, false, false, false},
			{true, false, false, false},
			{true, false, false, false},
			{true, false, false, false},
		}
		SynthContext.OpOutputs = []float32{0, 1, 1, 1}
	case 5:
		SynthContext.ModMatrix = [][]bool{
			{false, false, false, false},
			{true, false, false, false},
			{false, false, false, false},
			{false, false, false, false},
		}
		SynthContext.OpOutputs = []float32{0, 1, 1, 1}
	case 6:
		SynthContext.ModMatrix = [][]bool{
			{false, false, false, false},
			{false, false, false, false},
			{false, false, false, false},
			{false, false, false, false},
		}
		SynthContext.OpOutputs = []float32{1, 1, 1, 1}
	case 7:
		SynthContext.ModMatrix = [][]bool{
			{false, false, false, false},
			{true, false, false, false},
			{true, false, false, false},
			{false, true, true, false},
		}
		SynthContext.OpOutputs = []float32{0, 0, 0, 1}
	case 8:
		SynthContext.ModMatrix = [][]bool{
			{false, false, false, false},
			{true, false, false, false},
			{false, true, false, false},
			{false, true, false, false},
		}
		SynthContext.OpOutputs = []float32{0, 0, 1, 1}
	case 9:
		SynthContext.ModMatrix = [][]bool{
			{false, false, false, false},
			{true, false, false, false},
			{false, true, false, false},
			{false, false, false, false},
		}
		SynthContext.OpOutputs = []float32{0, 0, 1, 1}
	case 10:
		SynthContext.ModMatrix = [][]bool{
			{false, false, false, false},
			{true, false, false, false},
			{true, false, false, false},
			{false, false, false, false},
		}
		SynthContext.OpOutputs = []float32{0, 1, 1, 1}
	case 11:
		SynthContext.ModMatrix = [][]bool{
			{false, false, false, false},
			{false, false, false, false},
			{false, false, false, false},
			{true, true, true, false},
		}
		SynthContext.OpOutputs = []float32{0, 0, 0, 1}
	}
}

func GenerateWaveStr() *string {
	str := ""
	for _, n := range WaveOutput {
		str += strconv.Itoa(n) + " "
	}
	str += ";"
	return &str
}

var WaveSeqStr = ""
func GenerateWaveSeqStr() {
	str := ""
	tmpMac := SynthContext.Macro
	for i := 0; i < int(SynthContext.MacLen); i++ {
		SynthContext.Macro = int32(i)
		Synthesize()
		str += *GenerateWaveStr() + "\n"
	}
	SynthContext.Macro = tmpMac
	Synthesize()
	WaveSeqStr = str
}

type LowpassFilter struct {
    sampleRate float64
    cutoffFreq float64
    resonance  float64
    q          float64
    gain       float64
    a          []float64
    b          []float64
    x          []float64
    y          []float64
}

func NewLowpassFilter(sampleRate, cutoffFreq, resonance, gain float64) *LowpassFilter {
    // Calculate filter coefficients
    w0 := 2 * math.Pi * cutoffFreq / sampleRate
    alpha := math.Sin(w0) / (2 * resonance)
    cosw0 := math.Cos(w0)
    a0 := 1 + alpha
    a1 := -2 * cosw0 / a0
    a2 := (1 - alpha) / a0
    b0 := (1 - cosw0) / 2 / a0
    b1 := (1 - cosw0) / a0
    b2 := (1 - cosw0) / 2 / a0
    
    return &LowpassFilter{
        sampleRate: sampleRate,
        cutoffFreq: cutoffFreq,
        resonance:  resonance,
        q:          1 / (2 * resonance),
        gain:       gain,
        a:          []float64{-a1, -a2},
        b:          []float64{b0, b1, b2},
        x:          make([]float64, 4),
        y:          make([]float64, 4),
    }
}

func (lpf *LowpassFilter) Process(x float64) float64 {
    // Shift input and output samples
    lpf.x[3] = lpf.x[2]
    lpf.x[2] = lpf.x[1]
    lpf.x[1] = lpf.x[0]
    lpf.x[0] = x / lpf.gain

    lpf.y[3] = lpf.y[2]
    lpf.y[2] = lpf.y[1]
    lpf.y[1] = lpf.y[0]

    // Apply filter
    lpf.y[0] = (lpf.b[0]*lpf.x[0] + lpf.b[1]*lpf.x[1] + lpf.b[2]*lpf.x[2] -
        lpf.a[1]*lpf.y[1] - lpf.a[2]*lpf.y[2]) / lpf.a[0]

    lpf.x[3] = lpf.x[2]
    lpf.x[2] = lpf.x[1]
    lpf.x[1] = lpf.x[0]
    lpf.y[3] = lpf.y[2]
    lpf.y[2] = lpf.y[1]
    lpf.y[1] = lpf.y[0]

    return lpf.y[0] * lpf.gain
}

func lowpassFiltering222(cutoffFreq float64, resonance float64, sampleRate float64, input []float64) []float64 {
	gain := 1.0
    lpf := NewLowpassFilter(sampleRate, cutoffFreq, resonance, gain)
    output := make([]float64, len(input))

    for i := 0; i < len(input); i++ {
        output[i] = lpf.Process(input[i])
    }

    return output
}

type ButterworthFilter struct {
	xn     [2]float64
	yn     [2]float64
	a      [3]float64
	b      [3]float64
}

func NewButterworthFilter(cutoff, sampleRate, resonance float64) *ButterworthFilter {
	bwf := &ButterworthFilter{}
	omegaC := 2.0 * math.Pi * cutoff / sampleRate
	if resonance > 0.0 {
		d := 1.0 / resonance
		bwf.a[0] = 1.0 + omegaC*d
		bwf.a[1] = -2.0 * math.Cos(omegaC)
		bwf.a[2] = 1.0 - omegaC*d
		bwf.b[0] = (1.0 - math.Cos(omegaC)) / 2.0
		bwf.b[1] = 1.0 - math.Cos(omegaC)
		bwf.b[2] = (1.0 - math.Cos(omegaC)) / 2.0
	} else {
		bwf.a[0] = 1.0 + 2.0*math.Cos(omegaC) + math.Pow(omegaC, 2.0)
		bwf.a[1] = -2.0 * (1.0 + math.Pow(omegaC, 2.0))
		bwf.a[2] = 1.0 - 2.0*math.Cos(omegaC) + math.Pow(omegaC, 2.0)
		bwf.b[0] = math.Pow(omegaC, 2.0) / bwf.a[0]
		bwf.b[1] = 2.0 * bwf.b[0]
		bwf.b[2] = bwf.b[0]
	}
	return bwf
}

func (bwf *ButterworthFilter) Filter(sample float64) float64 {
	bwf.xn[0], bwf.xn[1] = bwf.xn[1], bwf.xn[0]
	bwf.xn[1] = sample
	bwf.yn[0], bwf.yn[1] = bwf.yn[1], bwf.yn[0]
	bwf.yn[1] = (bwf.b[0]*bwf.xn[1] + bwf.b[1]*bwf.xn[0] + bwf.b[2]*bwf.xn[0] - bwf.a[1]*bwf.yn[0] - bwf.a[2]*bwf.yn[1]) / bwf.a[0]
	return bwf.yn[1]
}


// Apply a lowpass filter to the input wavetable and return the filtered wavetable.
func lowpassFiltering(cutoff float64, resonance float64, pitch float64, wavetable []float64) []float64 {
    sampleRate := float64(len(wavetable)) / pitch
	filter := NewButterworthFilter(cutoff, sampleRate, resonance)
	filtered := make([]float64, len(wavetable))
	for i := 0; i < len(wavetable); i++ {
		filtered[i] = filter.Filter(wavetable[i])
	}
	for i := 0; i < len(wavetable); i++ {
		filtered[i] = filter.Filter(wavetable[i])
	}
	return filtered
}



func lowpassFiltering2(cutoffFrequency float64, resonance float64, sampleRate float64, input []float64) []float64 {
    // Calculate filter coefficients
	sampleRate = float64(len(input)) * sampleRate
	
	omegaC := 2.0 * math.Pi * cutoffFrequency / sampleRate
    s := math.Sin(omegaC)
    alpha := s / (2.0 * math.Sqrt(2.0))
    beta := 1.0 - alpha

    // Initialize filter state variables
    x1 := 0.0
    x2 := 0.0
    y1 := 0.0
    y2 := 0.0

    // Apply filter to input signal
    output := make([]float64, len(input))
    for i := range input {
        x0 := input[i] - resonance*y1
        y0 := alpha*(x0+x2) + beta*y1 + alpha*y2
        x2 = x1
        x1 = x0
        y2 = y1
        y1 = y0
        output[i] = y0
    }

    return output
}

func notetofreq(n float64) float64 {
	return 440 * math.Pow(2, (n-69)/12)
}

type biquadFilter struct {
	order int
	a0, a1, a2, b1, b2 float64 // Factors
	filterCutoff, Q, peakGain float64 // Cutoff, Q, and peak gain
	z1, z2 float64 // poles
	A, w0, w1, w2, d1, d2 []float64
	ep float64
}

func buildBQFilter(cutoff float64) *biquadFilter {
	filter := new(biquadFilter)
	filter.order = 4
	sampleRate := notetofreq(float64(SynthContext.Pitch)) * float64(SynthContext.WaveLen * SynthContext.Oversample)
	var norm float64
	K := math.Tan(math.Pi * cutoff/sampleRate)
	switch SynthContext.FilterType {
	case 0: // LPF BQ
		if SynthContext.Resonance == 0 {
			norm = 0
		} else {
			norm = 1 / (1 + K / float64(SynthContext.Resonance) + K*K)
		}
		filter.a0 = K * K * norm
		filter.a1 = 2 * filter.a0
		filter.a2 = filter.a0
		filter.b1 = 2 * (K * K - 1) * norm
		if SynthContext.Resonance == 0 {
			filter.b2 = 0
		} else {
			filter.b2 = (1 - K / float64(SynthContext.Resonance) + K * K) * norm
		} 
	case 1: // HPF BQ
	if SynthContext.Resonance == 0 {
			norm = 0
		} else {
			norm = 1 / (1 + K / float64(SynthContext.Resonance) + K*K)
		}
		filter.a0 = 1 * norm
		filter.a1 = -2 * filter.a0
		filter.a2 = filter.a0
		filter.b1 = 2 * (K * K - 1) * norm
		if SynthContext.Resonance == 0 {
			filter.b2 = 0
		} else {
			filter.b2 = (1 - K / float64(SynthContext.Resonance) + K * K) * norm
		}
	case 2: // BPF BQ
		if SynthContext.Resonance == 0 {
			norm = 0
		} else {
			norm = 1 / (1 + K / float64(SynthContext.Resonance) + K*K)
		}
		filter.a0 = K / float64(SynthContext.Resonance) * norm
		filter.a1 = 0
		filter.a2 = -filter.a0
		filter.b1 = 2 * (K * K - 1) * norm
		if SynthContext.Resonance == 0 {
			filter.b2 = 0
		} else {
			filter.b2 = (1 - K / float64(SynthContext.Resonance) + K * K) * norm
		}
	case 3: // BSF BQ
		if SynthContext.Resonance == 0 {
			norm = 0
		} else {
			norm = 1 / (1 + K / float64(SynthContext.Resonance) + K*K)
		}
		filter.a0 = (1 + K * K) * norm
		filter.a1 = 2 * (K * K - 1) * norm
		filter.a2 = filter.a0
		filter.b1 = filter.a1
		if SynthContext.Resonance == 0 {
			filter.b2 = 0
		} else {
			filter.b2 = (1 - K / float64(SynthContext.Resonance) + K * K) * norm
		}
	case 4: // AP BQ
		aa := (K - 1.0) / (K + 1.0)
		bb := -math.Cos(math.Pi * cutoff/sampleRate)
		filter.a0 = -aa
		filter.a1 = bb*(1.0 - aa)
		filter.a2 = 1.0
		filter.b1 = filter.a1
    	filter.b2 = filter.a0
	}
	return filter
}

const preWarm = 8

func (bqFilter *biquadFilter) processBQFilter(x float64) float64 {
	output := x * bqFilter.a0 + bqFilter.z1
	bqFilter.z1 = x * bqFilter.a1 + bqFilter.z2 - bqFilter.b1 * output
	bqFilter.z2 = x * bqFilter.a2 - bqFilter.b2 * output
	return output
}

func filter(wavetable []float64) []float64 {
	sampleRate := notetofreq(float64(SynthContext.Pitch)) * float64(SynthContext.WaveLen * SynthContext.Oversample)
	
	myCutoff := float64(SynthContext.Cutoff)
	if(SynthContext.FilterAdsrEnabled) {
		myCutoff = adsrFilter()
	}
	
	filterCutoff := 5 * math.Pow(10, float64(myCutoff) * 3);
	filterCutoff = math.Min(sampleRate/2,filterCutoff);
	
	filter := buildBQFilter(filterCutoff)
	outWave := make([]float64, len(wavetable))
	for i := 0; i < len(wavetable) * preWarm; i++ {
		filter.processBQFilter(wavetable[i % len(wavetable)])
	}
	for i := 0; i < len(wavetable); i++ {
		outWave[i] = filter.processBQFilter(wavetable[i % len(wavetable)])
	}
	return outWave
}

var FilterTypes = []string {
	"Biquad Lowpass",
	"Biquad Highpass",
	"Biquad Bandpass",
	"Biquad Bandstop",
	"Biquad Allpass",
}

func CreateFTIN163(macro bool) error {
	tmpLen := SynthContext.WaveLen
	tmpHei := SynthContext.WaveHei

	if(tmpLen > 240) {
		SynthContext.WaveLen = 240
	}
	if(tmpHei > 15) {
		SynthContext.WaveHei = 15
	}
	Synthesize()


	fpath, errZen := zenity.SelectFileSave(
		zenity.ConfirmOverwrite(),
		zenity.Filename("output"),
		zenity.FileFilters{
			{".fti files", []string{"*.fti"}, false},
		})
	if(errZen == zenity.ErrCanceled) {
		return errZen
	}
	if !strings.HasSuffix(fpath, ".fti") {
        fpath += ".fti"
    }
	file, err := os.Create(fpath)
	name := filepath.Base(fpath)
	println(name)
	name = name[:len(name)-len(filepath.Ext(fpath))]
	if (len(name) > 127) {
		name = name[:127]
	}

	// It's time to build the FUW file!!
	output := []byte {
		'F','T','I','2','.','4', // Header
		0x05, // Instrument type (N163 here)
		byte(len(name)), 0, 0, 0, // Length of name string
	}

	for _, character := range name {
		output = append(output, byte(character & 0xFF))
	}

	output = append(output, 0x05) // I dunno what it is

	output = append(output, 0x00) // We disable volume envelope...
	output = append(output, 0x00)
	output = append(output, 0x00)
	output = append(output, 0x00)

	waveMacroEnabled := byte(0)
	if(macro) {
		waveMacroEnabled = 1
	}
	output = append(output, waveMacroEnabled) // We enable waveform envelope!

	if(macro) {
		macroLen := int(math.Min(float64(SynthContext.MacLen), 64))
		output = append(output, byte(macroLen))

		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)

		output = append(output, 0xFF)
		output = append(output, 0xFF)
		output = append(output, 0xFF)
		output = append(output, 0xFF)
		output = append(output, 0xFF)
		output = append(output, 0xFF)
		output = append(output, 0xFF)
		output = append(output, 0xFF)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)

		

		for i := 0; i < macroLen; i++ {
			output = append(output, byte(i))
		}
		output = append(output, byte(SynthContext.WaveLen))
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)

		output = append(output, byte(macroLen))
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)

		tmpMacro := SynthContext.Macro
		for m := int32(0); m < int32(macroLen); m++ {
			SynthContext.Macro = m
			Synthesize()

			for _, sample := range WaveOutput {
				output = append(output, byte(sample & 0x0F))
			}
		}
		SynthContext.Macro = tmpMacro
		Synthesize()
	} else {
		output = append(output, byte(SynthContext.WaveLen))
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x01)
		output = append(output, 0x00)
		output = append(output, 0x00)
		output = append(output, 0x00)
		for _, sample := range WaveOutput {
			output = append(output, byte(sample & 0x0F))
		}
	}
	SynthContext.WaveLen = tmpLen
	SynthContext.WaveHei = tmpHei
	Synthesize()

	i, err := file.Write(output)
	i = int(i)
    if err != nil {
        return err
    }
    return nil
}

const FURNACE_FORMAT_VER uint16 = 143

func CreateFUW() error {
	fpath, errZen := zenity.SelectFileSave(
		zenity.ConfirmOverwrite(),
		zenity.Filename("output"),
		zenity.FileFilters{
			{".fuw files", []string{"*.fuw"}, false},
		})
	if(errZen == zenity.ErrCanceled) {
		return errZen
	}
	if !strings.HasSuffix(fpath, ".fuw") {
        fpath += ".fuw"
    }
	file, err := os.Create(fpath)

	var size uint32 = 1 + 4 + 4 + 4 +uint32(4 * len(WaveOutput)) 
	const HEADER_SIZE = 16 + 2 + 2 + 4 + 4 + 1 + 4 + 4 + 4

	// It's time to build the FUW file!!
	output := []byte {
		'-','F','u','r','n','a','c','e',' ','w','a','v','e','t','a','-', // Header, 16 bytes
		byte(FURNACE_FORMAT_VER & 0xFF), byte(FURNACE_FORMAT_VER >> 8), // Format version, 2 bytes
		'0', '0', // Reserved, 2 bytes
		'W','A','V','E', // WAVE chunk, 4 bytes
		byte(size & 0xFF), byte((size >> 8) & 0xFF), byte((size >> 16) & 0xFF), byte((size >> 24)), // Size of chunk, 4 bytes
		0, //empty string, 1 byte
		byte(SynthContext.WaveLen & 0xFF), byte((SynthContext.WaveLen >> 8) & 0xFF), byte((SynthContext.WaveLen >> 16) & 0xFF), byte((SynthContext.WaveLen >> 24)), // Wave length, 4 bytes
		0, 0, 0, 0, // Reserved, 4 bytes
		byte(SynthContext.WaveHei & 0xFF), byte((SynthContext.WaveHei >> 8) & 0xFF), byte((SynthContext.WaveHei >> 16) & 0xFF), byte((SynthContext.WaveHei >> 24)), // Wave height, 4 bytes
	}

	// Appending Data
	for _, sample := range WaveOutput {
		output = append(output, byte(sample & 0xFF))
		output = append(output, byte((sample >> 8) & 0xFF))
		output = append(output, byte((sample >> 16) & 0xFF))
		output = append(output, byte(sample >> 24))
	}

	i, err := file.Write(output)
	i = int(i)
    if err != nil {
        return err
    }
    return nil
}

func getSampleRate() int {
    return int(math.Floor((440 * float64(len(WaveOutput))) / 2.0))
}

func createWav(path string, macro bool, bits16 bool) error {
	file, err := os.Create(path)
	var bufLen int
if macro {
    bufLen = len(WaveOutput) * int(SynthContext.MacLen)
} else {
    bufLen = len(WaveOutput)
}

var frames = 1
if(macro) {
	frames = int(SynthContext.MacLen)
}

var bits byte = 8
if bits16 {
    bits = 16
}

chunkSize := 0
if(bits16) {
	chunkSize = 36 + (bufLen * frames) * 2
} else {
	chunkSize = 36 + (bufLen * frames)
}

subchunkSize := 0
if(bits16) {
	subchunkSize = (bufLen * frames) * 2
} else {
	subchunkSize = (bufLen * frames)
}

sampleRate := getSampleRate()
byteRate := 0
if(bits16) {
	byteRate = (sampleRate * 16) / 8
} else {
	byteRate = sampleRate
}

intBuffer := []byte{
	0x52, 0x49, 0x46, 0x46, // ChunkID: "RIFF" in ASCII form, big endian
	byte(chunkSize & 0xFF), byte((chunkSize >> 8) & 0xFF), byte((chunkSize >> 16) & 0xFF), byte(chunkSize >> 24), // ChunkSize - will be filled later, 
	0x57, 0x41, 0x56, 0x45, // Format: "WAVE" in ASCII form
	0x66, 0x6d, 0x74, 0x20, // Subchunk1ID: "fmt " in ASCII form
	0x10, 0x00, 0x00, 0x00, // Subchunk1Size: 16 for PCM
	0x01, 0x00,             // AudioFormat: PCM = 1
	0x01, 0x00,             // NumChannels: Mono = 1
	byte(sampleRate & 0xFF), byte((sampleRate >> 8) & 0xFF), byte((sampleRate >> 16) & 0xFF), byte(sampleRate >> 24), // SampleRate: 44100 Hz - little endian
	byte(byteRate & 0xFF), byte((byteRate >> 8) & 0xFF), byte((byteRate >> 16) & 0xFF), byte(byteRate >> 24), // ByteRate: 44100 * 1 * 16 / 8 - little endian
	byte(bits / 8), 0x00,             // BlockAlign: 1 * 16 / 8 - little endian
	bits, 0x00,             // BitsPerSample: 16 bits per sample
	0x64, 0x61, 0x74, 0x61, // Subchunk2ID: "data" in ASCII form
	byte(subchunkSize & 0xFF), byte((subchunkSize >> 8) & 0xFF), byte((subchunkSize >> 16) & 0xFF), byte(subchunkSize >> 24), // Subchunk2Size - will be filled later
}

    var output []uint8
    if macro {
        tmpMac := SynthContext.Macro
        for i := 0; i < int(SynthContext.MacLen); i++ {
            SynthContext.Macro = int32(i)
			Synthesize()
            for _, sample := range WaveOutput {
				var tmp float64;
				if(bits16) {

					if (SynthContext.WaveHei & 0x0001) == 1 {
						tmp = float64(sample) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
					} else {
						tmp = float64(sample) / (float64(SynthContext.WaveHei) / 2.0)
					}
					myOut := int16(math.Round((tmp-1)*float64((1<<(16-1))-1)))
					output = append(output, byte(myOut >> 8))
					output = append(output, byte(myOut & 0xFF))
					continue
				}
				if (SynthContext.WaveHei & 0x0001) == 1 {
					tmp = float64(sample) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
				} else {
					tmp = float64(sample) / (float64(SynthContext.WaveHei) / 2.0)
				}
				myOut := int16(math.Round((tmp-0)*float64((1<<(8-1))-1)))
				output = append(output, byte(myOut))                
            }
        }
		SynthContext.Macro = tmpMac
		Synthesize()
    } else {
        for _, sample := range WaveOutput {
			var tmp float64;
			if(bits16) {

				if (SynthContext.WaveHei & 0x0001) == 1 {
					tmp = float64(sample) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
				} else {
					tmp = float64(sample) / (float64(SynthContext.WaveHei) / 2.0)
				}
				myOut := int16(math.Round((tmp-1)*float64((1<<(16-1))-1)))
				output = append(output, byte(myOut & 0xFF))
				output = append(output, byte(myOut >> 8))
				continue
			}
			if (SynthContext.WaveHei & 0x0001) == 1 {
				tmp = float64(sample) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
			} else {
				tmp = float64(sample) / (float64(SynthContext.WaveHei) / 2.0)
			}
			myOut := int16(math.Round((tmp-0)*float64((1<<(8-1))-1)))
			output = append(output, byte(myOut))                
		}
    }

    for _, sample := range output {
        intBuffer = append(intBuffer, sample)
    }

    i, err := file.Write(intBuffer)
	i = int(i)
    if err != nil {
        return err
    }
    return nil
}

func SaveRaw(macro bool, mode int) error {
	path, errZen := zenity.SelectFileSave(
		zenity.ConfirmOverwrite(),
		zenity.Filename("output"),
		zenity.FileFilters{
			{".raw files", []string{"*.raw"}, false},
		})
	if(errZen == zenity.ErrCanceled) {
		return errZen
	}
	if !strings.HasSuffix(path, ".raw") {
        path += ".raw"
    }

	file, err := os.Create(path)
	if(err != nil) {
		return err
	}

	if(mode == 2) {
		if(len(WaveOutput) & 1 > 0) {
			zenity.Error("Only an even length is accepted for 4-bits export.")
			return nil
		}
	}

	bufLen := len(WaveOutput)
	if macro {
		bufLen = bufLen * int(SynthContext.MacLen)
	}
	var output []byte
	if macro {
        tmpMac := SynthContext.Macro
        for i := 0; i < int(SynthContext.MacLen); i++ {
            SynthContext.Macro = int32(i)
			Synthesize()
			if(mode == 2) {
				for i := 0; i < len(WaveOutput); i += 2 {
					var tmp1 float64
					var tmp2 float64
					if (SynthContext.WaveHei & 0x0001) == 1 {
						tmp1 = float64(WaveOutput[i]) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
						tmp2 = float64(WaveOutput[i + 1]) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
					} else {
						tmp1 = float64(WaveOutput[i]) / (float64(SynthContext.WaveHei) / 2.0)
						tmp2 = float64(WaveOutput[i + 1]) / (float64(SynthContext.WaveHei) / 2.0)
					}
					
					myOut1 := int16(math.Round((tmp1-0)*float64((1<<(4-1))-1)))
					myOut2 := int16(math.Round((tmp2-0)*float64((1<<(4-1))-1)))
					output = append(output, byte((myOut1 >> 4) | myOut2 & 0xF))
				}
				continue
			}
            for _, sample := range WaveOutput {
				var tmp float64;
				
				if (SynthContext.WaveHei & 0x0001) == 1 {
					tmp = float64(sample) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
				} else {
					tmp = float64(sample) / (float64(SynthContext.WaveHei) / 2.0)
				}
				switch mode {
				case 0: // Normalized RAW
				myOut := int16(math.Round((tmp-0)*float64((1<<(8-1))-1)))
				output = append(output, byte(myOut))
				case 1: // Non Normalized
				output = append(output, byte(sample))
				}
				                
            }
        }
		SynthContext.Macro = tmpMac
		Synthesize()
    } else {
        if(mode == 2) {
			for i := 0; i < len(WaveOutput); i += 2 {
				var tmp1 float64
				var tmp2 float64
				if (SynthContext.WaveHei & 0x0001) == 1 {
					tmp1 = float64(WaveOutput[i]) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
					tmp2 = float64(WaveOutput[i + 1]) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
				} else {
					tmp1 = float64(WaveOutput[i]) / (float64(SynthContext.WaveHei) / 2.0)
					tmp2 = float64(WaveOutput[i + 1]) / (float64(SynthContext.WaveHei) / 2.0)
				}
				
				myOut1 := int16(math.Round((tmp1-0)*float64((1<<(4-1))-1)))
				myOut2 := int16(math.Round((tmp2-0)*float64((1<<(4-1))-1)))
				output = append(output, byte((myOut1 << 4) | myOut2 & 0xF))
			}
		}
		for _, sample := range WaveOutput {
			var tmp float64;
			
			if (SynthContext.WaveHei & 0x0001) == 1 {
				tmp = float64(sample) / ((float64(SynthContext.WaveHei) / 2.0) + 0.5)
			} else {
				tmp = float64(sample) / (float64(SynthContext.WaveHei) / 2.0)
			}
			switch mode {
			case 0: // Normalized RAW
			myOut := int16(math.Round((tmp-0)*float64((1<<(8-1))-1)))
			output = append(output, byte(myOut))
			case 1: // Non Normalized
			output = append(output, byte(sample))
			}
							
		}
    }
	

	_, err2 := file.Write(output)
	if(err2 != nil) {
		return err2
	}
	return nil
}

func SaveTxt(macro bool) error {

	path, errZen := zenity.SelectFileSave(
		zenity.ConfirmOverwrite(),
		zenity.Filename("output"),
		zenity.FileFilters{
			{".txt files", []string{"*.txt"}, false},
		})
	if(errZen == zenity.ErrCanceled) {
		return errZen
	}
	if !strings.HasSuffix(path, ".txt") {
        path += ".txt"
    }

	str := ""
	file, err := os.Create(path)
	if(err != nil) {
		return err
	}
	if(macro) {
		GenerateWaveSeqStr()
		str = WaveSeqStr
	} else {
		str = *GenerateWaveStr()
	}
	_, err2 := file.WriteString(str)
	if(err2 != nil) {
		return err2
	}
	return nil
}

func SaveFile(macro bool, bits16 bool) {
	path, err := zenity.SelectFileSave(
		zenity.ConfirmOverwrite(),
		zenity.Filename("output"),
		zenity.FileFilters{
			{".WAV files", []string{"*.wav"}, false},
		})
	if(err == zenity.ErrCanceled) {
		return
	}
	if !strings.HasSuffix(path, ".wav") {
        path += ".wav"
    }
	createWav(path, macro, bits16)
}

const (
	toneHz   = 440
	sampleHz = 48000
	
)

var phase = 0.0

func phaseAcc(len int) float64 {
	freqTable := float64(sampleHz) / float64(len)
	playfreq:= float64(toneHz) / freqTable
	phase = math.Mod((phase + playfreq), float64(len))
	return math.Mod(phase, float64(len))
}

//export Wavetable
func Wavetable(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Int16)(unsafe.Pointer(&hdr))
	
	for i := 0; i < (n / 2); i ++ {
		sample := float64(0)
		ind := 0
		if(!SynthContext.SongPlaying) {
			buf[i] = (0)
			continue
		}
		if(SynthContext.SongPlaying && len(WaveOutput) > 0 && WaveOutput != nil){
			ind = int(phaseAcc(len(WaveOutput)))
		}
		if(SynthContext.SongPlaying && len(WaveOutput) > 0 && WaveOutput != nil) {
			sample = float64(WaveOutput[ind % len(WaveOutput)])
		}

		sample = sample * (255 / float64(SynthContext.WaveHei))
		s2 := C.Int16(sample) - 128
		buf[i] = C.Int16((s2) << 4)
	}
}

func InitAudio() {
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		println(err)
		return
	}

	spec := &sdl.AudioSpec{
		Freq:     sampleHz,
		Format:   sdl.AUDIO_S16,
		Channels: 1,
		Samples:  512,
		Callback: sdl.AudioCallback(C.Wavetable),
	}
	if err := sdl.OpenAudio(spec, nil); err != nil {
		println(err)
		return
	}
	sdl.PauseAudio(false)
}