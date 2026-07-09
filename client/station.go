package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Station struct {
	StationName        string
	Lat                float64
	Long               float64
	CrsCode            string
	ConstituentCountry string
}

func ReadStations() map[string]Station {
	b, err := os.ReadFile("stations/stations.json")
	if err != nil {
		log.Fatal(err)
	}
	var ss []Station
	err = json.Unmarshal(b, &ss)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("stations", len(ss))
	stationMap := make(map[string]Station)
	for _, s := range ss {
		//fmt.Println("station", s)
		stationMap[s.CrsCode] = s
	}
	return stationMap
}
