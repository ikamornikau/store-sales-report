import { describe, expect, it } from 'vitest'

import { addDays, daysBetween, defaultRange, formatMoney, validateRange } from '../format'

describe('date range', () => {
  it('allows a span of 92 days and rejects a longer one', () => {
    const from = '2026-01-01'
    expect(validateRange(from, addDays(from, 92))).toBeNull()
    expect(daysBetween(from, addDays(from, 92))).toBe(92)
    expect(validateRange(from, addDays(from, 93))).toBe('span')
  })

  it('rejects a reversed range and dates that are not real calendar days', () => {
    expect(validateRange('2026-04-03', '2026-01-01')).toBe('order')
    expect(validateRange('', '2026-01-01')).toBe('missing')
    expect(validateRange('2026-02-31', '2026-03-01')).toBe('missing')
  })

  it('defaults to a 30 day window ending today', () => {
    const range = defaultRange(new Date(2026, 10, 11))
    expect(range).toEqual({ from: '2026-10-13', to: '2026-11-11' })
    expect(validateRange(range.from, range.to)).toBeNull()
  })
})

describe('formatMoney', () => {
  it('formats cents in the store currency', () => {
    expect(formatMoney(154099, 'USD', 'en')).toBe('$1,540.99')
    expect(formatMoney(0, 'CAD', 'en')).toBe('CA$0.00')
  })
})
