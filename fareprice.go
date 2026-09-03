package main

import (
	"railprice/client"
	"strconv"
	"strings"
)

type FarePrice struct {
	Price       int
	CrossLondon bool
	RouteCode   string
	Route       string
	TicketCode  string
	TicketDesc  string
	Advance     bool
	IsDay       bool
	TicketClass string
	TicketType  string
	Restriction string
}

func newFarePrice(f client.Fare) FarePrice {
	priceMinorUnit, err := strconv.Atoi(f.Price)
	if err != nil {
		panic(err)
	}

	return FarePrice{
		Price:       priceMinorUnit,
		CrossLondon: flows[f.Flowid].Crosslondon == "1",
		RouteCode:   flows[f.Flowid].Route,
		Route:       strings.TrimSpace(routes[flows[f.Flowid].Route].AtbDesc),
		TicketCode:  f.Ticketcode,
		TicketDesc:  ticketTypes[f.Ticketcode].AtbDesc,
		Advance:     ticketTypes[f.Ticketcode].IsAdvance,
		IsDay:       ticketTypes[f.Ticketcode].IsDay,
		TicketClass: ticketTypes[f.Ticketcode].TicketClass,
		TicketType:  ticketTypes[f.Ticketcode].TicketType,
		Restriction: f.Restriction,
	}
}
