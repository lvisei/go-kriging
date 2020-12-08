package main

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"runtime/pprof"

	"github.com/liuvigongzuoshi/go-kriging/internal/canvas"
	"github.com/liuvigongzuoshi/go-kriging/internal/ordinary"
	"github.com/liuvigongzuoshi/go-kriging/pkg/json"
)

const dirPath = "testdata"

func main() {
	cpuProfile, _ := os.Create("./testdata/cpu_profile")
	if err := pprof.StartCPUProfile(cpuProfile); err != nil {
		log.Fatal(err)
	}
	//memProfile, _ := os.Create("./testdata/mem_profile")
	//if err := pprof.WriteHeapProfile(memProfile); err != nil {
	//	log.Fatal(err)
	//}
	defer func() {
		pprof.StopCPUProfile()
		cpuProfile.Close()
		//memProfile.Close()
	}()

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
	_ = ordinaryKriging.Train(ordinary.Spherical, 0, 100)
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

func generateGridData(ordinaryKriging *ordinary.Variogram, polygon ordinary.PolygonCoordinates) {
	defer timeCost()("插值生成成功图片耗时")
	gridMatrices := ordinaryKriging.Grid(polygon, 0.01)
	//colors := []color.RGBA{
	//	color.RGBA{R: 52, G: 146, B: 199, A: 255},
	//	color.RGBA{R: 104, G: 166, B: 179, A: 255},
	//	color.RGBA{R: 149, G: 189, B: 158, A: 255},
	//	color.RGBA{R: 191, G: 212, B: 138, A: 255},
	//	color.RGBA{R: 231, G: 237, B: 114, A: 255},
	//	color.RGBA{R: 250, G: 228, B: 90, A: 255},
	//	color.RGBA{R: 248, G: 179, B: 68, A: 255},
	//	color.RGBA{R: 247, G: 133, B: 50, A: 255},
	//	color.RGBA{R: 242, G: 86, B: 34, A: 255},
	//	color.RGBA{R: 232, G: 21, B: 19, A: 255},
	//}

	colors := []ordinary.GridLevelColor{
		ordinary.GridLevelColor{Color: ordinary.RGBA{40, 146, 199, 255}, Value: [2]float64{-30, -15}},
		ordinary.GridLevelColor{Color: ordinary.RGBA{96, 163, 181, 255}, Value: [2]float64{-15, -10}},
		ordinary.GridLevelColor{Color: ordinary.RGBA{140, 184, 164, 255}, Value: [2]float64{-10, -5}},
		ordinary.GridLevelColor{Color: ordinary.RGBA{177, 204, 145, 255}, Value: [2]float64{-5, 0}},
		ordinary.GridLevelColor{Color: ordinary.RGBA{215, 227, 125, 255}, Value: [2]float64{0, 5}},
		ordinary.GridLevelColor{Color: ordinary.RGBA{250, 250, 100, 255}, Value: [2]float64{5, 10}},
		ordinary.GridLevelColor{Color: ordinary.RGBA{252, 207, 81, 255}, Value: [2]float64{10, 15}},
		ordinary.GridLevelColor{Color: ordinary.RGBA{252, 164, 63, 255}, Value: [2]float64{15, 20}},
		ordinary.GridLevelColor{Color: ordinary.RGBA{242, 77, 31, 255}, Value: [2]float64{25, 30}},
		ordinary.GridLevelColor{Color: ordinary.RGBA{232, 16, 20, 255}, Value: [2]float64{30, 40}},
	}
	ctx := ordinaryKriging.Plot(gridMatrices, 500, 500, gridMatrices.Xlim, gridMatrices.Ylim, colors)

	subTitle := &canvas.TextConfig{
		Text:     "球面半变异函数模型",
		FontName: "data/fonts/source-han-sans-sc/regular.ttf",
		FontSize: 28,
		Color:    color.RGBA{R: 0, G: 0, B: 0, A: 255},
		OffsetX:  252,
		OffsetY:  40,
		AlignX:   0.5,
	}
	if err := ctx.DrawText(subTitle); err != nil {
		log.Fatalf("DrawText %v", err)
	}

	//writeFile("gridMatrices.json", gridMatrices)
	buffer, err := ctx.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		saveBuffer("grid.png", buffer)
	}

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

func readGeoJsonFile(filePath string) (ordinary.PolygonCoordinates, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
		return nil, err
	}
	var polygonGeometry ordinary.PolygonGeometry
	err = json.Unmarshal(content, &polygonGeometry)
	if err != nil {
		log.Fatalf("invalid json: %v", err)
		return nil, err
	}

	return polygonGeometry.Coordinates, nil
}

func timeCost() func(name string) {
	start := time.Now()
	return func(name string) {
		tc := time.Since(start)
		fmt.Printf("%v : time cost = %v s\n", name, tc.Seconds())
	}
}

func writeFile(fileName string, v interface{}) {
	filePath := fmt.Sprintf("%v/%v %v", dirPath, time.Now().Format("2006-01-02 15-04-05"), fileName)
	fmt.Printf("%#v\n", filePath)
	// fmt.Printf("%#v\n", v)
	content, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(filePath, content, os.ModePerm)
}

func saveBuffer(fileName string, content []byte) {
	filePath := fmt.Sprintf("%v/%v %v", dirPath, time.Now().Format("2006-01-02 15-04-05"), fileName)
	fmt.Printf("%#v\n", filePath)
	ioutil.WriteFile(filePath, content, os.ModePerm)
}
