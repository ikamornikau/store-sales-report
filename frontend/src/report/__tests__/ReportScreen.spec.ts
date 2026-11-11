import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount, type VueWrapper } from '@vue/test-utils'

import { addDays, formatMoney, MAX_RANGE_DAYS } from '../format'
import ReportScreen from '../ReportScreen.vue'
import type { CategoryRevenue, DailyRevenue, Store, Summary } from '../types'

const stores: Store[] = [
  {
    id: 'S-1004',
    name: 'Toronto Queen Street',
    currency: 'CAD',
    timezone: 'America/Toronto',
  },
  {
    id: 'S-1010',
    name: 'Manhattan Deli',
    currency: 'USD',
    timezone: 'America/New_York',
  },
]

const categories: CategoryRevenue[] = [
  { category: 'grocery', revenue_cents: 80000, transaction_count: 20 },
  { category: 'beverages', revenue_cents: 30000, transaction_count: 15 },
  { category: 'snacks', revenue_cents: 20000, transaction_count: 10 },
  { category: 'household', revenue_cents: 15000, transaction_count: 8 },
  { category: 'tobacco', revenue_cents: 7000, transaction_count: 4 },
  { category: 'lottery', revenue_cents: 2099, transaction_count: 2 },
]

const summary: Summary = {
  currency: 'USD',
  revenue_cents: 154099,
  transaction_count: 40,
  refunded_count: 2,
  voided_count: 1,
  average_basket_cents: 3852,
  items_sold: 1860,
}

const daily: DailyRevenue = {
  currency: 'USD',
  timezone: 'America/New_York',
  points: [
    { day: '2026-11-10', revenue_cents: 5000 },
    { day: '2026-11-11', revenue_cents: 7000 },
  ],
}

function json(body: unknown, status = 200) {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: () => Promise.resolve(body),
    text: () => Promise.resolve(typeof body === 'string' ? body : JSON.stringify(body)),
  }
}

function currencyFor(url: string): string {
  return url.includes('S-1004') ? 'CAD' : 'USD'
}

