import type {
  CategoryOrder,
  CategoryRevenuePage,
  DailyRevenue,
  SortDirection,
  Store,
  Summary,
} from './types'

export const CATEGORY_PAGE_SIZE = 5
export const REPORT_TIMEOUT_MS = 20_000

export type ReportFetcher = typeof fetch

export class ApiError extends Error {
  readonly status: number

  constructor(status: number, message: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

export interface CategoryQuery {
  from: string
  to: string
  orderBy: CategoryOrder
  direction: SortDirection
  limit: number
  offset: number
}

export function reportUrl(
  apiBaseUrl: string,
  path: string,
  query?: Record<string, string | number>,
): URL {
  const base = apiBaseUrl.endsWith('/') ? apiBaseUrl : `${apiBaseUrl}/`
  const url = new URL(path, base)
  if (query) {
    for (const [key, value] of Object.entries(query)) {
      url.searchParams.set(key, String(value))
    }
  }
  return url
}

export function storePath(storeId: string, suffix: string): string {
  return `api/v1/stores/${encodeURIComponent(storeId)}/${suffix}`
}

export function openRequest(timeoutMs = REPORT_TIMEOUT_MS): {
  signal: AbortSignal
  finish: () => void
  cancel: () => void
} {
  const controller = new AbortController()
  const timeoutId = setTimeout(() => {
    controller.abort(new DOMException('timeout', 'TimeoutError'))
  }, timeoutMs)
  const clear = () => clearTimeout(timeoutId)
  return {
    signal: controller.signal,
    finish: clear,
    cancel: () => {
      clear()
      controller.abort()
    },
  }
}

export function isAbortError(error: unknown): boolean {
  return isErrorName(error, 'AbortError')
}

export function isTimeoutError(error: unknown): boolean {
  return isErrorName(error, 'TimeoutError')
}

export async function fetchStores(
  apiBaseUrl: string,
  signal: AbortSignal,
  fetcher?: ReportFetcher,
): Promise<Store[]> {
  const stores = await getJson<Store[] | null>(
    reportUrl(apiBaseUrl, 'api/v1/stores/'),
    signal,
    fetcher,
  )
  return Array.isArray(stores) ? stores : []
}

export function fetchSummary(
  apiBaseUrl: string,
  storeId: string,
  from: string,
  to: string,
  signal: AbortSignal,
  fetcher?: ReportFetcher,
): Promise<Summary> {
  return getJson<Summary>(
    reportUrl(apiBaseUrl, storePath(storeId, 'summary'), { from, to }),
    signal,
    fetcher,
  )
}

export function fetchDailyRevenue(
  apiBaseUrl: string,
  storeId: string,
  from: string,
  to: string,
  signal: AbortSignal,
  fetcher?: ReportFetcher,
): Promise<DailyRevenue> {
  return getJson<DailyRevenue>(
    reportUrl(apiBaseUrl, storePath(storeId, 'revenue/daily'), { from, to }),
    signal,
    fetcher,
  )
}

export function fetchCategoryRevenue(
  apiBaseUrl: string,
  storeId: string,
  query: CategoryQuery,
  signal: AbortSignal,
  fetcher?: ReportFetcher,
): Promise<CategoryRevenuePage> {
  return getJson<CategoryRevenuePage>(
    reportUrl(apiBaseUrl, storePath(storeId, 'revenue/category'), {
      from: query.from,
      to: query.to,
      order_by: query.orderBy,
      direction: query.direction,
      limit: query.limit,
      offset: query.offset,
    }),
    signal,
    fetcher,
  )
}

async function getJson<T>(url: URL, signal: AbortSignal, fetcher?: ReportFetcher): Promise<T> {
  const doFetch = fetcher ?? fetch
  const response = await doFetch(url, {
    signal,
    headers: { Accept: 'application/json' },
    cache: 'no-store',
  })
  if (!response.ok) {
    const message = (await response.text()).trim()
    throw new ApiError(response.status, message || `HTTP ${response.status}`)
  }
  return (await response.json()) as T
}

function isErrorName(error: unknown, name: string): boolean {
  return error instanceof DOMException && error.name === name
}
