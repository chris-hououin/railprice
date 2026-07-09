package main

import (
	"railprice/client"
	"sort"
	"strings"
)

type sta struct {
	Nlc  string
	Crs  string
	Name string
	Lat  float64
	Long float64
}

func newSta(nlc string) sta {
	crs := locations[nlc].Crs
	return sta{
		nlc,
		crs,
		strings.TrimSpace(locations[nlc].CcDescOut),
		stations[crs].Lat,
		stations[crs].Long,
	}
}

func fetchAllStations() []sta {
	var stations []sta
	for n := range locations {
		stations = append(stations, newSta(n))
	}
	sort.Slice(stations, func(i, j int) bool {
		return stations[i].Name < stations[j].Name
	})
	return stations
}

func fetchNlc(nlc string) []string {
	stationsOrGroups := client.FetchClusterStations(nlc, clusterNlcMap)

	var stations []string
	for _, s := range stationsOrGroups {
		ss := client.FetchStationsFromGroup(s, locationGroups)
		stations = append(stations, ss...)
	}

	return stations
}
