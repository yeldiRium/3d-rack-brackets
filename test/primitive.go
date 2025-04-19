package test

import "github.com/ljanyst/ghostscad/primitive"

func RemoveParent(primitive primitive.Primitive) {
	primitive.SetParent(nil)
}
