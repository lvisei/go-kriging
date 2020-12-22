package ordinarykriging

import "image/color"

type ModelType string

const (
	Gaussian    ModelType = "gaussian"
	Exponential ModelType = "exponential"
	Spherical   ModelType = "spherical"
)

var (
	DefaultLegendColor = []color.Color{
		NewRGBA(40, 146, 199, 255),
		NewRGBA(96, 163, 181, 255),
		NewRGBA(140, 184, 164, 255),
		NewRGBA(177, 204, 145, 255),
		NewRGBA(215, 227, 125, 255),
		NewRGBA(250, 250, 100, 255),
		NewRGBA(252, 207, 81, 255),
		NewRGBA(252, 164, 63, 255),
		NewRGBA(242, 77, 31, 255),
		NewRGBA(232, 16, 20, 255),
	}
	DefaultGridLevelColor = []GridLevelColor{
		{Color: NewRGBA(40, 146, 199, 255), Value: [2]float64{-30, -15}},
		{Color: NewRGBA(96, 163, 181, 255), Value: [2]float64{-15, -10}},
		{Color: NewRGBA(140, 184, 164, 255), Value: [2]float64{-10, -5}},
		{Color: NewRGBA(177, 204, 145, 255), Value: [2]float64{-5, 0}},
		{Color: NewRGBA(215, 227, 125, 255), Value: [2]float64{0, 5}},
		{Color: NewRGBA(250, 250, 100, 255), Value: [2]float64{5, 10}},
		{Color: NewRGBA(252, 207, 81, 255), Value: [2]float64{10, 15}},
		{Color: NewRGBA(252, 164, 63, 255), Value: [2]float64{15, 20}},
		{Color: NewRGBA(247, 122, 45, 255), Value: [2]float64{20, 25}},
		{Color: NewRGBA(242, 77, 31, 255), Value: [2]float64{25, 30}},
		{Color: NewRGBA(232, 16, 20, 255), Value: [2]float64{30, 40}},
	}
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
	Data        [][]float64 `json:"data"`
	Width       float64     `json:"width"`
	Xlim        [2]float64  `json:"xLim"`
	Ylim        [2]float64  `json:"yLim"`
	Zlim        [2]float64  `json:"zLim"`
	NodataValue float64     `json:"nodataValue"`
}

type ContourRectangle struct {
	Contour     []float64  `json:"contour"`
	XWidth      int        `json:"xWidth"`
	YWidth      int        `json:"yWidth"`
	Xlim        [2]float64 `json:"xLim"`
	Ylim        [2]float64 `json:"yLim"`
	Zlim        [2]float64 `json:"zLim"`
	XResolution float64    `json:"xResolution"`
	YResolution float64    `json:"yResolution"`
}

type Point [2]float64 // example [103.614373, 27.00541]

type Ring []Point

type PolygonCoordinates []Ring

type PolygonGeometry struct {
	Type        string `json:"type" default:"Polygon"` // Polygon
	Coordinates []Ring `json:"coordinates,omitempty"`  // coordinates
}

func NewRGBA(r, g, b, a uint8) color.RGBA {
	_rgba := color.RGBA{R: r, G: g, B: b, A: a}
	return _rgba
}

type GridLevelColor struct {
	Value [2]float64 `json:"value"` // 值区间 [0, 5]
	Color color.RGBA `json:"color"` // RGBA颜色 {255, 255, 255, 255}
}

type PredictDate struct {
	X     int
	Y     int
	Value float64
}
