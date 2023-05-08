package kurumi

import "math"

var Interpolations = []string{
	"None",
	"Linear",
	"Cosine",
	"Cubic",
}

func noInterpolation(x float64, wt []uint8) float64 {
	t := x
	idx := math.Floor(t)
	len := len(wt)
	_, max := minandmax(wt)
	myMax := float64(max) / 2.0
	s0 := float64(wt[moduloInt(int(idx), len)])/myMax - 1
	return s0
}

func linearInterpolation(x float64, wt []uint8) float64 {
	t := x
	idx := int(math.Floor(t))
	len := len(wt)
	_, m := minandmax(wt)
	max := float64(m) / 2.0
	mu := (t - float64(idx))
	s0 := float64(wt[moduloInt(idx, len)])/max - 1.0
	s1 := float64(wt[moduloInt(idx+1, len)])/max - 1.0
	return s0 + mu*s1 - (mu * s0)
}

func cosineInterpolation(x float64, wt []uint8) float64 {
	t := x
	idx := int(math.Floor(t))
	len := len(wt)
	_, m := minandmax(wt)
	max := float64(m) / 2.0
	mu := (t - float64(idx))
	muCos := (1 - math.Cos(mu*math.Pi)/2)
	s0 := float64(wt[moduloInt(idx, len)])/max - 1.0
	s1 := float64(wt[moduloInt(idx+1, len)])/max - 1.0
	return s0 + muCos*s1 - (muCos * s0)
}

func cubicInterpolation(x float64, wt []uint8) float64 {
	t := x
	idx := int(math.Floor(t))
	len := len(wt)
	_, m := minandmax(wt)
	max := float64(m) / 2.0
	s0 := float64(wt[moduloInt(idx-1, len)])/max - 1.0
	s1 := float64(wt[moduloInt(idx, len)])/max - 1.0
	s2 := float64(wt[moduloInt(idx+1, len)])/max - 1.0
	s3 := float64(wt[moduloInt(idx+2, len)])/max - 1.0
	mu := (t - float64(idx))
	mu2 := mu * mu
	a0 := -0.5*s0 + 1.5*s1 - 1.5*s2 + 0.5*s3
	a1 := s0 - 2.5*s1 + 2*s2 - 0.5*s3
	a2 := -0.5*s0 + 0.5*s2
	a3 := s1
	return (a0*mu*mu2 + a1*mu2 + a2*mu + a3)
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
