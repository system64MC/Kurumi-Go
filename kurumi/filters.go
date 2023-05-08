package kurumi

import "math"

func notetofreq(n float64) float64 {
	return 440 * math.Pow(2, (n-69)/12)
}

type biquadFilter struct {
	order                     int
	a0, a1, a2, b1, b2        float64 // Factors
	filterCutoff, Q, peakGain float64 // Cutoff, Q, and peak gain
	z1, z2                    float64 // poles
	A, w0, w1, w2, d1, d2     []float64
	ep                        float64
}

func buildBQFilter(cutoff float64) *biquadFilter {
	filter := new(biquadFilter)
	filter.order = 4
	sampleRate := notetofreq(float64(SynthContext.Pitch)) * float64(SynthContext.WaveLen*SynthContext.Oversample)
	var norm float64
	K := math.Tan(math.Pi * cutoff / sampleRate)
	switch SynthContext.FilterType {
	case 0: // LPF BQ
		if SynthContext.Resonance == 0 {
			norm = 0
		} else {
			norm = 1 / (1 + K/float64(SynthContext.Resonance) + K*K)
		}
		filter.a0 = K * K * norm
		filter.a1 = 2 * filter.a0
		filter.a2 = filter.a0
		filter.b1 = 2 * (K*K - 1) * norm
		if SynthContext.Resonance == 0 {
			filter.b2 = 0
		} else {
			filter.b2 = (1 - K/float64(SynthContext.Resonance) + K*K) * norm
		}
	case 1: // HPF BQ
		if SynthContext.Resonance == 0 {
			norm = 0
		} else {
			norm = 1 / (1 + K/float64(SynthContext.Resonance) + K*K)
		}
		filter.a0 = 1 * norm
		filter.a1 = -2 * filter.a0
		filter.a2 = filter.a0
		filter.b1 = 2 * (K*K - 1) * norm
		if SynthContext.Resonance == 0 {
			filter.b2 = 0
		} else {
			filter.b2 = (1 - K/float64(SynthContext.Resonance) + K*K) * norm
		}
	case 2: // BPF BQ
		if SynthContext.Resonance == 0 {
			norm = 0
		} else {
			norm = 1 / (1 + K/float64(SynthContext.Resonance) + K*K)
		}
		filter.a0 = K / float64(SynthContext.Resonance) * norm
		filter.a1 = 0
		filter.a2 = -filter.a0
		filter.b1 = 2 * (K*K - 1) * norm
		if SynthContext.Resonance == 0 {
			filter.b2 = 0
		} else {
			filter.b2 = (1 - K/float64(SynthContext.Resonance) + K*K) * norm
		}
	case 3: // BSF BQ
		if SynthContext.Resonance == 0 {
			norm = 0
		} else {
			norm = 1 / (1 + K/float64(SynthContext.Resonance) + K*K)
		}
		filter.a0 = (1 + K*K) * norm
		filter.a1 = 2 * (K*K - 1) * norm
		filter.a2 = filter.a0
		filter.b1 = filter.a1
		if SynthContext.Resonance == 0 {
			filter.b2 = 0
		} else {
			filter.b2 = (1 - K/float64(SynthContext.Resonance) + K*K) * norm
		}
	case 4: // AP BQ
		aa := (K - 1.0) / (K + 1.0)
		bb := -math.Cos(math.Pi * cutoff / sampleRate)
		filter.a0 = -aa
		filter.a1 = bb * (1.0 - aa)
		filter.a2 = 1.0
		filter.b1 = filter.a1
		filter.b2 = filter.a0
	}
	return filter
}

const preWarm = 8

func (bqFilter *biquadFilter) processBQFilter(x float64) float64 {
	output := x*bqFilter.a0 + bqFilter.z1
	bqFilter.z1 = x*bqFilter.a1 + bqFilter.z2 - bqFilter.b1*output
	bqFilter.z2 = x*bqFilter.a2 - bqFilter.b2*output
	return output
}

func filter(wavetable []float64) []float64 {
	sampleRate := notetofreq(float64(SynthContext.Pitch)) * float64(SynthContext.WaveLen*SynthContext.Oversample)

	myCutoff := float64(SynthContext.Cutoff)
	if SynthContext.FilterAdsrEnabled {
		myCutoff = adsrFilter()
	}

	filterCutoff := 5 * math.Pow(10, float64(myCutoff)*3)
	filterCutoff = math.Min(sampleRate/2, filterCutoff)

	filter := buildBQFilter(filterCutoff)
	outWave := make([]float64, len(wavetable))
	for i := 0; i < len(wavetable)*preWarm; i++ {
		filter.processBQFilter(wavetable[i%len(wavetable)])
	}
	for i := 0; i < len(wavetable); i++ {
		outWave[i] = filter.processBQFilter(wavetable[i%len(wavetable)])
	}
	return outWave
}

var FilterTypes = []string{
	"Biquad Lowpass",
	"Biquad Highpass",
	"Biquad Bandpass",
	"Biquad Bandstop",
	"Biquad Allpass",
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
	xn [2]float64
	yn [2]float64
	a  [3]float64
	b  [3]float64
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

func smooth(arr []float64) []float64 {
	if SynthContext.SmoothWin == 0 {
		return arr
	}
	out := make([]float64, len(arr))
	for i := 0; i < len(arr); i++ {
		smp := 0.0
		for j := -SynthContext.SmoothWin; j <= SynthContext.SmoothWin; j++ {
			smp += arr[moduloFix(i+int(j), len(arr))]
		}
		avg := smp / (float64(SynthContext.SmoothWin*2) + 1)
		out[i] = avg
	}
	return out
}

func getSampleRate() int {
	return int(math.Floor((440 * float64(len(WaveOutput))) / 2.0))
}

func adsrFilter() float64 {
	macro := SynthContext.Macro
	// Attack
	if macro <= SynthContext.FilterAttack {
		if SynthContext.FilterAttack <= 0 {
			return float64(SynthContext.Cutoff)
		}
		return LinearInterpolation(0, float64(SynthContext.FilterStart), float64(SynthContext.FilterAttack), float64(SynthContext.Cutoff), float64(macro))
	}

	// Decay and Sustain
	if macro > SynthContext.FilterAttack && macro < (SynthContext.FilterAttack+SynthContext.FilterDecay) {
		if SynthContext.FilterDecay <= 0 {
			return float64(SynthContext.FilterSustain)
		}
		return LinearInterpolation(float64(SynthContext.FilterAttack), float64(SynthContext.Cutoff), float64(SynthContext.FilterAttack+SynthContext.FilterDecay), float64(SynthContext.FilterSustain), float64(macro))
	}
	return float64(SynthContext.FilterSustain)
}
