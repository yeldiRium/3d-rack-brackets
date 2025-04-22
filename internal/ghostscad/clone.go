package ghostscad

import "github.com/ljanyst/ghostscad/primitive"

func CloneTransform(transform *primitive.Transform) *primitive.Transform {
	return transform.Inverse().Inverse()
}
