
package main

import (
	"log"
	det "github.com/cybermaggedon/detector"
	
)

func main() {

	ii, err := det.LoadIndicatorsFromFile("ind3.json")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fsmc := det.CreateFsmCollection(ii)

	fsmc.Reset()
	fsmc.Dump()
	fsmc.Update(det.Token{Type:"tcp", Value:"80"})
	fsmc.Dump()
	fsmc.Update(det.Token{Type:"url", Value:"http://www.example.com/malware.dat"})
	fsmc.Dump()
	fsmc.Update(det.Token{Type:"end", Value:""})
	fsmc.Dump()
	hits := fsmc.GetHits()

	for _, hit := range hits {
		hit.Dump()
	}

	_ = fsmc
	
}

