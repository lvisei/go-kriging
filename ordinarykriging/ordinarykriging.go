// Package ordinarykriging
// Ordinary Kriging（OK）
// 普通克里金

package ordinarykriging

import (
	"errors"
	"image"
	"image/color"
	"math"
	"sort"
	"sync"

	"github.com/lvisei/go-kriging/canvas"
)

// Variogram ordinary kriging variogram
type Variogram struct {
	t []float64
	x []float64
	y []float64

	Nugget float64 `json:"nugget"`
	Range  float64 `json:"range"`
	Sill   float64 `json:"sill"`
	A      float64 `json:"A"`
	N      int     `json:"n"`

	K     []float64 `json:"K"`
	M     []float64 `json:"M"`
	model variogramModel
}

func NewOrdinary(t, x, y []float64) *Variogram {
	return &Variogram{t: t, x: x, y: y}
}

type variogramModel func(float64, float64, float64, float64, float64) float64

// krigingVariogramGaussian gaussian variogram models
func krigingVariogramGaussian(h, nugget, range_, sill, A float64) float64 {
	x := -(1.0 / A) * ((h / range_) * (h / range_))
	return nugget + ((sill-nugget)/range_)*
		(1.0-exp(x))
}

// krigingVariogramExponential exponential variogram models
func krigingVariogramExponential(h, nugget, range_, sill, A float64) float64 {
	x := -(1.0 / A) * (h / range_)
	return nugget + ((sill-nugget)/range_)*
		(1.0-exp(x))
}

// krigingVariogramSpherical spherical variogram models
func krigingVariogramSpherical(h, nugget, range_, sill, A float64) float64 {
	if h > range_ {
		return nugget + (sill-nugget)/range_
	} else {
		x := h / range_
		return nugget + ((sill-nugget)/range_)*
			(1.5*(x)-0.5*(pow3(x)))
	}
}

