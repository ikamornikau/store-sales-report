<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import {
  Chart as ChartJS,
  CategoryScale,
  Filler,
  LinearScale,
  LineElement,
  PointElement,
  type ActiveElement,
  type ChartEvent,
  type ChartOptions,
} from 'chart.js'
import { Line, type ChartComponentRef } from 'vue-chartjs'

import { formatDay, formatMoney, formatMoneyCompact, intlLocale } from './format'
import type { Locale } from './i18n'
import type { DailyPoint } from './types'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Filler)

// Hover marker is radius 5 with a 2px stroke centered on that radius.
const pointReach = 8

const props = defineProps<{
  points: DailyPoint[]
  currency: string
  locale: Locale
  dayLabel: string
  revenueLabel: string
}>()

const chartRef = ref<ChartComponentRef<'line'> | null>(null)
const hoverIndex = ref<number | null>(null)
const canDraw = canvasAvailable()

watch(
  () => props.points,
  () => {
    hoverIndex.value = null
  },
)

const activeIndex = computed(() => {
  const last = props.points.length - 1
  if (last < 0) return -1
  const chosen = hoverIndex.value ?? last
  return Math.min(last, Math.max(0, chosen))
})
const active = computed(() => props.points[activeIndex.value] ?? null)
const chartData = computed(() => ({
  labels: props.points.map((point) => formatDay(point.day, props.locale, 'short')),
  datasets: [
    {
      label: props.revenueLabel,
      data: props.points.map((point) => point.revenue_cents),
      borderColor: '#0c6b4d',
      backgroundColor: 'rgba(12, 107, 77, 0.16)',
      pointBackgroundColor: '#0c6b4d',
      pointBorderColor: '#fffdf8',
      pointBorderWidth: 2,
      pointRadius: 0,
      pointHoverRadius: 5,
      borderWidth: 2.5,
      clip: pointReach,
      fill: true,
      tension: 0,
    },
  ],
}))
const options = computed<ChartOptions<'line'>>(() => ({
  responsive: true,
  maintainAspectRatio: true,
  aspectRatio: 800 / 300,
  locale: intlLocale[props.locale],
  interaction: { mode: 'index', intersect: false },
  onHover,
  font: { family: "'Segoe UI', system-ui, sans-serif", size: 12 },
  color: '#5c564e',
  layout: { padding: { top: pointReach, right: pointReach } },
  plugins: { legend: { display: false } },
  scales: {
    x: {
      ticks: { maxTicksLimit: 6, maxRotation: 0, padding: pointReach },
      grid: { display: false },
      border: { display: false },
    },
    y: {
      beginAtZero: true,
      ticks: {
        maxTicksLimit: 4,
        padding: pointReach,
        callback: (value) => formatMoneyCompact(Number(value), props.currency, props.locale),
      },
      grid: { color: '#a89c8a' },
      border: { display: false },
    },
  },
}))

function canvasAvailable(): boolean {
  if (typeof document === 'undefined') return false
  const context = document.createElement('canvas').getContext('2d')
  return Boolean(context && typeof context.measureText === 'function')
}

function onHover(_event: ChartEvent, elements: ActiveElement[]) {
  hoverIndex.value = elements[0]?.index ?? null
}

function highlight(index: number) {
  const chart = chartRef.value?.chart
  if (!chart) return
  chart.setActiveElements([{ datasetIndex: 0, index }])
  chart.update()
}

function onKey(event: KeyboardEvent) {
  const last = props.points.length - 1
  if (last < 0) return
  let next = activeIndex.value
  if (event.key === 'ArrowRight') next = Math.min(last, next + 1)
  else if (event.key === 'ArrowLeft') next = Math.max(0, next - 1)
  else if (event.key === 'Home') next = 0
  else if (event.key === 'End') next = last
  else return
  event.preventDefault()
  hoverIndex.value = next
  highlight(next)
}
</script>

<template>
  <div class="chart" data-testid="chart">
    <p v-if="active" class="readout" data-testid="day-readout" aria-live="polite">
      <span>{{ formatDay(active.day, locale) }}</span>
      <strong>{{ formatMoney(active.revenue_cents, currency, locale) }}</strong>
    </p>
    <div class="plot" tabindex="0" role="group" :aria-label="revenueLabel" @keydown="onKey">
      <Line v-if="canDraw" ref="chartRef" :data="chartData" :options="options" />
    </div>
    <div class="sr-only">
      <table>
        <caption>
          {{
            revenueLabel
          }}
        </caption>
        <thead>
          <tr>
            <th scope="col">{{ dayLabel }}</th>
            <th scope="col">{{ revenueLabel }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="point in points" :key="point.day">
            <td>{{ formatDay(point.day, locale) }}</td>
            <td>{{ formatMoney(point.revenue_cents, currency, locale) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.chart {
  min-width: 0;
}

.readout {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  margin: 0 0 8px;
  font-variant-numeric: tabular-nums;
}

.plot {
  width: 100%;
  outline: none;
}

.plot:focus-visible {
  outline: 2px solid #0c6b4d;
  outline-offset: 3px;
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
