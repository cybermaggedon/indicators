
package main

import (
	"log"
	"fmt"
	det "github.com/cybermaggedon/detector"
)

func main() {

	ii, err := det.LoadIndicatorsFromFile("ind3.json")
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

