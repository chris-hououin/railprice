/* eslint-disable react/prop-types */
import { useState, useRef, useEffect } from 'react'
import './SearchBar.css'


export default function SearchBar({ stations, selectedStation, setSelectedStation }) {
  const [searchQuery, setSearchQuery] = useState('')
  const [showDropdown, setShowDropdown] = useState(false)
  const dropdownRef = useRef(null)

  // Filter stations based on search query using SearchText
  const filteredStations = stations.filter(station => station.SearchText?.includes(searchQuery.trim().toLowerCase()))

  // Handle station selection
  const handleSelectStation = (station) => {
    setSelectedStation(station.Nlc)
    setShowDropdown(false)
  }

  useEffect(() => {
    const getStationByNlc = (nlc) => {
      return stations.filter(station => station.Nlc === nlc)[0] || null
    }

    const s = getStationByNlc(selectedStation)
    if (s) {
      setSearchQuery(getStationByNlc(selectedStation).Name)
    }
  }, [stations, selectedStation]);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setShowDropdown(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  return (
    <div className="search-bar" ref={dropdownRef}>
      <input
        type="text"
        placeholder="Search station..."
        value={searchQuery}
        onChange={(e) => setSearchQuery(e.target.value)}
        onFocus={() => setShowDropdown(true)}
        className="search-input"
      />
      {showDropdown && (
        <div className="dropdown-list">
          {filteredStations.length > 0 ? (
            filteredStations.map(station => (
              <div
                key={station.Nlc}
                className="dropdown-item"
                onClick={() => handleSelectStation(station)}
                style={{
                  backgroundColor: selectedStation === station.Nlc ? '#e8f4f8' : 'white',
                  fontWeight: selectedStation === station.Nlc ? 'bold' : 'normal'
                }}
              >
                {station.Nlc && <span className="station-nlc">{station.Nlc}</span>}
                <span className="station-name">{station.Name}</span>
                {station.Crs && <span className="station-crs">{station.Crs}</span>}
              </div>
            ))) : (<div className="dropdown-empty">No stations found</div>)
          }
        </div>
      )}
    </div>
  )
}

