<template>
  <div class="trend-chart">
    <h3>{{ title }}</h3>
    <div ref="chartRef" class="chart"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import * as echarts from 'echarts'

defineProps<{
  data: any[]
  title: string
}>()

const chartRef = ref<HTMLElement | null>(null)

onMounted(() => {
  if (chartRef.value) {
    const chart = echarts.init(chartRef.value)
    chart.setOption({
      title: { text: 'Trend' },
      tooltip: {},
      xAxis: { type: 'category', data: [] },
      yAxis: { type: 'value' },
      series: [{ data: [], type: 'line' }]
    })
  }
})
</script>

<style scoped>
.trend-chart h3 {
  margin-bottom: 12px;
  color: #333;
}
.chart {
  width: 100%;
  height: 300px;
}
</style>
