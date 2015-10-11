package utils

import (
	"github.com/pyros2097/spike/math/shape"
)

// Allows a parent to set the area that is visible on a child actor to allow the child to cull when drawing itself. This must only
// be used for actors that are not rotated or scaled.
// When Group is given a culling rectangle with {@link Group#setCullingArea(Rectangle)}, it will automatically call
// {@link #setCullingArea(Rectangle)} on its children.
type Cullable interface {
	// param cullingArea The culling area in the child actor's coordinates.
	SetCullingArea(cullingArea *shape.Rectangle)
}
