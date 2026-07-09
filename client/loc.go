package client

import (
	"fmt"
	"strings"
)

type Location struct {
	UicCode           string
	end               string
	start             string
	quote             string
	adminArea         string
	Nlc               string
	desc              string
	Crs               string
	resv              string
	ersCountry        string
	ersCode           string
	FareGroup         string
	county            string
	pte               string
	zoneNo            string
	zoneInd           string
	region            string
	hierarchy         string
	CcDescOut         string
	ccDescRtn         string
	atbDescOut        string
	atbDescRtn        string
	specialFacilities string
	lulDirection      string
	lulUtsMode        string
	lulZone1          string
	lulZone2          string
	lulZone3          string
	lulZone4          string
	lulZone5          string
	lulZone6          string
	lulUtsLondonStn   string
	utsCode           string
	utsACode          string
	utsPtrBias        string
	utsOffset         string
	utsNorth          string
	utsEast           string
	utsSouth          string
	UtsWest           string
}

func newLoc(line string) Location {
	uicCode := line[2:9]
	end := line[9:17]
	start := line[17:25]
	quote := line[25:33]
	adminArea := line[33:36]
	nlc := line[36:40]
	desc := line[40:56]
	crs := line[56:59]
	resv := line[59:64]
	ersCountry := line[64:66]
	ersCode := line[66:69]
	fareGroup := strings.TrimSpace(line[69:75])
	county := line[75:77]
	pte := line[77:79]
	zoneNo := line[79:83]
	zoneInd := line[83:85]
	region := line[85:86]
	hierarchy := line[86:87]
	ccDescOut := line[87:128]
	ccDescRtn := line[128:144]
	atbDescOut := line[144:204]
	atbDescRtn := line[204:234]
	specialFacilities := line[234:260]
	lulDirection := line[260:261]
	lulUtsMode := line[261:262]
	lulZone1 := line[262:263]
	lulZone2 := line[263:264]
	lulZone3 := line[264:265]
	lulZone4 := line[265:266]
	lulZone5 := line[266:267]
	lulZone6 := line[267:268]
	lulUtsLondonStn := line[268:269]
	utsCode := line[269:272]
	utsACode := line[272:275]
	utsPtrBias := line[275:276]
	utsOffset := line[276:277]
	utsNorth := line[277:280]
	utsEast := line[280:283]
	utsSouth := line[283:286]
	utsWest := line[286:289]
	l := Location{
		uicCode,
		end,
		start,
		quote,
		adminArea,
		nlc,
		desc,
		crs,
		resv,
		ersCountry,
		ersCode,
		fareGroup,
		county,
		pte,
		zoneNo,
		zoneInd,
		region,
		hierarchy,
		ccDescOut,
		ccDescRtn,
		atbDescOut,
		atbDescRtn,
		specialFacilities,
		lulDirection,
		lulUtsMode,
		lulZone1,
		lulZone2,
		lulZone3,
		lulZone4,
		lulZone5,
		lulZone6,
		lulUtsLondonStn,
		utsCode,
		utsACode,
		utsPtrBias,
		utsOffset,
		utsNorth,
		utsEast,
		utsSouth,
		utsWest,
	}
	return l
}

func determineLocationRecordType(line string) RecordType {
	if strings.HasPrefix(line, "/!!") {
		return CommentRecord
	} else if strings.HasPrefix(line, "RL") {
		return LocationRecord
	}

	return CommentRecord
}

func parseLocationLine(line string) interface{} {
	recordType := determineLocationRecordType(line)
	if recordType == LocationRecord {
		return newLoc(line)
	} else if recordType == CommentRecord {
		return nil
	}
	return nil
}
func parseLocationLines(lines []string) map[string]Location {
	locations := make(map[string]Location)
	for l := range lines {
		parsed := parseLocationLine(lines[l])
		switch parsed.(type) {
		case nil:
			break
		case Location:
			l := parsed.(Location)
			locations[l.Nlc] = l
		}
	}
	fmt.Println("locations", len(locations))
	return locations
}

func ReadLocations(filename string) map[string]Location {
	locationLines, _ := readFileLines(filename + "/" + filename + ".LOC")
	locations := parseLocationLines(locationLines)
	return locations
}

func MapStationGroups(locations map[string]Location) map[string][]string {
	mappedStationGroups := make(map[string][]string)
	for _, loc := range locations {
		_, exists := mappedStationGroups[loc.FareGroup]
		if !exists {
			mappedStationGroups[loc.FareGroup] = []string{loc.Nlc}
		} else {
			mappedStationGroups[loc.FareGroup] = append(mappedStationGroups[loc.FareGroup], loc.Nlc)
		}
	}
	return mappedStationGroups
}

func FetchStationsFromGroup(groupNlc string, locationGroup map[string][]string) []string {
	_, exists := locationGroup[groupNlc]

	if exists {
		return locationGroup[groupNlc]
	} else {
		// It is a station NLC in a group
		return []string{groupNlc}
	}
}
