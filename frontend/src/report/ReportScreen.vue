<script setup lang="ts">
import { computed, onUnmounted, ref, watch } from 'vue'

import {
  ApiError,
  CATEGORY_PAGE_SIZE,
  fetchCategoryRevenue,
  fetchDailyRevenue,
  fetchStores,
  fetchSummary,
  isAbortError,
  isTimeoutError,
  openRequest,
  type CategoryQuery,
  type ReportFetcher,
} from './api'
import CategoryTable from './CategoryTable.vue'
import DailyRevenue from './DailyRevenue.vue'
import { defaultRange, validateRange, type RangeProblem } from './format'
import { isLocale, persistLocale, readLocale, messages, type Locale } from './i18n'
import ReportFilters from './ReportFilters.vue'
import ReportHeader from './ReportHeader.vue'
import SummaryFigures from './SummaryFigures.vue'
import type {
  CategoryOrder,
  CategoryRevenuePage,
  DailyRevenue as DailyRevenueData,
  SortDirection,
  Store,
  Summary,
} from './types'

const props = defineProps<{
  apiBaseUrl: string
  locale?: Locale
  fetcher?: ReportFetcher
}>()

const emit = defineEmits<{
  'update:locale': [locale: Locale]
}>()

const initialRange = defaultRange()
const activeLocale = ref<Locale>(initialLocale())
const stores = ref<Store[]>([])
const storeId = ref('')
const from = ref(initialRange.from)
const to = ref(initialRange.to)
const orderBy = ref<CategoryOrder>('revenue')
const direction = ref<SortDirection>('desc')
const offset = ref(0)
const shownOffset = ref(0)
const storesStatus = ref<'idle' | 'loading' | 'ready' | 'error'>('idle')
const reportStatus = ref<'idle' | 'loading' | 'ready' | 'error'>('idle')
const categoryLoading = ref(false)
const storesError = ref('')
const reportError = ref('')
const categoryError = ref('')
const rangeProblem = ref<RangeProblem | null>(null)
const summary = ref<Summary | null>(null)
const daily = ref<DailyRevenueData | null>(null)
const categories = ref<CategoryRevenuePage | null>(null)

const t = computed(() => messages[activeLocale.value])
const sortedStores = computed(() =>
  [...stores.value].sort((left, right) => left.name.localeCompare(right.name, activeLocale.value)),
)
const selectedStore = computed(
  () => stores.value.find((store) => store.id === storeId.value) ?? null,
)
const busy = computed(
  () =>
    storesStatus.value === 'loading' || reportStatus.value === 'loading' || categoryLoading.value,
)

let storesToken = 0
let reportToken = 0
let categoryToken = 0
let cancelRequest: (() => void) | null = null

watch(
  () => props.locale,
  (value) => {
    if (value && isLocale(value)) activeLocale.value = value
  },
)

watch(storeId, (id) => {
  if (!id) return
  offset.value = 0
  void loadReport()
})

watch([from, to], () => {
  offset.value = 0
  const problem = validateRange(from.value, to.value)
  rangeProblem.value = problem
  if (problem) {
    cancelActive()
    reportToken += 1
    categoryToken += 1
    reportStatus.value = 'idle'
    summary.value = null
    daily.value = null
    categories.value = null
    reportError.value = ''
    return
  }
  void loadReport()
})

watch(
  () => props.apiBaseUrl,
  () => {
    storeId.value = ''
    stores.value = []
    summary.value = null
    daily.value = null
    categories.value = null
    reportStatus.value = 'idle'
    void loadStores()
  },
  { immediate: true },
)

onUnmounted(() => {
  storesToken += 1
  reportToken += 1
  categoryToken += 1
  cancelActive()
})

function initialLocale(): Locale {
  if (props.locale && isLocale(props.locale)) return props.locale
  return readLocale() ?? 'en'
}

function setLocale(next: Locale) {
  activeLocale.value = next
  persistLocale(next)
  emit('update:locale', next)
}

function beginRequest() {
  cancelActive()
  const request = openRequest()
  cancelRequest = request.cancel
  return request
}

function cancelActive() {
  cancelRequest?.()
  cancelRequest = null
}

function categoryQuery(): CategoryQuery {
  return {
    from: from.value,
    to: to.value,
    orderBy: orderBy.value,
    direction: direction.value,
    limit: CATEGORY_PAGE_SIZE,
    offset: offset.value,
  }
}

