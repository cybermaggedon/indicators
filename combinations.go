package detector

// A Combination represents a set of terms.
type Combination struct {
	set map[*Term]bool
}

// Create an empty Combination
func NewCombination() Combination {
	c := Combination{}
	c.set = make(map[*Term]bool)
	return c
}

// Return the number of elements in the set
func (s *Combination) Size() int {
	return len(s.set)
}

// Return true if the set is empty
func (s *Combination) IsEmpty() bool {
	return len(s.set) == 0
}

// Add a Term to the set
func (c *Combination) Add(l *Term) {
	c.set[l] = true
}

// Delete a term from the set
func (c *Combination) Delete(l *Term) {
	delete(c.set, l)
}

// Return true if the set contains the element.
func (c *Combination) Contains(l *Term) bool {
	_, ok := c.set[l]
	return ok
}

// Return a chan which iterates over all elements in the set
func (c *Combination) Iter() chan *Term {
	ch := make(chan *Term)
	go func() {
		for v, _ := range c.set {
			ch <- v
		}
		close(ch)
	}()
	return ch
}

// Returns an identical copy of the Combination.  Altering the new
// Combination does not affect the old
func (c *Combination) Copy() Combination {
	c2 := NewCombination()
	for k, _ := range c.set {
		c2.set[k] = true
	}
	return c2
}

// Converts the Combination to an array of Terms.
func (c *Combination) ToArray() []*Term {
	keys := []*Term{}
	for v, _ := range c.set {
		keys = append(keys, v)
	}
	return keys
}

// Combinations is an array of Combination
type Combinations []Combination