// Train using gaussian processes with bayesian priors
func (variogram *Variogram) Train(model ModelType, sigma2 float64, alpha float64) (*Variogram, error) {
	variogram.Nugget = 0.0
	variogram.Range = 0.0
	variogram.Sill = 0.0
	variogram.A = float64(1) / float64(3)
	variogram.N = 0.0

	switch model {
	case Gaussian:
		variogram.model = krigingVariogramGaussian
		break
	case Exponential:
		variogram.model = krigingVariogramExponential
		break
	case Spherical:
		variogram.model = krigingVariogramSpherical
		break
	}

	// Lag distance/semivariance
	var i, j, k, l, n int
	n = len(variogram.t)

	var distance DistanceList
	distance = make([][2]float64, (n*n-n)/2)

	i = 0
	k = 0
	for ; i < n; i++ {
		for j = 0; j < i; {
			distance[k] = [2]float64{}
			distance[k][0] = math.Sqrt(pow2(variogram.x[i]-variogram.x[j]) + pow2(variogram.y[i]-variogram.y[j]))
			distance[k][1] = math.Abs(variogram.t[i] - variogram.t[j])
			j++
			k++
		}
	}
	sort.Sort(distance)
	variogram.Range = distance[(n*n-n)/2-1][0]

	// Bin lag distance
	var lags int
	if ((n*n - n) / 2) > 30 {
		lags = 30
	} else {
		lags = (n*n - n) / 2
	}

	tolerance := variogram.Range / float64(lags)

	lag := make([]float64, lags)
	semi := make([]float64, lags)
	if lags < 30 {
		for l = 0; l < lags; l++ {
			lag[l] = distance[l][0]
			semi[l] = distance[l][1]
		}
	} else {
		i = 0
		j = 0
		k = 0
		l = 0
		for i < lags && j < ((n*n-n)/2) {
			for {
				if distance[j][0] > (float64(i+1) * tolerance) {
					break
				}
				lag[l] += distance[j][0]
				semi[l] += distance[j][1]
				j++
				k++
				if j >= ((n*n - n) / 2) {
					break
				}
			}

			if k > 0 {
				lag[l] = lag[l] / float64(k)
				semi[l] = semi[l] / float64(k)
				l++
			}
			i++
			k = 0
		}
		if l < 2 {
			return nil, errors.New("not enough points")
		}
	}

	// Feature transformation
	n = l
	variogram.Range = lag[n-1] - lag[0]
	X := make([]float64, 2*n)
	for i := 0; i < len(X); i++ {
		X[i] = 1
	}
	Y := make([]float64, n)
	var A = variogram.A
	for i = 0; i < n; i++ {
		switch model {
		case Gaussian:
			X[i*2+1] = 1.0 - exp(-(1.0/A)*pow2(lag[i]/variogram.Range))
			break
		case Exponential:
			X[i*2+1] = 1.0 - exp(-(1.0/A)*lag[i]/variogram.Range)
			break
		case Spherical:
			X[i*2+1] = 1.5*(lag[i]/variogram.Range) - 0.5*pow3(lag[i]/variogram.Range)
			break
		}
		Y[i] = semi[i]
	}

	// Least squares
	var Xt = matrixTranspose(X, n, 2)
	var Z = matrixMultiply(Xt, X, 2, n, 2)
	Z = matrixAdd(Z, matrixDiag(float64(1)/alpha, 2), 2, 2)
	var cloneZ = make([]float64, len(Z))
	copy(cloneZ, Z)
	if matrixChol(Z, 2) {
		matrixChol2inv(Z, 2)
	} else {
		// TODO false
		Z, _ = matrixInverse(cloneZ, 2)
	}

	var W = matrixMultiply(matrixMultiply(Z, Xt, 2, 2, n), Y, 2, n, 1)

	// Variogram parameters
	variogram.Nugget = W[0]
	variogram.Sill = W[1]*variogram.Range + variogram.Nugget
	variogram.N = len(variogram.x)

	// Gram matrix with prior
	n = len(variogram.x)
	K := make([]float64, n*n)
	for i = 0; i < n; i++ {
		for j = 0; j < i; j++ {
			K[i*n+j] = variogram.model(
				math.Sqrt(pow2(variogram.x[i]-variogram.x[j])+pow2(variogram.y[i]-variogram.y[j])),
				variogram.Nugget,
				variogram.Range,
				variogram.Sill,
				variogram.A)
			K[j*n+i] = K[i*n+j]
		}
		K[i*n+i] = variogram.model(0, variogram.Nugget,
			variogram.Range,
			variogram.Sill,
			variogram.A)
	}

	// Inverse penalized Gram matrix projected to target vector
	var C = matrixAdd(K, matrixDiag(sigma2, n), n, n)
	var cloneC = make([]float64, len(C))
	copy(cloneC, C)
	if matrixChol(C, n) {
		matrixChol2inv(C, n)
	} else {
		// TODO false
		C, _ = matrixInverse(cloneC, n)
	}

	// Copy unprojected inverted matrix as K 复制未投影的逆矩阵为K
	copy(K, C)
	var M = matrixMultiply(C, variogram.t, n, n, 1)
	variogram.K = K
	variogram.M = M

	return variogram, nil
}

// Predict model prediction
func (variogram *Variogram) Predict(x, y float64) float64 {
	k := make([]float64, variogram.N)
	for i := 0; i < variogram.N; i++ {
		x_ := x - variogram.x[i]
		y_ := y - variogram.y[i]
		h := math.Sqrt(pow2(x_) + pow2(y_))
		k[i] = variogram.model(
			h,
			variogram.Nugget, variogram.Range,
			variogram.Sill, variogram.A,
		)
	}

	return matrixMultiply(k, variogram.M, 1, variogram.N, 1)[0]
}

func (variogram *Variogram) Variance(x, y float64) {

}

