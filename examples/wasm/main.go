package main

import (
	"encoding/json"
	"fmt"
	"github.com/liuvigongzuoshi/go-kriging/internal/ordinary"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"log"
	"syscall/js"
)

func main() {
	fmt.Println("Instantiate, kriging WebAssembly!")
	done := make(chan int, 0)
	js.Global().Set("RunOrdinaryKriging", js.FuncOf(RunOrdinaryKrigingFunc))
	js.Global().Set("RunOrdinaryKrigingTrain", js.FuncOf(RunOrdinaryKrigingTrainFunc))
	<-done
}

func RunOrdinaryKrigingFunc(this js.Value, args []js.Value) interface{} {
	done := make(chan *ordinary.GridMatrices, 1)
	values := make([]float64, args[0].Length())
	for i := 0; i < len(values); i++ {
		values[i] = args[0].Index(i).Float()
	}
	lons := make([]float64, args[1].Length())
	for i := 0; i < len(lons); i++ {
		lons[i] = args[1].Index(i).Float()
	}
	lats := make([]float64, args[2].Length())
	for i := 0; i < len(lats); i++ {
		lats[i] = args[2].Index(i).Float()
	}
	model := args[3].String()
	sigma2 := args[4].Float()
	alpha := args[5].Float()

	geoJsonString := args[6].String()
	polygon, err := readGeoJsonBytes([]byte(geoJsonString))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		gridMatrices := RunOrdinaryKriging(values, lons, lats, model, sigma2, alpha, polygon)
		done <- gridMatrices
	}()

	gridMatrices := <-done

	gridBuffer, err := json.Marshal(gridMatrices)
	if err != nil {
		log.Fatal(err)
	}

	//js.Global().Get("console").Call("log", "wasm geoJsonString: ", js.ValueOf(geoJsonString))
	//js.Global().Get("console").Call("log", "wasm gridMatrices: ", string(gridBuffer))
	return string(gridBuffer)
}

func RunOrdinaryKrigingTrainFunc(this js.Value, args []js.Value) interface{} {
	values := make([]float64, args[0].Length())
	for i := 0; i < len(values); i++ {
		values[i] = args[0].Index(i).Float()
	}
	lons := make([]float64, args[1].Length())
	for i := 0; i < len(lons); i++ {
		lons[i] = args[1].Index(i).Float()
	}
	lats := make([]float64, args[2].Length())
	for i := 0; i < len(lats); i++ {
		lats[i] = args[2].Index(i).Float()
	}
	model := args[3].String()
	sigma2 := args[4].Float()
	alpha := args[5].Float()

	variogram := RunOrdinaryKrigingTrain(values, lons, lats, model, sigma2, alpha)
	variogramBuffer, err := json.Marshal(variogram)
	if err != nil {
		log.Fatal(err)
	}

	//js.Global().Get("console").Call("log", "wasm variogram: ", string(variogramBuffer))
	return string(variogramBuffer)
}

func RunOrdinaryKrigingTrain(values, lons, lats []float64, model string, sigma2 float64, alpha float64) *ordinary.Variogram {
	ordinaryKriging := ordinary.NewOrdinary(values, lons, lats)
	variogram := ordinaryKriging.Train(ordinary.ModelType(model), sigma2, alpha)
	return variogram
}

func RunOrdinaryKriging(values, lons, lats []float64, model string, sigma2 float64, alpha float64, polygon ordinary.Polygon) *ordinary.GridMatrices {
	ordinaryKriging := ordinary.NewOrdinary(values, lons, lats)
	_ = ordinaryKriging.Train(ordinary.ModelType(model), sigma2, alpha)
	return ordinaryKriging.Grid(polygon, 0.01)
}

func readGeoJsonBytes(content []byte) (ordinary.Polygon, error) {
	fc := geojson.NewFeatureCollection()
	err := json.Unmarshal(content, &fc)
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
