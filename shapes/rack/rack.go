package rack

import "github.com/ljanyst/ghostscad/primitive"

func MakeRack() primitive.Primitive {
	return primitive.NewSphere(15)
}
