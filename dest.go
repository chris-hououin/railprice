package main

import (
	"encoding/json"
	"sort"
)

type OriginPrices struct {
	Orig   sta
	Prices []FarePrice
}

type Origins struct {
	Dest  sta
	Origs map[string]*OriginPrices
}

func (originPrices *OriginPrices) AddItem(item FarePrice) []FarePrice {
	originPrices.Prices = append(originPrices.Prices, item)
	return originPrices.Prices
}

func mapOriginPrices(myFares []ReversibleFare) map[string]*OriginPrices {
	var myPrices = make(map[string]*OriginPrices)
	for _, f := range myFares {
		//fmt.Println(flows[f.fare.Flowid], f)
		var origs []string
		if f.reversed {
			origs = fetchNlcs(flows[f.fare.Flowid].Dest)
		} else {
			origs = fetchNlcs(flows[f.fare.Flowid].Orig)
		}

		for _, o := range origs {
			_, exists := myPrices[o]
			if exists {
				myPrices[o].AddItem(newFarePrice(f.fare))
			} else {
				newOriginPrices := OriginPrices{
					Orig:   newSta(o),
					Prices: []FarePrice{newFarePrice(f.fare)},
				}
				myPrices[o] = &newOriginPrices
			}
		}
	}
	return myPrices
}

func fetchFlowIdsForDestinations(myDestinationClusterIds []string) map[string]ReversibleFlow {
	myFlowIds := make(map[string]ReversibleFlow)

	for _, id := range myDestinationClusterIds {
		myDestFlows, exists := destFlows[id]
		if exists {
			for _, flow := range myDestFlows {
				myFlowIds[flow.Flowid] = ReversibleFlow{flow.Flowid, false}
			}
		}

		myOrigFlows, exists := origFlows[id]
		if exists {
			for _, flow := range myOrigFlows {
				if flow.Direction == "R" {
					myFlowIds[flow.Flowid] = ReversibleFlow{flow.Flowid, true}
				}
			}
		}
	}
	return myFlowIds
}

func processDestinationToFile(station string) {

	myClusterIds := fetchAllNlcsForStation(station)
	//fmt.Println("My clusters", time.Since(start), myClusterIds)

	myFlowIds := fetchFlowIdsForDestinations(myClusterIds)
	//fmt.Println("My flows", time.Since(start), len(myFlowIds))

	myFares := fetchFares(myFlowIds)
	//fmt.Println("My fares", time.Since(start), len(myFares))

	sort.Slice(myFares, func(i, j int) bool {
		ii := myFares[i].fare.Price
		jj := myFares[j].fare.Price
		return ii < jj
	})
	//fmt.Println("Sorted fares", time.Since(start), len(myFares))

	origins := Origins{
		Dest:  newSta(station),
		Origs: mapOriginPrices(myFares),
	}
	//fmt.Println("mapped fares", time.Since(start), len(myFares))

	j, _ := json.Marshal(origins)
	//fmt.Println("Marshal", time.Since(start))

	saveCompressed(j, "out/dest/"+station+".json.br")
	//fmt.Println("Save", time.Since(start))
}
