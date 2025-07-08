package rack

import (
	"bufio"
	"fmt"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/ljanyst/ghostscad/primitive"

	"github.com/yeldiRium/3d-rack-brackets/internal/shapes"
)

const (
	rackFootLength         = 170 + rackSpineInlayWidth
	rackFootThicknessFront = 15
	rackFootThicknessBack  = 10
	rackFootWidth          = rackSpineWidth
)

type RackFoot struct {
	primitive.ParentImpl
	prefix string

	name     string
	contents *primitive.List

	anchors         map[string]shapes.Anchor
	anchorTransform *primitive.Transform
}

func NewRackFoot(name string) *RackFoot {
	footBox := primitive.NewRotation(
		mgl64.Vec3{0, 90, 0},
		primitive.NewLinearExtrusion(
			rackFootWidth,
			primitive.NewPolygon([]mgl64.Vec2{
				{0, 0},
				{rackFootThicknessFront, 0},
				{rackFootThicknessBack, rackFootLength},
				{0, rackFootLength},
			}),
		),
	)

	rackFoot := &RackFoot{
		name:     name,
		contents: primitive.NewList(),
	}
	rackFoot.contents.Add(footBox)
	rackFoot.anchors = map[string]shapes.Anchor{
		"top": shapes.NewAnchor(
			"top",
			rackFoot,
			primitive.NewTranslation(mgl64.Vec3{
				0,
				(rackSpineThickness / 2) + rackSpineInlayWidth,
				0,
			}),
			mgl64.Vec3{0, 0, 1},
		),
	}

	return rackFoot
}

func (foot *RackFoot) Anchors() map[string]shapes.Anchor {
	return foot.anchors
}

func (foot *RackFoot) SetAnchorTransform(transform primitive.Transform) error {
	if foot.anchorTransform != nil {
		// TODO check if the preexisting anchorTransform might be identical to
		// transform. If so, don't return an error.
		return fmt.Errorf("trying to set conflicting anchor transforms")
	}

	foot.anchorTransform = &transform
	return nil
}

func (foot *RackFoot) GetAnchorTransform() *primitive.Transform {
	return foot.anchorTransform
}

func (foot *RackFoot) Disable() primitive.Primitive {
	foot.prefix = "*"
	return foot
}

func (foot *RackFoot) ShowOnly() primitive.Primitive {
	foot.prefix = "!"
	return foot
}

func (foot *RackFoot) Highlight() primitive.Primitive {
	foot.prefix = "#"
	return foot
}

func (foot *RackFoot) Transparent() primitive.Primitive {
	foot.prefix = "%"
	return foot
}

func (foot *RackFoot) Prefix() string {
	return foot.prefix
}

func (foot *RackFoot) Render(w *bufio.Writer) {
	if foot.anchorTransform == nil {
		panic("cannot render foot without resolving its anchors")
	}
	foot.anchorTransform.Add(foot.contents)
	foot.anchorTransform.Render(w)
}
