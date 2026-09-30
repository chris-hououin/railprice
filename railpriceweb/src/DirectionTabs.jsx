import './DirectionTabs.css'

export default function DirectionTabs({ isDestination, onChange }) {
  return (
    <div className="direction-tabs" aria-label="Direction selector">
      <button
        type="button"
        className={`direction-tab ${!isDestination ? 'selected' : ''}`}
        onClick={() => onChange(false)}
      >
        Origin
      </button>
      <button
        type="button"
        className={`direction-tab ${isDestination ? 'selected' : ''}`}
        onClick={() => onChange(true)}
      >
        Destination
      </button>
    </div>
  )
}
