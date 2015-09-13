// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package scene2d

/** Determines how touch input events are distributed to an actor and any children.
 * @author Nathan Sweet */
type Touchable int

const (
	/** All touch input events will be received by the actor and any children. */
	TouchableEnabled Touchable = iota

	/** No touch input events will be received by the actor or any children. */
	TouchableDisabled

	/** No touch input events will be received by the actor, but children will still receive events. Note that events on the
	 * children will still bubble to the parent. */
	TouchableChildrenOnly
)
