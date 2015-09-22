// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package collision

import (
	"github.com/pyros2097/spike/math/utils"
	. "github.com/pyros2097/spike/math/vector"
)

// Encapsulates a 3D sphere with a center and a radius
type Sphere struct {
	Radius float32
	Center *Vector3
}

// Constructs a sphere with the given center and radius
// param center The center
// param radius The radius
func NewSphere(center *Vector3, radius float32) *Sphere {
	return &Sphere{radius, NewVector3Copy(center)}
}

// @param sphere the other sphere
// return whether this and the other sphere overlap
func (self *Sphere) Overlaps(sphere *Sphere) bool {
	return self.Center.Dst2V(sphere.Center) < (self.Radius+sphere.Radius)*(self.Radius+sphere.Radius)
}

func (self *Sphere) Volume() float32 {
	return utils.PI_4_3 * self.Radius * self.Radius * self.Radius
}

func (self *Sphere) SurfaceArea() float32 {
	return 4 * utils.PI * self.Radius * self.Radius
}

// public int hashCode () {
// 	final int prime = 71;
// 	int result = 1;
// 	result = prime * result + self.Center.hashCode();
// 	result = prime * result + NumberUtils.floatToRawIntBits(self.Radius);
// 	return result;
// }

// public boolean equals (Object o) {
// 	if (this == o) return true;
// 	if (o == null || o.getClass() != self.getClass()) return false;
// 	Sphere s = (Sphere)o;
// 	return self.Radius == s.radius && self.Center.equals(s.center);
// }
