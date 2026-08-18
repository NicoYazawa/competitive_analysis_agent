<template>
  <div class="chart-container">
    <div class="chart-header">
      <div>
        <div class="chart-title">{{ title }}</div>
        <div v-if="subtitle" class="chart-subtitle">{{ subtitle }}</div>
      </div>
      <div v-if="legend?.length" class="chart-legend">
        <div v-for="item in legend" :key="item.label" class="legend-item">
          <div class="legend-dot" :style="{ background: item.color }"></div>
          {{ item.label }}
        </div>
      </div>
    </div>
    <div ref="chartRef" class="chart-area"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import * as echarts from 'echarts'

export interface ChartLegend {
  label: string
  color: string
}

const props = defineProps<{
  title: string
  subtitle?: string
  legend?: ChartLegend[]
  xData?: string[]
  series?: {
    name: string
    data: number[]
    color: string
    smooth?: boolean
  }[]
  type?: 'line' | 'bar'
}>()

const chartRef = ref<HTMLElement | null>(null)
let chartInstance: echarts.ECharts | null = null

const initChart = () => {
  if (!chartRef.value) return

  chartInstance = echarts.init(chartRef.value)

  const option: echarts.EChartsOption = {
    tooltip: {
      trigger: 'axis',
      backgroundColor: '#111111',
      borderColor: '#111111',
      textStyle: {
        color: '#ffffff',
        fontSize: 12,
        fontFamily: 'ui-monospace, JetBrains Mono, monospace'
      },
      axisPointer: {
        type: 'cross',
        crossStyle: {
          color: '#e5e5e5'
        }
      }
    },
    grid: {
      left: 50,
      right: 16,
      top: 20,
      bottom: 30
    },
    xAxis: {
      type: 'category',
      data: props.xData || [],
      boundaryGap: false,
      axisLine: {
        lineStyle: { color: '#e5e5e5' }
      },
      axisLabel: {
        color: '#6b6b6b',
        fontSize: 10,
        fontFamily: 'ui-monospace, JetBrains Mono, monospace'
      },
      axisTick: { show: false }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisLabel: {
        color: '#6b6b6b',
        fontSize: 10,
        fontFamily: 'ui-monospace, JetBrains Mono, monospace'
      },
      splitLine: {
        lineStyle: { color: '#e5e5e5', type: 'dashed' }
      }
    },
    series: (props.series || []).map(s => ({
      name: s.name,
      type: props.type || 'line',
      data: s.data,
      smooth: s.smooth ?? true,
      lineStyle: { color: s.color, width: 2 },
      itemStyle: { color: s.color },
      areaStyle: props.type !== 'bar' ? {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: s.color + '26' },
          { offset: 1, color: s.color + '00' }
        ])
      } : undefined,
      symbol: 'circle',
      symbolSize: 4
    }))
  }

  chartInstance.setOption(option)
}

const handleResize = () => {
  chartInstance?.resize()
}

onMounted(() => {
  initChart()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  chartInstance?.dispose()
})

watch(() => [props.series, props.xData], () => {
  initChart()
}, { deep: true })
</script>

<style scoped>
.chart-container {
  background: #ffffff;
  border: 1px solid #e5e5e5;
  border-radius: 12px;
  padding: 24px;
  position: relative;
}
.chart-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 20px;
  gap: 12px;
}
.chart-title {
  font-size: 15px;
  font-weight: 600;
  color: #111111;
}
.chart-subtitle {
  font-size: 12px;
  color: #6b6b6b;
  margin-top: 2px;
}
.chart-legend {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}
.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #6b6b6b;
}
.legend-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}
.chart-area {
  position: relative;
  height: 200px;
  width: 100%;
}
</style>
