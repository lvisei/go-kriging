package main

import (
	"encoding/csv"
	"fmt"
	"github.com/liuvigongzuoshi/go-kriging/internal/ordinary"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const pngDirPath = "testdata"

func main() {
	data, err := readCsvFile("examples/csv/data/2292.csv")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%#v\n", data)

	defer timeCost()()

	ordinaryKriging := &ordinary.Variogram{T: data["values"], X: data["lats"], Y: data["lons"]}
	ordinaryKriging.Train(ordinary.Exponential, 0, 100)

	//ordinaryKriging.Grid(0.01)
	//ordinaryKriging.RectangleGrid(0.01)
	xWidth, yWidth := 400, 400
	krigingValue, rangeMaxPM, colorperiod := ordinaryKriging.GeneratePngGrid(xWidth, yWidth)
	pngPath := fmt.Sprintf("%v/%v.png", pngDirPath, time.Now().Format("2006-01-02 15:04:05"))
	img := ordinaryKriging.GeneratePng(krigingValue, rangeMaxPM, colorperiod, xWidth, yWidth)
	err = os.MkdirAll(filepath.Dir(pngPath), os.ModePerm)
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

func readCsvFile(filePath string) (map[string][]float64, error) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
		return nil, err
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
		return nil, err
	}

	var values, lats, lons []float64

	for i := 1; i < len(records); i++ {
		var value, lat, lon float64
		if lat, err = strconv.ParseFloat(records[i][1], 64); err != nil {
			return nil, err
		}
		lats = append(lats, lat)
		if lon, err = strconv.ParseFloat(records[i][2], 64); err != nil {
			return nil, err
		}
		lons = append(lons, lon)
		if value, err = strconv.ParseFloat(records[i][5], 64); err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return map[string][]float64{"values": values, "lats": lats, "lons": lons}, nil
}

func timeCost() func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		fmt.Printf("time cost = %v s\n", tc.Seconds())
	}
}
