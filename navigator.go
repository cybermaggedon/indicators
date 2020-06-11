package detector

import (
	"strconv"
)

// The term tree is a simple tree, it is useful to collect information to
// allow some other forms of navigation in the tree.
type Navigator struct {

	// Maps state name to logic element
	LogicState map[string]*Term

	// Maps logic element to state name
	StateName map[*Term]string

	// Element to parent.  root element has parent = nil
	Parent map[*Term]*Term

	// Set of basic states
	BasicStates Combination

	// Array of all match terms i.e. only have key/value pairs, not
	// AND/OR/NOT.
	Terms []*Term
}

// Constructs a navigator from a term tree.
func (i *Indicator) BuildNavigator() *Navigator {

	// Allocate Navigator.
	n := &Navigator{
		LogicState: make(map[string]*Term),
		StateName:  make(map[*Term]string),
		Parent:     make(map[*Term]*Term),
	}

	// New state IDs start at 1.
	next_id := 1

	// Walk the term tree, collecting information
	i.Walk(func(l *Term, state interface{}, par *Term) error {
		state_id := "s" + strconv.Itoa(next_id)
		next_id++
		n.StateName[l] = state_id
		n.LogicState[state_id] = l
		n.Parent[l] = par
		return nil
	})

	// Collect state information
	n.BasicStates, n.Terms = i.DiscoverStates(n)

	return n

}
