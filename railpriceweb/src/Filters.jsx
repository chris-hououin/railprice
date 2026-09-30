import { Button, ToggleButton } from 'react-aria-components'
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
      <div className="ticket-type-tabs" aria-label="Ticket type filters">
        {ticketTypes.map(type => (
          <ToggleButton
            key={type.key}
            className="ticket-tab"
            isSelected={ticketTypeFilter.includes(type.key)}
            onChange={() => toggleTicketType(type.key)}
          >
            {type.label}
          </ToggleButton>
        ))}
      </div>

      <Button
        className={`cross-london-btn ${crossLondonFilter ? 'active' : ''}`}
        onPress={() => setCrossLondonFilter(!crossLondonFilter)}
        aria-label="Toggle cross London routes"
        title="Cross London"
      >
        ✠
      </Button>
    </div>
  )
}
