/* eslint-disable react/prop-types */
import { useEffect, useMemo, useRef, useState } from 'react'
import { Input } from 'react-aria-components'
import './SearchBar.css'

export default function SearchBar({ stations, selectedStation, setSelectedStation }) {
  const [searchQuery, setSearchQuery] = useState('')
  const [showDropdown, setShowDropdown] = useState(false)
  const dropdownRef = useRef(null)

  useEffect(() => {
    const selected = stations.find(station => station.Nlc === selectedStation)
    setSearchQuery(selected ? selected.Name : '')
  }, [stations, selectedStation])

  const filteredStations = useMemo(() => {
    const query = searchQuery.trim().toLowerCase()
    let matches = !query ? stations : stations.filter(station => station.SearchText?.includes(query))
    return matches.sort((a, b) => {
      if (a.Nlc === selectedStation) return -1
      if (b.Nlc === selectedStation) return 1
      return 0
    })
  }, [searchQuery, selectedStation, stations])

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setShowDropdown(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const handleSelectStation = (station) => {
    setSelectedStation(station.Nlc)
    setSearchQuery(station.Name)
    setShowDropdown(false)
  }

  return (
    <div className="search-bar" ref={dropdownRef}>
      <div className="search-field">
        <Input
          className="search-input"
          placeholder="Search station..."
          value={searchQuery}
          onChange={(event) => {
            setSearchQuery(event.target.value)
            setShowDropdown(true)
          }}
          onFocus={() => setShowDropdown(true)}
          aria-label="Search station"
        />
      </div>
      {showDropdown && (
        <div className="dropdown-list" role="listbox" aria-label="Station results">
          {filteredStations.length > 0 ? (
            filteredStations.map(station => (
              <div
                key={station.Nlc}
                className={`dropdown-item ${selectedStation === station.Nlc ? 'selected' : ''}`}
                onMouseDown={(event) => {
                  event.preventDefault()
                  handleSelectStation(station)
                }}
                role="option"
                aria-selected={selectedStation === station.Nlc}
              >
                {station.Nlc && <span className="station-nlc">{station.Nlc}</span>}
                <span className="station-name">{station.Name}</span>
                {station.Crs && <span className="station-crs">{station.Crs}</span>}
              </div>
            ))
          ) : (
            <div className="dropdown-empty">No stations found</div>
          )}
        </div>
      )}
    </div>
  )
}

