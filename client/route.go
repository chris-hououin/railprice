package client

import (
	"strings"
)

type Route struct {
	routeCode string
	end       string
	start     string
	quote     string
	desc      string
	AtbDesc   string
}

func newRoute(line string) Route {
	routeCode := line[2:7]
	end := line[7:15]
	start := line[15:23]
	quote := line[23:31]
	desc := line[31:47]
	atbDesc := line[47:82]
	r := Route{
		routeCode,
		end,
		start,
		quote,
		desc,
		atbDesc,
	}
	return r
}

func determineRouteRecordType(line string) RecordType {
	if strings.HasPrefix(line, "/!!") {
		return CommentRecord
	} else if strings.HasPrefix(line, "RR") {
		return RouteRecord
	}
	return CommentRecord
}

func parseRouteLine(line string) interface{} {
	recordType := determineRouteRecordType(line)
	if recordType == RouteRecord {
		return newRoute(line)
	} else if recordType == CommentRecord {
		return nil
	}
	return nil
}

func parseRouteLines(lines []string) map[string]Route {
	routes := make(map[string]Route)
	for l := range lines {
		parsed := parseRouteLine(lines[l])
		switch parsed.(type) {
		case nil:
			break
		case Route:
			r := parsed.(Route)
			routes[r.routeCode] = r
		}
	}

	return routes
}

func ReadRoutes(filename string) map[string]Route {
	routeLines, _ := readFileLines(filename + "/" + filename + ".RTE")
	routes := parseRouteLines(routeLines)
	return routes
}
