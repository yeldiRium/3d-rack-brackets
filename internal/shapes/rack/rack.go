package rack

import (
	"fmt"

	"github.com/ljanyst/ghostscad/primitive"
)

const (
	SCREW_RADIUS_M6 = 3.0

	RACK_SPINE_WIDTH          = 15.875
	RACK_SPINE_THICKNESS      = 10.0
	RACK_SEGMENT_HEIGHT       = 44.45
	RACK_SEGMENT_HOLE_SPACING = 6.35
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
