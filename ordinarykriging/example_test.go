package ordinarykriging_test

import (
	"fmt"

	"github.com/liuvigongzuoshi/go-kriging/ordinarykriging"
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

func ExampleVariogram_Contour_Exponential() {
	ordinaryKriging := ordinarykriging.NewOrdinary(values, lats, lons)
	ordinaryKriging.Train(ordinarykriging.Exponential, 0, 100)
	contourRectangle := ordinaryKriging.Contour(200, 200)
	fmt.Printf("%#v", contourRectangle.Contour[:10])
	// Output:
	// []float64{31.062802427639, 31.67443506838088, 32.27805611994289, 32.87380457354172, 33.461820447530236, 34.042244827188064, 34.61521990152406, 35.18088899653797, 35.73939660436969, 36.29088840795865}

}

func ExampleVariogram_Contour_Spherical() {
	ordinaryKriging := ordinarykriging.NewOrdinary(values, lats, lons)
	ordinaryKriging.Train(ordinarykriging.Spherical, 0, 100)
	contourRectangle := ordinaryKriging.Contour(200, 200)
	fmt.Printf("%#v", contourRectangle.Contour[:10])
	// Output:
	// []float64{31.062802427637955, 31.35568613698695, 31.649070507173384, 31.942636980745615, 32.23631166432933, 32.53011270795399, 32.82405871070559, 33.11816872323164, 33.41246224930053, 33.70695924637588}

}

func ExampleVariogram_Contour_Gaussian() {
	ordinaryKriging := ordinarykriging.NewOrdinary(values, lats, lons)
	ordinaryKriging.Train(ordinarykriging.Gaussian, 0, 100)
	contourRectangle := ordinaryKriging.Contour(200, 200)
	fmt.Printf("%#v", contourRectangle.Contour[:10])
	// Output:
	// []float64{31.062802438363697, 31.19418275387546, 31.32895545000455, 31.467084521380848, 31.60853364267539, 31.753266184588494, 31.901245229805934, 32.052433588854164, 32.206793815838, 32.36428822416533}

}
