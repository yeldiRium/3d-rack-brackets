package rack

import (
	"bufio"
	"fmt"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/ljanyst/ghostscad/primitive"

	"github.com/yeldiRium/3d-rack-brackets/shapes"
)

const (
	SCREW_RADIUS_M6 = 3.0

	RACK_SPINE_WIDTH          = 15.875
	RACK_SPINE_THICKNESS      = 10.0
	RACK_SEGMENT_HEIGHT       = 44.45
	RACK_SEGMENT_HOLE_SPACING = 6.35
)

type RackSegment struct {
	primitive.ParentImpl
	prefix string

	name string
	contents *primitive.List

	anchors         map[string]shapes.Anchor
	anchorTransform *primitive.Transform
}

func NewRackSegment(name string) *RackSegment {
	spine := primitive.NewCube(mgl64.Vec3{RACK_SPINE_WIDTH, RACK_SPINE_THICKNESS, RACK_SEGMENT_HEIGHT})
	cutout := primitive.NewCylinder(RACK_SPINE_THICKNESS+1, SCREW_RADIUS_M6)
	orientedCutout := primitive.NewRotation(mgl64.Vec3{90, 0, 0}, cutout)

	firstCutout := primitive.NewTranslation(mgl64.Vec3{0, 0, (RACK_SEGMENT_HEIGHT / 2) - RACK_SEGMENT_HOLE_SPACING}, orientedCutout)
	secondCutout := orientedCutout
	thirdCutout := primitive.NewTranslation(mgl64.Vec3{0, 0, -(RACK_SEGMENT_HEIGHT / 2) + RACK_SEGMENT_HOLE_SPACING}, orientedCutout)

	spineWithCutouts := primitive.NewDifference(spine, firstCutout, secondCutout, thirdCutout)

	rackSegment := &RackSegment{
		name: name,
		contents: primitive.NewList(),
	}
	rackSegment.contents.Add(spineWithCutouts)
	rackSegment.anchors = map[string]shapes.Anchor{
		"top": shapes.NewAnchor("top", rackSegment, primitive.NewTranslation(mgl64.Vec3{0, 0, RACK_SEGMENT_HEIGHT / 2}), mgl64.Vec3{0, 0, 1}),
		"bottom":  shapes.NewAnchor("bottom", rackSegment, primitive.NewTranslation(mgl64.Vec3{0, 0, -RACK_SEGMENT_HEIGHT / 2}), mgl64.Vec3{0, 0, -1}),
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

func (rackSegment *RackSegment) Disable() primitive.Primitive {
	rackSegment.prefix = "*"
	return rackSegment
}

func (rackSegment *RackSegment) ShowOnly() primitive.Primitive {
	rackSegment.prefix = "!"
	return rackSegment
}

func (rackSegment *RackSegment) Highlight() primitive.Primitive {
	rackSegment.prefix = "#"
	return rackSegment
}

func (rackSegment *RackSegment) Transparent() primitive.Primitive {
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

type Rack struct {
	primitive.ParentImpl
	primitive.List
	Segments []*RackSegment
}

func MakeRack(heightUnits uint8) *Rack {
	rack := &Rack{
		Segments: make([]*RackSegment, 0, heightUnits),
	}

	if heightUnits == 0 {
		return rack
	}

	var previousSegment *RackSegment

	for i := uint8(0); i < heightUnits; i++ {
		nextSegment := NewRackSegment(fmt.Sprintf("segment-%d", i))

		if previousSegment != nil {
			previousSegment.Anchors()["bottom"].Connect(nextSegment.Anchors()["top"], 0)
		}

		previousSegment = nextSegment
		rack.Add(nextSegment)
		rack.Segments = append(rack.Segments, nextSegment)
	}

	return rack
}
