package client

import (
	"fmt"
	"strings"
)

func determineRecordType(line string) RecordType {
	if strings.HasPrefix(line, "/!!") {
		return CommentRecord
	} else if strings.HasPrefix(line, "RF") {
		return FlowRecord
	} else if strings.HasPrefix(line, "RT") {
		return FareRecord
	} else {
		panic("Unrecognized record type: " + line)
	}
}

func parseLine(line string) interface{} {
	recordType := determineRecordType(line)
	if recordType == FlowRecord {
		return newFlow(line)
	} else if recordType == FareRecord {
		return newFare(line)
	} else if recordType == CommentRecord {
		return nil
	}
	return nil
}

type Flow struct {
	Orig        string
	Dest        string
	Route       string
	status      string
	usage       string
	Direction   string
	end         string
	start       string
	toc         string
	Crosslondon string
	discount    string
	publication string
	Flowid      string
}

type Fare struct {
	Flowid      string
	Ticketcode  string
	Price       string
	Restriction string
}

func newFlow(line string) Flow {
	orig := line[2:6]
	dest := line[6:10]
	route := line[10:15]
	status := line[15:18]
	usage := line[18:19]
	direction := line[19:20]
	end := line[20:28]
	start := line[28:36]
	toc := line[36:39]
	crosslondon := line[39:40]
	discount := line[40:41]
	publication := line[41:42]
	flowid := line[42:49]
	f := Flow{
		orig,
		dest,
		route,
		status,
		usage,
		direction,
		end,
		start,
		toc,
		crosslondon,
		discount,
		publication,
		flowid,
	}
	return f
}

func newFare(line string) Fare {
	flowid := line[2:9]
	ticketcode := line[9:12]
	price := line[12:20]
	restriction := line[20:22]
	f := Fare{
		flowid,
		ticketcode,
		price,
		restriction,
	}
	return f
}

func parseFlowLines(lines []string) (map[string]Flow, []Fare) {
	flows := make(map[string]Flow)
	var fares []Fare
	for l := range lines {
		parsed := parseLine(lines[l])
		switch parsed.(type) {
		case nil:
			break
		case Flow:
			f := parsed.(Flow)
			flows[f.Flowid] = f
		case Fare:
			f := parsed.(Fare)
			fares = append(fares, f)
		}
	}
	fmt.Println("flows", len(flows))
	fmt.Println("fares", len(fares))

	return flows, fares
}

func mapFares(fares []Fare) map[string][]Fare {
	mappedFares := make(map[string][]Fare)
	for _, fare := range fares {
		_, exists := mappedFares[fare.Flowid]
		if !exists {
			mappedFares[fare.Flowid] = []Fare{fare}
		} else {
			mappedFares[fare.Flowid] = append(mappedFares[fare.Flowid], fare)
		}
	}
	return mappedFares
}

func MapFlows(flows map[string]Flow) (map[string][]Flow, map[string][]Flow) {
	origFlows := make(map[string][]Flow)
	destFlows := make(map[string][]Flow)
	for _, flow := range flows {
		_, origExists := origFlows[flow.Orig]
		if !origExists {
			origFlows[flow.Orig] = []Flow{flow}
		} else {
			origFlows[flow.Orig] = append(origFlows[flow.Orig], flow)
		}

		_, destExists := destFlows[flow.Dest]
		if !destExists {
			destFlows[flow.Dest] = []Flow{flow}
		} else {
			destFlows[flow.Dest] = append(destFlows[flow.Dest], flow)
		}
	}

	return origFlows, destFlows
}

func ReadFlows(filename string) (map[string]Flow, map[string][]Fare) {
	flowLines, _ := readFileLines(filename + "/" + filename + ".FFL")
	fmt.Println("FFL", len(flowLines))
	flows, fares := parseFlowLines(flowLines)
	return flows, mapFares(fares)
}
