// Function to interpolate color using rank distribution
export const getMarkerColor = (price, possiblePrices) => {
  if (!possiblePrices || possiblePrices.length === 0) return 'rgb(128, 128, 128)' // Gray fallback

  // Find the rank (position) of this price in the sorted list
  const rank = possiblePrices.findIndex(p => p === price)
  const ratio = possiblePrices.length > 1 ? rank / (possiblePrices.length - 1) : 0.5
  const clampedRatio = Math.max(0, Math.min(1, ratio))

  // Rainbow gradient: Red -> Orange -> Yellow -> Green -> Blue -> Indigo -> Violet
  let r, g, b

  if (clampedRatio < 1 / 7) {
    // Red to Orange
    const localRatio = clampedRatio / (1 / 7)
    r = 255
    g = Math.round(0 + 165 * localRatio)
    b = 0
  } else if (clampedRatio < 2 / 7) {
    // Orange to Yellow
    const localRatio = (clampedRatio - 1 / 7) / (1 / 7)
    r = 255
    g = Math.round(165 + 90 * localRatio)
    b = 0
  } else if (clampedRatio < 3 / 7) {
    // Yellow to Green
    const localRatio = (clampedRatio - 2 / 7) / (1 / 7)
    r = Math.round(255 - 255 * localRatio)
    g = 255
    b = 0
  } else if (clampedRatio < 4 / 7) {
    // Green to Blue
    const localRatio = (clampedRatio - 3 / 7) / (1 / 7)
    r = 0
    g = 255
    b = Math.round(0 + 255 * localRatio)
  } else if (clampedRatio < 5 / 7) {
    // Blue to Indigo
    const localRatio = (clampedRatio - 4 / 7) / (1 / 7)
    r = Math.round(0 + 75 * localRatio)
    g = Math.round(255 - 100 * localRatio)
    b = 255
  } else if (clampedRatio < 6 / 7) {
    // Indigo to Violet
    const localRatio = (clampedRatio - 5 / 7) / (1 / 7)
    r = Math.round(75 + 180 * localRatio)
    g = Math.round(155 - 155 * localRatio)
    b = Math.round(255 - 100 * localRatio)
  } else {
    // Violet
    r = 255
    g = 0
    b = 155
  }

  return `rgb(${r}, ${g}, ${b})`
}

export const getMarkerColor2 = (price, possiblePrices) => {
  const rank = possiblePrices.findIndex(p => p === price)
  const total = possiblePrices.length

  const hPart = 12
  const hSize = total / hPart
  const sv = (rank % hSize) / hSize;
  const svSplit = 0.8
  const h = Math.floor(rank / hSize) * (360 / hPart)
  const s = sv < svSplit ? 0.2 + sv / svSplit * 0.8 : 1
  const v = sv > svSplit ? 0.8 + ((1 - ((sv - svSplit) / (1 - svSplit))) * 0.2) : 1

  if (price === 4400) {
    console.log("four", rank, total, hPart, sv, h,s,v)
  }

  return hsv2rgb(h,s,v)
}

const hsv2rgb = (h,s,v) => {
  let f= (n,k=(n+h/60)%6) => v - v*s*Math.max( Math.min(k,4-k,1), 0);
  return `rgb(${f(5)*255}, ${f(3)*255}, ${f(1)*255})`;
}