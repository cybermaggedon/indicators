
package main

import (
	"log"
	"fmt"
)

func main() {

	ii, err := LoadIndicatorsFromFile("ind3.json")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	ind := ii.Indicators[1]

	fsm := ind.GenerateFsm()

	fsmm := fsm.Mapify()

	for k, v := range *fsmm {
		fmt.Printf("%s -> %s:%s -> %s\n",
			k.State, k.Token.Type, k.Token.Value, v)
	}
	
}

