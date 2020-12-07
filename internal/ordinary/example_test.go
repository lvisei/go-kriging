package ordinary_test

import (
	"fmt"
	"github.com/liuvigongzuoshi/go-kriging/internal/ordinary"
)

var (
	values = FloatList{45.986076009952846, 46.223032113384235, 52.821454425024626, 89.19253247046487, 31.062802427638776}
	lats   = FloatList{117.99598607600996, 117.99622303211338, 118.00282145442502, 118.03919253247047, 117.98106280242764}
	lons   = FloatList{31.995986076009952, 31.99622303211338, 32.002821454425025, 32.03919253247046, 31.981062802427637}
)

type FloatList []float64

func (t FloatList) min() float64 {
	min := float64(0)
	for i := 0; i < len(t); i++ {
		if min == 0 || min > t[i] {
			min = t[i]
		}
	}

	return min
}

func (t FloatList) max() float64 {
	max := float64(0)
	for i := 0; i < len(t); i++ {
		if max < t[i] {
			max = t[i]
		}
	}

	return max
}

func Example_OrdinaryKriging_Exponential() {
	ordinaryKriging := ordinary.NewOrdinary(values, lats, lons)
	ordinaryKriging.Train(ordinary.Exponential, 0, 100)
	krigingValue, _, _ := ordinaryKriging.GeneratePngGrid(200, 200)
	fmt.Printf("%#v", krigingValue[:10])
	// Output:
	// []float64{31.832037963867943, 32.2958798424233, 32.83227510797089, 33.38759955944839, 33.946936215772, 34.50483331478148, 35.05891641620112, 35.60803760024491, 36.15161335849867, 36.68934894716563}

}

func Example_OrdinaryKriging_Spherical() {
	ordinaryKriging := ordinary.NewOrdinary(values, lats, lons)
	ordinaryKriging.Train(ordinary.Spherical, 0, 100)
	krigingValue, _, _ := ordinaryKriging.GeneratePngGrid(200, 200)
	fmt.Printf("%#v", krigingValue[:10])
	// Output:
	// []float64{31.42123794010535, 31.644499557807247, 31.90817908629488, 32.185341522126315, 32.46845610328156, 32.75473696943352, 33.042940816984725, 33.33243905472105, 33.622885866814954, 33.91408001569965}

}

func Example_OrdinaryKriging_Gaussian() {
	ordinaryKriging := ordinary.NewOrdinary(values, lats, lons)
	ordinaryKriging.Train(ordinary.Gaussian, 0, 100)
	krigingValue, _, _ := ordinaryKriging.GeneratePngGrid(200, 200)
	fmt.Printf("%#v", krigingValue[:10])
	// Output:
	// []float64{31.3237814741247, 31.456757078398557, 31.593073987805035, 31.73269631764414, 31.875587884114843, 32.02171222008403, 32.17103259071844, 32.32351200912366, 32.47911325177369, 32.63779887397147}

}
