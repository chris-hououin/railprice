package main

import (
	"fmt"
	"os"
	"railprice/client"
	"sync"
	"time"
)

var flows map[string]client.Flow
var origFlows map[string][]client.Flow
var destFlows map[string][]client.Flow
var fares map[string][]client.Fare
var clusters []client.Cluster
var clusterNlcMap map[string][]string
var locations map[string]client.Location
var locationGroups map[string][]string
var stations map[string]client.Station
var routes map[string]client.Route
var ticketTypes map[string]client.TicketType
var ticketAdvances []client.TicketAdvance
var ticketValidityRecords map[string]client.TicketValidity

func readFiles(filename string) {
	flows, fares = client.ReadFlows(filename)
	origFlows, destFlows = client.MapFlows(flows)
	clusters = client.ReadClusters(filename)
	clusterNlcMap = client.MapClusters(clusters)
	locations = client.ReadLocations(filename)
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

	//for i, s := range fetchAllStations() {
	//	fmt.Println(i, s.Nlc)
	//	fetchPricesToFile(s.Nlc)
	//}

	start := time.Now()

	sem := make(chan struct{}, WORKERS)
	var wg sync.WaitGroup
	for i, s := range fetchAllStations() {
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
