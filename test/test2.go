
package main

import (
	"log"
//	"fmt"
)

func main() {

	ii, err := LoadIndicatorsFromFile("ind3.json")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fsmc := CreateFsmCollection(ii)

	fsmc.Reset()
	fsmc.Dump()
	fsmc.Update(Token{Type:"tcp", Value:"80"})
	fsmc.Dump()
	fsmc.Update(Token{Type:"url", Value:"http://www.example.com/malware.dat"})
	fsmc.Dump()
	fsmc.Update(Token{Type:"end", Value:""})
	fsmc.Dump()
	hits := fsmc.GetHits()

	for _, hit := range hits {
		hit.Dump()
	}

	/*
	for k, v := range *fsmm {
		fmt.Printf("%s -> %s:%s -> %s\n",
			k.State, k.Token.Type, k.Token.Value, v)
	}
*/

	_ = fsmc
	
}

