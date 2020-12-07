package ordinary

type ModelType string

const (
	Gaussian    ModelType = "gaussian"
	Exponential ModelType = "exponential"
	Spherical   ModelType = "spherical"
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
	Data  [][]float64 `json:"data"`
	Width float64     `json:"width"`
	Xlim  [2]float64  `json:"xLim"`
	Ylim  [2]float64  `json:"yLim"`
	Zlim  [2]float64  `json:"zLim"`
}

type PolygonCoordinates [][][2]float64

type PolygonGeometry struct {
	Type        string             `json:"type"`
	Coordinates PolygonCoordinates `json:"coordinates,omitempty"`
}
