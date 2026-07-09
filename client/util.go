package client

import (
	"os"
	"strings"
	"time"
)

type RecordType int

const (
	CommentRecord RecordType = iota
	FlowRecord
	FareRecord
	ClusterRecord
	LocationRecord
	RouteRecord
	TicketTypeRecord
	TicketAdvanceRecord
	TicketValidityRecord
)

func readFileLines(path string) ([]string, error) {
	dat, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(dat), "\n")
	return lines, nil
}

func FormatDate(date string) string {
	t, _ := time.Parse("02012006", date)
	return t.Format("20060102")
}
