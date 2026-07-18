import {MapContainer, Marker, Popup, TileLayer, useMap} from 'react-leaflet'
import L from 'leaflet'
import {useEffect, useRef, useState} from 'react'
import {useQuery} from '@tanstack/react-query'
import './App.css'
import {
  formatMoney,
  formatTicketType,
  getDestPrices,
  getFilteredPrices,
  getPossibleDestPrices,
  getPossiblePrices,
  getPriceList
} from "./utils.jsx";
import {getMarkerColor2} from "./GetMarkerColor.jsx";
import PriceList from './PriceList.jsx'
import SearchBar from './SearchBar.jsx'

const ticketTypes = [
  { key: 'A', label: 'Advance' },
  { key: 'S', label: 'Single' },
  { key: 'D', label: 'Day Return' },
  { key: 'P', label: 'Period Return' },
  { key: 'N', label: 'Season' }
]

// Function to create colored marker icon
const createColoredMarker = (color) => {
  return L.divIcon({
    html: `<div style="background-color: ${color}; width: 20px; height: 20px; border-radius: 50%; border: 2px solid white; display: flex; align-items: center; justify-content: center;"></div>`,
    iconSize: [20, 20],
    iconAnchor: [10, 10],
    popupAnchor: [0, -10],
    className: 'map-marker-icon'
  })
}

// Custom component to capture map instance
/* eslint-disable react/prop-types */
function MapController({ mapRef }) {
  const map = useMap()

  useEffect(() => {
    mapRef.current = map
  }, [map, mapRef])

  useEffect(() => {
    const text = `Powered by National Rail Enquiries. This website is a work of fiction. <a href="https://github.com/chris-hououin/railprice" target="_blank">Build ${document.querySelector('meta[name="build-number"]')?.getAttribute('content')}</a>`
    map?.attributionControl?.addAttribution(text)
  }, [])

  return null
}
/* eslint-enable react/prop-types */

function App() {
  const position = [51.505, -0.09] // Default center (London)
  const mapRef = useRef(null)
  const markerRefs = useRef({})

  const [selectedStation, setSelectedStation] = useState('7022')
  const [ticketTypeFilter, setTicketTypeFilter] = useState(['P'])
  const [crossLondonFilter, setCrossLondonFilter] = useState(false)


  const toggleTicketType = (key) => {
    setTicketTypeFilter(prev =>
      prev.includes(key) ? prev.filter(t => t !== key) : [...prev, key]
    )
  }

  const { data: stations = [] } = useQuery({
    queryKey: ['stations'],
    queryFn: async () => {
      const response = await fetch('stations.json')
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      return response.json().then(ss => ss.map(s => ({...s, SearchText: (s.Nlc + s.Crs + s.Name).toLowerCase()})))
    }
  })

  const { data: prices = {}, isLoading: loading, error } = useQuery({
    queryKey: ['prices' + selectedStation],
    queryFn: async () => {
      const response = await fetch(`orig/${selectedStation}.json.br`)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      if (response.headers.get('Content-Encoding') === 'br') {
        return response.json()
      } else {
        const compressedData = await response.blob();

        const ds = new DecompressionStream('brotli');
        const decompressedStream = compressedData.stream().pipeThrough(ds);
        const decompressedResponse = new Response(decompressedStream);

        return decompressedResponse.json();
      }
    }
  })

  console.log("prices", prices);

  const filteredPrices = getFilteredPrices(prices, ticketTypeFilter, crossLondonFilter);
  const possiblePrices = getPossiblePrices(filteredPrices);
  const destPrices = getDestPrices(filteredPrices);
  const possibleDestPrices = getPossibleDestPrices(destPrices);
  const priceList = getPriceList(filteredPrices);

  return (
    <div className="app-container">
      {error && <div className="error-toast">{error}</div>}
      {loading && <div className="loading-toast">Loading prices...</div>}
      <SearchBar
        stations={stations}
        selectedStation={selectedStation}
        setSelectedStation={setSelectedStation}
      />
      <div className="filters">
        <div className="ticket-type-tabs">
          {ticketTypes.map(type => (
            <button
              key={type.key}
              className={`ticket-tab ${ticketTypeFilter.includes(type.key) ? 'active' : ''}`}
              onClick={() => toggleTicketType(type.key)}
            >
              {type.label}
            </button>
          ))}
        </div>

        <button
          className={`cross-london-btn ${crossLondonFilter ? 'active' : ''}`}
          onClick={() => setCrossLondonFilter(!crossLondonFilter)}
          title="Cross London"
        >
          ✠
        </button>
      </div>

      <MapContainer center={position} zoom={13} className="map">
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright" target="_blank">OpenStreetMap</a> contributors'
          url="https://tile.openstreetmap.org/{z}/{x}/{y}.png"
          opacity={0.4}
        />
        <TileLayer
          attribution='&copy; <a href="https://www.openrailwaymap.org/" target="_blank">OpenRailwayMap</a>'
          url="https://tiles.openrailwaymap.org/standard/{z}/{x}/{y}.png"
          opacity={0.6}
        />
        {destPrices.map(
          dest => (
            !(dest.dest.Lat === 0 && dest.dest.Long === 0) &&
              <Marker
                key={dest.dest.Nlc}
                position={[dest.dest.Lat, dest.dest.Long]}
                icon={createColoredMarker(getMarkerColor2(dest.minPrice, possibleDestPrices))}
                ref={m => { markerRefs.current[dest.dest.Nlc] = m}}
              >
                <Popup>
                  <div><a style={{ cursor: "pointer"}} onClick={() => setSelectedStation(dest.dest.Nlc)}><strong>{dest.dest.Name}</strong></a></div>
                  {dest.prices.map((price, i) => (
                    <div key={'' + dest.dest.nlc + i}>
                      {price.Advance ? '🅰' : ''}
                      {formatTicketType(price.TicketType)}
                      {price.TicketCode}
                      {formatMoney(price.Price)}
                      {price.Route} {price.RouteCode}
                      {price.CrossLondon ? '✠' : ''}
                    </div>
                  ))}
                </Popup>
              </Marker>
          )
        )}
        <MapController mapRef={mapRef} />
      </MapContainer>

      <PriceList
        priceList={priceList}
        possiblePrices={possiblePrices}
        mapRef={mapRef}
        markerRefs={markerRefs}
      />

    </div>
  )
}

export default App
