<template>
  <div class="dashboard">
    <h1>Market Overview</h1>
    <div class="metrics">
      <ScoreCard title="Total Products" :value="0" icon="📦" />
      <ScoreCard title="Active Competitors" :value="0" icon="🏢" />
      <ScoreCard title="Price Changes" :value="0" icon="📊" />
      <ScoreCard title="Supply Alerts" :value="0" icon="⚠️" />
    </div>
    <div class="chart-container">
      <TrendChart :data="[]" title="Market Trends" />
    </div>
  </div>
</template>

<script setup lang="ts">
import ScoreCard from '@/components/ScoreCard.vue'
import TrendChart from '@/components/TrendChart.vue'
</script>

<style scoped>
.dashboard {
  padding: 24px;
}
.metrics {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}
.chart-container {
  background: white;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}
</style>
