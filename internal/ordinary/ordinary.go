// Package ordinary
// Ordinary Kriging（OK）
// 普通克里金

package ordinary

import (
	"image"
	"image/color"
	"math"
	"sort"
)

// Variogram ordinary kriging variogram
type Variogram struct {
	T []float64
	X []float64
	Y []float64

	nugget float64
	range_ float64
	sill   float64
	A      float64
	n      int

	K     []float64
	M     []float64
	model variogramModel
}

type variogramModel func(float64, float64, float64, float64, float64) float64

//type VariogramModel interface {
//	model(float64, float64, float64, float64, float64) float64
//}

// krigingVariogramGaussian gaussian variogram models
func krigingVariogramGaussian(h, nugget, range_, sill, A float64) float64 {
	return nugget + ((sill-nugget)/range_)*
		(1.0-math.Exp(-(1.0/A)*math.Pow(h/range_, 2)))
}

// krigingVariogramExponential exponential variogram models
func krigingVariogramExponential(h, nugget, range_, sill, A float64) float64 {
	return nugget + ((sill-nugget)/range_)*
		(1.0-math.Exp(-(1.0/A)*(h/range_)))
}

// krigingVariogramSpherical spherical variogram models
func krigingVariogramSpherical(h, nugget, range_, sill, A float64) float64 {
	if h > range_ {
		return nugget + (sill-nugget)/range_
	} else {
		return nugget + ((sill-nugget)/range_)*
			(1.5*(h/range_)-0.5*math.Pow(h/range_, 3))
	}
}

// Train using gaussian processes with bayesian priors
func (variogram *Variogram) Train(model ModelType, sigma2 float64, alpha float64) *Variogram {
	variogram.nugget = 0.0
	variogram.range_ = 0.0
	variogram.sill = 0.0
	variogram.A = float64(1) / float64(3)
	variogram.n = 0.0

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
	n = len(variogram.T)

	var distance DistanceList
	distance = make([][2]float64, (n*n-n)/2)

	i = 0
	k = 0
	for ; i < n; i++ {
		for j = 0; j < i; {
			distance[k] = [2]float64{}
			distance[k][0] = math.Pow(
				math.Pow(variogram.X[i]-variogram.X[j], 2)+
					math.Pow(variogram.Y[i]-variogram.Y[j], 2), 0.5)
			distance[k][1] = math.Abs(variogram.T[i] - variogram.T[j])
			j++
			k++
		}
	}
	sort.Sort(distance)
	variogram.range_ = distance[(n*n-n)/2-1][0]

	// Bin lag distance
	var lags int
	if ((n*n - n) / 2) > 30 {
		lags = 30
	} else {
		lags = (n*n - n) / 2
	}

	tolerance := variogram.range_ / float64(lags)

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
			// Error: Not enough points
			return variogram
		}
	}

	// Feature transformation
	n = l
	variogram.range_ = lag[n-1] - lag[0]
	X := make([]float64, 2*n)
	for i := 0; i < len(X); i++ {
		X[i] = 1
	}
	Y := make([]float64, n)
	var A = variogram.A
	for i = 0; i < n; i++ {
		switch model {
		case Gaussian:
			X[i*2+1] = 1.0 - math.Exp(-(1.0/A)*math.Pow(lag[i]/variogram.range_, 2))
			break
		case Exponential:
			X[i*2+1] = 1.0 - math.Exp(-(1.0/A)*lag[i]/variogram.range_)
			break
		case Spherical:
			X[i*2+1] = 1.5*(lag[i]/variogram.range_) - 0.5*math.Pow(lag[i]/variogram.range_, 3)
			break
		}
		Y[i] = semi[i]
	}

	// Least squares
	var Xt = krigingMatrixTranspose(X, n, 2)
	var Z = krigingMatrixMultiply(Xt, X, 2, n, 2)
	Z = krigingMatrixAdd(Z, krigingMatrixDiag(float64(1)/alpha, 2), 2, 2)
	var cloneZ = Z
	if krigingMatrixChol(Z, 2) {
		krigingMatrixChol2inv(Z, 2)
	} else {
		krigingMatrixSolve(cloneZ, 2)
		Z = cloneZ
	}
	var W = krigingMatrixMultiply(krigingMatrixMultiply(Z, Xt, 2, 2, n), Y, 2, n, 1)

	// Variogram parameters
	variogram.nugget = W[0]
	variogram.sill = W[1]*variogram.range_ + variogram.nugget
	variogram.n = len(variogram.X)

	// Gram matrix with prior
	n = len(variogram.X)
	K := make([]float64, n*n)
	for i = 0; i < n; i++ {
		for j = 0; j < i; j++ {
			K[i*n+j] = variogram.model(math.Pow(math.Pow(variogram.X[i]-variogram.X[j], 2)+
				math.Pow(variogram.Y[i]-variogram.Y[j], 2), 0.5),
				variogram.nugget,
				variogram.range_,
				variogram.sill,
				variogram.A)
			K[j*n+i] = K[i*n+j]
		}
		K[i*n+i] = variogram.model(0, variogram.nugget,
			variogram.range_,
			variogram.sill,
			variogram.A)
	}

	// Inverse penalized Gram matrix projected to target vector
	var C = krigingMatrixAdd(K, krigingMatrixDiag(sigma2, n), n, n)
	var cloneC = C
	if krigingMatrixChol(C, n) {
		krigingMatrixChol2inv(C, n)
	} else {
		krigingMatrixSolve(cloneC, n)
		C = cloneC
	}

	// Copy unprojected inverted matrix as K
	K = C
	M := krigingMatrixMultiply(C, variogram.T, n, n, 1)
	variogram.K = K
	variogram.M = M

	return variogram
}

