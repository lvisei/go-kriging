package ordinarykriging

import (
	"math"
)

// matrixTranspose The matrix is reversed, and the horizontal matrix becomes the vertical matrix
// 矩阵颠倒，横向矩阵变成纵向矩阵
func matrixTranspose(X []float64, n, m int) []float64 {
	Z := make([]float64, m*n)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			Z[j*n+i] = X[i*m+j]
		}
	}

	return Z
}

// matrixMultiply naive matrix multiplication
// 矩阵相乘, 横向矩阵乘纵向矩阵
func matrixMultiply(X, Y []float64, n, m, p int) []float64 {
	Z := make([]float64, n*p)
	for i := 0; i < n; i++ {
		for j := 0; j < p; j++ {
			Z[i*p+j] = 0
			for k := 0; k < m; k++ {
				Z[i*p+j] += X[i*m+k] * Y[k*p+j]
			}
		}
	}
	return Z
}

// matrixAdd
// 矩阵相加
func matrixAdd(X, Y []float64, n, m int) []float64 {
	Z := make([]float64, n*m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			Z[i*m+j] = X[i*m+j] + Y[i*m+j]
		}
	}

	return Z
}

// matrixDiag matrix algebra
func matrixDiag(c float64, n int) []float64 {
	Z := make([]float64, n*n)
	for i := 0; i < n; i++ {
		Z[i*n+i] = c
	}
	return Z
}

// matrixChol cholesky decomposition
// Cholesky 分解
func matrixChol(X []float64, n int) bool {
	p := make([]float64, n)

	for i := 0; i < n; i++ {
		p[i] = X[i*n+i]
	}

	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			p[i] -= X[i*n+j] * X[i*n+j]
		}
		if p[i] <= 0 {
			return false
		}
		p[i] = math.Sqrt(p[i])
		for j := i + 1; j < n; j++ {
			for k := 0; k < i; k++ {
				X[j*n+i] -= X[j*n+k] * X[i*n+k]
				X[j*n+i] /= p[i]
			}
		}
	}

	for i := 0; i < n; i++ {
		X[i*n+i] = p[i]
	}
	return true
}

// matrixChol2inv inversion of cholesky decomposition
// cholesky 分解的求逆
func matrixChol2inv(X []float64, n int) {
	var i, j, k int
	var sum float64

	for i = 0; i < n; i++ {
		X[i*n+i] = 1 / X[i*n+i]
		for j = i + 1; j < n; j++ {
			sum = 0
			for k = i; k < j; k++ {
				sum -= X[j*n+k] * X[k*n+i]
			}
			X[j*n+i] = sum / X[j*n+j]
		}
	}

	for i = 0; i < n; i++ {
		for j = i + 1; j < n; j++ {
			X[i*n+j] = 0
		}
	}

	for i = 0; i < n; i++ {
		X[i*n+i] *= X[i*n+i]
		for k = i + 1; k < n; k++ {
			X[i*n+i] += X[k*n+i] * X[k*n+i]
		}

		for j = i + 1; j < n; j++ {
			for k = j; k < n; k++ {
				X[i*n+j] += X[k*n+i] * X[k*n+j]
			}
		}
	}

	for i = 0; i < n; i++ {
		for j = 0; j < i; j++ {
			X[i*n+j] = X[j*n+i]
		}
	}
}
