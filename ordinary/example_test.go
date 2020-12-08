package ordinary_test

import (
	"fmt"

	"github.com/liuvigongzuoshi/go-kriging/ordinary"
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
	contourRectangle := ordinaryKriging.Contour(200, 200)
	fmt.Printf("%#v", contourRectangle.Contour[:10])
	// Output:
	// []float64{31.062802427638875, 31.674435068380763, 32.278056119942775, 32.873804573541605, 33.461820447530116, 34.04224482718795, 34.615219901523936, 35.18088899653787, 35.73939660436959, 36.29088840795853}

}

func Example_OrdinaryKriging_Spherical() {
	ordinaryKriging := ordinary.NewOrdinary(values, lats, lons)
	ordinaryKriging.Train(ordinary.Spherical, 0, 100)
	contourRectangle := ordinaryKriging.Contour(200, 200)
	fmt.Printf("%#v", contourRectangle.Contour[:10])
	// Output:
	// []float64{31.062802427638914, 31.355686136987902, 31.64907050717433, 31.94263698074654, 32.23631166433025, 32.5301127079549, 32.82405871070649, 33.11816872323253, 33.41246224930141, 33.70695924637675}

}

func Example_OrdinaryKriging_Gaussian() {
	ordinaryKriging := ordinary.NewOrdinary(values, lats, lons)
	ordinaryKriging.Train(ordinary.Gaussian, 0, 100)
	contourRectangle := ordinaryKriging.Contour(200, 200)
	fmt.Printf("%#v", contourRectangle.Contour[:10])
	// Output:
	// []float64{31.062802429121923, 31.194182744748254, 31.328955440993354, 31.467084512481577, 31.60853363388739, 31.75326617591121, 31.90124522123483, 32.05243358038846, 32.206793807474796, 32.36428821590181}

}
