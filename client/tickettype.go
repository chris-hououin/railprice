package client

import (
	"strings"
)

type TicketType struct {
	ticketCode          string
	end                 string
	start               string
	quote               string
	desc                string
	TicketClass         string
	TicketType          string
	ticketGroup         string
	lastValidDay        string
	maxPassengers       string
	minPassengers       string
	maxAdults           string
	minAdults           string
	maxChildren         string
	minChildren         string
	restrictedDate      string
	restrictedTrain     string
	restrictedArea      string
	validityCode        string
	AtbDesc             string
	lulLondonIssue      string
	reservationRequired string
	capriCode           string
	lul93               string
	utsCode             string
	timeRestriction     string
	freePassLul         string
	packageMkr          string
	fareMultiplier      string
	discountCategory    string
	IsAdvance           bool
	IsDay               bool
}

func newTicketType(line string, ticketAdvances []TicketAdvance, ticketValidityRecords map[string]TicketValidity) TicketType {
	ticketCode := line[1:4]
	end := line[4:12]
	start := line[12:20]
	quote := line[20:28]
	desc := line[28:43]
	ticketClass := line[43:44]
	ticketType := line[44:45]
	ticketGroup := line[45:46]
	lastValidDay := line[46:54]
	maxPassengers := line[54:57]
	minPassengers := line[57:60]
	maxAdults := line[60:63]
	minAdults := line[63:66]
	maxChildren := line[66:69]
	minChildren := line[69:72]
	restrictedDate := line[72:73]
	restrictedTrain := line[73:74]
	restrictedArea := line[74:75]
	validityCode := line[75:77]
	atbDesc := strings.TrimSpace(line[77:97])
	lulLondonIssue := line[97:98]
	reservationRequired := line[98:99]
	capriCode := line[99:102]
	lul93 := line[102:103]
	utsCode := line[103:105]
	timeRestriction := line[105:106]
	freePassLul := line[106:107]
	packageMkr := line[107:108]
	fareMultiplier := line[108:111]
	discountCategory := line[111:113]
	isAdvance := determineAdvance(ticketCode, atbDesc, reservationRequired, ticketAdvances)
	isDay := determineDay(validityCode, ticketValidityRecords)
	t := TicketType{
		ticketCode,
		end,
		start,
		quote,
		desc,
		ticketClass,
		ticketType,
		ticketGroup,
		lastValidDay,
		maxPassengers,
		minPassengers,
		maxAdults,
		minAdults,
		maxChildren,
		minChildren,
		restrictedDate,
		restrictedTrain,
		restrictedArea,
		validityCode,
		atbDesc,
		lulLondonIssue,
		reservationRequired,
		capriCode,
		lul93,
		utsCode,
		timeRestriction,
		freePassLul,
		packageMkr,
		fareMultiplier,
		discountCategory,
		isAdvance,
		isDay,
	}
	return t
}

func determineDay(validityCode string, ticketValidityRecords map[string]TicketValidity) bool {
	v := ticketValidityRecords[validityCode]
	return v.OutDays <= 1 && v.OutMonths == 0 && v.RetDays <= 1 && v.RetMonths == 0
}

func determineAdvance(ticketCode string, atbDesc string, reservationRequired string, ticketAdvances []TicketAdvance) bool {
	for _, ticketAdvance := range ticketAdvances {
		if ticketAdvance.ticketCode == ticketCode {
			return true
		}
		if reservationRequired != "N" {
			return true
		}
		if strings.Contains(strings.ToUpper(atbDesc), "ADVANCE") {
			return true
		}
		if strings.Contains(strings.ToUpper(atbDesc), "ADV ") {
			return true
		}
	}
	return false
}

func determineTicketTypeRecordType(line string) RecordType {
	if strings.HasPrefix(line, "/!!") {
		return CommentRecord
	}
	return TicketTypeRecord
}

func parseTicketTypeLine(line string, ticketAdvances []TicketAdvance, ticketValidityRecords map[string]TicketValidity) interface{} {
	recordType := determineTicketTypeRecordType(line)
	if recordType == TicketTypeRecord {
		return newTicketType(line, ticketAdvances, ticketValidityRecords)
	} else if recordType == CommentRecord {
		return nil
	}
	return nil
}

func parseTicketTypeLines(lines []string, ticketAdvances []TicketAdvance, ticketValidityRecords map[string]TicketValidity) map[string]TicketType {
	ticketTypes := make(map[string]TicketType)
	for l := range lines {
		parsed := parseTicketTypeLine(lines[l], ticketAdvances, ticketValidityRecords)
		switch parsed.(type) {
		case nil:
			break
		case TicketType:
			t := parsed.(TicketType)
			ticketTypes[t.ticketCode] = t
		}
	}

	return ticketTypes
}

func ReadTicketType(filename string, ticketAdvances []TicketAdvance, ticketValidityRecords map[string]TicketValidity) map[string]TicketType {
	ticketTypeLines, _ := readFileLines(filename + "/" + filename + ".TTY")
	ticketTypes := parseTicketTypeLines(ticketTypeLines, ticketAdvances, ticketValidityRecords)
	return ticketTypes
}
