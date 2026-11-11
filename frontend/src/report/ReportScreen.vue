<script setup lang="ts">
import { onMounted, ref } from 'vue'

const props = defineProps<{
  apiBaseUrl: string
}>()

const text = ref('')

onMounted(async () => {
  try {
    const response = await fetch(new URL('/api/v1/stores/', props.apiBaseUrl))
    text.value = await response.text()
  } catch (error: unknown) {
    text.value = error instanceof Error ? error.message : 'Failed to load stores'
  }
})
</script>

<template>
  <pre>{{ text }}</pre>
</template>
