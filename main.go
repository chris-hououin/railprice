package main

import (
	"fmt"
	"os"
	"railprice/client"
	"sync"
	"time"
)

var allStations []sta
var locations map[string]client.Location
var flows map[string]client.Flow
var origFlows map[string][]client.Flow
var destFlows map[string][]client.Flow
var fares map[string][]client.Fare
var clusters []client.Cluster
var clusterNlcMap map[string][]string
var locationGroups map[string][]string
var routes map[string]client.Route
var stations map[string]client.Station
var ticketTypes map[string]client.TicketType
var ticketAdvances []client.TicketAdvance
var ticketValidityRecords map[string]client.TicketValidity

func readFiles(filename string) {
	locations = client.ReadLocations(filename)
	allStations = fetchAllStations()
	flows, fares = client.ReadFlows(filename)
	origFlows, destFlows = client.MapFlows(flows)
	clusters = client.ReadClusters(filename)
	clusterNlcMap = client.MapClusters(clusters)
	locationGroups = client.MapStationGroups(locations)
	routes = client.ReadRoutes(filename)
	ticketAdvances = client.ReadTicketAdvance(filename)
	ticketValidityRecords = client.ReadTicketValidity(filename)
	ticketTypes = client.ReadTicketType(filename, ticketAdvances, ticketValidityRecords)
	stations = client.ReadStations()
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: go run main.go <filename>")
		return
	}

	WORKERS := 8
	filename := os.Args[1]
	readFiles(filename)

	err := os.MkdirAll("out", os.ModePerm)
	if err != nil {
		panic(err)
	}

	storeStations(allStations)

	start := time.Now()

	sem := make(chan struct{}, WORKERS)
	var wg sync.WaitGroup
	for i, s := range allStations {
		wg.Go(func() {
			sem <- struct{}{}
			defer func() { <-sem }()
			fmt.Println(i, s.Nlc)
			fetchPricesToFile(s.Nlc)
		})
	}
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Println("Took", elapsed)
}
