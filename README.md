# go-kriging

[![GoDoc](https://godoc.org/github.com/liuvigongzuoshi/go-kriging?status.svg)](https://pkg.go.dev/github.com/liuvigongzuoshi/go-kriging)

Golang Multi-Goroutine spatial interpolation algorithm library for geospatial prediction and mapping via ordinary kriging.

Based on [oeo4b/kriging.js](https://github.com/oeo4b/kriging.js)  refactoring and optimized the algorithm and added some new features. 

## Fitting a Model

The train method with the new ordinaryKriging fits your input to whatever variogram model you specify - gaussian, exponential or spherical - and returns a variogram variable.


```go
import "github.com/liuvigongzuoshi/go-kriging/ordinarykriging"

func main() {
  sigma2 := 0
  alpha := 100
  ordinaryKriging := ordinarykriging.NewOrdinary(values, x, y)
  variogram = ordinaryKriging.Train(ordinarykriging.Spherical, sigma2, alpha)  
}
```

## Predicting New Values

Values can be predicted for new coordinate pairs by using the predict method with the new ordinaryKriging.

```go
import "github.com/liuvigongzuoshi/go-kriging/ordinarykriging"

func main() {
  // ...
  // Pair of new coordinates to predict
  xnew := 0.5481
  ynew := 0.4455
  tpredicted := ordinaryKriging.predict(xnew, ynew)
}
```

## Variogram and Probability Model

The various variogram models can be interpreted as kernel functions for 2-dimensional coordinates a, b and parameters nugget, range, sill and A. Reparameterized as a linear function, with w = [nugget, (sill-nugget)/range], this becomes:

- Gaussian: k(a,b) = w[0] + w[1] * ( 1 - exp{ -( ||a-b|| / range )2 / A } )
- Exponential: k(a,b) = w[0] + w[1] * ( 1 - exp{ -( ||a-b|| / range ) / A } )
- Spherical: k(a,b) = w[0] + w[1] * ( 1.5 * ( ||a-b|| / range ) - 0.5 * ( ||a-b|| / range )3 )

The variance parameter α of the prior distribution for w should be manually set, according to:

- w ~ N(w|0, αI)

Using the fitted kernel function hyperparameters and setting K as the Gram matrix, the prior and likelihood for the gaussian process become:

- y ~ N(y|0, K)
- t|y ~ N(t|y, σ2I)

The variance parameter σ2 of the likelihood reflects the error in the gaussian process and should be manually set.


## Other

[kriging-wasm example](https://github.com/liuvigongzuoshi/kriging-wasm) - Test example used by wasm compiled with go-kriging algorithm code.

[go-kriging-service](https://github.com/liuvigongzuoshi/go-kriging-service) - Call the REST service written by the go-kriging algorithm package, which supports concurrent calls by multiple users, and has a simple logging and fault-tolerant recovery mechanism.
