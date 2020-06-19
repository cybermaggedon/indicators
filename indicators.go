package detector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Represents an indicator set.  Is JSON serialisable
type Indicators struct {
	Description string       `json:"description,omitempty"`
	Version     string       `json:"version,omitempty"`
	Indicators  []*Indicator `json:"indicators,omitempty"`
}

// Add an indicator
func (ii *Indicators) Add(i *Indicator) {
	ii.Indicators = append(ii.Indicators, i)
}

// Get an indicator by ID.  O(n)
func (ii *Indicators) Get(id string) *Indicator {
	for _, v := range ii.Indicators {
		if v.Id == id {
			return v
		}
	}
	return nil
}

// An indicator descriptor describes the results of a hit.
type Descriptor struct {
	Description string  `json:"description,omitempty"`
	Category    string  `json:"category,omitempty"`
	Author      string  `json:"author,omitempty"`
	Source      string  `json:"source,omitempty"`
	Type        string  `json:"type,omitempty"`
	Value       string  `json:"value,omitempty"`
	Probability float32 `json:"probability,omitempty"`
}

// An indicator
type Indicator struct {
	Id         string     `json:"id,omitempty"`
	Descriptor Descriptor `json:"descriptor,omitempty"`
	Term
}

// Loads indicators from a byte array
func LoadIndicators(data []byte) (*Indicators, error) {
	var ii Indicators
	err := json.Unmarshal(data, &ii)
	if err != nil {
		return nil, err
	}

	// Having loaded indicators, set probability to 1.0 for anything
	// without a probability.  Can't tell the difference between
	// setting 0 and a missing field, but in practice specifying a
	// probability of 0 doesn't make sense. 
	for i, _ := range ii.Indicators {
		if ii.Indicators[i].Descriptor.Probability == 0.0 {
			ii.Indicators[i].Descriptor.Probability = 1.0
		}
	}

	return &ii, nil
}

// Loads indicators from a file.
func LoadIndicatorsFromFile(path string) (*Indicators, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadIndicators(data)
}

// Dumps indicators.
func (ii *Indicators) Dump() {
	fmt.Println("Id:", ii.Version)
	fmt.Println("Description:", ii.Description)
	for _, i := range ii.Indicators {
		i.Dump()
	}
}

// Dumps an indicator.
func (i *Indicator) Dump() {
	fmt.Println()
	fmt.Println("  Id:", i.Id)
	fmt.Println("  Description:", i.Descriptor.Description)
	fmt.Println("  Category:", i.Descriptor.Category)
	fmt.Println("  Author:", i.Descriptor.Author)
	fmt.Println("  Source:", i.Descriptor.Source)
	fmt.Println("  Type:", i.Descriptor.Type)
	fmt.Println("  Value:", i.Descriptor.Value)
	fmt.Println("  Probability:", i.Descriptor.Probability)
	i.Term.Dump(0)
}
