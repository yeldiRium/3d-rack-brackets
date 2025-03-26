package rack

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/ljanyst/ghostscad/primitive"
)

const (
	SCREW_RADIUS_M6 = 3.0
)

func MakeRackSegment() primitive.Primitive {
	thickness := 10.0
	rackSegmentWidth := 15.875
	rackSegmentHeight := 44.45
	distanceEdgeToHole := 6.35

	spine := primitive.NewCube(mgl64.Vec3{rackSegmentWidth, thickness, rackSegmentHeight})
	cutout := primitive.NewCylinder(thickness + 1, SCREW_RADIUS_M6)
	orientedCutout := primitive.NewRotation(mgl64.Vec3{90, 0, 0}, cutout)

	firstCutout := primitive.NewTranslation(mgl64.Vec3{0, 0, (rackSegmentHeight / 2) - distanceEdgeToHole}, orientedCutout)
	secondCutout := orientedCutout
	thirdCutout := primitive.NewTranslation(mgl64.Vec3{0, 0, - (rackSegmentHeight / 2) + distanceEdgeToHole}, orientedCutout)

	rackSegment := primitive.NewDifference(spine, firstCutout, secondCutout, thirdCutout)

	orientedRackSegment := primitive.NewRotation(mgl64.Vec3{-90, 0, 0}, rackSegment)

	return orientedRackSegment
}
