package kurumi

import "math"

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

// Buffers used for output in the Matrix
var samples = make([]float64, 4)

type ModulationMode int

const (
	fmMod   ModulationMode = 0
	orMod   ModulationMode = 1
	xorMod  ModulationMode = 2
	andMod  ModulationMode = 3
	nandMod ModulationMode = 4
	addMod  ModulationMode = 5
	subMod  ModulationMode = 6
	mulMod  ModulationMode = 7
)

func logicMod(x float64, modValue float64, opId int) float64 {
	op := &SynthContext.Operators[opId]
	switch ModulationMode(op.ModMode) {
	case fmMod: // FM
		return op.oscillate(x+modValue+op.getFB()) * op.getVolume()
	case orMod: // OR
		a := int(math.Round((modValue + 1) * 32767.5))
		b := int(math.Round(((op.oscillate(x) * op.getVolume()) + (1 * op.getVolume())) * 32767.5))
		return (float64(a|b) / 32767.5) - (1 * op.getVolume())
	case xorMod: // XOR
		a := int(math.Round((modValue + 1) * 32767.5))
		b := int(math.Round(((op.oscillate(x) * op.getVolume()) + (1 * op.getVolume())) * 32767.5))
		return (float64(a^b) / 32767.5) - (1 * op.getVolume())
	case andMod: // AND
		a := int(math.Round((modValue + 1) * 32767.5))
		b := int(math.Round(((op.oscillate(x) * op.getVolume()) + (1 * op.getVolume())) * 32767.5))
		return (float64(a&b) / 32767.5) - (1 * op.getVolume())
	case nandMod: // NAND
		a := int(math.Round((modValue + 1) * 32767.5))
		b := int(math.Round(((op.oscillate(x) * op.getVolume()) + (1 * op.getVolume())) * 32767.5))
		return float64(^int(a&b))/32767.5 - (1 * op.getVolume())
	case addMod: // ADD
		return modValue + (op.oscillate(x) * op.getVolume())
	case subMod: // SUB
		return op.oscillate(x)*op.getVolume() - modValue
	case mulMod: // MUL
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
