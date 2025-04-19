package shapes

import (
	"fmt"
	"testing"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/ljanyst/ghostscad/primitive"
	"github.com/stretchr/testify/assert"
	"github.com/yeldiRium/3d-rack-brackets/test"
)

type FooAnchored struct {
	primitive.ParentImpl
	primitive.Cube
	name string
	anchors         map[string]Anchor
	anchorTransform *primitive.Transform
}

func NewFoo(name string, edge float64) *FooAnchored {
	foo := &FooAnchored{
		name: name,
		Cube: *primitive.NewCube(mgl64.Vec3{edge, edge, edge}),
	}
	foo.anchors = map[string]Anchor{
		"top": NewAnchor("top", foo, *primitive.NewTranslation(mgl64.Vec3{0, 0, edge / 2}), mgl64.Vec3{0, 0, 1}),
		"bottom": NewAnchor("bottom", foo, *primitive.NewTranslation(mgl64.Vec3{0, 0, -edge / 2}), mgl64.Vec3{0, 0, -1}),
		"right":  NewAnchor("right", foo, *primitive.NewTranslation(mgl64.Vec3{edge / 2, 0, 0}), mgl64.Vec3{1, 0, 0}),
	}
	return foo
}

func (foo *FooAnchored) Anchors() map[string]Anchor {
	return foo.anchors
}

func (foo *FooAnchored) SetAnchorTransform(transform primitive.Transform) error {
	if foo.anchorTransform != nil {
		// TODO check if the preexisting anchorTransform might be identical to
		// transform. If so, don't return an error.
		return fmt.Errorf("trying to set conflicting anchor transforms")
	}

	foo.anchorTransform = &transform
	return nil
}

func (foo *FooAnchored) GetAnchorTransform() *primitive.Transform {
	return foo.anchorTransform
}

func TestResolveAnchors(t *testing.T) {
	t.Run("resolves the anchor transform from one foo to another correctly.", func(t *testing.T) {
		fooOne := NewFoo("fooOne", 7)
		fooTwo := NewFoo("fooTwo", 2)

		err := fooOne.Anchors()["bottom"].Connect(fooTwo.Anchors()["right"], 45)
		assert.NoError(t, err)

		err = ResolveAnchors(fooOne)
		assert.NoError(t, err)

		expectedTransformFooOne := primitive.NewTranslation(mgl64.Vec3{0, 0, 0})
		expectedTransformFooTwo := primitive.NewTranslation(mgl64.Vec3{})
		expectedTransformFooTwo.Append(primitive.NewTranslation(mgl64.Vec3{0, 0, -3.5}))
		expectedTransformFooTwo.Append(primitive.NewRotationByAxis(45, mgl64.Vec3{0, 0, -1}))
		expectedTransformFooTwo.Append(primitive.NewRotation(mgl64.Vec3{0, -90, 0}))
		expectedTransformFooTwo.Append(primitive.NewTranslation(mgl64.Vec3{1, 0, 0}))

		test.RemoveParent(fooOne.anchorTransform.Items)
		test.RemoveParent(fooTwo.anchorTransform.Items)
		test.RemoveParent(expectedTransformFooOne.Items)
		test.RemoveParent(expectedTransformFooTwo.Items)
		assert.Equal(t, expectedTransformFooOne, fooOne.anchorTransform)
		assert.Equal(t, expectedTransformFooTwo, fooTwo.anchorTransform)
	})

	t.Run("resolves a transform correctly for opposite normals", func(t *testing.T) {
		fooOne := NewFoo("fooOne", 7)
		fooTwo := NewFoo("fooTwo", 2)

		err := fooOne.Anchors()["bottom"].Connect(fooTwo.Anchors()["top"], 0)
		assert.NoError(t, err)

		err = ResolveAnchors(fooOne)
		assert.NoError(t, err)

		expectedTransformFooOne := primitive.NewTranslation(mgl64.Vec3{0, 0, 0})
		expectedTransformFooTwo := primitive.NewTranslation(mgl64.Vec3{})
		expectedTransformFooTwo.Append(primitive.NewTranslation(mgl64.Vec3{0, 0, -3.5}))
		expectedTransformFooTwo.Append(primitive.NewRotationByAxis(0, mgl64.Vec3{0, 0, -1}))
		expectedTransformFooTwo.Append(primitive.NewRotation(mgl64.Vec3{0, 180, 0})) // TODO: fix rotation
		expectedTransformFooTwo.Append(primitive.NewTranslation(mgl64.Vec3{0, 0, -1}))

		test.RemoveParent(fooOne.anchorTransform.Items)
		test.RemoveParent(fooTwo.anchorTransform.Items)
		test.RemoveParent(expectedTransformFooOne.Items)
		test.RemoveParent(expectedTransformFooTwo.Items)
		assert.Equal(t, expectedTransformFooOne, fooOne.anchorTransform)
		assert.Equal(t, expectedTransformFooTwo, fooTwo.anchorTransform)
	})
}
