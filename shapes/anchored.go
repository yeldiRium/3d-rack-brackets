package shapes

import (
	"errors"
	"math"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/ljanyst/ghostscad/primitive"
	"github.com/yeldiRium/3d-rack-brackets/ghostscad"
)

var (
	ErrAnchorAlreadyConnected = errors.New("anchor already has different connection")
)

type Anchor interface {
	Name() string
	Parent() Anchored
	Connect(target Anchor, angle float64) error
	Connection() *anchorConnection
	Translation() *primitive.Transform
	Normal() mgl64.Vec3
}

type anchor struct {
	// name identifies the anchor and is mainly used for debugging
	name string

	// parent is the shape that the anchor is a part of. The anchor is used to
	// connect the parent to another primitive.
	parent Anchored

	// translation is a translation from the anchored parent's origin to the an-
	// chor. Thus applying the inverse transform to the anchor results in the an-
	// chored parent's origin.
	translation *primitive.Transform

	// normal is the direction in which the anchor connects. When connecting to
	// another anchor, they will be rotated so that their normals are opposite
	// each other.
	normal mgl64.Vec3

	// connectedAnchor is a second anchor that is connected to this one. Each an-
	// chor can always be tied to at most one other anchor.
	connection *anchorConnection
}

func NewAnchor(name string, parent Anchored, translation *primitive.Transform, normal mgl64.Vec3) *anchor {
	return &anchor{
		name:        name,
		parent:      parent,
		translation: translation,
		normal:      normal,
	}
}

func (anchor *anchor) Name() string {
	return anchor.name
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

func (anchor *anchor) Translation() *primitive.Transform {
	return anchor.translation
}

func (anchor *anchor) Normal() mgl64.Vec3 {
	return anchor.normal
}

type AnchorConnection interface {
	Target() Anchor
	Angle() float64
	WasResolved() bool
}

type anchorConnection struct {
	target      Anchor
	angle       float64
	wasResolved bool
}

func (c *anchorConnection) Target() Anchor {
	return c.target
}

func (c *anchorConnection) Angle() float64 {
	return c.angle
}

func (c *anchorConnection) WasResolved() bool {
	return c.wasResolved
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
			if connection == nil || connection.wasResolved {
				continue
			}

			targetAnchor := connection.Target()
			angle := connection.Angle()
			targetAnchored := targetAnchor.Parent()

			matchAnchorOrientationRotation := calculateRotationFromVec3ToVec3(anchor.Normal(), targetAnchor.Normal())

			matchAnchorOrientation := primitive.NewRotation(matchAnchorOrientationRotation)
			rotateAroundConnection := primitive.NewRotationByAxis(angle, anchor.Normal())
			moveByStartAnchor := anchor.Translation()
			moveByTargetAnchor := targetAnchor.Translation().Inverse()

			targetTransformation := ghostscad.CloneTransform(currentTransform)
			targetTransformation.Append(moveByStartAnchor)
			targetTransformation.Append(rotateAroundConnection) 
			targetTransformation.Append(matchAnchorOrientation)
			targetTransformation.Append(moveByTargetAnchor)

			err := targetAnchored.SetAnchorTransform(*targetTransformation)
			if err != nil {
				return err
			}
			connection.wasResolved = true
			targetAnchor.Connection().wasResolved = true

			anchoredQueue = append(anchoredQueue, targetAnchored)
		}
	}

	return nil
}

// TODO: fix rotation for parallel vectors
func calculateRotationFromVec3ToVec3(from, to mgl64.Vec3) mgl64.Vec3 {
	axis := from.Cross(to).Normalize()
	if math.IsNaN(axis[0]) {
		// If any of the axis' values are NaN, from and to are parallel and we can
		// choose any orthogonal
		axis = findOrthogonal(from)

		// TODO: rotate by 180deg or 0deg, dependening on orientation of from and to
	}
	angle := math.Acos(from.Normalize().Dot(to.Normalize()))

	rotationMatrix := mgl64.HomogRotate3D(angle, axis)
	eulerAngles := eulerAngles(rotationMatrix)

	return mgl64.Vec3{
		mgl64.RadToDeg(eulerAngles[0]),
		mgl64.RadToDeg(eulerAngles[1]),
		mgl64.RadToDeg(eulerAngles[2]),
	}
}

// to find an orthogonal u for v, we need to set u dot v = 0
// u1 * v1 + u2 * v2 + u3 * v3 = 0
// setting u1 = 1 and u2 = 1 we get
// v1 + v2 + u3 * v3 = 0
// u3 = (-v1 - v2) / v3
func findOrthogonal(v mgl64.Vec3) mgl64.Vec3 {
	z := (-v[0] - v[1]) / v[2]
	return mgl64.Vec3{1, 1, z}
}

// eulerAngles takes a radians rotation matrix and calculates its radians euler angles
func eulerAngles(rotationMatrix mgl64.Mat4) mgl64.Vec3 {
	sy := math.Sqrt(rotationMatrix[0]*rotationMatrix[0] + rotationMatrix[1]*rotationMatrix[1])

	singular := sy < 1e-6

	var x, y, z float64
	if !singular {
		x = math.Atan2(rotationMatrix[6], rotationMatrix[10])
		y = math.Atan2(-rotationMatrix[2], sy)
		z = math.Atan2(rotationMatrix[1], rotationMatrix[0])
	} else {
		x = math.Atan2(-rotationMatrix[9], rotationMatrix[5])
		y = math.Atan2(-rotationMatrix[2], sy)
		z = 0
	}

	return mgl64.Vec3{x, y, z}
}