// Grid gridded matrices or contour paths
// 根据 PolygonCoordinates 生成裁剪过的矩阵网格数据
// 这里 polygon 是一个三维数组，可以变相的支持的多个面，但不符合 Polygon 规范
// PolygonCoordinates [[[x,y]],[[x,y]]] 两个面
func (variogram *Variogram) Grid(polygon PolygonCoordinates, width float64) *GridMatrices {
	n := len(polygon)
	if n == 0 {
		return &GridMatrices{}
	}

	var nodataValue float64 = -9999

	// Boundaries of polygon space
	xlim := [2]float64{polygon[0][0][0], polygon[0][0][0]}
	ylim := [2]float64{polygon[0][0][1], polygon[0][0][1]}

	// Polygons
	for i := 0; i < n; i++ {
		// Vertices
		for j := 0; j < len(polygon[i]); j++ {
			if polygon[i][j][0] < xlim[0] {
				xlim[0] = polygon[i][j][0]
			}
			if polygon[i][j][0] > xlim[1] {
				xlim[1] = polygon[i][j][0]
			}
			if polygon[i][j][1] < ylim[0] {
				ylim[0] = polygon[i][j][1]
			}
			if polygon[i][j][1] > ylim[1] {
				ylim[1] = polygon[i][j][1]
			}
		}
	}

	// Alloc for O(N^2) space
	x := int(math.Ceil((xlim[1] - xlim[0]) / width))
	y := int(math.Ceil((ylim[1] - ylim[0]) / width))

	A := make([][]float64, x+1)
	for i := 0; i <= x; i++ {
		A[i] = make([]float64, y+1)
	}

	for i := 0; i < n; i++ {
		currentPolygon := polygon[i]
		var lxlim [2]float64 // Local dimensions
		var lylim [2]float64 // Local dimensions
		// Range for currentPolygon
		lxlim[0] = currentPolygon[0][0]
		lxlim[1] = lxlim[0]
		lylim[0] = currentPolygon[0][1]
		lylim[1] = lylim[0]
		for j := 1; j < len(currentPolygon); j++ { // Vertices
			if currentPolygon[j][0] < lxlim[0] {
				lxlim[0] = currentPolygon[j][0]
			}
			if currentPolygon[j][0] > lxlim[1] {
				lxlim[1] = currentPolygon[j][0]
			}
			if currentPolygon[j][1] < lylim[0] {
				lylim[0] = currentPolygon[j][1]
			}
			if currentPolygon[j][1] > lylim[1] {
				lylim[1] = currentPolygon[j][1]
			}
		}

		var a, b [2]int
		// Loop through polygon subspace
		a[0] = int(math.Floor(((lxlim[0] - math.Mod(lxlim[0]-xlim[0], width)) - xlim[0]) / width))
		a[1] = int(math.Ceil(((lxlim[1] - math.Mod(lxlim[1]-xlim[1], width)) - xlim[0]) / width))
		b[0] = int(math.Floor(((lylim[0] - math.Mod(lylim[0]-ylim[0], width)) - ylim[0]) / width))
		b[1] = int(math.Ceil(((lylim[1] - math.Mod(lylim[1]-ylim[1], width)) - ylim[0]) / width))

		var wg sync.WaitGroup
		predictCh := make(chan *PredictDate, (b[1]-b[0])*(a[1]-a[0]))
		var parallelPredict = func(j, k int, polygon []Point, xTarget, yTarget float64) {
			predictDate := &PredictDate{X: j, Y: k}
			predictDate.Value = variogram.Predict(xTarget,
				yTarget,
			)
			predictCh <- predictDate
			defer wg.Done()
		}

		var xTarget, yTarget float64
		for j := a[0]; j <= a[1]; j++ {
			xTarget = xlim[0] + float64(j)*width
			for k := b[0]; k <= b[1]; k++ {
				yTarget = ylim[0] + float64(k)*width

				if pipFloat64(currentPolygon, xTarget, yTarget) {
					wg.Add(1)
					go parallelPredict(j, k, currentPolygon, xTarget, yTarget)
				} else {
					A[j][k] = nodataValue
				}
				//if pipFloat64(currentPolygon, xTarget, yTarget) {
				//	A[j][k] = variogram.Predict(xTarget,
				//		yTarget,
				//	)
				//}
			}
		}

		go func() {
			wg.Wait()
			close(predictCh)
		}()

		for predictDate := range predictCh {
			if predictDate.Value != 0 {
				j := predictDate.X
				k := predictDate.Y
				A[j][k] = predictDate.Value
			}

		}
	}

	gridMatrices := &GridMatrices{
		Xlim:        xlim,
		Ylim:        ylim,
		Zlim:        [2]float64{minFloat64(variogram.t), maxFloat64(variogram.t)},
		Width:       width,
		Data:        A,
		NodataValue: nodataValue,
	}
	return gridMatrices
}

