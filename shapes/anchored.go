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
	Connect(target Anchor, angle float64) error
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

func (anchor *anchor) Connect(target Anchor, angle float64) error {
	if target == nil {
		panic("target must not be nil")
	}

	inverseAngle := math.Mod(360-angle, 360)
	if anchor.connection != nil &&
		anchor.connection.target == target &&
		anchor.connection.angle == angle &&
		target.Connection() != nil &&
		target.Connection().Target() == anchor &&
		target.Connection().Angle() == inverseAngle {
		return nil
	}

	if anchor.connection != nil && (anchor.connection.target != target || anchor.connection.angle != angle) {
		return ErrAnchorAlreadyConnected
	}
	targetConnection := target.Connection()
	if targetConnection != nil {
		if targetConnection.Target() != anchor || targetConnection.Angle() != inverseAngle {
			return ErrAnchorAlreadyConnected
		}
	}

	anchor.connection = &anchorConnection{
		target: target,
		angle:  angle,
	}
	target.Connect(anchor, inverseAngle)

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
	Angle() float64
}

type anchorConnection struct {
	target Anchor
	angle  float64
}

func (c *anchorConnection) Target() Anchor {
	return c.target
}

func (c *anchorConnection) Angle() float64 {
	return c.angle
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
	GetAnchorTransform() *primitive.Transform
}

func ResolveAnchors(start Anchored) error {
	start.SetAnchorTransform(*primitive.NewTranslation(mgl64.Vec3{}))
	anchoredQueue := []Anchored{start}

	processedAnchoreds := map[Anchored]bool{}

	for len(anchoredQueue) > 0 {
		currentAnchored := anchoredQueue[0]
		anchoredQueue = anchoredQueue[1:]

		if processedAnchoreds[currentAnchored] {
			continue
		}

		processedAnchoreds[currentAnchored] = true

		currentTransform := currentAnchored.GetAnchorTransform()
		if currentTransform == nil {
			panic("an already processed anchored unexpectedly has no transform. this should never happen")
		}

		for _, anchor := range currentAnchored.Anchors() {
			connection := anchor.Connection()
			if connection == nil {
				continue
			}
			targetAnchor := connection.Target()
			angle := connection.Angle()
			targetAnchored := targetAnchor.Parent()

			matchAnchorOrientationRotation := calculateRotationFromVec3ToVec3(anchor.Normal(), targetAnchor.Normal())

			matchAnchorOrientation := primitive.NewRotation(matchAnchorOrientationRotation)
			rotateAroundConnection := primitive.NewRotationByAxis(angle, anchor.Normal())
			moveByStartAnchor := anchor.Translation()
			moveByTargetAnchor := targetAnchor.Translation()

			targetTransformation := matchAnchorOrientation
			targetTransformation.Append(rotateAroundConnection)
			targetTransformation.Append(&moveByStartAnchor)
			targetTransformation.Append(&moveByTargetAnchor)

			err := targetAnchored.SetAnchorTransform(*targetTransformation)
			if err != nil {
				return err
			}

			anchoredQueue = append(anchoredQueue, targetAnchored)
		}
	}

	return nil
}

func calculateRotationFromVec3ToVec3(from, to mgl64.Vec3) mgl64.Vec3 {
	// TODO implement this. maybe using mat3s is better than vec3s? hm.
}
