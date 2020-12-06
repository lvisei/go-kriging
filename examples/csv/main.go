package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/liuvigongzuoshi/go-kriging/internal/ordinary"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strconv"
	"time"
)

const dirPath = "testdata"

func main() {
	cpuProfile, _ := os.Create("./testdata/cpu_profile")
	pprof.StartCPUProfile(cpuProfile)
	defer pprof.StopCPUProfile()
	//memProfile, _ := os.Create("./testdata/mem_profile")
	//pprof.WriteHeapProfile(memProfile)
	//defer memProfile.Close()
	//var wg sync.WaitGroup
	//wg.Add(1)
	data, err := readCsvFile("examples/csv/data/2045.csv")
	if err != nil {
		log.Fatal(err)
	}
	polygon, err := readGeoJsonFile("examples/csv/data/yn.json")
	if err != nil {
		log.Fatal(err)
	}
	defer timeCost()("训练模型加插值总耗时")

	ordinaryKriging := ordinary.NewOrdinary(data["values"], data["lons"], data["lats"])
	_ = ordinaryKriging.Train(ordinary.Exponential, 0, 100)
	//writeFile("variogram.json", variogram)

	//go func() {
	//	defer wg.Done()
	//	generateGridData(ordinaryKriging, polygon)
	//}()
	generateGridData(ordinaryKriging, polygon)

	//go func() {
	//	defer wg.Done()
	//	generatePng(ordinaryKriging)
	//}()

	//wg.Wait()
}

func generateGridData(ordinaryKriging *ordinary.Variogram, polygon ordinary.Polygon) {
	defer timeCost()("插值耗时")
	gridMatrices := ordinaryKriging.Grid(polygon, 0.01)
	_ = gridMatrices
	writeFile("gridMatrices.json", gridMatrices)
}

func generatePng(ordinaryKriging *ordinary.Variogram) {
	defer timeCost()("生成插值图片耗时")
	xWidth, yWidth := 800, 800
	krigingValue, rangeMaxPM, colorperiod := ordinaryKriging.GeneratePngGrid(xWidth, yWidth)
	pngPath := fmt.Sprintf("%v/%v.png", dirPath, time.Now().Format("2006-01-02 15:04:05"))
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

func generateRectangleGrid() {
	//bound := polygon.Bound()
	//bbox := [4]float64{bound.Min.Lat(), bound.Min.Lon(), bound.Max.Lat(), bound.Max.Lon()}
	//gridDate := ordinaryKriging.RectangleGrid(bbox, 0.01)
	//writeFile("gridDate.json", gridDate)

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
		if lon, err = strconv.ParseFloat(records[i][1], 64); err != nil {
			return nil, err
		}
		lons = append(lons, lon)
		if lat, err = strconv.ParseFloat(records[i][2], 64); err != nil {
			return nil, err
		}
		lats = append(lats, lat)
		if value, err = strconv.ParseFloat(records[i][5], 64); err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	data := map[string][]float64{"values": values, "lats": lats, "lons": lons}

	//fmt.Printf("values %#v\n lons %#v\n lats %#v\n", values, lons, lats)

	return data, nil
}

func readGeoJsonFile(filePath string) (ordinary.Polygon, error) {
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

	p := make(ordinary.Polygon, 0, len(polygon))
	for _, ring := range polygon {
		points := make([][2]float64, 0, len(ring))
		for _, point := range ring {
			points = append(points, [2]float64{point.X(), point.Y()})
		}
		p = append(p, points)
	}

	return p, nil
}

func timeCost() func(name string) {
	start := time.Now()
	return func(name string) {
		tc := time.Since(start)
		fmt.Printf("%v : time cost = %v s\n", name, tc.Seconds())
	}
}

func writeFile(fileName string, v interface{}) {
	filePath := fmt.Sprintf("%v/%v %v", dirPath, time.Now().Format("2006-01-02 15:04:05"), fileName)
	fmt.Printf("%#v\n", filePath)
	// fmt.Printf("%#v\n", v)
	content, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(filePath, content, os.ModePerm)
}
