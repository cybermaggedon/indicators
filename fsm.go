package detector

import (
	"fmt"
	"sort"
)

// FIXME: Strategy thing is missing.

// Represents a type/value pair.
type Token struct {
	Type  string
	Value string
}

// Represents an FSM transition.
type FsmTransition struct {

	// The 'current' state.
	Current string

	// The token which activates this transition
	Token []Token

	// The state to transition to
	Next string
}

// An Fsm is an array of transitions.  This is how the FSM is represented
// when we're working on it.  Once complete, the FSM is converted to an
// FsmMap which is a structure optimised for navigation.
type Fsm struct {
	Transitions []FsmTransition
}

func (fsm *Fsm) Dump() {
	for _, v := range fsm.Transitions {
		for _, w := range v.Token {
			fmt.Printf("%s -> %s:%s -> %s\n",
				v.Current, w.Type, w.Value, v.Next)
		}
	}
}

// Describes an event causing an FSM transition, consists of a current
// state, and a token.
type FsmEvent struct {
	State string
	Token Token
}

// An FSM optimised for navigation, maps FSM events to new states.
type FsmMap map[FsmEvent]string

func (fsm *FsmMap) Dump() {
	for event, next := range *fsm {
		fmt.Printf("%s -> %s:%s -> %s\n",
			event.State, event.Token.Type, event.Token.Value, next)
	}
}

// Converts an Fsm to an FsmMap.
func (fsm *Fsm) Mapify() *FsmMap {
	fsm_map := FsmMap{}

	for _, trans := range fsm.Transitions {
		for _, token := range trans.Token {
			key := FsmEvent{
				State: trans.Current,
				Token: token,
			}
			fsm_map[key] = trans.Next
		}
	}

	return &fsm_map

}

// Returns array of activators, which is the tokens which take the FSM out
// of the 'init' state.
func (fsm *Fsm) GetActivators() []Token {
	activs := map[Token]bool{}
	for _, transition := range fsm.Transitions {
		if transition.Current == "init" {
			for _, token := range transition.Token {
				activs[token] = true
			}
		}
	}

	lst := make([]Token, 0, len(activs))
	for v, _ := range activs {
		lst = append(lst, v)
	}

	return lst

}

// Flattens an FSM, compacts all edges to the same current/next transition.
// (current, [token1], next), (current, [token2], next) -> (current,
// [token1, token2], next)
func (fsm *Fsm) Flatten() {

	type Path struct {
		Current string
		Next    string
	}

	// Convert to map
	transitions := make(map[Path][]Token)

	for _, v := range fsm.Transitions {

		for _, term := range v.Token {

			key := Path{v.Current, v.Next}
			if _, ok := transitions[key]; !ok {
				transitions[key] = []Token{}
			}
			transitions[key] = append(transitions[key], term)
		}

	}

	// Back to list
	fsm.Transitions = []FsmTransition{}

	for k, v := range transitions {
		fsm.Transitions = append(fsm.Transitions, FsmTransition{k.Current, v, k.Next})
	}

}

// Removes unnavigable transitions from the FSM.  Two types of unnavigable
// transitions are discovered:
// - those which cannot lead to 'hit' are re-labeled as 'fail'.
// - those which cannot be discovered from 'init' are removed altogether.
func (fsm *Fsm) RemoveInvalidTransitions(n *Navigator) {

	valid_hit_states := make(map[string]bool)
	valid_hit_states["hit"] = true

	last_run_size := 1

	for {
		for _, v := range fsm.Transitions {
			if _, ok := valid_hit_states[v.Next]; ok {
				valid_hit_states[v.Current] = true
			}
		}
		if len(valid_hit_states) == last_run_size {
			break
		}
		last_run_size = len(valid_hit_states)
	}

	valid_trav_states := make(map[string]bool)
	valid_trav_states["init"] = true

	last_run_size = 1

	for {
		for _, v := range fsm.Transitions {
			if _, ok := valid_trav_states[v.Current]; ok {
				valid_trav_states[v.Next] = true
			}
		}
		if len(valid_trav_states) == last_run_size {
			break
		}
		last_run_size = len(valid_trav_states)
	}

	transitions := []FsmTransition{}
	for _, v := range fsm.Transitions {
		if _, ok := valid_trav_states[v.Current]; !ok {
			continue
		}
		if _, ok := valid_hit_states[v.Next]; !ok {
			v.Next = "fail"
		}
		transitions = append(transitions, v)
	}
	fsm.Transitions = transitions

}

// Converts a term combination to a name for the state
// e.g. [s2, s5, s12] -> "s2-5-12"
func NameCombinationState(comb *Combination, n *Navigator, l *Term) string {

	// If root of tree is active, this is a 'hit' state.
	if comb.Contains(l) {
		return "hit"
	}

	// Empty set is "init"
	if comb.IsEmpty() {
		return "init"
	}

	// Create an array of states.
	// FIXME: Use Copy method.
	states := make([]string, 0, comb.Size())
	for k := range comb.Iter() {
		states = append(states, n.state_name[k][1:])
	}

	// Sort strings, so that this function is always stable.
	sort.Strings(states)

	// Start with the letter 's'.
	cstate := "s"
	sep := ""

	// Iterate, append state names, minus the 's'.
	for _, v := range states {
		cstate = cstate + sep + v
		sep = "-"
	}

	return cstate

}

