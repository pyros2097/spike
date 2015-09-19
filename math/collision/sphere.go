// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package shape

// Encapsulates a 3D sphere with a center and a radius
type Sphere struct {
	radius float32
	center *Vector3
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
	return self.center.Dst2(sphere.center) < (self.radius+sphere.radius)*(self.radius+sphere.radius)
}

func (self *Sphere) Volume() float32 {
	return PI_4_3 * self.radius * self.radius * self.radius
}

func (self *Sphere) SurfaceArea() float32 {
	return 4 * PI * self.radius * self.radius
}

// public int hashCode () {
// 	final int prime = 71;
// 	int result = 1;
// 	result = prime * result + self.center.hashCode();
// 	result = prime * result + NumberUtils.floatToRawIntBits(self.radius);
// 	return result;
// }

// public boolean equals (Object o) {
// 	if (this == o) return true;
// 	if (o == null || o.getClass() != self.getClass()) return false;
// 	Sphere s = (Sphere)o;
// 	return self.radius == s.radius && self.center.equals(s.center);
// }
