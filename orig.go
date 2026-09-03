package main

import (
	"encoding/json"
	"sort"
)

type DestinationPrices struct {
	Dest   sta
	Prices []FarePrice
}

type Destinations struct {
	Orig  sta
	Dests map[string]*DestinationPrices
}

func (destinationPrices *DestinationPrices) AddItem(item FarePrice) []FarePrice {
	destinationPrices.Prices = append(destinationPrices.Prices, item)
	return destinationPrices.Prices
}

func mapDestinationPrices(myFares []ReversibleFare) map[string]*DestinationPrices {
	var myPrices = make(map[string]*DestinationPrices)
	for _, f := range myFares {
		//fmt.Println(flows[f.fare.Flowid], f)
		var dests []string
		if f.reversed {
			dests = fetchNlc(flows[f.fare.Flowid].Orig)
		} else {
			dests = fetchNlc(flows[f.fare.Flowid].Dest)
		}

		for _, d := range dests {
			_, exists := myPrices[d]
			if exists {
				myPrices[d].AddItem(newFarePrice(f.fare))
			} else {
				newDestinationPrices := DestinationPrices{
					Dest:   newSta(d),
					Prices: []FarePrice{newFarePrice(f.fare)},
				}
				myPrices[d] = &newDestinationPrices
			}
		}
	}
	return myPrices
}

func fetchFlowIds(myClusterIds []string) map[string]ReversibleFlow {
	myFlowIds := make(map[string]ReversibleFlow)

	for _, id := range myClusterIds {
		myOrigFlows, exists := origFlows[id]
		if exists {
			for _, flow := range myOrigFlows {
				myFlowIds[flow.Flowid] = ReversibleFlow{flow.Flowid, false}
			}
		}

		myDestFlows, exists := destFlows[id]
		if exists {
			for _, flow := range myDestFlows {
				if flow.Direction == "R" {
					myFlowIds[flow.Flowid] = ReversibleFlow{flow.Flowid, true}
				}
			}
		}
	}
	return myFlowIds
}

func fetchPricesToFile(station string) {

	myClusterIds := fetchAllNlcsForStation(station)
	//fmt.Println("My clusters", time.Since(start), myClusterIds)

	myFlowIds := fetchFlowIds(myClusterIds)
	//fmt.Println("My flows", time.Since(start), len(myFlowIds))

	myFares := fetchFares(myFlowIds)
	//fmt.Println("My fares", time.Since(start), len(myFares))

	sort.Slice(myFares, func(i, j int) bool {
		ii := myFares[i].fare.Price
		jj := myFares[j].fare.Price
		return ii < jj
	})
	//fmt.Println("Sorted fares", time.Since(start), len(myFares))

	destinations := Destinations{
		Orig:  newSta(station),
		Dests: mapDestinationPrices(myFares),
	}
	//fmt.Println("mapped fares", time.Since(start), len(myFares))

	j, _ := json.Marshal(destinations)
	//fmt.Println("Marshal", time.Since(start))

	saveCompressed(j, "out/orig/"+station+".json.br")
	//fmt.Println("Save", time.Since(start))
}
