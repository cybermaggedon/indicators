
// Detector library includes classes for loading indicators, building FSMs from indicators,
// using FSMs to analyze terms and make hit decisions.  Call CreateFsmCollection to create
// a collection from indicators.  Call Reset method to start scanning, Update methods with
// all tokens to scan, and then GetHits to find out the indicators which hit.
package detector

import (
	"fmt"
)

// A collection of indicators and derived FSMs.
type FsmCollection struct {

	// Array of all FSMs
	Fsms []*FsmMap

	// Map from FSM to Indicator
	Indicators map[*FsmMap]*Indicator

	// Map from activator tokens to array  of FSMs.  The tokens include all tokens which
	// lead out of the 'init' state for each FSM.  By 'active', I'm refering to any
	// FSM which has left the 'init' state.
	Activators map[Token][]*FsmMap

	// The current state of all active FSMs, maps FSM to the current state string.
	State map[*FsmMap]string
}

// Dump an FSM collection showing all tracked states.
func (c *FsmCollection) Dump() {
	fmt.Println("State:")
	for fsm, state := range c.State {
		fmt.Println("  ", c.Indicators[fsm].Id, " in state ", state)
	}
}

// Resets an FSM collection so that all FSMs revert to the inactive state.  This would be
// called to forget existing scanning history when scanning something new.
func (c *FsmCollection) Reset() {
	c.State = map[*FsmMap]string{}
}

// Update an FSM collection for a new token.
func (c *FsmCollection) Update(token Token) {

	// If the token is an activator, activate all relevant FSMs to the init state.
	// The next code segment will apply the transition from init to the next state.
	if fsms, ok := c.Activators[token]; ok {
		for _, fsm := range fsms {
			if _, ok = c.State[fsm]; !ok {
				c.State[fsm] = "init"
			}
		}
	}

	// Iterate over all active FSMs, moving to the next state if necessary.
	for fsm, state := range c.State {
		event := FsmEvent{State: state, Token: token}
		if newstate, ok := (*fsm)[event]; ok {
			c.State[fsm] = newstate
		}
	}

}

// Returns all active FSM hits.  This would be called once scanning is complete to
// return hits.
func (c *FsmCollection) GetHits() []*Indicator {

	hits := []*Indicator{}

	for fsm, state := range c.State {
		if state == "hit" {
			hits = append(hits, c.Indicators[fsm])
		}
	}

	return hits

}

// Create an FSM collection from a set of indicators.
func CreateFsmCollection(ii *Indicators) *FsmCollection {

	// Initialise the FSM collection to null state, and allocate all maps.
	fsmc := FsmCollection{}
	fsmc.Indicators = map[*FsmMap]*Indicator{}
	fsmc.Activators = map[Token][]*FsmMap{}
	fsmc.State = map[*FsmMap]string{}

	// Iterate over indicators
	for _, ind := range ii.Indicators {

		// Generate the FSM for this indicator.
		fsm := ind.GenerateFsm()

		// Convert the FSM to its 'map' form.
		fsmm := fsm.Mapify()

		// Add mapping from FSM to corresponding indicator.
		fsmc.Indicators[fsmm] = ind


		// Get activator terms 
		activs := fsm.GetActivators()

		// Add activator terms to the activator map.
		for _, activ := range activs {
			if _, ok := fsmc.Activators[activ]; !ok {
				fsmc.Activators[activ] = nil
			}
			fsmc.Activators[activ] = append(fsmc.Activators[activ], fsmm)
		}

		// Append FSM to FSM list.
		fsmc.Fsms = append(fsmc.Fsms, fsmm)

	}

	return &fsmc

}

