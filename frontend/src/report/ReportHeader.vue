<script setup lang="ts">
import { computed } from 'vue'

import { localeLabel, locales, messages, type Locale } from './i18n'

const props = defineProps<{
  locale: Locale
}>()

const emit = defineEmits<{
  'update:locale': [locale: Locale]
}>()

const t = computed(() => messages[props.locale])
</script>

<template>
  <header class="top">
    <div>
      <h1>{{ t.title }}</h1>
      <p class="lede">{{ t.subtitle }}</p>
    </div>
    <div class="languages" role="group" :aria-label="t.language">
      <button
        v-for="code in locales"
        :key="code"
        type="button"
        :data-locale="code"
        :aria-pressed="locale === code"
        :aria-label="localeLabel[code]"
        @click="emit('update:locale', code)"
      >
        {{ localeLabel[code] }}
      </button>
    </div>
  </header>
</template>

<style scoped>
.top,
.top *,
.top *::before,
.top *::after {
  box-sizing: border-box;
}

.top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 24px;
}

h1,
p {
  margin: 0;
}

h1 {
  font-size: 1.85rem;
  line-height: 1.15;
  letter-spacing: -0.03em;
}

.lede {
  max-width: 62ch;
  margin-top: 8px;
  color: var(--muted);
}

.languages {
  display: flex;
  flex-shrink: 0;
  padding: 3px;
  border: 1px solid var(--line);
  border-radius: 999px;
  background: var(--card);
}

.languages button {
  min-height: 34px;
  padding: 0 12px;
  border: 0;
  border-radius: 999px;
  background: transparent;
  color: inherit;
  font: inherit;
  cursor: pointer;
}

.languages button[aria-pressed='true'] {
  background: var(--accent);
  color: #fff;
}

button:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}
</style>
