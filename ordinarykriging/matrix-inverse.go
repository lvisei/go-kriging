package ordinarykriging

import (
	"gonum.org/v1/gonum/mat"
	"math"
)

// gaussJordanInversion inversion via gauss-jordan elimination
// 矩阵求逆 高斯-若尔当消元法
func gaussJordanInversion(X []float64, n int) bool {
	// m 是列数, n 是行数
	var m = n
	var b = make([]float64, n*n)
	var indexC = make([]int, n)
	var indexR = make([]int, n)
	var ipiV = make([]int, n)
	var iCol, iRow int
	var dum, pivinv float64

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				b[i*n+j] = 1
			} else {
				b[i*n+j] = 0
			}
		}
	}

	for j := 0; j < n; j++ {
		ipiV[j] = 0
	}

	for i := 0; i < n; i++ {
		var big float64 = 0
		for j := 0; j < n; j++ {
			if ipiV[j] != 1 {
				for k := 0; k < n; k++ {
					if ipiV[k] == 0 {
						absoluteValue := math.Abs(X[j*n+k])
						if absoluteValue >= big {
							big = absoluteValue
							iRow = j
							iCol = k
						}
					}
				}
			}
		}
		ipiV[iCol]++
		if iRow != iCol {
			for l := 0; l < n; l++ {
				X[iRow*n+l], X[iCol*n+l] = X[iCol*n+l], X[iRow*n+l]
			}
			for l := 0; l < m; l++ {
				b[iRow*n+l], b[iCol*n+l] = b[iCol*n+l], b[iRow*n+l]
			}
		}
		indexR[i] = iRow
		indexC[i] = iCol

		if X[iCol*n+iCol] == 0 {
			return false
		} // Singular

		pivinv = 1 / X[iCol*n+iCol]
		X[iCol*n+iCol] = 1
		for l := 0; l < n; l++ {
			X[iCol*n+l] *= pivinv
		}
		for l := 0; l < m; l++ {
			b[iCol*n+l] *= pivinv
		}

		for ll := 0; ll < n; ll++ {
			if ll != iCol {
				dum = X[ll*n+iCol]
				X[ll*n+iCol] = 0
				for l := 0; l < n; l++ {
					X[ll*n+l] -= X[iCol*n+l] * dum
				}
				for l := 0; l < m; l++ {
					b[ll*n+l] -= b[iCol*n+l] * dum
				}
			}
		}
	}

	for l := n - 1; l >= 0; l-- {
		if indexR[l] != indexC[l] {
			for k := 0; k < n; k++ {
				X[k*n+indexR[l]], X[k*n+indexC[l]] = X[k*n+indexC[l]], X[k*n+indexR[l]]
			}
		}
	}

	return true
}

// matrixInverseByCol 矩阵求逆 列主元消去法
// https://github.com/zourongcsu/matrxi-algorithm/blob/master/inverse_matrix.cpp
func matrixInverseByCol(a [][]float64) ([][]float64, bool) {
	// eMat 返回 n 阶单位矩阵
	var eMat = func(n int) ([][]float64, bool) {
		sol := make([][]float64, n)
		for i := 0; i < n; i++ {
			sol[i] = make([]float64, n)
		}

		// 判断阶数
		if n < 1 {
			return nil, false
		}

		//分配元素
		for i := 0; i < n; i++ {
			sol[i][i] = 1.0
		}

		return sol, true
	}
	// maxAbs 向量第一个绝对值最大值及其位置
	var maxAbs = func(a []float64) (float64, int, bool) {
		var sol float64
		var ii int
		var err = false

		n := len(a)
		ii = 0
		sol = a[ii]
		for i := 1; i < n; i++ {
			if math.Abs(sol) < math.Abs(a[i]) {
				ii = i
				sol = a[i]
			}
		}

		err = true
		return sol, ii, err
	}

	//判断是否方阵
	if len(a) != len(a[0]) {
		return nil, false
	}

	n := len(a)
	sol, isOk := eMat(n)
	if !isOk {
		return nil, false
	}

	//lint:ignore SA4006 for temp
	temp1 := make([]float64, n)

	//主元消去
	for i := 0; i < n; i++ {
		//求第i列的主元素并调整行顺序
		aCol := make([]float64, n-i)
		for iCol := i; iCol < n; iCol++ {
			aCol[iCol-i] = a[iCol][i]
		}
		_, ii, _ := maxAbs(aCol)
		if ii+i != i {
			temp1 = a[ii+i]
			a[ii+i] = a[i]
			a[i] = temp1
			temp1 = sol[ii+i]
			sol[ii+i] = sol[i]
			sol[i] = temp1
		}

		//列消去
		//本行主元置一
		mul := a[i][i]
		for j := 0; j < n; j++ {
			a[i][j] = a[i][j] / mul
			sol[i][j] = sol[i][j] / mul
		}
		//其它列置零
		for j := 0; j < n; j++ {
			if j != i {
				mul = a[j][i] / a[i][i]
				for k := 0; k < n; k++ {
					a[j][k] = a[j][k] - a[i][k]*mul
					sol[j][k] = sol[j][k] - sol[i][k]*mul
				}
			}
		}
	}

	return sol, true
}

func matrixInverseByCol_(x []float64, n int) ([]float64, bool) {
	a := make([][]float64, n)
	for i := 0; i < n; i++ {
		a[i] = x[i*n : i*n+n]
	}
	aa, ok := matrixInverseByCol(a)
	if !ok {
		return x, false
	}

	var xx []float64
	for i := 0; i < n; i++ {
		xx = append(xx, aa[i]...)
	}

	return xx, false
}

func matrixInverse(x []float64, n int) ([]float64, bool) {
	a := mat.NewDense(n, n, x)
	var ia mat.Dense

	// Take the inverse of a and place the result in ia.
	err := ia.Inverse(a)
	if err != nil {
		return ia.RawMatrix().Data, false
	}

	return ia.RawMatrix().Data, true
}
