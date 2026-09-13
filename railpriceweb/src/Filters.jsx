import './Filters.css'

const ticketTypes = [
  { key: 'A', label: 'Advance' },
  { key: 'S', label: 'Single' },
  { key: 'D', label: 'Day Return' },
  { key: 'P', label: 'Period Return' },
  { key: 'N', label: 'Season' }
]

export default function Filters({
  ticketTypeFilter,
  setTicketTypeFilter,
  crossLondonFilter,
  setCrossLondonFilter
}) {
  const toggleTicketType = (key) => {
    setTicketTypeFilter(prev =>
      prev.includes(key) ? prev.filter(t => t !== key) : [...prev, key]
    )
  }
  return (
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
  )
}
