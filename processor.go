package main

import (
	"bytes"
	"io"
	"os"
	"railprice/client"
	"slices"

	"github.com/andybalholm/brotli"
)

type ReversibleFlow struct {
	flowId   string
	reversed bool
}

type ReversibleFare struct {
	fare     client.Fare
	reversed bool
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

func fetchFares(myFlowIds map[string]ReversibleFlow) []ReversibleFare {
	var myFares []ReversibleFare
	for _, flowId := range myFlowIds {
		for _, fare := range fares[flowId.flowId] {
			myFares = append(myFares, ReversibleFare{fare, myFlowIds[flowId.flowId].reversed})
		}
	}
	return myFares
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
