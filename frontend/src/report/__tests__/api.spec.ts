import { describe, expect, it } from 'vitest'

import { reportUrl, storePath } from '../api'

describe('report urls', () => {
  it('keeps the trailing slash on the stores collection', () => {
    expect(reportUrl('https://reports.example', 'api/v1/stores/').href).toBe(
      'https://reports.example/api/v1/stores/',
    )
  })

  it('preserves a path prefix and encodes the store id', () => {
    expect(storePath('S 1001/west', 'summary')).toBe('api/v1/stores/S%201001%2Fwest/summary')
    expect(
      reportUrl('https://reports.example/prefix', storePath('S-1010', 'revenue/category'), {
        from: '2026-10-13',
        to: '2026-11-11',
        order_by: 'count',
        direction: 'asc',
        limit: 5,
        offset: 5,
      }).href,
    ).toBe(
      'https://reports.example/prefix/api/v1/stores/S-1010/revenue/category?from=2026-10-13&to=2026-11-11&order_by=count&direction=asc&limit=5&offset=5',
    )
  })
})
