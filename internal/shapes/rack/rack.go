package rack

import (
	"fmt"

	"github.com/ljanyst/ghostscad/primitive"
)

const (
	screwRadiusM6 = 3.0

	rackSpineWidth          = 15.875
	rackSpineThickness      = 10.0
	rackSegmentHeight       = 44.45
	rackSegmentHoleSpacing = 6.35
)

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
			if err := previousSegment.Anchors()["bottom"].Connect(nextSegment.Anchors()["top"], 0); err != nil {
				panic("failed to connect rack segments. this should not happen")
			}
		}

		previousSegment = nextSegment
		rack.Add(nextSegment)
		rack.Segments = append(rack.Segments, nextSegment)
	}

	return rack
}
