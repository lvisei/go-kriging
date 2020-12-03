package ordinary

func minFloat64(t []float64) float64 {
	min := float64(0)
	for i := 0; i < len(t); i++ {
		if min == 0 || min > t[i] {
			min = t[i]
		}
	}

	return min
}

func maxFloat64(t []float64) float64 {
	max := float64(0)
	for i := 0; i < len(t); i++ {
		if max < t[i] {
			max = t[i]
		}
	}

	return max
}

func pipFloat64(t [][2]float64, x, y float64) bool {
	c := false
	for i, j := 0, len(t)-1; i < len(t); j, i = i, i+1 {
		if ((t[i][1] > y) != (t[j][1] > y)) && (x < (t[j][0]-t[i][0])*(y-t[i][1])/(t[j][1]-t[i][1])+t[i][0]) {
			c = !c
		}
	}

	return c
}
