import {ArrowRight, ArrowRightLeft, CalendarDays, ExternalLink, FastForward} from "lucide-react";

const isEmptyObject = obj => {
  return Object.keys(obj).length === 0
}

export const formatMoney = (minor) => {
  const pounds = Math.floor(minor / 100)
  const pence = minor % 100
  return `£${pounds}.${pence.toString().padStart(2, '0')}`
}

export const formatTicketType = (type) => {
  switch (type) {
    case 'S': return <ArrowRight size={14} />;
    case 'R': return <ArrowRightLeft size={14} />;
    case 'N': return <CalendarDays size={14} />;
  }
}

export const advanceIcon = <FastForward size="14" />;

export const formatPriceLine = (station, i, dest, price) =>
  <div key={'' + dest.dest.nlc + i} style={{ display: 'flex', fontFamily: 'monospace'}}>
    <div style={{ display: 'inline-flex', alignItems: 'center', marginRight: 'auto'}}>
      {formatTicketType(price.TicketType)}
      {price.Advance && advanceIcon}
      &nbsp;
      {price.TicketCode}
      {price.Restriction.trim() && `(${price.Restriction.trim()})`}
      &nbsp;
      {price.CrossLondon && '✠'}
      {price.Route}
    </div>
    <div style={{ display: 'inline-flex', alignItems: 'center', textAlign: 'right'}}>
      &nbsp;
      {formatMoney(price.Price)}
      <a
        href={`https://www.brfares.com/!fares?orig=${station}&dest=${dest.dest.Nlc}`}
        target={`_blank`}
      >
        <ExternalLink size={11} />
      </a>
    </div>
  </div>

const parseFilters = (s) => {
  if (s && s.trim()) {
    return s.split(',')
  }
  return []
}

export const getPriceTicketType = (price) => {
  if (price.Advance) {
    return 'A'
  } else if (price.TicketType === 'S') {
    return 'S';
  } else if (price.TicketType === 'R') {
    if (price.IsDay) {
      return 'D'
    } else {
      return 'P'
    }
  } else {
    return 'N'
  }
}

export const getFilteredPrices = (prices, ticketTypeFilter, crossLondonFilter, iRoute, xRoute, iTicket, xTicket) => {
  console.log("getFilteredPrices")
  if (isEmptyObject(prices)) {
    return {};
  }

  const iRouteList = parseFilters(iRoute)
  const xRouteList = parseFilters(xRoute)
  const iTicketList = parseFilters(iTicket)
  const xTicketList = parseFilters(xTicket)

  console.log("filters", ticketTypeFilter, crossLondonFilter, iRouteList, xRouteList, iTicketList, xTicketList)

  const filteredDestPrices = { Orig: prices.Orig, Dests: {} }
  Object.entries(prices.Dests).forEach(([nlc, dest]) => {
    const filteredPrices = []
    dest.Prices.forEach(p => {
      if (!ticketTypeFilter.includes(getPriceTicketType(p))) {
        // filter out
      } else if (crossLondonFilter && p.CrossLondon !== true) {
        // filter out
      } else if (iRouteList.length > 0 && !iRouteList.includes(p.RouteCode)) {
        // filter out
      } else if (xRouteList.length > 0 && xRouteList.includes(p.RouteCode)) {
        // filter out
      } else if (iTicketList.length > 0 && !iTicketList.includes(p.TicketCode)) {
        // filter out
      } else if (xTicketList.length > 0 && xTicketList.includes(p.TicketCode)) {
        // filter out
      } else {
        filteredPrices.push(p)
      }
    })

    if (filteredPrices.length > 0) {
      filteredDestPrices.Dests[nlc] = {
        Dest: dest.Dest,
        Prices: filteredPrices
      };
    }
  })
  return filteredDestPrices;
}

export const getPossiblePrices = (prices) => {
  if (isEmptyObject(prices)) {
    return new Set();
  }
  const possiblePrices = new Set();
  Object.entries(prices.Dests).forEach(([nlc, dest]) => {
    if ('Prices' in dest) {
      dest.Prices.forEach(price => {
        possiblePrices.add(price.Price)
      })
    } else {
      console.log(dest);
    }
  })
  const sortedPossiblePrices = Array.from(possiblePrices).sort((a, b) => a - b);
  console.log("sortedPossiblePrices", sortedPossiblePrices);
  return sortedPossiblePrices;
}

export const getPossibleTicketTypes = (prices) => {
  if (isEmptyObject(prices)) {
    return new Set();
  }
  const ticketTypes = new Set();
  Object.entries(prices.Dests).forEach(([nlc, dest]) => {
    if ('Prices' in dest) {
      dest.Prices.forEach(price => {
        ticketTypes.add(price.TicketCode)
      })
    } else {
      console.log(dest);
    }
  })
  console.log("ticketTypes", ticketTypes);
  return Array.from(ticketTypes);
}

export const getPossibleDestPrices = (destPrices) => {
  const possibleDestPrices = new Set();
  destPrices.forEach(destPrice => {
    possibleDestPrices.add(destPrice.minPrice)
  });
  const sortedPossibleDestPrices = Array.from(possibleDestPrices).sort((a, b) => a - b);
  console.log("sortedPossibleDestPrices", sortedPossibleDestPrices);
  return sortedPossibleDestPrices;
}

export const getDestPrices = (prices) => {
  if (isEmptyObject(prices)) {
    return [];
  }
  const destPrices = [];
  Object.entries(prices.Dests).forEach(([nlc, dest]) => {
    const destPriceList = []
    dest.Prices.forEach(price => {
      destPriceList.push(price.Price);
    });
    const minPrice = Math.min.apply(Math, destPriceList);
    destPrices.push({
      dest: dest.Dest,
      minPrice: minPrice,
      prices: dest.Prices
    });
  });
  console.log("destPrices", destPrices)
  return destPrices;
}

export const getPriceList = (prices) => {
  if (isEmptyObject(prices)) {
    return [];
  }
  const priceList = [];
  Object.entries(prices.Dests).forEach(([nlc, dest]) => {
    dest.Prices.forEach(price => {
      priceList.push({
        dest: dest.Dest,
        price: price
      });
    });
  });
  const sortedPriceList = priceList;
  sortedPriceList.sort((a, b) => a.price.Price - b.price.Price);
  console.log("priceList", priceList)
  return sortedPriceList;
}