// Contour contour paths
// 根据宽高度生成轮廓数据
func (variogram *Variogram) Contour(xWidth, yWidth int) *ContourRectangle {
	xlim := [2]float64{minFloat64(variogram.x), maxFloat64(variogram.x)}
	ylim := [2]float64{minFloat64(variogram.y), maxFloat64(variogram.y)}
	zlim := [2]float64{minFloat64(variogram.t), maxFloat64(variogram.t)}
	xl := xlim[1] - xlim[0]
	yl := ylim[1] - ylim[0]
	gridW := xl / float64(xWidth)
	gridH := yl / float64(yWidth)
	var contour []float64

	var xTarget, yTarget float64

	for j := 0; j < yWidth; j++ {
		yTarget = ylim[0] + float64(j)*gridW
		for k := 0; k < xWidth; k++ {
			xTarget = xlim[0] + float64(k)*gridH
			contour = append(contour, variogram.Predict(xTarget, yTarget))
		}
	}

	contourRectangle := &ContourRectangle{
		Contour:     contour,
		XWidth:      xWidth,
		YWidth:      yWidth,
		Xlim:        xlim,
		Ylim:        ylim,
		Zlim:        zlim,
		XResolution: 1,
		YResolution: 1,
	}

	return contourRectangle
}

// ContourWithBBox contour paths
// 根据 bbox 生成轮廓数据
func (variogram *Variogram) ContourWithBBox(bbox [4]float64, width float64) *ContourRectangle {
	// x方向
	xlim := [2]float64{bbox[0], bbox[2]}
	ylim := [2]float64{bbox[1], bbox[3]}
	zlim := [2]float64{minFloat64(variogram.t), maxFloat64(variogram.t)}

	// xy 方向地理跨度
	geoXWidth := xlim[1] - xlim[0]
	geoYWidth := ylim[1] - ylim[0]

	xWidth := int(math.Ceil(width))
	// 让图像的xy比例与地理的xy比例保持一致
	yWidth := int(math.Ceil(float64(xWidth) / (geoXWidth / geoYWidth)))

	// 地理跨度/图像跨度=当前地图图上分辨率
	var xResolution = geoXWidth * 1.0 / float64(xWidth)
	var yResolution = geoYWidth * 1.0 / float64(yWidth)

	var xTarget, yTarget float64
	var contour []float64

	for j := 0; j < yWidth; j++ {
		yTarget = bbox[1] + float64(j)*yResolution
		for k := 0; k < xWidth; k++ {
			xTarget = bbox[0] + float64(k)*xResolution
			contour = append(contour, variogram.Predict(xTarget, yTarget))
		}
	}
	contourRectangle := &ContourRectangle{
		Contour:     contour,
		XWidth:      xWidth,
		YWidth:      yWidth,
		Xlim:        xlim,
		Ylim:        ylim,
		Zlim:        zlim,
		XResolution: xResolution,
		YResolution: yResolution,
	}

	return contourRectangle
}

