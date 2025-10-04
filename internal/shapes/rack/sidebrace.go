package rack

import (
	"bufio"
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/ljanyst/ghostscad/primitive"

	"github.com/yeldiRium/3d-rack-brackets/internal/shapes"
)

const (
	sideBracePadding         = 10
	sideBraceInnerPadding    = 2
	sideBraceAttachmentDepth = 20
	sideBraceWidth           = 3.0
)

type SideBrace struct {
	primitive.ParentImpl
	prefix string

	name     string
	contents *primitive.List

	anchors         map[string]shapes.Anchor
	anchorTransform *primitive.Transform
}

// NewSideBrace constructs a side brace.
// heightUnit is the number of the segment the brace belongs to.
//
//	0 is the lowest brace.
//	This is used to calculate the connection point to the rack foot.
func NewSideBrace(name string, totalHeight, heightUnit uint8) *SideBrace {
	footOffsetY := float64(totalHeight-heightUnit-1)*rackSegmentHeight + rackFootSpacerHeight
	footOffsetZ := math.Sqrt(float64(totalHeight-heightUnit)/float64(totalHeight)) * rackFootLength * 2 / 3

	scaledAttachmentDepth := sideBraceAttachmentDepth * math.Pow(0.8, float64(totalHeight-heightUnit-1))

	shape := primitive.NewPolygon([]mgl64.Vec2{
		{0, 0},
		{rackSegmentHeight, 0},
		{rackSegmentHeight, rackSpineThickness},
		{rackSegmentHeight - sideBracePadding, rackSpineThickness},
		{rackSegmentHeight + footOffsetY, footOffsetZ - scaledAttachmentDepth},
		{rackSegmentHeight + footOffsetY, footOffsetZ},
		{sideBracePadding, rackSpineThickness},
		{0, rackSpineThickness},
	})

	cutoutDepth := math.Pow(float64(totalHeight-heightUnit)/float64(totalHeight), 1.05) * rackFootLength / 3
	cutout := primitive.NewDifference(primitive.NewPolygon([]mgl64.Vec2{
		{rackSegmentHeight/2 - sideBraceInnerPadding, rackSpineThickness},
		{rackSegmentHeight + footOffsetY, footOffsetZ - scaledAttachmentDepth/2},
		{rackSegmentHeight/2 + sideBraceInnerPadding, rackSpineThickness},
	}), primitive.NewPolygon([]mgl64.Vec2{
		{0, rackSpineThickness + cutoutDepth},
		{0, rackFootLength},
		{rackSegmentHeight + footOffsetY, rackFootLength},
		{rackSegmentHeight + footOffsetY, rackSpineThickness + cutoutDepth},
	}))

	finalShape := primitive.NewDifference(shape, cutout)

	extrusion := primitive.NewLinearExtrusion(
		sideBraceWidth,
		finalShape,
	)

	sideBrace := &SideBrace{
		name:     name,
		contents: primitive.NewList(),
	}
	sideBrace.contents.Add(extrusion)
	sideBrace.anchors = map[string]shapes.Anchor{
		"segmentattach": shapes.NewAnchor(
			"segmentattach",
			sideBrace,
			primitive.NewTranslation(mgl64.Vec3{
				rackSegmentHeight / 2,
				rackSpineThickness / 2,
				-sideBraceWidth / 2,
			}),
			mgl64.Vec3{0, 0, 1},
		),
	}

	return sideBrace
}

func (sideBrace *SideBrace) Anchors() map[string]shapes.Anchor {
	return sideBrace.anchors
}

func (sideBrace *SideBrace) SetAnchorTransform(transform primitive.Transform) error {
	if sideBrace.anchorTransform != nil {
		// TODO check if the preexisting anchorTransform might be identical to
		// transform. If so, don't return an error.
		return fmt.Errorf("trying to set conflicting anchor transforms")
	}

	sideBrace.anchorTransform = &transform
	return nil
}

func (sideBrace *SideBrace) GetAnchorTransform() *primitive.Transform {
	return sideBrace.anchorTransform
}

func (sideBrace *SideBrace) Disable() primitive.Primitive {
	sideBrace.prefix = "*"
	return sideBrace
}

func (sideBrace *SideBrace) ShowOnly() primitive.Primitive {
	sideBrace.prefix = "!"
	return sideBrace
}

func (sideBrace *SideBrace) Highlight() primitive.Primitive {
	sideBrace.prefix = "#"
	return sideBrace
}

func (sideBrace *SideBrace) Transparent() primitive.Primitive {
	sideBrace.prefix = "%"
	return sideBrace
}

func (sideBrace *SideBrace) Prefix() string {
	return sideBrace.prefix
}

func (sideBrace *SideBrace) Render(w *bufio.Writer) {
	if sideBrace.anchorTransform == nil {
		panic("cannot render side brace without resolving its anchors")
	}
	sideBrace.anchorTransform.Add(sideBrace.contents)
	sideBrace.anchorTransform.Render(w)
}
