package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/liuvigongzuoshi/go-kriging/internal/ordinary"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

const dirPath = "testdata"

func main() {
	data, err := readCsvFile("examples/csv/data/2292.csv")
	if err != nil {
		log.Fatal(err)
	}
	polygon, err := readGeoJsonFile("examples/csv/data/yn.json")
	if err != nil {
		log.Fatal(err)
	}
	defer timeCost()()

	ordinaryKriging := &ordinary.Variogram{T: data["values"], X: data["lats"], Y: data["lons"]}
	_ = ordinaryKriging.Train(ordinary.Exponential, 0, 100)
	//writeFile("variogram.json", variogram)

	gridMatrices := ordinaryKriging.Grid(polygon, 0.01)
	writeFile("gridMatrices.json", gridMatrices)

	//bound := polygon.Bound()
	//bbox := [4]float64{bound.Min.Lat(), bound.Min.Lon(), bound.Max.Lat(), bound.Max.Lon()}
	//gridDate := ordinaryKriging.RectangleGrid(bbox, 0.01)
	//writeFile("gridDate.json", gridDate)

	//xWidth, yWidth := 800, 800
	//krigingValue, rangeMaxPM, colorperiod := ordinaryKriging.GeneratePngGrid(xWidth, yWidth)
	//pngPath := fmt.Sprintf("%v/%v.png", dirPath, time.Now().Format("2006-01-02 15:04:05"))
	//img := ordinaryKriging.GeneratePng(krigingValue, rangeMaxPM, colorperiod, xWidth, yWidth)
	//err = os.MkdirAll(filepath.Dir(pngPath), os.ModePerm)
	//if err != nil {
	//	return
	//}
	//file, err := os.Create(pngPath)
	//if err != nil {
	//	return
	//}
	//defer file.Close()
	//png.Encode(file, img)
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

func readGeoJsonFile(filePath string) (orb.Polygon, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
		return nil, err
	}
	fc := geojson.NewFeatureCollection()
	err = json.Unmarshal(content, &fc)
	if err != nil {
		log.Fatalf("invalid json: %v", err)
		return nil, err
	}
	polygon := fc.Features[0].Geometry.(orb.Polygon)

	return polygon, nil
}

func timeCost() func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		fmt.Printf("time cost = %v s\n", tc.Seconds())
	}
}

func writeFile(fileName string, v interface{}) {
	filePath := fmt.Sprintf("%v/%v %v", dirPath, time.Now().Format("2006-01-02 15:04:05"), fileName)
	fmt.Printf("%#v\n", filePath)
	content, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(filePath, content, os.ModePerm)
}
