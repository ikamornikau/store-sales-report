<script setup lang="ts">
import { computed } from 'vue'

import { formatCount, formatMoney } from './format'
import { fill, messages, type Locale } from './i18n'
import type { Summary } from './types'

const props = defineProps<{
  summary: Summary
  locale: Locale
}>()

const t = computed(() => messages[props.locale])
const hasExclusions = computed(
  () => (props.summary.refunded_count ?? 0) + (props.summary.voided_count ?? 0) > 0,
)

function money(cents: number | undefined): string {
  return formatMoney(cents ?? 0, props.summary.currency, props.locale)
}

function count(value: number | undefined): string {
  return formatCount(value ?? 0, props.locale)
}
</script>

<template>
  <section class="panel" aria-labelledby="figures-heading">
    <h2 id="figures-heading" class="sr-only">{{ t.keyFigures }}</h2>
    <dl class="figures">
      <div class="figure">
        <dt>{{ t.revenue }}</dt>
        <dd data-testid="revenue">{{ money(summary.revenue_cents) }}</dd>
      </div>
      <div class="figure">
        <dt>{{ t.transactions }}</dt>
        <dd data-testid="transactions">{{ count(summary.transaction_count) }}</dd>
      </div>
      <div class="figure">
        <dt>{{ t.averageBasket }}</dt>
        <dd data-testid="average-basket">{{ money(summary.average_basket_cents) }}</dd>
      </div>
      <div class="figure">
        <dt>{{ t.itemsSold }}</dt>
        <dd data-testid="items-sold">{{ count(summary.items_sold) }}</dd>
      </div>
    </dl>
    <p v-if="hasExclusions" class="excluded">
      {{
        fill(t.excluded, {
          refunded: count(summary.refunded_count),
          voided: count(summary.voided_count),
        })
      }}
    </p>
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

.figures {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  margin: 0;
}

.figure {
  min-width: 0;
  padding: 4px 8px 2px 0;
}

.figure dt {
  color: var(--muted);
  font-size: 0.75rem;
  font-weight: 600;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}

.figure dd {
  margin: 8px 0 0;
  font-size: clamp(1.2rem, 1.7vw, 1.7rem);
  font-weight: 600;
  letter-spacing: -0.03em;
  font-variant-numeric: tabular-nums;
  overflow-wrap: anywhere;
}

.figure:first-child dd {
  color: var(--accent);
}

.excluded {
  margin-top: 14px;
  color: var(--muted);
}

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}
</style>
