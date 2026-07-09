package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"railprice/client"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/andybalholm/brotli"
)

type DestinationPrice struct {
	Price       int
	CrossLondon bool
	RouteCode   string
	Route       string
	TicketCode  string
	TicketDesc  string
	Advance     bool
	IsDay       bool
	TicketClass string
	TicketType  string
	Restriction string
}

type DestinationPrices struct {
	Dest   sta
	Prices []DestinationPrice
}

type ReversibleFlow struct {
	flowId   string
	reversed bool
}

type ReversibleFare struct {
	fare     client.Fare
	reversed bool
}

func (destinationPrices *DestinationPrices) AddItem(item DestinationPrice) []DestinationPrice {
	destinationPrices.Prices = append(destinationPrices.Prices, item)
	return destinationPrices.Prices
}

type Destinations struct {
	Orig  sta
	Dests map[string]*DestinationPrices
}

func newDestinationPrice(f client.Fare) DestinationPrice {
	priceMinorUnit, err := strconv.Atoi(f.Price)
	if err != nil {
		panic(err)
	}

	return DestinationPrice{
		Price:       priceMinorUnit,
		CrossLondon: flows[f.Flowid].Crosslondon == "1",
		RouteCode:   flows[f.Flowid].Route,
		Route:       strings.TrimSpace(routes[flows[f.Flowid].Route].AtbDesc),
		TicketCode:  f.Ticketcode,
		TicketDesc:  ticketTypes[f.Ticketcode].AtbDesc,
		Advance:     ticketTypes[f.Ticketcode].IsAdvance,
		IsDay:       ticketTypes[f.Ticketcode].IsDay,
		TicketClass: ticketTypes[f.Ticketcode].TicketClass,
		TicketType:  ticketTypes[f.Ticketcode].TicketType,
		Restriction: f.Restriction,
	}
}

func fetchAllNlcsForStation(station string) []string {
	stations := []string{station}

	myStationIds := slices.Clone(stations)
	for _, s := range stations {
		myStationIds = append(myStationIds, locations[s].FareGroup)
	}

	var myClusterIds []string
	myClusterIds = slices.Clone(myStationIds)
	for _, c := range clusters {
		if slices.Contains(myStationIds, c.ClusterNlc) {
			myClusterIds = append(myClusterIds, c.ClusterId)
		}
	}

	return myClusterIds
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

func fetchFares(myFlowIds map[string]ReversibleFlow) []ReversibleFare {
	var myFares []ReversibleFare
	for _, flowId := range myFlowIds {
		for _, fare := range fares[flowId.flowId] {
			myFares = append(myFares, ReversibleFare{fare, myFlowIds[flowId.flowId].reversed})
		}
	}
	return myFares
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

	saveCompressed(j, "out/"+station+".json.br")
	//fmt.Println("Save", time.Since(start))
}

func saveCompressed(j []byte, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := brotli.NewWriterOptions(f, brotli.WriterOptions{Quality: brotli.DefaultCompression})
	defer w.Close()

	_, err = io.Copy(w, bytes.NewBuffer(j))
	if err != nil {
		panic(err)
	}

	err = w.Flush()
	if err != nil {
		panic(err)
	}
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
				myPrices[d].AddItem(newDestinationPrice(f.fare))
			} else {
				newDestinationPrices := DestinationPrices{
					Dest:   newSta(d),
					Prices: []DestinationPrice{newDestinationPrice(f.fare)},
				}
				myPrices[d] = &newDestinationPrices
			}
		}
	}
	return myPrices
}
