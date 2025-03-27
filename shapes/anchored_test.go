package shapes

import (
	"testing"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/ljanyst/ghostscad/primitive"
	"github.com/stretchr/testify/assert"
)

type FooAnchored struct {
	primitive.ParentImpl
	primitive.Cube
	anchors map[string]Anchor
	anchorTransform *primitive.Transform
}

func NewFoo(edge float64) *FooAnchored {
	foo := &FooAnchored{}
	foo.Cube = *primitive.NewCube(mgl64.Vec3{edge, edge, edge})
	foo.anchors = map[string]Anchor{
		"bottom": NewAnchor(foo, *primitive.NewTranslation(mgl64.Vec3{0, 0, -edge / 2}), mgl64.Vec3{0, 0, -1}),
		"right": NewAnchor(foo, *primitive.NewTranslation(mgl64.Vec3{edge / 2, 0, 0}), mgl64.Vec3{1, 0, 0}),
	}
	return foo
}

func (foo *FooAnchored) Anchors() map[string]Anchor {
	return foo.anchors
}

func (foo *FooAnchored) SetAnchorTransform(t primitive.Transform) error {
	panic("unimplemented")
}

func TestResolveAnchors(t *testing.T) {
	t.Run("resolves the anchor transform from one foo to another correctly.", func(t *testing.T) {
		fooOne := NewFoo(7)
		fooTwo := NewFoo(2)

		err := fooOne.Anchors()["bottom"].Connect(fooTwo.Anchors()["right"], 45)
		assert.NoError(t, err)

		ResolveAnchors(fooOne)

		expectedTransformFooOne := primitive.NewTranslation(mgl64.Vec3{0, 0, 0})
		expectedTransformFooOne.Append(primitive.NewRotation(mgl64.Vec3{0, 0, 0}))
		expectedTransformFooTwo := primitive.NewTranslation(mgl64.Vec3{4.5, 0, 0})
		expectedTransformFooTwo.Append(primitive.NewRotation(mgl64.Vec3{0, 0, 0})) // TODO: fix rotation

		assert.Equal(t, expectedTransformFooOne, fooOne.anchorTransform)
		assert.Equal(t, expectedTransformFooTwo, fooTwo.anchorTransform)
	})
}
