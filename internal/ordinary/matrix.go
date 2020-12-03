package ordinary

import "math"

// krigingMatrixTranspose 矩阵颠倒，横向矩阵变成纵向矩阵
func krigingMatrixTranspose(X []float64, n, m int) []float64 {
	Z := make([]float64, m*n)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			Z[j*n+i] = X[i*m+j]
		}
	}

	return Z
}

// krigingMatrixMultiply naive matrix multiplication 矩阵相乘, 横向矩阵*纵向矩阵
func krigingMatrixMultiply(X, Y []float64, n, m, p int) []float64 {
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

func krigingMatrixAdd(X, Y []float64, n, m int) []float64 {
	Z := make([]float64, n*m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			Z[i*m+j] = X[i*m+j] + Y[i*m+j]
		}
	}

	return Z
}

// krigingMatrixDiag matrix algebra
func krigingMatrixDiag(c float64, n int) []float64 {
	Z := make([]float64, n*n)
	for i := 0; i < n; i++ {
		Z[i*n+i] = c
	}
	return Z
}

// krigingMatrixChol cholesky decomposition
func krigingMatrixChol(X []float64, n int) bool {
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

// krigingMatrixChol2inv inversion of cholesky decomposition
func krigingMatrixChol2inv(X []float64, n int) {
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

// krigingMatrixSolve inversion via gauss-jordan elimination
func krigingMatrixSolve(X []float64, n int) bool {
	var m = n
	var b = make([]float64, n*n)
	var indxc = make([]int, n)
	var indxr = make([]int, n)
	var ipiv = make([]int, n)
	var i, icol, irow, j, k, l, ll int
	var big, dum, pivinv, temp float64

	for i = 0; i < n; i++ {
		for j = 0; j < n; j++ {
			if i == j {
				b[i*n+j] = 1
			} else {
				b[i*n+j] = 0
			}
		}
	}

	for j = 0; j < n; j++ {
		ipiv[j] = 0
	}

	for i = 0; i < n; i++ {
		big = 0
		for j = 0; j < n; j++ {
			if ipiv[j] != 1 {
				for k = 0; k < n; k++ {
					if ipiv[k] == 0 {
						if math.Abs(X[j*n+k]) >= big {
							big = math.Abs(X[j*n+k])
							irow = j
							icol = k
						}
					}
				}
			}
		}
		ipiv[icol]++
		if irow != icol {
			for l = 0; l < n; l++ {
				temp = X[irow*n+l]
				X[irow*n+l] = X[icol*n+l]
				X[icol*n+l] = temp
			}
			for l = 0; l < m; l++ {
				temp = b[irow*n+l]
				b[irow*n+l] = b[icol*n+l]
				b[icol*n+l] = temp
			}
		}
		indxr[i] = irow
		indxc[i] = icol

		if X[icol*n+icol] == 0 {
			return false
		} // Singular

		pivinv = 1 / X[icol*n+icol]
		X[icol*n+icol] = 1
		for l = 0; l < n; l++ {
			X[icol*n+l] *= pivinv
		}
		for l = 0; l < m; l++ {
			b[icol*n+l] *= pivinv
		}

		for ll = 0; ll < n; ll++ {
			if ll != icol {
				dum = X[ll*n+icol]
				X[ll*n+icol] = 0
				for l = 0; l < n; l++ {
					X[ll*n+l] -= X[icol*n+l] * dum
				}
				for l = 0; l < m; l++ {
					b[ll*n+l] -= b[icol*n+l] * dum
				}
			}
		}
	}
	for l = n - 1; l >= 0; l-- {
		if indxr[l] != indxc[l] {
			for k = 0; k < n; k++ {
				temp = X[k*n+indxr[l]]
				X[k*n+indxr[l]] = X[k*n+indxc[l]]
				X[k*n+indxc[l]] = temp
			}
		}
	}

	return true
}