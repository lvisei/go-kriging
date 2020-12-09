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

	"github.com/liuvigongzuoshi/go-kriging/canvas"
	"github.com/liuvigongzuoshi/go-kriging/ordinary"
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

	data, err := readCsvFile("examples/csv/data/2045.csv")
	if err != nil {
		log.Fatal(err)
	}
	polygon, err := readGeoJsonFile("examples/csv/data/yn.json")
	if err != nil {
		log.Fatal(err)
	}
	defer timeCost()("训练模型加插值总耗时")

	ordinaryKriging := ordinary.NewOrdinary(data["values"], data["x"], data["y"])
	_ = ordinaryKriging.Train(ordinary.Spherical, 0, 100)

	gridPlot(ordinaryKriging, polygon)

	//var wg sync.WaitGroup
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	gridPlot(ordinaryKriging, polygon)
	//}()
	//go func() {
	//	defer wg.Done()
	//	contourRectanglePng(ordinaryKriging)
	//}()

	//wg.Wait()
}

func gridPlot(ordinaryKriging *ordinary.Variogram, polygon ordinary.PolygonCoordinates) {
	defer timeCost()("插值生成网格图片耗时")
	gridMatrices := ordinaryKriging.Grid(polygon, 0.01)
	ctx := ordinaryKriging.Plot(gridMatrices, 500, 500, gridMatrices.Xlim, gridMatrices.Ylim, ordinary.DefaultGridLevelColor)

	subTitle := &canvas.TextConfig{
		Text:     "球面半变异函数模型",
		FontName: "testdata/fonts/source-han-sans-sc/regular.ttf",
		FontSize: 28,
		Color:    color.RGBA{R: 0, G: 0, B: 0, A: 255},
		OffsetX:  252,
		OffsetY:  40,
		AlignX:   0.5,
	}
	if err := ctx.DrawText(subTitle); err != nil {
		log.Fatalf("DrawText %v", err)
	}

	buffer, err := ctx.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		saveBuffer("grid.png", buffer)
	}

	//writeFile("gridMatrices.json", gridMatrices)
}

func contourRectanglePng(ordinaryKriging *ordinary.Variogram) {
	defer timeCost()("插值生成矩形图片耗时")
	xWidth, yWidth := 800, 800
	contourRectangle := ordinaryKriging.Contour(xWidth, yWidth)
	pngPath := fmt.Sprintf("%v/%v.png", dirPath, time.Now().Format("2006-01-02 15:04:05"))
	ctx := ordinaryKriging.PlotRectangleContour(contourRectangle, 500, 500, contourRectangle.Xlim, contourRectangle.Ylim, ordinary.DefaultLegendColor)
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

	buffer, err := ctx.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		saveBuffer("rectangle.png", buffer)
	}
}

func ContourWithBBoxPng(bbox [4]float64) {
	//contourRectangle := ordinaryKriging.ContourWithBBox(bbox, 0.01)
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

	data := map[string][]float64{"values": values, "x": lons, "y": lats}

	//fmt.Printf("values %#v\n lons %#v\n lats %#v\n", values, lons, lats)
	// writeFile("data.json", data)

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
