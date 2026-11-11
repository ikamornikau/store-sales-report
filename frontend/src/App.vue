<script setup lang="ts">
import { ReportScreen } from '@/report'

const apiBaseUrl = window.location.origin
</script>

<template>
  <ReportScreen :api-base-url="apiBaseUrl" />
</template>
