package detector

import (
	"fmt"
)

// An indicator term, can be one of or, and, not or type/value pair.
type Term struct {
	Type  string  `json:"type,omitempty"`
	Value string  `json:"value,omitempty"`
	And   []*Term `json:"and,omitempty"`
	Or    []*Term `json:"or,omitempty"`
	Not   *Term   `json:"not,omitempty"`
}

// Dump a term
func (l *Term) Dump(indent int) {
	if l.Type != "" && l.Value != "" {
		for v := 0; v < indent+2; v++ {
			fmt.Print("  ")
		}
		fmt.Println(l.Type, ":", l.Value)
		return
	}
	if len(l.And) > 0 {
		for v := 0; v < indent+2; v++ {
			fmt.Print("  ")
		}
		fmt.Println("And")
		for _, v := range l.And {
			v.Dump(indent + 1)
		}
		return
	}
	if len(l.Or) > 0 {
		for v := 0; v < indent+2; v++ {
			fmt.Print("  ")
		}
		fmt.Println("Or")
		for _, v := range l.Or {
			v.Dump(indent + 1)
		}
		return
	}
	if l.Not != nil {
		for v := 0; v < indent+2; v++ {
			fmt.Print("  ")
		}
		fmt.Println("Not")
		l.Not.Dump(indent + 1)
		return
	}
}

// Returns true if this is an AND expression
func (l *Term) IsAnd() bool {
	return len(l.And) > 0
}

// Returns true if this is an OR expression
func (l *Term) IsOr() bool {
	return len(l.Or) > 0
}

// Returns true if this is a NOT expression
func (l *Term) IsNot() bool {
	return l.Not != nil
}

// Returns true if this is a type/value expression
func (l *Term) IsMatchTerm() bool {
	return l.Type != "" && l.Value != ""
}

// Works out state change based on 'end' state.
func (l *Term) RecordEnd(state *Combination, n *Navigator) {

	if l.IsAnd() {
		for _, v := range l.And {
			v.RecordEnd(state, n)
		}
	}
	if l.IsOr() {
		for _, v := range l.Or {
			v.RecordEnd(state, n)
		}
	}
	if l.IsNot() {
		if state.Contains(l) {
			return
		}
		l.Not.RecordEnd(state, n)
		if state.Contains(l.Not) {
			return
		}
		state.Add(l)
		if n.parent[l] != nil {
			n.parent[l].Evaluate(state, n)
		}
	}
	if l.IsMatchTerm() {
		return
	}
}

// Evaluates a term's state when child state has changed
func (l *Term) Evaluate(state *Combination, n *Navigator) {

	if state.Contains(l) {
		return
	}

	parent := n.parent[l]

	if l.IsAnd() {
		count := 0
		for _, v := range l.And {
			if state.Contains(v) {
				count++
			}
		}
		if count == len(l.And) {
			state.Add(l)
			if parent != nil {
				parent.Evaluate(state, n)
			}
		}
		return
	}

	if l.IsOr() {
		count := 0
		for _, v := range l.Or {
			if state.Contains(v) {
				count++
			}
		}
		if count > 0 {
			state.Add(l)
			if parent != nil {
				parent.Evaluate(state, n)
			}
		}
		return
	}

	if l.IsNot() {
		return
	}

	if l.IsMatchTerm() {
		return
	}

}

// Works out state change based on activating a term - used when a match
// term becomes true
func (l *Term) Activate(state *Combination, n *Navigator) {

	if !l.IsMatchTerm() {
		panic("Activate called on non-term")
	}

	if state.Contains(l) {
		return
	}

	state.Add(l)
	if n.parent[l] != nil {
		n.parent[l].Evaluate(state, n)
	}

}

// A walk observer is a callback used to study every term in a term tree,
// breadth-first
type WalkObserver func(*Term, interface{}, *Term) error

// Walks a term tree, calling a callback for every term
func (l *Term) WalkState(wo WalkObserver, state interface{}, parent *Term) error {
	for _, v := range l.And {
		err := v.WalkState(wo, state, l)
		if err != nil {
			return err
		}
	}
	for _, v := range l.Or {
		err := v.WalkState(wo, state, l)
		if err != nil {
			return err
		}
	}
	if l.Not != nil {
		err := l.Not.WalkState(wo, state, l)
		if err != nil {
			return err
		}
	}
	err := wo(l, state, parent)
	return err
}

// Wrapper for WalkState passing nil defaults
func (l *Term) Walk(wo WalkObserver) error {
	return l.WalkState(wo, nil, nil)
}
