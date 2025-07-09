package rack

import (
	"fmt"

	"github.com/ljanyst/ghostscad/primitive"
)

const (
	screwRadiusM6 = 3.0

	rackSpineWidth         = 15.875
	rackSpineThickness     = 10.0
	rackSpineInlayWidth    = 2.0
)

type Rack struct {
	primitive.ParentImpl
	primitive.List
	Foot     *RackFoot
}

func MakeRack(heightUnits uint8) *Rack {
	rack := &Rack{}

	if heightUnits == 0 {
		return rack
	}

	var previousSegment *RackSegment

	for i := uint8(0); i < heightUnits; i++ {
		nextSegment := NewRackSegment(fmt.Sprintf("segment-%d", i))
		nextBrace := NewSideBrace(fmt.Sprintf("sidebrace-%d", i), heightUnits, i)

		if previousSegment != nil {
			if err := previousSegment.Anchors()["bottom"].Connect(nextSegment.Anchors()["top"], 0); err != nil {
				panic("failed to connect rack segments. this should not happen")
			}
		}
		if err := nextSegment.Anchors()["left"].Connect(nextBrace.Anchors()["segmentattach"], 0); err != nil {
			panic("failed to attach side brace to rack segment")
		}

		previousSegment = nextSegment
		rack.Add(nextSegment)
		rack.Add(nextBrace)
	}

	foot := NewRackFoot("foot")
	if err := foot.Anchors()["top"].Connect(previousSegment.Anchors()["bottom"], 0); err != nil {
		panic("failed to connect rack segments. this should not happen")
	}
	rack.Foot = foot
	rack.Add(foot)

	return rack
}
