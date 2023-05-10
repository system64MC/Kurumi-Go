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

type FilterType int

const (
	LowPass  FilterType = 0
	HighPass FilterType = 1
	BandPass FilterType = 2
	BandStop FilterType = 3
	AllPass  FilterType = 4
)

func buildBQFilter(cutoff float64) *biquadFilter {
	filter := new(biquadFilter)
	filter.order = 4
	sampleRate := notetofreq(float64(SynthContext.Pitch)) * float64(SynthContext.WaveLen*SynthContext.Oversample)
	var norm float64
	K := math.Tan(math.Pi * cutoff / sampleRate)
	switch FilterType(SynthContext.FilterType) {
	case LowPass: // LPF BQ
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
	case HighPass: // HPF BQ
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
	case BandPass: // BPF BQ
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
	case BandStop: // BSF BQ
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
	case AllPass: // AP BQ
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