describe('ReportScreen', () => {
  let wrapper: VueWrapper | undefined
  let summaryStatus = 200
  let omitAverage = false

  beforeEach(() => {
    vi.stubGlobal('localStorage', memoryStorage())
    vi.spyOn(HTMLCanvasElement.prototype, 'getContext').mockReturnValue(null)
  })

  afterEach(() => {
    wrapper?.unmount()
    wrapper = undefined
    summaryStatus = 200
    omitAverage = false
    vi.restoreAllMocks()
    vi.unstubAllGlobals()
  })

  it('shows completed-sale figures, the daily chart, and the first category page', async () => {
    installFetch()
    wrapper = mount(ReportScreen, { props: { apiBaseUrl: 'https://reports.example' } })
    await settle()

    const urls = requestedUrls()
    expect(urls[0]).toBe('https://reports.example/api/v1/stores/')
    expect(urls.filter((url) => url.includes('/summary'))).toHaveLength(1)
    expect(urls.some((url) => url.includes('/api/v1/stores/S-1010/summary?from='))).toBe(true)
    expect(urls.filter((url) => url.includes('/revenue/daily?'))).toHaveLength(1)
    expect(
      urls.some(
        (url) =>
          url.includes('order_by=revenue') && url.includes('offset=0') && url.includes('limit=5'),
      ),
    ).toBe(true)
    expect(wrapper.get('[data-testid="revenue"]').text()).toBe(formatMoney(154099, 'USD', 'en'))
    expect(wrapper.get('[data-testid="transactions"]').text()).toBe('40')
    expect(wrapper.get('[data-testid="average-basket"]').text()).toBe(
      formatMoney(3852, 'USD', 'en'),
    )
    expect(wrapper.get('[data-testid="items-sold"]').text()).toBe('1,860')
    expect(wrapper.text()).toContain('Refunded transactions (2)')
    expect(wrapper.text()).toContain('voided transactions (1)')
    expect(wrapper.text()).toContain('America/New_York')
    expect(wrapper.get('[data-testid="day-readout"]').text()).toContain(
      formatMoney(7000, 'USD', 'en'),
    )
    expect(wrapper.text()).toContain('Grocery')
    expect(wrapper.text()).not.toContain('Lottery')
    expect(wrapper.text()).toContain('Showing 1–5 of 6')

    const init = vi.mocked(fetch).mock.calls[0]?.[1]
    expect(init).toMatchObject({ cache: 'no-store', headers: { Accept: 'application/json' } })
  })

  it('sorts and pages category revenue through the API', async () => {
    installFetch()
    wrapper = mount(ReportScreen, { props: { apiBaseUrl: 'https://reports.example' } })
    await settle()

    await wrapper.get('[data-sort="count"]').trigger('click')
    await settle()
    expect(latestUrl('/revenue/category')).toContain('order_by=count')
    expect(latestUrl('/revenue/category')).toContain('direction=desc')
    expect(latestUrl('/revenue/category')).toContain('offset=0')

    await wrapper.get('[data-sort="count"]').trigger('click')
    await settle()
    expect(latestUrl('/revenue/category')).toContain('direction=asc')

    await wrapper.get('[data-testid="next-page"]').trigger('click')
    await settle()
    expect(latestUrl('/revenue/category')).toContain('offset=5')
    expect(wrapper.text()).toContain('Lottery')
    expect(wrapper.get('[data-testid="next-page"]').attributes('disabled')).toBeDefined()

    await wrapper.get('[data-testid="previous-page"]').trigger('click')
    await settle()
    expect(latestUrl('/revenue/category')).toContain('offset=0')
  })

  it('switches store and language, and remembers the language', async () => {
    installFetch()
    wrapper = mount(ReportScreen, { props: { apiBaseUrl: 'https://reports.example' } })
    await settle()

    await wrapper.get('select[name="store"]').setValue('S-1004')
    await settle()
    expect(latestUrl('/summary')).toContain('/S-1004/summary')
    expect(wrapper.get('[data-testid="revenue"]').text()).toBe(formatMoney(154099, 'CAD', 'en'))

    await wrapper.get('[data-locale="fr"]').trigger('click')
    expect(wrapper.text()).toContain('Magasin')
    expect(wrapper.text()).toContain('Épicerie')
    expect(wrapper.get('[data-testid="revenue"]').text()).toBe(formatMoney(154099, 'CAD', 'fr'))

    wrapper.unmount()
    wrapper = mount(ReportScreen, { props: { apiBaseUrl: 'https://reports.example' } })
    await settle()
    expect(wrapper.get('[data-locale="fr"]').attributes('aria-pressed')).toBe('true')
    expect(wrapper.text()).toContain('Ventes du magasin')
  })

  it('follows a locale passed in by the host app', async () => {
    installFetch()
    localStorage.setItem('store-sales-report.locale', 'en')
    wrapper = mount(ReportScreen, {
      props: { apiBaseUrl: 'https://reports.example', locale: 'es' },
    })
    await settle()

    expect(wrapper.text()).toContain('Tienda')
    expect(wrapper.text()).toContain('Comestibles')
  })

  it('uses the host fetcher instead of the global fetch', async () => {
    const globalFetch = vi.fn<typeof fetch>()
    vi.stubGlobal('fetch', globalFetch)
    const fetcher = vi.fn<typeof fetch>(async (input: RequestInfo | URL) => {
      return answer(String(input)) as unknown as Response
    })
    wrapper = mount(ReportScreen, {
      props: { apiBaseUrl: 'https://reports.example', fetcher },
    })
    await settle()

    expect(fetcher).toHaveBeenCalled()
    expect(globalFetch).not.toHaveBeenCalled()
    expect(wrapper.get('[data-testid="revenue"]').text()).toBe(formatMoney(154099, 'USD', 'en'))
  })

  it('does not request a report when the date range is longer than 92 days', async () => {
    installFetch()
    wrapper = mount(ReportScreen, { props: { apiBaseUrl: 'https://reports.example' } })
    await settle()
    const before = requestedUrls().filter((url) => url.includes('/summary')).length
    const fromInput = wrapper.get('input[name="from"]')
    const toInput = wrapper.get('input[name="to"]')
    const fromValue = (fromInput.element as HTMLInputElement).value
    const toValue = (toInput.element as HTMLInputElement).value
    expect(fromInput.attributes('min')).toBe(addDays(toValue, -MAX_RANGE_DAYS))
    expect(fromInput.attributes('max')).toBe(toValue)
    expect(toInput.attributes('min')).toBe(fromValue)
    expect(toInput.attributes('max')).toBe(addDays(fromValue, MAX_RANGE_DAYS))

    await fromInput.setValue('2020-01-01')
    await settle()

    expect(requestedUrls().filter((url) => url.includes('/summary')).length).toBe(before)
    expect(wrapper.get('[data-testid="range-error"]').text()).toContain('92')
    expect(wrapper.find('[data-testid="revenue"]').exists()).toBe(false)
  })

  it('shows a retry when the report request fails, including a slow query', async () => {
    installFetch()
    summaryStatus = 500
    wrapper = mount(ReportScreen, { props: { apiBaseUrl: 'https://reports.example' } })
    await settle()
    expect(wrapper.text()).toContain('The report could not be loaded.')

    summaryStatus = 200
    await wrapper.get('[data-testid="retry-report"]').trigger('click')
    await settle()
    expect(wrapper.get('[data-testid="revenue"]').text()).toBe(formatMoney(154099, 'USD', 'en'))
  })

  it('treats a missing average basket as zero', async () => {
    installFetch()
    omitAverage = true
    wrapper = mount(ReportScreen, { props: { apiBaseUrl: 'https://reports.example' } })
    await settle()
    expect(wrapper.get('[data-testid="average-basket"]').text()).toBe(formatMoney(0, 'USD', 'en'))
  })

  it('explains when the report times out or the store list cannot be loaded', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new DOMException('timeout', 'TimeoutError')))
    wrapper = mount(ReportScreen, { props: { apiBaseUrl: 'https://reports.example' } })
    await settle()
    expect(wrapper.text()).toContain('The report took too long to load.')

    wrapper.unmount()
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('offline')))
    wrapper = mount(ReportScreen, { props: { apiBaseUrl: 'https://reports.example' } })
    await settle()
    expect(wrapper.text()).toContain('Stores could not be loaded.')
    expect(requestedUrls().some((url) => url.includes('/summary'))).toBe(false)
  })

  function installFetch() {
    vi.stubGlobal(
      'fetch',
      vi.fn(async (input: RequestInfo | URL) => answer(String(input))),
    )
  }

  function answer(url: string) {
    if (url.includes('/revenue/category')) {
      const parsed = new URL(url)
      const offset = Number(parsed.searchParams.get('offset') ?? '0')
      const limit = Number(parsed.searchParams.get('limit') ?? '5')
      return json({
        currency: currencyFor(url),
        total: categories.length,
        items: categories.slice(offset, offset + limit),
      })
    }
    if (url.includes('/revenue/daily')) {
      return json({ ...daily, currency: currencyFor(url) })
    }
    if (url.includes('/summary')) {
      if (summaryStatus !== 200) return json('failed to get transaction summary', summaryStatus)
      const body: Summary = { ...summary, currency: currencyFor(url) }
      if (omitAverage) delete body.average_basket_cents
      return json(body)
    }
    if (url.endsWith('/api/v1/stores/')) return json(stores)
    return json('not found', 404)
  }

  function requestedUrls(): string[] {
    return vi.mocked(fetch).mock.calls.map((call) => String(call[0]))
  }

  function latestUrl(fragment: string): string {
    const matches = requestedUrls().filter((url) => url.includes(fragment))
    return matches[matches.length - 1] ?? ''
  }

  async function settle() {
    for (let attempt = 0; attempt < 8; attempt += 1) await flushPromises()
  }
})

function memoryStorage(): Storage {
  const values = new Map<string, string>()
  return {
    get length() {
      return values.size
    },
    clear() {
      values.clear()
    },
    getItem(key: string) {
      return values.get(key) ?? null
    },
    key(index: number) {
      return [...values.keys()][index] ?? null
    },
    removeItem(key: string) {
      values.delete(key)
    },
    setItem(key: string, value: string) {
      values.set(key, value)
    },
  }
}