async function loadStores() {
  const token = ++storesToken
  reportToken += 1
  categoryToken += 1
  const request = beginRequest()
  storesStatus.value = 'loading'
  storesError.value = ''
  try {
    const list = await fetchStores(props.apiBaseUrl, request.signal, props.fetcher)
    if (token !== storesToken) return
    stores.value = list
    storesStatus.value = 'ready'
    const first = sortedStores.value[0]
    if (!first) {
      storeId.value = ''
      return
    }
    if (storeId.value === first.id) {
      offset.value = 0
      void loadReport()
      return
    }
    storeId.value = first.id
  } catch (error) {
    if (token !== storesToken || isAbortError(error)) return
    storesStatus.value = 'error'
    storesError.value = isTimeoutError(error) ? t.value.timeout : t.value.storesError
  } finally {
    request.finish()
  }
}

async function loadReport() {
  if (!storeId.value) return
  const problem = validateRange(from.value, to.value)
  rangeProblem.value = problem
  if (problem) return

  const token = ++reportToken
  categoryToken += 1
  const request = beginRequest()
  const store = storeId.value
  const query = categoryQuery()
  reportStatus.value = 'loading'
  reportError.value = ''
  categoryError.value = ''
  summary.value = null
  daily.value = null
  categories.value = null
  try {
    const [summaryResult, dailyResult, categoryResult] = await Promise.all([
      fetchSummary(props.apiBaseUrl, store, query.from, query.to, request.signal, props.fetcher),
      fetchDailyRevenue(
        props.apiBaseUrl,
        store,
        query.from,
        query.to,
        request.signal,
        props.fetcher,
      ),
      fetchCategoryPage(store, query, request.signal),
    ])
    if (token !== reportToken) return
    summary.value = summaryResult
    daily.value = {
      ...dailyResult,
      points: Array.isArray(dailyResult.points) ? dailyResult.points : [],
    }
    categories.value = normalizeCategories(categoryResult)
    shownOffset.value = offset.value
    reportStatus.value = 'ready'
  } catch (error) {
    if (token !== reportToken || isAbortError(error)) return
    reportStatus.value = 'error'
    reportError.value = describeError(error)
  } finally {
    request.finish()
  }
}

async function loadCategories() {
  if (!storeId.value || validateRange(from.value, to.value)) return
  const token = ++categoryToken
  const request = beginRequest()
  const store = storeId.value
  const query = categoryQuery()
  categoryLoading.value = true
  categoryError.value = ''
  try {
    const page = await fetchCategoryPage(store, query, request.signal)
    if (token !== categoryToken) return
    categories.value = normalizeCategories(page)
    shownOffset.value = offset.value
  } catch (error) {
    if (token !== categoryToken || isAbortError(error)) return
    categoryError.value = describeError(error)
  } finally {
    if (token === categoryToken) categoryLoading.value = false
    request.finish()
  }
}

async function fetchCategoryPage(
  store: string,
  query: CategoryQuery,
  signal: AbortSignal,
): Promise<CategoryRevenuePage> {
  let page = normalizeCategories(
    await fetchCategoryRevenue(props.apiBaseUrl, store, query, signal, props.fetcher),
  )
  const nextOffset = clampedOffset(page.total, query.offset)
  if (nextOffset !== query.offset) {
    offset.value = nextOffset
    page = normalizeCategories(
      await fetchCategoryRevenue(
        props.apiBaseUrl,
        store,
        { ...query, offset: nextOffset },
        signal,
        props.fetcher,
      ),
    )
  }
  return page
}

function normalizeCategories(page: CategoryRevenuePage): CategoryRevenuePage {
  return {
    currency: page.currency,
    total: page.total ?? 0,
    items: Array.isArray(page.items) ? page.items : [],
  }
}

function clampedOffset(total: number, current: number): number {
  if (total <= 0) return 0
  if (current < total) return current
  return Math.floor((total - 1) / CATEGORY_PAGE_SIZE) * CATEGORY_PAGE_SIZE
}

function describeError(error: unknown): string {
  if (isTimeoutError(error)) return t.value.timeout
  if (error instanceof ApiError) {
    if (error.status === 404) return t.value.storeNotFound
    if (error.status === 400 && error.message) return error.message
    return t.value.loadError
  }
  return t.value.loadError
}

function toggleSort(next: CategoryOrder) {
  if (orderBy.value === next) {
    direction.value = direction.value === 'asc' ? 'desc' : 'asc'
  } else {
    orderBy.value = next
    direction.value = 'desc'
  }
  offset.value = 0
  void loadCategories()
}

function goTo(nextOffset: number) {
  offset.value = Math.max(0, nextOffset)
  void loadCategories()
}
</script>

