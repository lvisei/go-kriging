package ordinary

import "image/color"

type ModelType string

const (
	Gaussian    ModelType = "gaussian"
	Exponential ModelType = "exponential"
	Spherical   ModelType = "spherical"
)

var (
	DefaultLegendColor = []color.Color{
		RGBA{40, 146, 199, 255},
		RGBA{96, 163, 181, 255},
		RGBA{140, 184, 164, 255},
		RGBA{177, 204, 145, 255},
		RGBA{215, 227, 125, 255},
		RGBA{250, 250, 100, 255},
		RGBA{252, 207, 81, 255},
		RGBA{252, 164, 63, 255},
		RGBA{242, 77, 31, 255},
		RGBA{232, 16, 20, 255},
	}
	DefaultGridLevelColor = []GridLevelColor{
		{Color: RGBA{40, 146, 199, 255}, Value: [2]float64{-30, -15}},
		{Color: RGBA{96, 163, 181, 255}, Value: [2]float64{-15, -10}},
		{Color: RGBA{140, 184, 164, 255}, Value: [2]float64{-10, -5}},
		{Color: RGBA{177, 204, 145, 255}, Value: [2]float64{-5, 0}},
		{Color: RGBA{215, 227, 125, 255}, Value: [2]float64{0, 5}},
		{Color: RGBA{250, 250, 100, 255}, Value: [2]float64{5, 10}},
		{Color: RGBA{252, 207, 81, 255}, Value: [2]float64{10, 15}},
		{Color: RGBA{252, 164, 63, 255}, Value: [2]float64{15, 20}},
		{Color: RGBA{242, 77, 31, 255}, Value: [2]float64{25, 30}},
		{Color: RGBA{232, 16, 20, 255}, Value: [2]float64{30, 40}},
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
	Data  [][]float64 `json:"data"`
	Width float64     `json:"width"`
	Xlim  [2]float64  `json:"xLim"`
	Ylim  [2]float64  `json:"yLim"`
	Zlim  [2]float64  `json:"zLim"`
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

type Point [2]float64 // [103.614373, 27.00541]

type Ring []Point // [[103.614373, 27.00541],[104.174357, 26.635252],[104.356163, 28.018448],[103.614373, 27.00541]]

type PolygonCoordinates []Ring // [[[103.614373, 27.00541],[104.174357, 26.635252],[104.356163, 28.018448],[103.614373, 27.00541]]]

type PolygonGeometry struct {
	Type        string `json:"type" example:"Polygon"` // Polygon
	Coordinates []Ring `json:"coordinates,omitempty"`  // coordinates
}

type RGBA [4]uint8

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
	Value [2]float64 `json:"value" example:"0, 5"`               // [0, 5]
	Color RGBA       `json:"color" example:"255, 255, 255, 255"` // [255, 255, 255, 255]
}
