import type { Locale } from './i18n'

export const MAX_RANGE_DAYS = 92
export const DEFAULT_SPAN_DAYS = 29

export const intlLocale: Record<Locale, string> = {
  en: 'en-US',
  es: 'es-ES',
  fr: 'fr-FR',
}

const ISO_DATE = /^\d{4}-\d{2}-\d{2}$/

export type RangeProblem = 'missing' | 'order' | 'span'

export function isIsoDate(value: string): boolean {
  if (!ISO_DATE.test(value)) return false
  const [year, month, day] = value.split('-').map(Number)
  if (!year || !month || !day) return false
  const date = new Date(Date.UTC(year, month - 1, day))
  return (
    date.getUTCFullYear() === year && date.getUTCMonth() === month - 1 && date.getUTCDate() === day
  )
}

export function isoDate(date: Date): string {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

export function addDays(iso: string, days: number): string {
  if (!isIsoDate(iso)) return iso
  const [year, month, day] = iso.split('-').map(Number)
  const date = new Date(Date.UTC(year!, month! - 1, day!))
  date.setUTCDate(date.getUTCDate() + days)
  return date.toISOString().slice(0, 10)
}

export function daysBetween(from: string, to: string): number {
  const start = Date.parse(`${from}T00:00:00Z`)
  const end = Date.parse(`${to}T00:00:00Z`)
  return Math.round((end - start) / 86_400_000)
}

export function defaultRange(today = new Date()): { from: string; to: string } {
  const to = isoDate(today)
  return { from: addDays(to, -DEFAULT_SPAN_DAYS), to }
}

export function validateRange(from: string, to: string): RangeProblem | null {
  if (!isIsoDate(from) || !isIsoDate(to)) return 'missing'
  if (from > to) return 'order'
  if (daysBetween(from, to) > MAX_RANGE_DAYS) return 'span'
  return null
}

export function formatCount(value: number, locale: Locale): string {
  return new Intl.NumberFormat(intlLocale[locale]).format(value)
}

export function formatMoney(cents: number, currency: string, locale: Locale): string {
  const amount = cents / 100
  try {
    return new Intl.NumberFormat(intlLocale[locale], {
      style: 'currency',
      currency,
    }).format(amount)
  } catch {
    return `${amount.toFixed(2)} ${currency}`
  }
}

export function formatMoneyCompact(cents: number, currency: string, locale: Locale): string {
  const amount = cents / 100
  try {
    return new Intl.NumberFormat(intlLocale[locale], {
      style: 'currency',
      currency,
      notation: 'compact',
      maximumFractionDigits: 1,
    }).format(amount)
  } catch {
    return `${amount.toFixed(0)} ${currency}`
  }
}

export function formatDay(
  iso: string,
  locale: Locale,
  style: 'short' | 'medium' = 'medium',
): string {
  if (!isIsoDate(iso)) return iso
  const [year, month, day] = iso.split('-').map(Number)
  return new Intl.DateTimeFormat(intlLocale[locale], {
    timeZone: 'UTC',
    month: 'short',
    day: 'numeric',
    year: style === 'medium' ? 'numeric' : undefined,
  }).format(new Date(Date.UTC(year!, month! - 1, day!)))
}
