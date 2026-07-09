package client

import (
	"strings"
)

type TicketAdvance struct {
	ticketCode      string
	restrictionCode string
	restrictionFlag string
	tocId           string
	endDate         string
	start           string
	checkType       string
	apData          string
	bookingTime     string
}

func newTicketAdvance(line string) TicketAdvance {
	ticketCode := line[0:3]
	restrictionCode := line[3:5]
	restrictionFlag := line[5:6]
	tocId := line[6:8]
	endDate := line[8:16]
	start := line[16:24]
	checkType := line[24:25]
	apData := line[25:33]
	bookingTime := line[33:37]
	t := TicketAdvance{
		ticketCode,
		restrictionCode,
		restrictionFlag,
		tocId,
		endDate,
		start,
		checkType,
		apData,
		bookingTime,
	}
	return t
}

func determineTicketAdvanceRecordType(line string) RecordType {
	if strings.HasPrefix(line, "/!!") {
		return CommentRecord
	}
	return TicketAdvanceRecord
}

func parseTicketAdvanceLine(line string) interface{} {
	recordType := determineTicketAdvanceRecordType(line)
	if recordType == TicketAdvanceRecord {
		return newTicketAdvance(line)
	} else if recordType == CommentRecord {
		return nil
	}
	return nil
}

func parseTicketAdvanceLines(lines []string) []TicketAdvance {
	var ticketAdvances []TicketAdvance
	for l := range lines {
		parsed := parseTicketAdvanceLine(lines[l])
		switch parsed.(type) {
		case nil:
			break
		case TicketAdvance:
			t := parsed.(TicketAdvance)
			ticketAdvances = append(ticketAdvances, t)
		}
	}

	return ticketAdvances
}

func ReadTicketAdvance(filename string) []TicketAdvance {
	ticketAdvanceLines, _ := readFileLines(filename + "/" + filename + ".TAP")
	ticketAdvances := parseTicketAdvanceLines(ticketAdvanceLines)
	return ticketAdvances
}