// Predict model prediction
func (variogram *Variogram) Predict(x, y float64) float64 {
	k := make([]float64, variogram.n)
	for i := 0; i < variogram.n; i++ {
		k[i] = variogram.model(math.Pow(math.Pow(x-variogram.X[i], 2)+
			math.Pow(y-variogram.Y[i], 2), 0.5),
			variogram.nugget, variogram.range_,
			variogram.sill, variogram.A)
	}

	return krigingMatrixMultiply(k, variogram.M, 1, variogram.n, 1)[0]
}

func (variogram *Variogram) Variance(x, y float64) {

}

// Grid gridded matrices or contour paths
func (variogram *Variogram) Grid(polygons [][][2]float64, width float64) GridMatrices {
	n := len(polygons)
	if n == 0 {
		return GridMatrices{}
	}

	// Boundaries of polygons space
	xlim := [2]float64{polygons[0][0][0], polygons[0][0][0]}
	ylim := [2]float64{polygons[0][0][1], polygons[0][0][1]}

	// Polygons
	for i := 0; i < n; i++ {
		// Vertices
		for j := 0; j < len(polygons[i]); j++ {
			if polygons[i][j][0] < xlim[0] {
				xlim[0] = polygons[i][j][0]
			}
			if polygons[i][j][0] > xlim[1] {
				xlim[1] = polygons[i][j][0]
			}
			if polygons[i][j][1] < ylim[0] {
				ylim[0] = polygons[i][j][1]
			}
			if polygons[i][j][1] > ylim[1] {
				ylim[1] = polygons[i][j][1]
			}
		}
	}

	// Alloc for O(n^2) space
	var xtarget, ytarget float64
	var a, b [2]int
	var lxlim [2]float64 // Local dimensions
	var lylim [2]float64 // Local dimensions
	var x = int(math.Ceil((xlim[1] - xlim[0]) / width))
	var y = int(math.Ceil((ylim[1] - ylim[0]) / width))

	A := make([][]float64, x+1)
	for i := 0; i <= x; i++ {
		A[i] = make([]float64, y+1)
	}
	for i := 0; i < n; i++ {
		// Range for polygons[i]
		lxlim[0] = polygons[i][0][0]
		lxlim[1] = lxlim[0]
		lylim[0] = polygons[i][0][1]
		lylim[1] = lylim[0]
		for j := 1; j < len(polygons[i]); j++ { // Vertices
			if polygons[i][j][0] < lxlim[0] {
				lxlim[0] = polygons[i][j][0]
			}
			if polygons[i][j][0] > lxlim[1] {
				lxlim[1] = polygons[i][j][0]
			}
			if polygons[i][j][1] < lylim[0] {
				lylim[0] = polygons[i][j][1]
			}
			if polygons[i][j][1] > lylim[1] {
				lylim[1] = polygons[i][j][1]
			}
		}

		// Loop through polygon subspace
		a[0] = int(math.Floor(((lxlim[0] - math.Mod(lxlim[0]-xlim[0], width)) - xlim[0]) / width))
		a[1] = int(math.Ceil(((lxlim[1] - math.Mod(lxlim[1]-xlim[1], width)) - xlim[0]) / width))
		b[0] = int(math.Floor(((lylim[0] - math.Mod(lylim[0]-ylim[0], width)) - ylim[0]) / width))
		b[1] = int(math.Ceil(((lylim[1] - math.Mod(lylim[1]-ylim[1], width)) - ylim[0]) / width))
		for j := a[0]; j <= a[1]; j++ {
			for k := b[0]; k <= b[1]; k++ {
				xtarget = xlim[0] + float64(j)*width
				ytarget = ylim[0] + float64(k)*width

				if pipFloat64(polygons[i], xtarget, ytarget) {
					A[j][k] = variogram.Predict(xtarget,
						ytarget,
					)
				}
			}
		}
	}

	gridMatrices := GridMatrices{
		xlim:  xlim,
		ylim:  ylim,
		zlim:  [2]float64{minFloat64(variogram.T), maxFloat64(variogram.T)},
		width: width, data: A,
	}
	return gridMatrices
}

