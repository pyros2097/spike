// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package utils

// Interface for disposable resources.
type Disposable interface {
	// Releases all resources of this object.
	Dispose()
}
