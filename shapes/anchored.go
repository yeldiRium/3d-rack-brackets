package shapes

import (
	"errors"
	"math"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/ljanyst/ghostscad/primitive"
)

var (
	ErrAnchorAlreadyConnected = errors.New("anchor already has different connection")
)

type Anchor interface {
	Parent() Anchored
	Connect(target Anchor, rotation float64) error
	Connection() *anchorConnection
	Translation() primitive.Transform
	Normal() mgl64.Vec3
}

type anchor struct {
	// parent is the shape that the anchor is a part of. The anchor is used to
	// connect the parent to another primitive.
	parent Anchored

	// translation is a translation from the anchored parent's origin to the an-
	// chor. Thus applying the inverse transform to the anchor results in the an-
	// chored parent's origin.
	translation primitive.Transform

	// normal is the direction in which the anchor connects. When connecting to
	// another anchor, they will be rotated so that their normals are opposite
	// each other.
	normal mgl64.Vec3

	// connectedAnchor is a second anchor that is connected to this one. Each an-
	// chor can always be tied to at most one other anchor.
	connection *anchorConnection
}

func (anchor *anchor) Parent() Anchored {
	return anchor.parent
}

func (anchor *anchor) Connect(target Anchor, rotation float64) error {
	if target == nil {
		panic("target must not be nil")
	}

	inverseRotation := math.Mod(360-rotation, 360)
	if anchor.connection != nil &&
		anchor.connection.target == target &&
		anchor.connection.rotation == rotation &&
		target.Connection() != nil &&
		target.Connection().Target() == anchor &&
		target.Connection().Rotation() == inverseRotation {
		return nil
	}

	if anchor.connection != nil && (anchor.connection.target != target || anchor.connection.rotation != rotation) {
		return ErrAnchorAlreadyConnected
	}
	targetConnection := target.Connection()
	if targetConnection != nil {
		if targetConnection.Target() != anchor || targetConnection.Rotation() != inverseRotation {
			return ErrAnchorAlreadyConnected
		}
	}

	anchor.connection = &anchorConnection{
		target:   target,
		rotation: rotation,
	}
	target.Connect(anchor, inverseRotation)

	return nil
}

func (anchor *anchor) Connection() *anchorConnection {
	return anchor.connection
}

func (anchor *anchor) Translation() primitive.Transform {
	return anchor.translation
}

func (anchor *anchor) Normal() mgl64.Vec3 {
	return anchor.normal
}

type AnchorConnection interface {
	Target() Anchor
	Rotation() float64
}

type anchorConnection struct {
	target   Anchor
	rotation float64
}

func (c *anchorConnection) Target() Anchor {
	return c.target
}

func (c *anchorConnection) Rotation() float64 {
	return c.rotation
}

func NewAnchor(parent Anchored, translation primitive.Transform, normal mgl64.Vec3) *anchor {
	anchor := &anchor{}
	anchor.parent = parent
	anchor.translation = translation
	anchor.normal = normal

	return anchor
}

type Anchored interface {
	Anchors() map[string]Anchor
	SetAnchorTransform(t primitive.Transform) error
}

func ResolveAnchors(start Anchored) error {
	// TODO: implement graph algo that traverses anchors, resolves tranforms, de-
	// tects obvious collections and writes the resolved transforms into the an-
	// choreds.
	return nil
}
