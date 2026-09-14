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

export const formatPriceLine = (station, i, point, price) =>
  <div key={'' + point.point.nlc + i} style={{ display: 'flex', fontFamily: 'monospace'}}>
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
        href={`https://www.brfares.com/!fares?orig=${station}&dest=${point.point.Nlc}`}
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

const mapOrigPrices = (prices) => {
  if (isEmptyObject(prices)) {
    return {};
  }

  const pointPrices = {
    Station: prices.Orig,
    Points: Object.entries(prices.Dests).reduce((acc, [nlc, dest]) => {
      acc[nlc] = {
        Point: dest.Dest,
        Prices: dest.Prices
      };
      return acc;
    }, {})
  }

  console.log("pointPrices", pointPrices);
  return pointPrices;
}

const mapDestPrices = (prices) => {
  if (isEmptyObject(prices)) {
    return {};
  }

  const pointPrices = {
    Station: prices.Dest,
    Points: Object.entries(prices.Origs).reduce((acc, [nlc, orig]) => {
      acc[nlc] = {
        Point: orig.Orig,
        Prices: orig.Prices
      };
      return acc;
    }, {})
  }

  console.log("pointPrices", pointPrices);
  return pointPrices;
}

export const mapPrices = (direction, prices) => {
  switch (direction) {
    case 'orig':
      return mapOrigPrices(prices);
    case 'dest':
      return mapDestPrices(prices);
    default:
      console.error("direction", direction);
      return {};
  }
}

export const getFilteredPrices = (prices, ticketTypeFilter, crossLondonFilter, iRoute, xRoute, iTicket, xTicket) => {
  if (isEmptyObject(prices)) {
    return {};
  }

  const iRouteList = parseFilters(iRoute)
  const xRouteList = parseFilters(xRoute)
  const iTicketList = parseFilters(iTicket)
  const xTicketList = parseFilters(xTicket)

  console.log("filters", ticketTypeFilter, crossLondonFilter, iRouteList, xRouteList, iTicketList, xTicketList)

  const filteredPointPrices = { Station: prices.Station, Points: {} }
  Object.entries(prices.Points).forEach(([nlc, point]) => {
    const filteredPrices = []
    point.Prices.forEach(p => {
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
      filteredPointPrices.Points[nlc] = {
        Point: point.Point,
        Prices: filteredPrices
      };
    }
  })
  return filteredPointPrices;
}

export const getPossiblePrices = (prices) => {
  if (isEmptyObject(prices)) {
    return new Set();
  }
  const possiblePrices = new Set();
  Object.entries(prices.Points).forEach(([nlc, point]) => {
    if ('Prices' in point) {
      point.Prices.forEach(price => {
        possiblePrices.add(price.Price)
      })
    } else {
      console.log(point);
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
  Object.entries(prices.Points).forEach(([nlc, point]) => {
    if ('Prices' in point) {
      point.Prices.forEach(price => {
        ticketTypes.add(price.TicketCode)
      })
    } else {
      console.log(point);
    }
  })
  console.log("ticketTypes", ticketTypes);
  return Array.from(ticketTypes);
}

export const getPossiblePointPrices = (pointPrices) => {
  const possiblePointPrices = new Set();
  pointPrices.forEach(pointPrice => {
    possiblePointPrices.add(pointPrice.minPrice)
  });
  const sortedPossiblePointPrices = Array.from(possiblePointPrices).sort((a, b) => a - b);
  console.log("sortedPossiblePointPrices", sortedPossiblePointPrices);
  return sortedPossiblePointPrices;
}

export const getPointPrices = (prices) => {
  if (isEmptyObject(prices)) {
    return [];
  }
  const pointPrices = [];
  Object.entries(prices.Points).forEach(([nlc, point]) => {
    const pointPriceList = []
    point.Prices.forEach(price => {
      pointPriceList.push(price.Price);
    });
    const minPrice = Math.min.apply(Math, pointPriceList);
    pointPrices.push({
      point: point.Point,
      minPrice: minPrice,
      prices: point.Prices
    });
  });
  console.log("pointPrices", pointPrices)
  return pointPrices;
}

export const getPriceList = (prices) => {
  if (isEmptyObject(prices)) {
    return [];
  }
  const priceList = [];
  Object.entries(prices.Points).forEach(([nlc, point]) => {
    point.Prices.forEach(price => {
      priceList.push({
        point: point.Point,
        price: price
      });
    });
  });
  const sortedPriceList = priceList;
  sortedPriceList.sort((a, b) => a.price.Price - b.price.Price);
  console.log("priceList", priceList)
  return sortedPriceList;
}