package kurumi

import (
	"math"
	"math/rand"
)

func sine(op *Operator, x float64) float64 {
	return math.Sin((x * op.getMult() * 2 * math.Pi) + (float64(op.Phase)*2*math.Pi + (op.getPhase() * math.Pi * 2)))
}
func rectSine(op *Operator, x float64) float64 {
	return math.Max(sine(op, x), 0)
}
func absSine(op *Operator, x float64) float64 {
	return math.Abs(sine(op, x))
}
func quarterSine(op *Operator, x float64) float64 {
	if math.Mod((x*op.getMult()+float64(op.Phase)+op.getPhase()), 0.5) <= 0.25 {
		return absSine(op, x)
	}
	return 0
}
func squishedSine(op *Operator, x float64) float64 {
	if sine(op, x) > 0 {
		return math.Sin((x * op.getMult() * 4 * math.Pi) + (float64(op.Phase)*4*math.Pi + (op.getPhase() * math.Pi * 4)))
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
	a := moduloF64(x*math.Pi*2*op.getMult()+(float64(op.Phase)*math.Pi*2+op.getPhase()*math.Pi*2), math.Pi*2)
	if a >= (math.Pi * width * 2) {
		return -1
	}
	return 1
}

func rectSquare(op *Operator, x float64) float64 {
	a := square(op, x)
	if a < 0 {
		a = 0
	}
	return a
}

func saw(op *Operator, x float64) float64 {
	return math.Atan(math.Tan(x*math.Pi*op.getMult()+(float64(op.Phase)*math.Pi+(op.getPhase()*math.Pi)))) / (math.Pi * 0.5)
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
	if math.Mod((x*op.getMult()+(float64(op.Phase)+op.getPhase())), 0.5) <= 0.25 {
		return triangle(op, x)
	}
	return 0
}
func squishedTriangle(op *Operator, x float64) float64 {
	if sine(op, x) > 0 {
		return math.Asin(math.Sin((x*op.getMult()*4*math.Pi)+(float64(op.Phase)*math.Pi*4+(op.getPhase()*math.Pi*4)))) / (math.Pi / 2)
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

var lsfr = uint16(0b01001_1010_1011_1010)

func lsfrShift() {
	// lsfr = lsfr >> 1
	lsfr = (lsfr << 1) | (((lsfr >> 13) ^ (lsfr >> 14)) & 1)
}

func noise1bitLsfr(op *Operator, x float64) float64 {
	lsfrShift()
	return float64(lsfr&1)*2 - 1
	//return (rand.Float64() * 2) - 1
}

func noise8bitLsfr(op *Operator, x float64) float64 {
	lsfrShift()
	//return float64(lsfr & 1) * 2 - 1
	return float64(lsfr&0xFF)/float64(0x7F) - 1
	//return (rand.Float64() * 2) - 1
}

func noiseRandom(op *Operator, x float64) float64 {
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

func getWTSample(op *Operator, x float64) float64 {
	if op.Morphing {
		a := interpolate(x*float64(len(op.Wavetable))*op.getMult()+(float64(op.Phase)*float64(len(op.Wavetable))+op.getPhase()*float64(len(op.Wavetable))), op, op.Wavetable)
		b := interpolate(x*float64(len(op.MorphWave))*op.getMult()+(float64(op.Phase)*float64(len(op.MorphWave))+op.getPhase()*float64(len(op.MorphWave))), op, op.MorphWave)
		c := float64(SynthContext.Macro) / float64(op.MorphTime)
		return lerp(a, b, c)
	}
	return interpolate(x*float64(len(op.Wavetable))*op.getMult()+(float64(op.Phase)*float64(len(op.Wavetable))+op.getPhase()*float64(len(op.Wavetable))), op, op.Wavetable)
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
	noise1bitLsfr,
	noise8bitLsfr,
	noiseRandom,
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
	"Noise (1 bit, LFSR)",
	"Noise (8 bits, LFSR)",
	"Noise (Random)",
	"Custom",
	"Rect. Custom",
	"Abs. Custom",
	"Cubed Custom",
}
