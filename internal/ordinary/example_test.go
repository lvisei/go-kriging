package ordinary_test

import (
	"fmt"
	"github.com/liuvigongzuoshi/go-kriging/internal/ordinary"
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

func Example_OrdinaryKriging() {
	values := FloatList{45.986076009952846, 46.223032113384235, 52.821454425024626, 89.19253247046487, 31.062802427638776}
	lats := FloatList{117.99598607600996, 117.99622303211338, 118.00282145442502, 118.03919253247047, 117.98106280242764}
	lons := FloatList{31.995986076009952, 31.99622303211338, 32.002821454425025, 32.03919253247046, 31.981062802427637}

	ordinaryKriging := &ordinary.Variogram{T: values, X: lats, Y: lons}
	ordinaryKriging.Train(ordinary.Exponential, 0, 100)
	krigingValue, _, _ := ordinaryKriging.GeneratePngGrid(200, 200)
	fmt.Printf("%#v", krigingValue[:10])
	// Output:
	// []float64{32.58034918188755, 33.47952397605678, 34.50902699390895, 35.56247811084377, 36.61093522661021, 37.644407996545574, 38.65908202680655, 39.653617264358715, 40.627823344662545, 41.58210725922872}

}
