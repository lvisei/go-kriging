package ordinary

type ModelType int

const (
	Gaussian    ModelType = iota // value --> 0
	Exponential                  // value --> 1
	Spherical                    // value --> 2
)

type DistanceList [][2]float64

func (t DistanceList) Len() int {
	return len(t)
}

func (t DistanceList) Less(i, j int) bool {
	return t[i][0] < t[j][0]
}

func (t DistanceList) Swap(i, j int) {
	tmp := t[i]
	t[i] = t[j]
	t[j] = tmp
}

type GridMatrices struct {
	data  [][]float64
	width float64
	xlim  [2]float64
	ylim  [2]float64
	zlim  [2]float64
}
