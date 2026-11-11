<script setup lang="ts">
import { ReportScreen } from '@/report'

const apiBaseUrl = window.location.origin
</script>

<template>
  <ReportScreen :api-base-url="apiBaseUrl" />
</template>

<style>
html,
body,
#app {
  margin: 0;
  min-height: 100%;
  background: #f3f0e8;
}
</style>
