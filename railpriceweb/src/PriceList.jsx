import {advanceIcon, formatMoney, formatTicketType} from './utils.jsx'
import {useEffect, useState} from "react";
import './PriceList.css'

export default function PriceList({ priceList, possiblePrices, mapRef, markerRefs, showPriceList }) {

  const [currentPage, setCurrentPage] = useState(1)
  const itemsPerPage = 500

  // Pagination for price list
  const totalPages = Math.max(1, Math.ceil(priceList.length / itemsPerPage))
  const startIndex = (currentPage - 1) * itemsPerPage
  const endIndex = Math.min(priceList.length, startIndex + itemsPerPage)
  const paginatedPriceList = priceList.slice(startIndex, endIndex)

  useEffect(() => {
    // if the data length changes reset to page 1 when current page becomes invalid
    if (currentPage > totalPages) {
      setCurrentPage(1)
    }
  }, [priceList.length, totalPages, currentPage])


  // Handle list item click
  const handleListItemClick = (price, idx) => {
    if (mapRef.current && (price.point.Lat || price.point.Long)) {
      // Center map on the marker
      mapRef.current.setView([price.point.Lat, price.point.Long], 12)

      // Open popup
      setTimeout(() => {
        const marker = markerRefs.current[idx]
        if (marker) {
          marker.openPopup()
        }
      }, 100)
    }
  }

  if (!showPriceList) {
    return null
  }

  return (
    <div className="price-list">
      <div className="price-list-title">
        <div>Prices ({priceList.length})</div>
        <div className="title-pager">
          <button className="pager-btn" onClick={() => setCurrentPage(p => Math.max(1, p - 1))} disabled={currentPage === 1} aria-label="Previous page">
            ‹
          </button>
          <span className="pagination-info">{currentPage} / {totalPages}</span>
          <button className="pager-btn" onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))} disabled={currentPage === totalPages} aria-label="Next page">
            ›
          </button>
        </div>
      </div>
      <div className="price-list-content">
        {paginatedPriceList.map((price, idx) => (
          <div key={startIndex + idx} className="price-item" onClick={() => handleListItemClick(price, price.point.Nlc)}>
            <div className="price-item-details">
              <div style={{display: 'flex'}}>
                <div className="price-item-desc" style={{ display: 'inline-flex', alignItems: 'center', marginRight: 'auto' }}>
                  {formatTicketType(price.price.TicketType)}
                  {price.price.Advance && advanceIcon}
                  &nbsp;
                  {price.price.TicketCode}
                  {' '}
                  {price.price.TicketDesc}
                </div>
                <div>{formatMoney(price.price.Price)}</div>

              </div>
              <div className="price-item-name">
                {price.point.Name} {price.point.Crs.trim() && `(${price.point.Crs})`}
              </div>
              <div className="price-item-route">
                {price.price.RouteCode} {price.price.CrossLondon && '✠'} {price.price.Route}
                {price.price.Restriction.trim() && ` (${price.price.Restriction})`}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}


