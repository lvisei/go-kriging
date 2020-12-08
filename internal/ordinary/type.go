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

type RGBA [4]uint32

func (c RGBA) RGBA() (r, g, b, a uint32) {
	r = uint32(c[0])
	r |= r << 8
	g = uint32(c[1])
	g |= g << 8
	b = uint32(c[2])
	b |= b << 8
	a = uint32(c[3])
	a |= a << 8
	return
}

type GridLevelColor struct {
	Value [2]float64 `json:"value"`

	Color RGBA `json:"color"`
}
