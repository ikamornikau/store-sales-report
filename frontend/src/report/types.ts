export interface Store {
  id: string
  name: string
  currency: string
  timezone: string
}

export interface Summary {
  currency: string
  revenue_cents: number
  transaction_count: number
  refunded_count: number
  voided_count: number
  average_basket_cents?: number
  items_sold: number
}

export interface DailyPoint {
  day: string
  revenue_cents: number
}

export interface DailyRevenue {
  currency: string
  timezone: string
  points: DailyPoint[]
}

export interface CategoryRevenue {
  category: string
  revenue_cents: number
  transaction_count: number
}

export interface CategoryRevenuePage {
  currency: string
  total: number
  items: CategoryRevenue[]
}

export type CategoryOrder = 'revenue' | 'count'
export type SortDirection = 'asc' | 'desc'
