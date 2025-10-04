package rack

import (
	"bufio"
	"fmt"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/ljanyst/ghostscad/primitive"

	"github.com/yeldiRium/3d-rack-brackets/internal/shapes"
)

const (
	rackSegmentHeight      = 44.45
	rackSegmentHoleSpacing = 6.35
)

type RackSegment struct {
	primitive.ParentImpl

	prefix string

	name     string
	contents *primitive.List

	anchors         map[string]shapes.Anchor
	anchorTransform *primitive.Transform
}

func NewRackSegment(name string) *RackSegment {
	spine := primitive.NewCube(mgl64.Vec3{rackSpineWidth, rackSpineThickness, rackSegmentHeight})
	cutout := primitive.NewCylinder(rackSpineThickness+1, screwRadiusM6)
	orientedCutout := primitive.NewRotation(mgl64.Vec3{90, 0, 0}, cutout)

	firstCutout := primitive.NewTranslation(mgl64.Vec3{0, 0, (rackSegmentHeight / 2) - rackSegmentHoleSpacing}, orientedCutout)
	secondCutout := orientedCutout
	thirdCutout := primitive.NewTranslation(mgl64.Vec3{0, 0, -(rackSegmentHeight / 2) + rackSegmentHoleSpacing}, orientedCutout)

	spineWithCutouts := primitive.NewDifference(spine, firstCutout, secondCutout, thirdCutout)

	rackSegment := &RackSegment{
		name:     name,
		contents: primitive.NewList(),
	}
	rackSegment.contents.Add(spineWithCutouts)
	rackSegment.anchors = map[string]shapes.Anchor{
		"top":    shapes.NewAnchor("top", rackSegment, primitive.NewTranslation(mgl64.Vec3{0, 0, rackSegmentHeight / 2}), mgl64.Vec3{0, 0, 1}),
		"left":   shapes.NewAnchor("left", rackSegment, primitive.NewTranslation(mgl64.Vec3{rackSpineWidth / 2, 0, 0}), mgl64.Vec3{1, 0, 0}),
		"bottom": shapes.NewAnchor("bottom", rackSegment, primitive.NewTranslation(mgl64.Vec3{0, 0, -rackSegmentHeight / 2}), mgl64.Vec3{0, 0, -1}),
	}

	return rackSegment
}

func (rackSegment *RackSegment) Anchors() map[string]shapes.Anchor {
	return rackSegment.anchors
}

func (rackSegment *RackSegment) SetAnchorTransform(transform primitive.Transform) error {
	if rackSegment.anchorTransform != nil {
		// TODO check if the preexisting anchorTransform might be identical to
		// transform. If so, don't return an error.
		return fmt.Errorf("trying to set conflicting anchor transforms")
	}

	rackSegment.anchorTransform = &transform

	return nil
}

func (rackSegment *RackSegment) GetAnchorTransform() *primitive.Transform {
	return rackSegment.anchorTransform
}

func (rackSegment *RackSegment) Disable() primitive.Primitive { //nolint:ireturn
	rackSegment.prefix = "*"

	return rackSegment
}

func (rackSegment *RackSegment) ShowOnly() primitive.Primitive { //nolint:ireturn
	rackSegment.prefix = "!"

	return rackSegment
}

func (rackSegment *RackSegment) Highlight() primitive.Primitive { //nolint:ireturn
	rackSegment.prefix = "#"

	return rackSegment
}

func (rackSegment *RackSegment) Transparent() primitive.Primitive { //nolint:ireturn
	rackSegment.prefix = "%"

	return rackSegment
}

func (rackSegment *RackSegment) Prefix() string {
	return rackSegment.prefix
}

func (rackSegment *RackSegment) Render(w *bufio.Writer) {
	if rackSegment.anchorTransform == nil {
		panic("cannot render racksegment without resolving its anchors")
	}
	rackSegment.anchorTransform.Add(rackSegment.contents)
	rackSegment.anchorTransform.Render(w)
}
