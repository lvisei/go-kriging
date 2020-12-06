package ordinary_test

import (
	"fmt"
	"github.com/liuvigongzuoshi/go-kriging/internal/ordinary"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const pngDirPath = "testdata"

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
	var values, lats, lons FloatList
	values = append(values, randFloats(0, 100, count)...)
	lats = append(lats, randFloats(117.95, 118.05, count)...)
	lons = append(lons, randFloats(31.95, 32.05, count)...)
	return values, lats, lons
}

func TestGenerateData(t *testing.T) {
	values, lats, lons := generateData(100)
	fmt.Printf("values %v\n", values)
	fmt.Printf("lats %v\n", lats)
	fmt.Printf("lons %v\n", lons)
}

func TestTrain(t *testing.T) {
	values, lats, lons := generateData(100)
	ordinaryKriging := ordinary.NewOrdinary(values, lats, lons)
	ordinaryKriging.Train(ordinary.Exponential, 0, 100)
}

func TestPredict(t *testing.T) {
	values, lats, lons := generateData(100)
	ordinaryKriging := ordinary.NewOrdinary(values, lats, lons)
	ordinaryKriging.Train(ordinary.Exponential, 0, 100)
	ordinaryKriging.GeneratePngGrid(200, 2000)
}

func TestOrdinaryKriging(t *testing.T) {
	values, lats, lons := generateData(100)

	ordinaryKriging := ordinary.NewOrdinary(values, lats, lons)
	ordinaryKriging.Train(ordinary.Exponential, 0, 100)

	xWidth, yWidth := 500, 500
	krigingValue, rangeMaxPM, colorperiod := ordinaryKriging.GeneratePngGrid(xWidth, yWidth)
	pngPath := fmt.Sprintf("%v/%v.png", pngDirPath, time.Now().Format("2006-01-02 15:04:05"))
	img := ordinaryKriging.GeneratePng(krigingValue, rangeMaxPM, colorperiod, xWidth, yWidth)

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
