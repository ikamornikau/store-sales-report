<script setup lang="ts">
import { computed } from 'vue'

import { addDays, isIsoDate, MAX_RANGE_DAYS, type RangeProblem } from './format'
import { fill, messages, type Locale } from './i18n'
import type { Store } from './types'

const props = defineProps<{
  stores: Store[]
  storeId: string
  from: string
  to: string
  locale: Locale
  rangeProblem: RangeProblem | null
}>()

const emit = defineEmits<{
  'update:storeId': [id: string]
  'update:from': [value: string]
  'update:to': [value: string]
}>()

const t = computed(() => messages[props.locale])
const storeIdModel = computed({
  get: () => props.storeId,
  set: (id: string) => emit('update:storeId', id),
})
const fromModel = computed({
  get: () => props.from,
  set: (value: string) => emit('update:from', value),
})
const toModel = computed({
  get: () => props.to,
  set: (value: string) => emit('update:to', value),
})
const fromMin = computed(() =>
  isIsoDate(props.to) ? addDays(props.to, -MAX_RANGE_DAYS) : undefined,
)
const fromMax = computed(() => (isIsoDate(props.to) ? props.to : undefined))
const toMin = computed(() => (isIsoDate(props.from) ? props.from : undefined))
const toMax = computed(() =>
  isIsoDate(props.from) ? addDays(props.from, MAX_RANGE_DAYS) : undefined,
)
const rangeMessage = computed(() => {
  if (props.rangeProblem === 'span') return fill(t.value.rangeError, { days: MAX_RANGE_DAYS })
  if (props.rangeProblem === 'order') return t.value.orderError
  if (props.rangeProblem === 'missing') return t.value.datesRequired
  return ''
})
</script>

<template>
  <form class="filters" novalidate @submit.prevent>
    <label class="field">
      <span>{{ t.store }}</span>
      <select v-model="storeIdModel" name="store">
        <option v-for="store in stores" :key="store.id" :value="store.id">
          {{ store.name }} ({{ store.currency }})
        </option>
      </select>
    </label>
    <label class="field">
      <span>{{ t.from }}</span>
      <input
        v-model="fromModel"
        name="from"
        type="date"
        :min="fromMin"
        :max="fromMax"
        :aria-invalid="rangeProblem !== null"
      />
    </label>
    <label class="field">
      <span>{{ t.to }}</span>
      <input
        v-model="toModel"
        name="to"
        type="date"
        :min="toMin"
        :max="toMax"
        :aria-invalid="rangeProblem !== null"
      />
    </label>
    <p class="hint">{{ fill(t.rangeHint, { days: MAX_RANGE_DAYS }) }}</p>
    <p v-if="rangeProblem" class="alert" role="alert" data-testid="range-error">
      {{ rangeMessage }}
    </p>
  </form>
</template>

<style scoped>
.filters,
.filters *,
.filters *::before,
.filters *::after {
  box-sizing: border-box;
}

.filters {
  display: grid;
  grid-template-columns: minmax(240px, 1.5fr) minmax(160px, 0.7fr) minmax(160px, 0.7fr);
  gap: 14px 16px;
  margin-top: 24px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
}

.field span {
  color: var(--muted);
  font-size: 0.75rem;
  font-weight: 600;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}

.field select,
.field input {
  width: 100%;
  height: 42px;
  padding: 0 12px;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: #fff;
  color: inherit;
  font: inherit;
}

.field select {
  -webkit-appearance: none;
  appearance: none;
  padding-right: 40px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='8' viewBox='0 0 12 8'%3E%3Cpath fill='none' stroke='%235c564e' stroke-width='1.6' stroke-linecap='round' stroke-linejoin='round' d='M1 1.5 6 6.5 11 1.5'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 16px center;
  background-size: 12px 8px;
}

p {
  margin: 0;
}

.hint,
.alert {
  grid-column: 1 / -1;
}

.hint {
  margin-top: -4px;
  color: var(--muted);
}

.alert {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 12px 14px;
  border-radius: 10px;
  background: var(--danger-bg);
  color: var(--danger);
}

select:focus-visible,
input:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}
</style>
