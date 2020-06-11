package detector

import (
	"strconv"
)

// The term tree is a simple tree, it is useful to collect information to
// allow some other forms of navigation in the tree.
type Navigator struct {

	// Maps state name to logic element
	logic_state map[string]*Term

	// Maps logic element to state name
	state_name map[*Term]string

	// Element to parent.  root element has parent = nil
	parent map[*Term]*Term

	// Set of basic states
	basic_states Combination

	// Array of all match terms i.e. only have key/value pairs, not
	// AND/OR/NOT.
	terms []*Term
}

// Constructs a navigator from a term tree.
func (i *Indicator) BuildNavigator() *Navigator {

	// Allocate Navigator.
	n := &Navigator{
		logic_state: make(map[string]*Term),
		state_name:  make(map[*Term]string),
		parent:      make(map[*Term]*Term),
	}

	// New state IDs start at 1.
	next_id := 1

	// Walk the term tree, collecting information
	i.Walk(func(l *Term, state interface{}, par *Term) error {
		state_id := "s" + strconv.Itoa(next_id)
		next_id++
		n.state_name[l] = state_id
		n.logic_state[state_id] = l
		n.parent[l] = par
		return nil
	})

	// Collect state information
	n.basic_states, n.terms = i.DiscoverStates(n)

	return n

}
