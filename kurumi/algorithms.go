package kurumi

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
