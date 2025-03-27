package rack

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/ljanyst/ghostscad/primitive"
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
	primitive.List
}

func NewRackSegment() *RackSegment {
	spine := primitive.NewCube(mgl64.Vec3{RACK_SPINE_WIDTH, RACK_SPINE_THICKNESS, RACK_SEGMENT_HEIGHT})
	cutout := primitive.NewCylinder(RACK_SPINE_THICKNESS+1, SCREW_RADIUS_M6)
	orientedCutout := primitive.NewRotation(mgl64.Vec3{90, 0, 0}, cutout)

	firstCutout := primitive.NewTranslation(mgl64.Vec3{0, 0, (RACK_SEGMENT_HEIGHT / 2) - RACK_SEGMENT_HOLE_SPACING}, orientedCutout)
	secondCutout := orientedCutout
	thirdCutout := primitive.NewTranslation(mgl64.Vec3{0, 0, -(RACK_SEGMENT_HEIGHT / 2) + RACK_SEGMENT_HOLE_SPACING}, orientedCutout)

	spineWithCutouts := primitive.NewDifference(spine, firstCutout, secondCutout, thirdCutout)

	rackSegment := &RackSegment{}
	rackSegment.Add(spineWithCutouts)

	return rackSegment
}

type Rack struct {
	primitive.ParentImpl
	primitive.List
}

func MakeRack(heightUnits uint8) *Rack {
	rack := &Rack{}

	// Since the origin point of each segment is it's center, we need to calculate
	// height as the distance between the top and the bottom center. Otherwise
	// the translations don't fit or the calculation below has to be more complic-
	// ated.
	finalHeight := RACK_SEGMENT_HEIGHT * float64(heightUnits - 1)
	zStart := -finalHeight / 2

	for i := uint8(0); i < heightUnits; i++ {
		rack.Add(primitive.NewTranslation(
			mgl64.Vec3{0, 0, zStart + float64(i) * RACK_SEGMENT_HEIGHT},
			NewRackSegment(),
		))
	}

	return rack
}