// Plot plotting on the canvas
// 绘制裁剪过的矩阵网格数据到 canvas 上
func (variogram *Variogram) Plot(gridMatrices *GridMatrices, width, height int, xlim, ylim [2]float64, colors []GridLevelColor) *canvas.Canvas {
	// Create canvas
	ctx := canvas.NewCanvas(width, height)
	// Starting boundaries
	range_ := [...]float64{xlim[1] - xlim[0], ylim[1] - ylim[0], gridMatrices.Zlim[1] - gridMatrices.Zlim[0]}

	n := len(gridMatrices.Data)
	m := len(gridMatrices.Data[0])
	// 计算色块宽高度
	wx := math.Ceil(gridMatrices.Width * float64(width) / (xlim[1] - xlim[0]))
	wy := math.Ceil(gridMatrices.Width * float64(height) / (ylim[1] - ylim[0]))

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if gridMatrices.Data[i][j] == gridMatrices.NodataValue {
				continue
			} else {
				x := float64(width) * (float64(i)*gridMatrices.Width + gridMatrices.Xlim[0] - xlim[0]) / range_[0]
				y := float64(height) * (1 - (float64(j)*gridMatrices.Width+gridMatrices.Ylim[0]-ylim[0])/range_[1])
				distance := gridMatrices.Data[i][j] - gridMatrices.Zlim[0]
				z := (distance) / range_[2]
				if z < 0 {
					z = 0.0
				} else if z > 1 {
					z = 1.0
				}

				colorIndex := -1
				for index, item := range colors {
					if distance >= item.Value[0] && distance <= item.Value[1] {
						colorIndex = index
						break
					}
				}
				if colorIndex == -1 {
					continue
				}
				color := colors[colorIndex].Color
				ctx.DrawRect(math.Round(x-wx/2), math.Round(y-wy/2), wx, wy, color)
			}
		}
	}

	return ctx
}

// PlotRectangleGrid
// 绘制矩形网格到数据 canvas 上
func (variogram *Variogram) PlotRectangleGrid(contourRectangle *ContourRectangle, width, height int, xlim, ylim [2]float64, colors []color.Color) *canvas.Canvas {
	// Create canvas
	ctx := canvas.NewCanvas(width, height)
	// Starting boundaries
	range_ := [...]float64{xlim[1] - xlim[0], ylim[1] - ylim[0], contourRectangle.Zlim[1] - contourRectangle.Zlim[0]}
	n := contourRectangle.XWidth
	m := contourRectangle.YWidth
	// 计算色块宽高度
	wx := math.Ceil(contourRectangle.XResolution * float64(width) / (xlim[1] - xlim[0]))
	wy := math.Ceil(contourRectangle.YResolution * float64(height) / (ylim[1] - ylim[0]))

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			index := i*n + j
			x := (float64(width) * (float64(j)*contourRectangle.XResolution + contourRectangle.Xlim[0] - xlim[0])) / range_[0]
			y := float64(height) * (1 - (float64(i)*contourRectangle.YResolution+contourRectangle.Ylim[0]-ylim[0])/range_[1])
			z := (contourRectangle.Contour[index] - contourRectangle.Zlim[0]) / range_[2]
			if z < 0 {
				z = 0.0
			} else if z > 1 {
				z = 1.0
			}

			colorIndex := int(math.Floor((float64(len(colors)) - 1) * z))
			color := colors[colorIndex]
			ctx.DrawRect(math.Round(x-wx/2), math.Round(y-wy/2), wx, wy, color)
		}
	}

	return ctx
}

// PlotPng
func (variogram *Variogram) PlotPng(rectangleGrids *ContourRectangle) *image.RGBA {
	contour := rectangleGrids.Contour
	xWidth := rectangleGrids.XWidth
	yWidth := rectangleGrids.YWidth
	zlim := rectangleGrids.Zlim
	colorperiod := (zlim[1] - zlim[0]) / 5
	img := image.NewRGBA(image.Rect(0, 0, xWidth, yWidth))

	for i := 0; i < xWidth*yWidth; i++ {
		zi := contour[i]
		var color color.RGBA

		if zi <= zlim[1] && zi > zlim[1]-colorperiod {
			color.R = 189
			color.G = 0
			color.B = 36
			color.A = 128
		} else if zi <= zlim[1]-colorperiod && zi > zlim[1]-2*colorperiod {
			color.R = 240
			color.G = 59
			color.B = 32
			color.A = 128
		} else if zi <= zlim[1]-2*colorperiod && zi > zlim[1]-3*colorperiod {
			color.R = 253
			color.G = 141
			color.B = 60
			color.A = 128
		} else if zi <= zlim[1]-3*colorperiod && zi > zlim[1]-4*colorperiod {
			color.R = 254
			color.G = 204
			color.B = 92
			color.A = 128
		} else {
			color.R = 255
			color.G = 255
			color.B = 78
			color.A = 128
		}

		x := i % xWidth
		y := i / yWidth
		img.Set(x, y, color)
	}

	return img
}