// getGridInfo gridded matrices or contour paths
func (variogram *Variogram) RectangleGrid(bbox [4]float64, width float64) map[string]interface{} {
	var grid []float64

	// x方向
	xlim := [2]float64{bbox[0], bbox[2]}
	ylim := [2]float64{bbox[1], bbox[3]}
	zlim := [2]float64{minFloat64(variogram.T), maxFloat64(variogram.T)}

	// xy 方向地理跨度
	geoxWidth := xlim[1] - xlim[0]
	geoyWidth := ylim[1] - ylim[0]

	// 如果x_width设置，初始基于200计算
	var xWidth, yWidth int
	if width != 0 {
		xWidth = 200
	} else {
		xWidth = int(math.Ceil(width))
	}
	//让图像的xy比例与地理的xy比例保持一致
	yWidth = int(math.Ceil(float64(xWidth) / (geoxWidth / geoyWidth)))

	//地理跨度/图像跨度=当前地图图上分辨率
	var xResolution = geoxWidth * 1.0 / float64(xWidth)
	var yResolution = geoyWidth * 1.0 / float64(yWidth)

	var xTarget, yTarget float64

	for j := 0; j < yWidth; j++ {
		for k := 0; k < xWidth; k++ {
			xTarget = bbox[0] + float64(k)*xResolution
			yTarget = bbox[1] + float64(j)*yResolution
			grid = append(grid, variogram.Predict(xTarget, yTarget))
		}
	}
	gridDate := map[string]interface{}{
		"grid":        grid,
		"n":           xWidth,
		"m":           yWidth,
		"xlim":        xlim,
		"ylim":        ylim,
		"zlim":        zlim,
		"xResolution": xResolution,
		"yResolution": yResolution,
	}

	return gridDate

}

func (variogram *Variogram) Contour() {

}

// Plot plotting on the Canvas
func (variogram *Variogram) Plot() {

}

// GeneratePngGrid
func (variogram *Variogram) GeneratePngGrid(xWidth, yWidth int) ([]float64, float64, float64) {
	rangeMaxX := maxFloat64(variogram.X)
	rangeMinX := minFloat64(variogram.X)
	rangeMaxY := maxFloat64(variogram.Y)
	rangeMinY := minFloat64(variogram.Y)
	rangeMaxT := maxFloat64(variogram.T)
	rangeMinT := minFloat64(variogram.T)
	colorperiod := (rangeMaxT - rangeMinT) / 5
	var xl = rangeMaxX - rangeMinX
	var yl = rangeMaxY - rangeMinY
	var gridX = xl / float64(xWidth)
	var gridY = yl / float64(yWidth)
	var gridPoint [][2]float64
	var krigingValue []float64
	var gX = rangeMinX

	for i := 0; i < xWidth; i++ {
		gX = gX + gridX
		gY := rangeMinY
		for j := 0; j < yWidth; j++ {
			gY = gY + gridY
			var pp = [2]float64{gX, gY}
			krigingValue = append(krigingValue, variogram.Predict(gX, gY))
			gridPoint = append(gridPoint, pp)
		}
	}

	return krigingValue, rangeMaxT, colorperiod
}

// GeneratePng
func (variogram *Variogram) GeneratePng(krigingValue []float64, rangeMaxPM float64, colorperiod float64, xWidth, yWidth int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, xWidth, yWidth))

	for i := 0; i < xWidth*yWidth; i++ {
		zi := krigingValue[i]
		var color color.RGBA

		if zi <= rangeMaxPM && zi > rangeMaxPM-colorperiod {
			color.R = 189
			color.G = 0
			color.B = 36
			color.A = 128
		} else if zi <= rangeMaxPM-colorperiod && zi > rangeMaxPM-2*colorperiod {
			color.R = 240
			color.G = 59
			color.B = 32
			color.A = 128
		} else if zi <= rangeMaxPM-2*colorperiod && zi > rangeMaxPM-3*colorperiod {
			color.R = 253
			color.G = 141
			color.B = 60
			color.A = 128
		} else if zi <= rangeMaxPM-3*colorperiod && zi > rangeMaxPM-4*colorperiod {
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
