<script setup lang="ts">
import { computed } from 'vue'

import RevenueChart from './RevenueChart.vue'
import { fill, messages, type Locale } from './i18n'
import type { DailyPoint } from './types'

const props = defineProps<{
  points: DailyPoint[]
  currency: string
  timezone: string
  locale: Locale
}>()

const t = computed(() => messages[props.locale])
</script>

<template>
  <section class="panel" aria-labelledby="daily-heading">
    <div class="section-head">
      <h2 id="daily-heading">{{ t.revenuePerDay }}</h2>
      <p>{{ fill(t.dayCaption, { timezone }) }}</p>
    </div>
    <p v-if="points.length === 0" class="empty">{{ t.emptyChart }}</p>
    <RevenueChart
      v-else
      :points="points"
      :currency="currency"
      :locale="locale"
      :day-label="t.day"
      :revenue-label="t.revenue"
    />
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

.section-head p,
.empty {
  color: var(--muted);
}

.section-head p {
  text-align: right;
}

@media (max-width: 1200px) {
  .section-head {
    align-items: flex-start;
    flex-direction: column;
  }

  .section-head p {
    text-align: left;
  }
}
</style>
