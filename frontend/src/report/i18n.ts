export const locales = ['en', 'es', 'fr'] as const

export type Locale = (typeof locales)[number]

export interface Messages {
  title: string
  subtitle: string
  store: string
  from: string
  to: string
  language: string
  rangeHint: string
  keyFigures: string
  revenue: string
  transactions: string
  averageBasket: string
  itemsSold: string
  revenuePerDay: string
  day: string
  dayCaption: string
  revenueByCategory: string
  category: string
  previous: string
  next: string
  page: string
  showing: string
  emptyCategories: string
  emptyChart: string
  loadingStores: string
  loadingReport: string
  loadError: string
  storesError: string
  retry: string
  rangeError: string
  orderError: string
  datesRequired: string
  noStores: string
  excluded: string
  timeout: string
  storeNotFound: string
  categories: {
    beverages: string
    snacks: string
    tobacco: string
    grocery: string
    lottery: string
    household: string
  }
}

export const localeLabel: Record<Locale, string> = {
  en: 'English',
  es: 'Español',
  fr: 'Français',
}

export const messages: Record<Locale, Messages> = {
  en: {
    title: 'Store sales',
    subtitle: 'Completed sales only. Refunded and voided transactions are excluded.',
    store: 'Store',
    from: 'From',
    to: 'To',
    language: 'Language',
    rangeHint: 'Up to {days} days.',
    keyFigures: 'Key figures',
    revenue: 'Revenue',
    transactions: 'Transactions',
    averageBasket: 'Average basket',
    itemsSold: 'Items sold',
    revenuePerDay: 'Revenue per day',
    day: 'Day',
    dayCaption: 'Days follow the store clock ({timezone}).',
    revenueByCategory: 'Revenue by category',
    category: 'Category',
    previous: 'Previous',
    next: 'Next',
    page: 'Page {page} of {pages}',
    showing: 'Showing {from}–{to} of {total}',
    emptyCategories: 'No category sales in this range.',
    emptyChart: 'No daily revenue in this range.',
    loadingStores: 'Loading stores…',
    loadingReport: 'Loading the report. Large ranges can take a moment.',
    loadError: 'The report could not be loaded.',
    storesError: 'Stores could not be loaded.',
    retry: 'Retry',
    rangeError: 'Choose a range of {days} days or less.',
    orderError: 'The start date must be on or before the end date.',
    datesRequired: 'Choose a start date and an end date.',
    noStores: 'No stores are available.',
    excluded:
      'Refunded transactions ({refunded}) and voided transactions ({voided}) are not included.',
    timeout: 'The report took too long to load. Try again.',
    storeNotFound: 'That store could not be found.',
    categories: {
      beverages: 'Beverages',
      snacks: 'Snacks',
      tobacco: 'Tobacco',
      grocery: 'Grocery',
      lottery: 'Lottery',
      household: 'Household',
    },
  },
  es: {
    title: 'Ventas de la tienda',
    subtitle:
      'Solo ventas completadas. Las transacciones reembolsadas y anuladas quedan excluidas.',
    store: 'Tienda',
    from: 'Desde',
    to: 'Hasta',
    language: 'Idioma',
    rangeHint: 'Hasta {days} días.',
    keyFigures: 'Cifras clave',
    revenue: 'Ingresos',
    transactions: 'Transacciones',
    averageBasket: 'Ticket medio',
    itemsSold: 'Artículos vendidos',
    revenuePerDay: 'Ingresos por día',
    day: 'Día',
    dayCaption: 'Los días siguen el reloj de la tienda ({timezone}).',
    revenueByCategory: 'Ingresos por categoría',
    category: 'Categoría',
    previous: 'Anterior',
    next: 'Siguiente',
    page: 'Página {page} de {pages}',
    showing: 'Mostrando {from}–{to} de {total}',
    emptyCategories: 'No hay ventas por categoría en este periodo.',
    emptyChart: 'No hay ingresos diarios en este periodo.',
    loadingStores: 'Cargando tiendas…',
    loadingReport: 'Cargando el informe. Los periodos largos pueden tardar un momento.',
    loadError: 'No se pudo cargar el informe.',
    storesError: 'No se pudieron cargar las tiendas.',
    retry: 'Reintentar',
    rangeError: 'Elige un periodo de {days} días o menos.',
    orderError: 'La fecha inicial debe ser anterior o igual a la final.',
    datesRequired: 'Elige una fecha de inicio y una de fin.',
    noStores: 'No hay tiendas disponibles.',
    excluded:
      'No se incluyen las transacciones reembolsadas ({refunded}) ni las anuladas ({voided}).',
    timeout: 'El informe tardó demasiado en cargar. Inténtalo de nuevo.',
    storeNotFound: 'No se encontró esa tienda.',
    categories: {
      beverages: 'Bebidas',
      snacks: 'Aperitivos',
      tobacco: 'Tabaco',
      grocery: 'Comestibles',
      lottery: 'Lotería',
      household: 'Hogar',
    },
  },
  fr: {
    title: 'Ventes du magasin',
    subtitle:
      'Ventes finalisées uniquement. Les transactions remboursées et annulées sont exclues.',
    store: 'Magasin',
    from: 'Du',
    to: 'Au',
    language: 'Langue',
    rangeHint: "Jusqu'à {days} jours.",
    keyFigures: 'Chiffres clés',
    revenue: "Chiffre d'affaires",
    transactions: 'Transactions',
    averageBasket: 'Panier moyen',
    itemsSold: 'Articles vendus',
    revenuePerDay: "Chiffre d'affaires par jour",
    day: 'Jour',
    dayCaption: "Les jours suivent l'horloge du magasin ({timezone}).",
    revenueByCategory: "Chiffre d'affaires par catégorie",
    category: 'Catégorie',
    previous: 'Précédent',
    next: 'Suivant',
    page: 'Page {page} sur {pages}',
    showing: '{from}–{to} sur {total}',
    emptyCategories: 'Aucune vente par catégorie sur cette période.',
    emptyChart: "Aucun chiffre d'affaires quotidien sur cette période.",
    loadingStores: 'Chargement des magasins…',
    loadingReport: 'Chargement du rapport. Les longues périodes peuvent prendre un moment.',
    loadError: 'Impossible de charger le rapport.',
    storesError: 'Impossible de charger les magasins.',
    retry: 'Réessayer',
    rangeError: 'Choisissez une période de {days} jours ou moins.',
    orderError: 'La date de début doit être antérieure ou égale à la date de fin.',
    datesRequired: 'Choisissez une date de début et une date de fin.',
    noStores: 'Aucun magasin disponible.',
    excluded:
      'Les transactions remboursées ({refunded}) et annulées ({voided}) ne sont pas incluses.',
    timeout: 'Le rapport a mis trop de temps à charger. Réessayez.',
    storeNotFound: 'Ce magasin est introuvable.',
    categories: {
      beverages: 'Boissons',
      snacks: 'En-cas',
      tobacco: 'Tabac',
      grocery: 'Épicerie',
      lottery: 'Loterie',
      household: 'Maison',
    },
  },
}

const STORAGE_KEY = 'store-sales-report.locale'

export function isLocale(value: string): value is Locale {
  return (locales as readonly string[]).includes(value)
}

export function readLocale(): Locale | null {
  try {
    const value = localStorage.getItem(STORAGE_KEY)
    return value && isLocale(value) ? value : null
  } catch {
    return null
  }
}

export function persistLocale(locale: Locale): void {
  try {
    localStorage.setItem(STORAGE_KEY, locale)
  } catch {
    // Storage can be unavailable in private mode. The choice still applies to this view.
  }
}

export function fill(template: string, vars: Record<string, string | number>): string {
  return template.replace(/\{(\w+)\}/g, (match, key: string) => {
    const value = vars[key]
    return value === undefined ? match : String(value)
  })
}

export function categoryLabel(locale: Locale, category: string): string {
  const labels: Record<string, string> = messages[locale].categories
  return labels[category] ?? category
}
