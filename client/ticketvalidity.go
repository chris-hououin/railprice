package client

import (
	"strconv"
	"strings"
)

type TicketValidity struct {
	validityCode   string
	endDate        string
	startDate      string
	desc           string
	OutDays        int
	OutMonths      int
	RetDays        int
	RetMonths      int
	retAfterDays   string
	retAfterMonths string
	retAfterDay    string
	breakOut       string
	breakRtn       string
	outDesc        string
	rtnDesc        string
}

func newTicketValidity(line string) TicketValidity {
	validityCode := line[0:2]
	endDate := line[2:10]
	startDate := line[10:18]
	desc := line[18:38]
	outDays, _ := strconv.Atoi(line[38:40])
	outMonths, _ := strconv.Atoi(line[40:42])
	retDays, _ := strconv.Atoi(line[42:44])
	retMonths, _ := strconv.Atoi(line[44:46])
	retAfterDays := line[46:48]
	retAfterMonths := line[48:50]
	retAfterDay := line[50:52]
	breakOut := line[52:53]
	breakRtn := line[53:54]
	outDesc := line[54:68]
	rtnDesc := line[54:68]
	t := TicketValidity{
		validityCode,
		endDate,
		startDate,
		desc,
		outDays,
		outMonths,
		retDays,
		retMonths,
		retAfterDays,
		retAfterDay,
		retAfterMonths,
		breakOut,
		breakRtn,
		outDesc,
		rtnDesc,
	}
	return t
}

func determineTicketValidityRecord(line string) RecordType {
	if strings.HasPrefix(line, "/!!") {
		return CommentRecord
	}
	return TicketValidityRecord
}

func parseTicketValidityRecordLine(line string) interface{} {
	recordType := determineTicketValidityRecord(line)
	if recordType == TicketValidityRecord {
		return newTicketValidity(line)
	} else if recordType == CommentRecord {
		return nil
	}
	return nil
}

func parseTicketValidityRecordLines(lines []string) map[string]TicketValidity {
	ticketValidityRecords := make(map[string]TicketValidity)
	for l := range lines {
		parsed := parseTicketValidityRecordLine(lines[l])
		switch parsed.(type) {
		case nil:
			break
		case TicketValidity:
			t := parsed.(TicketValidity)
			ticketValidityRecords[t.validityCode] = t
		}
	}

	return ticketValidityRecords
}

func ReadTicketValidity(filename string) map[string]TicketValidity {
	ticketValidityLines, _ := readFileLines(filename + "/" + filename + ".TVL")
	ticketValidityRecords := parseTicketValidityRecordLines(ticketValidityLines)
	return ticketValidityRecords
}