// Takes a set of terms in the form of a channel, and returns all term
// combinations, including all terms and the empty set.
// e.g. [s2, s5] -> {[], [s2], [s5], [s2, s5]}
func GetCombinations(l chan *Term) Combinations {

	// Basic design is, this is recursive.  This recursion takes the
	// first element from the channel, and delegates the rest of the
	// channel to a further recursive call of this function.

	// Get first element from channel.
	elt, ok := <-l

	// If list is empty, return a list containing only the empty set.
	if !ok {
		return Combinations{NewCombination()}
	}

	// Recurse the next level down.  Return the combination for
	// everything else in this set.
	subset := GetCombinations(l)

	// For this recursion start a new array.
	combs := Combinations{}

	// Iterate over combations from the subset.
	for _, comb := range subset {

		// Add all elements from the set
		combs = append(combs, comb)

		// Also add all elements with this element added.
		comb2 := comb.Copy()
		comb2.Add(elt)
		combs = append(combs, comb2)

	}

	return combs

}

// Return new state which describes what happens when this particular token
// is term is observed.
func (i *Indicator) ExerciseToken(current Combination, term *Term, n *Navigator) Combination {

	// Start with a copy of the 'current' state.
	next := current.Copy()

	// term=nil tests the 'end' state.
	if term == nil {
		// See what happens when 'end' is observed.
		i.RecordEnd(&next, n)
	} else {
		// See what happens when this term is observed.
		term.Activate(&next, n)
	}

	// Return the new state.
	return next

}

// Extract all FSM transition by exercising all possible terms in all
// possible basic state combinations.
func (i *Indicator) ExtractTransitions(basic_combis Combinations, terms []*Term, n *Navigator) *Fsm {

	// Start with transitions as an empty array
	transitions := []FsmTransition{}

	// Get root of term tree
	root := &i.Term

	// Iterate over all state combinations
	for _, comb := range basic_combis {

		// Get state name for the combination state.
		cur_state := NameCombinationState(&comb, n, root)

		// Iterate over all terms
		for _, term := range terms {

			// Get the state which results from observing this
			// token in the current state
			newstate := i.ExerciseToken(comb, term, n)

			// The new state combination is reduced by taking
			// out all states which aren't basic states.
			newstate2 := NewCombination()
			for v := range newstate.Iter() {
				if n.basic_states.Contains(v) {
					newstate2.Add(v)
				}
			}

			// Convert this new state to a state name.
			next_state := NameCombinationState(&newstate2, n, root)

			// If the term causes no state transition, we can
			// move on.
			if cur_state == next_state {
				continue
			}

			// Create transation and add to transition array
			transition := FsmTransition{
				Current: cur_state,
				Token:   []Token{Token{term.Type, term.Value}},
				Next:    next_state,
			}
			transitions = append(transitions, transition)

		}

	}

	// Same again, but see what happens when applying 'end'
	// FIXME: Could just add {"end", ""} to the token array.  So we
	// wouldn't need this code.
	for _, comb := range basic_combis {

		newstate := i.ExerciseToken(comb, nil, n)

		cur_state := NameCombinationState(&comb, n, root)
		next_state := NameCombinationState(&newstate, n, root)

		if cur_state == next_state {
			continue
		}

		transition := FsmTransition{
			Current: cur_state,
			Token:   []Token{Token{"end", ""}},
			Next:    next_state,
		}

		transitions = append(transitions, transition)

	}

	// Create and return an FSM.
	fsm := &Fsm{
		Transitions: transitions,
	}
	return fsm

}

// Find all 'basic states' and terms.  The basic states are the places in
// the tree where state information can stored: Children of AND and children
// of NOT.  NOT nodes are never themselves basic state nodes.
func (i *Indicator) DiscoverStates(n *Navigator) (Combination, []*Term) {

	// Initialise
	basic_states := NewCombination()
	terms := []*Term{}

	// Walk the term tree
	i.Walk(func(l *Term, state interface{}, par *Term) error {

		// Collect match terms
		if l.IsMatchTerm() {
			terms = append(terms, l)
		}

		// Skip the root element
		if n.parent[l] == nil {
			return nil
		}

		// Skip a NOT elemetn
		if l.IsNot() {
			return nil
		}

		// If we got this far, and parent is AND or NOT, this is a
		// basic state.
		if n.parent[l].IsAnd() {
			basic_states.Add(l)
		}
		if n.parent[l].IsNot() {
			basic_states.Add(l)
		}
		return nil
	})

	return basic_states, terms

}

// Dump a term tree.
func (l *Term) DumpTree(n *Navigator, indent int) {

	for i := 0; i < indent; i++ {
		fmt.Print("  ")
	}
	fmt.Print(n.state_name[l], ": ")

	if l.IsAnd() {
		fmt.Println("AND")
		for _, v := range l.And {
			v.DumpTree(n, indent+1)
		}
	}
	if l.IsOr() {
		fmt.Println("OR")
		for _, v := range l.Or {
			v.DumpTree(n, indent+1)
		}
	}
	if l.IsNot() {
		fmt.Println("NOT")
		l.Not.DumpTree(n, indent+1)
	}
	if l.IsMatchTerm() {
		fmt.Println(l.Type, ":", l.Value)
	}

}

// Generate an FSM from an indicator
func (i *Indicator) GenerateFsm() *Fsm {

	n := i.BuildNavigator()

	// Get all combinations of the basic states.
	combinations := GetCombinations(n.basic_states.Iter())

	// Get all transitions.
	fsm := i.ExtractTransitions(combinations, n.terms, n)

	// Flatten FSM
	fsm.Flatten()

	// Remove invalid transition
	fsm.RemoveInvalidTransitions(n)

	return fsm

}
