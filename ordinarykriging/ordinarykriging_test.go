package ordinarykriging_test

import (
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/lvisei/go-kriging/ordinarykriging"
)

const pngDirPath = "tempdata"

var (
	randomValues, randomLats, randomLons = generateData(100)
)

func randFloats(min, max float64, n int) []float64 {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	res := make([]float64, n)
	for i := range res {
		res[i] = min + r.Float64()*(max-min)
	}
	return res
}

func generateData(count int) (FloatList, FloatList, FloatList) {
	var randomValues, randomLats, randomLons FloatList
	randomValues = append(randomValues, randFloats(0, 100, count)...)
	randomLats = append(randomLats, randFloats(117.95, 118.05, count)...)
	randomLons = append(randomLons, randFloats(31.95, 32.05, count)...)
	return randomValues, randomLats, randomLons
}

func TestVariogram_Predict_100(t *testing.T) {
	var randomValues, randomLats, randomLons FloatList
	randomLons = append(randomLons, randFloats(-90, 90, 100)...)
	randomLats = append(randomLats, randFloats(-180, 180, 100)...)
	for _, lon := range randomLons {
		var value float64
		if lon > 0 {
			value = 100
		} else {
			value = 0
		}
		randomValues = append(randomValues, value)
	}
	ordinaryKriging := ordinarykriging.NewOrdinary(randomValues, randomLats, randomLons)
	if _, err := ordinaryKriging.Train(ordinarykriging.Exponential, 0, 10); err != nil {
		t.Fatal("variogram is null", err)
	}
	if ordinaryKriging.Predict(180, 0) < 50 {
		t.Fatal("unexpected result (<50)")
	}
	if ordinaryKriging.Predict(-180, 0) > 50 {
		t.Fatal("unexpected result (>50)")
	}
}

func TestVariogram_Predict_1000(t *testing.T) {
	var randomValues, randomLats, randomLons FloatList
	randomLons = append(randomLons, randFloats(-90, 90, 1000)...)
	randomLats = append(randomLats, randFloats(-180, 180, 1000)...)
	for _, lon := range randomLons {
		var value float64
		if lon > 0 {
			value = 100
		} else {
			value = 0
		}
		randomValues = append(randomValues, value)
	}
	ordinaryKriging := ordinarykriging.NewOrdinary(randomValues, randomLats, randomLons)
	if _, err := ordinaryKriging.Train(ordinarykriging.Exponential, 0, 10); err != nil {
		t.Fatal("variogram is null", err)
	}
	if ordinaryKriging.Predict(180, 0) < 50 {
		t.Fatal("unexpected result (<50)")
	}
	if ordinaryKriging.Predict(-180, 0) > 50 {
		t.Fatal("unexpected result (>50)")
	}
}

func TestVariogram_Plot(t *testing.T) {
	ordinaryKriging := ordinarykriging.NewOrdinary(randomValues, randomLats, randomLons)
	ordinaryKriging.Train(ordinarykriging.Exponential, 0, 100)

	contourRectangle := ordinaryKriging.Contour(500, 500)
	pngPath := fmt.Sprintf("%v/%v.png", pngDirPath, time.Now().Format("2006-01-02 15:04:05"))
	img := ordinaryKriging.PlotPng(contourRectangle)

	err := os.MkdirAll(filepath.Dir(pngPath), os.ModePerm)
	if err != nil {
		return
	}
	file, err := os.Create(pngPath)
	if err != nil {
		return
	}
	defer file.Close()
	png.Encode(file, img)
}

func BenchmarkVariogram_Train_Exponential(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ordinaryKriging := ordinarykriging.NewOrdinary(randomValues, randomLats, randomLons)
		ordinaryKriging.Train(ordinarykriging.Exponential, 0, 100)
	}
}

func BenchmarkVariogram_Train_Spherical(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ordinaryKriging := ordinarykriging.NewOrdinary(randomValues, randomLats, randomLons)
		ordinaryKriging.Train(ordinarykriging.Spherical, 0, 100)
	}
}

func BenchmarkVariogram_Train_Gaussian(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ordinaryKriging := ordinarykriging.NewOrdinary(randomValues, randomLats, randomLons)
		ordinaryKriging.Train(ordinarykriging.Gaussian, 0, 100)
	}
}

func BenchmarkVariogram_Contour(b *testing.B) {
	ordinaryKriging := ordinarykriging.NewOrdinary(randomValues, randomLats, randomLons)
	ordinaryKriging.Train(ordinarykriging.Exponential, 0, 100)
	for n := 0; n < b.N; n++ {
		ordinaryKriging.Contour(600, 600)
	}
}
