<template>
  <div class="pricing">
    <h1>Pricing Strategy</h1>
    <PriceTable :data="[]" />
  </div>
</template>

<script setup lang="ts">
import PriceTable from '@/components/PriceTable.vue'
</script>

<style scoped>
.pricing {
  padding: 24px;
}
</style>
