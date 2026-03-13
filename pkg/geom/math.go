package geom

func SmoothStep(edge0, edge1, x float64) float64 {
	t := Clamp((x-edge0)/(edge1-edge0), 0.0, 1.0)

	return t * t * (3 - 2*t)
}

func Clamp(x, low, high float64) float64 {
	if x < low {
		return low
	}

	if x > high {
		return high
	}

	return x
}
