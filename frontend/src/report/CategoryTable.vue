<script setup lang="ts">
import { computed } from 'vue'

import { CATEGORY_PAGE_SIZE } from './api'
import { formatCount, formatMoney } from './format'
import { categoryLabel, fill, messages, type Locale } from './i18n'
import type { CategoryOrder, CategoryRevenuePage, SortDirection } from './types'

const props = defineProps<{
  categories: CategoryRevenuePage
  locale: Locale
  orderBy: CategoryOrder
  direction: SortDirection
  offset: number
  loading: boolean
  error: string
}>()

const emit = defineEmits<{
  sort: [order: CategoryOrder]
  page: [offset: number]
  retry: []
}>()

const t = computed(() => messages[props.locale])
const page = computed(() => Math.floor(props.offset / CATEGORY_PAGE_SIZE) + 1)
const pages = computed(() =>
  Math.max(1, Math.ceil((props.categories.total ?? 0) / CATEGORY_PAGE_SIZE)),
)
const showingFrom = computed(() => (props.categories.total > 0 ? props.offset + 1 : 0))
const showingTo = computed(() => props.offset + props.categories.items.length)
const hasNext = computed(
  () => props.offset + props.categories.items.length < props.categories.total,
)
const rowPads = computed(() => {
  const count = props.categories.items.length
  if (count === 0) return 0
  return Math.max(0, CATEGORY_PAGE_SIZE - count)
})

function ariaSort(column: CategoryOrder): 'ascending' | 'descending' | 'none' {
  if (props.orderBy !== column) return 'none'
  return props.direction === 'asc' ? 'ascending' : 'descending'
}

function money(cents: number): string {
  return formatMoney(cents, props.categories.currency, props.locale)
}

function count(value: number): string {
  return formatCount(value, props.locale)
}
</script>

<template>
  <section class="panel" aria-labelledby="category-heading">
    <div class="section-head">
      <h2 id="category-heading">{{ t.revenueByCategory }}</h2>
    </div>
    <p v-if="error" class="alert" role="alert">
      <span>{{ error }}</span>
      <button type="button" @click="emit('retry')">{{ t.retry }}</button>
    </p>
    <div class="table-wrap" :aria-busy="loading">
      <table data-testid="category-table" aria-labelledby="category-heading">
        <thead>
          <tr>
            <th scope="col">{{ t.category }}</th>
            <th scope="col" class="num" :aria-sort="ariaSort('revenue')">
              <button
                type="button"
                class="sort"
                data-sort="revenue"
                @click="emit('sort', 'revenue')"
              >
                {{ t.revenue }}
                <span v-if="orderBy === 'revenue'" aria-hidden="true">{{
                  direction === 'asc' ? '↑' : '↓'
                }}</span>
              </button>
            </th>
            <th scope="col" class="num" :aria-sort="ariaSort('count')">
              <button type="button" class="sort" data-sort="count" @click="emit('sort', 'count')">
                {{ t.transactions }}
                <span v-if="orderBy === 'count'" aria-hidden="true">{{
                  direction === 'asc' ? '↑' : '↓'
                }}</span>
              </button>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="categories.items.length === 0">
            <td colspan="3" class="empty-cell">{{ t.emptyCategories }}</td>
          </tr>
          <tr v-for="item in categories.items" :key="item.category">
            <th scope="row">{{ categoryLabel(locale, item.category) }}</th>
            <td class="num">{{ money(item.revenue_cents) }}</td>
            <td class="num">{{ count(item.transaction_count) }}</td>
          </tr>
          <tr v-for="index in rowPads" :key="`pad-${index}`" class="pad" aria-hidden="true">
            <td colspan="3">&nbsp;</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div v-if="categories.total > 0" class="pager">
      <p>
        {{
          fill(t.showing, {
            from: showingFrom,
            to: showingTo,
            total: count(categories.total),
          })
        }}
      </p>
      <div class="pager-buttons">
        <button
          type="button"
          data-testid="previous-page"
          :disabled="offset === 0 || loading"
          @click="emit('page', offset - CATEGORY_PAGE_SIZE)"
        >
          {{ t.previous }}
        </button>
        <span>{{ fill(t.page, { page, pages }) }}</span>
        <button
          type="button"
          data-testid="next-page"
          :disabled="!hasNext || loading"
          @click="emit('page', offset + CATEGORY_PAGE_SIZE)"
        >
          {{ t.next }}
        </button>
      </div>
    </div>
  </section>
</template>

<style scoped>
.panel,
.panel *,
.panel *::before,
.panel *::after {
  box-sizing: border-box;
}

.panel {
  margin-top: 22px;
  padding: 18px;
  border: 1px solid var(--line);
  border-radius: 14px;
  background: var(--card);
}

h2,
p {
  margin: 0;
}

h2 {
  font-size: 1.15rem;
  letter-spacing: -0.02em;
}

.section-head {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: 16px;
  margin-bottom: 14px;
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

.alert button,
.pager button {
  min-height: 34px;
  padding: 0 12px;
  border: 1px solid currentColor;
  border-radius: 8px;
  background: #fff;
  color: inherit;
  font: inherit;
  cursor: pointer;
}

.alert button {
  color: var(--danger);
  border-color: #e4b4b4;
}

.pager button:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

button:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}

.table-wrap {
  overflow-x: auto;
}

.table-wrap[aria-busy='true'] tbody {
  opacity: 0.45;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  padding: 11px 12px;
  border-bottom: 1px solid var(--line);
  text-align: left;
}

th.num,
td.num {
  text-align: right;
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}

.sort {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  margin-left: auto;
  padding: 0;
  border: 0;
  background: transparent;
  color: inherit;
  font: inherit;
  font-weight: 600;
  cursor: pointer;
}

tbody th {
  font-weight: 600;
}

.empty-cell {
  color: var(--muted);
  text-align: center;
}

.pad td {
  color: transparent;
  user-select: none;
}

.pager {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-top: 14px;
}

.pager-buttons {
  display: flex;
  align-items: center;
  gap: 10px;
}

@media (max-width: 1200px) {
  .section-head {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