<template>
  <section class="screen" :lang="activeLocale">
    <div
      v-if="busy"
      class="progress"
      role="progressbar"
      :aria-label="storesStatus === 'loading' ? t.loadingStores : t.loadingReport"
    />
    <div class="sheet">
      <ReportHeader :locale="activeLocale" @update:locale="setLocale" />

      <p v-if="storesStatus === 'loading'" class="status" role="status">{{ t.loadingStores }}</p>
      <p v-else-if="storesStatus === 'error'" class="alert" role="alert">
        <span>{{ storesError }}</span>
        <button type="button" @click="loadStores">{{ t.retry }}</button>
      </p>
      <p v-else-if="storesStatus === 'ready' && stores.length === 0" class="empty">
        {{ t.noStores }}
      </p>

      <ReportFilters
        v-if="stores.length > 0"
        :stores="sortedStores"
        :store-id="storeId"
        :from="from"
        :to="to"
        :locale="activeLocale"
        :range-problem="rangeProblem"
        @update:store-id="storeId = $event"
        @update:from="from = $event"
        @update:to="to = $event"
      />

      <template v-if="selectedStore && !rangeProblem">
        <div class="report-head">
          <h2>{{ selectedStore.name }}</h2>
          <p>{{ selectedStore.timezone }} · {{ selectedStore.currency }}</p>
        </div>

        <p v-if="reportStatus === 'loading'" class="status" role="status">{{ t.loadingReport }}</p>
        <p v-else-if="reportStatus === 'error'" class="alert" role="alert">
          <span>{{ reportError }}</span>
          <button type="button" data-testid="retry-report" @click="loadReport">
            {{ t.retry }}
          </button>
        </p>

        <template v-else-if="summary && daily && categories">
          <SummaryFigures :summary="summary" :locale="activeLocale" />
          <DailyRevenue
            :points="daily.points"
            :currency="daily.currency"
            :timezone="daily.timezone || selectedStore.timezone"
            :locale="activeLocale"
          />
          <CategoryTable
            :categories="categories"
            :locale="activeLocale"
            :order-by="orderBy"
            :direction="direction"
            :offset="shownOffset"
            :loading="categoryLoading"
            :error="categoryError"
            @sort="toggleSort"
            @page="goTo"
            @retry="loadCategories"
          />
        </template>
      </template>
    </div>
  </section>
</template>

<style scoped>
.screen {
  --ink: #1f1a14;
  --muted: #5c564e;
  --line: #e3dcd0;
  --paper: #f3f0e8;
  --card: #fffdf8;
  --accent: #0c6b4d;
  --danger: #6f2424;
  --danger-bg: #fdecec;
  position: relative;
  box-sizing: border-box;
  width: 100%;
  background: var(--paper);
  color: var(--ink);
  font-family: 'Segoe UI', system-ui, sans-serif;
  font-size: 16px;
  line-height: 1.45;
  border-top: 4px solid var(--accent);
}

.screen *,
.screen *::before,
.screen *::after {
  box-sizing: border-box;
}

.sheet {
  box-sizing: border-box;
  width: 100%;
  max-width: 1440px;
  margin: 0 auto;
  padding: 32px 40px 24px;
}

.progress {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  overflow: hidden;
  background: #d7ebe3;
}

.progress::before {
  content: '';
  display: block;
  width: 35%;
  height: 100%;
  background: var(--accent);
  animation: slide 1.1s ease-in-out infinite;
}

@keyframes slide {
  from {
    transform: translateX(-120%);
  }
  to {
    transform: translateX(320%);
  }
}

@media (prefers-reduced-motion: reduce) {
  .progress::before {
    animation: none;
    width: 100%;
  }
}

h2,
p {
  margin: 0;
}

h2 {
  font-size: 1.15rem;
  letter-spacing: -0.02em;
}

.status,
.empty,
.report-head p {
  color: var(--muted);
}

.status {
  margin-top: 22px;
}

.alert {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-top: 16px;
  padding: 12px 14px;
  border-radius: 10px;
  background: var(--danger-bg);
  color: var(--danger);
}

.alert button {
  min-height: 34px;
  padding: 0 12px;
  border: 1px solid #e4b4b4;
  border-radius: 8px;
  background: #fff;
  color: var(--danger);
  font: inherit;
  cursor: pointer;
}

button:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}

.report-head {
  margin-top: 28px;
}

.report-head p {
  margin-top: 4px;
}

@media (max-width: 1200px) {
  .sheet {
    padding: 28px 24px 20px;
  }
}
</style>
