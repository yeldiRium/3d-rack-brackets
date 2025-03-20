package rack

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/ljanyst/ghostscad/primitive"
)

func square(i float64) float64 {
	return i * i
}

func MakeRack() primitive.Primitive {
	amount := 15
	spheres := make([]primitive.Primitive, amount * amount)

	diameter := 5.0
	distance := 10.0
	start := -float64(amount-1) / 2 * distance
	for i := 0; i < amount; i++ {
		for j := 0; j < amount; j++ {
			jitter := math.Sin(math.Sqrt(square(float64(i)) + square(float64(j)))) * 15

			sphere := primitive.NewSphere(diameter)
			translation := mgl64.Vec3{
				start + distance*float64(i), start + distance*float64(j), jitter,
			}
			translatedSphere := primitive.NewTranslation(translation, sphere)

			spheres[i * amount + j] = translatedSphere
		}
	}

	rotation := mgl64.Vec3{0, 0, 0}
	return primitive.NewRotation(rotation, primitive.NewList(spheres...))
}
