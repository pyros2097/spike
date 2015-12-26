package math

// This class represents a cumulative distribution.
// It can be used in scenarios where there are values with different probabilities
// and it's required to pick one of those respecting the probability.
// For example one could represent the frequency of the alphabet letters using a cumulative distribution
// and use it to randomly pick a letter respecting their probabilities (useful when generating random words).
// Another example could be point generation on a mesh surface: one could generate a cumulative distribution using
// triangles areas as interval size, in this way triangles with a large area will be picked more often than triangles with a smaller one.
// See http://en.wikipedia.org/wiki/Cumulative_distribution_function for a detailed explanation.
type CumulativeDistribution struct {
	values []*CumulativeValue
}

type CumulativeValue struct {
	value     interface{}
	frequency float32
	interval  float32
}

func NewCumulativeDistribution() *CumulativeDistribution {
	return &CumulativeDistribution{values: make([]*CumulativeValue, 10)}
}

// Adds a value with a given interval size to the distribution
func (self *CumulativeDistribution) Add(value interface{}, intervalSize float32) {
	// append(self.values, &CumulativeValue(value, 0, intervalSize))
}

// Adds a value with interval size equal to zero to the distributio
func (self *CumulativeDistribution) AddValue(value interface{}) {
	// append(self.values, &CumulativeValue(value, 0, 0))
}

// Generate the cumulative distribution
func (self *CumulativeDistribution) Generate() {
	var sum float32
	for i := 0; i < len(self.values); i++ {
		sum += self.values[i].interval
		self.values[i].frequency = sum
	}
}

// Generate the cumulative distribution in [0,1] where each interval will get a frequency between [0,1]
func (self *CumulativeDistribution) GenerateNormalized() {
	var sum, intervalSum float32
	for i := 0; i < len(self.values); i++ {
		sum += self.values[i].interval
	}
	for i := 0; i < len(self.values); i++ {
		intervalSum += self.values[i].interval / sum
		self.values[i].frequency = intervalSum
	}
}

// Generate the cumulative distribution in [0,1] where each value will have the same frequency and interval siz
func (self *CumulativeDistribution) GenerateUniform() {
	freq := 1 / len(self.values)
	for i := 0; i < len(self.values); i++ {
		//reset the interval to the normalized frequency
		self.values[i].interval = float32(freq)
		self.values[i].frequency = float32((i + 1) * freq)
	}
}

// Finds the value whose interval contains the given probability
// Binary search algorithm is used to find the value.
// param probability
// return the value whose interval contains the probability
func (self *CumulativeDistribution) ValueP(probability float32) interface{} {
	value := &CumulativeValue{}
	imax := len(self.values) - 1
	imin := 0
	var imid int
	for imin <= imax {
		imid = imin + ((imax - imin) / 2)
		value = self.values[imid]
		if probability < value.frequency {
			imax = imid - 1
		} else if probability > value.frequency {
			imin = imid + 1
		} else {
			break
		}
	}
	return value.value
}

// return the value whose interval contains a random probability in [0,1]
func (self *CumulativeDistribution) Value() interface{} {
	return "" // TODO:
	// return self.Value(random())
}

// return the amount of values
func (self *CumulativeDistribution) Size() int {
	return len(self.values)
}

// return the interval size for the value at the given position
func (self *CumulativeDistribution) GetInterval(index int) float32 {
	return self.values[index].interval
}

// return the value at the given position
func (self *CumulativeDistribution) GetValue(index int) interface{} {
	return self.values[index].value
}

// Set the interval size on the passed in object.
// The object must be present in the distribution.
func (self *CumulativeDistribution) SetInterval(obj interface{}, intervalSize float32) {
	for _, value := range self.values {
		if value.value == obj {
			value.interval = intervalSize
			return
		}
	}
}

// Sets the interval size for the value at the given index
func (self *CumulativeDistribution) SetIntervalIndex(index int, intervalSize float32) {
	self.values[index].interval = intervalSize
}

// Removes all the values from the distribution
func (self *CumulativeDistribution) Clear() {
	self.values = make([]*CumulativeValue, 10)
}
