<template>
  <div class="pricing fade-in">
    <div class="page-header">
      <h1>定价策略</h1>
      <p>竞品价格追踪与智能定价建议</p>
    </div>

    <!-- KPI Row -->
    <div class="kpi-grid" style="margin-bottom: 32px;">
      <div class="kpi-card">
        <div class="kpi-label">监控SKU数</div>
        <div class="kpi-value">{{ kpis.skuCount.toLocaleString() }}</div>
        <div class="kpi-delta neutral">覆盖 {{ kpis.skuCount }} 个品类</div>
      </div>
      <div class="kpi-card warn">
        <div class="kpi-label">价格低于均值预警</div>
        <div class="kpi-value">{{ kpis.underpriceAlerts }}</div>
        <div class="kpi-delta warning">需关注</div>
      </div>
      <div class="kpi-card">
        <div class="kpi-label">平均利润空间</div>
        <div class="kpi-value">{{ kpis.avgProfitMargin }}<span>%</span></div>
        <div class="kpi-delta up">↑ {{ kpis.profitDelta }}%</div>
      </div>
      <div class="kpi-card">
        <div class="kpi-label">本月调价建议</div>
        <div class="kpi-value">{{ kpis.pendingAdjustments }}</div>
        <div class="kpi-delta neutral">待审核</div>
      </div>
    </div>

    <!-- Price Trend Chart -->
    <div class="chart-container" style="margin-bottom: 32px;">
      <div class="chart-header">
        <div>
          <div class="chart-title">价格走势对比</div>
          <div class="chart-subtitle">近60天 · 热门SKU价格追踪</div>
        </div>
        <div class="chart-legend">
          <div class="legend-item">
            <div class="legend-dot" style="background: #2a78d6"></div>
            Anker 737
          </div>
          <div class="legend-item">
            <div class="legend-dot" style="background: #eb6834"></div>
            Jackery 1000
          </div>
          <div class="legend-item">
            <div class="legend-dot" style="background: #1baf7a"></div>
            UGREEN 100W
          </div>
        </div>
      </div>
      <TrendChart
        title=""
        :legend="[]"
        :xData="store.pricingTrend60Days.xData"
        :series="[
          { name: 'Anker 737', data: store.pricingTrend60Days.anker, color: '#2a78d6' },
          { name: 'Jackery 1000', data: store.pricingTrend60Days.jackery, color: '#eb6834' },
          { name: 'UGREEN 100W', data: store.pricingTrend60Days.ugreen, color: '#1baf7a' }
        ]"
        style="height: 240px;"
      />
    </div>

    <!-- Price Table + Suggestions -->
    <div class="grid-2">
      <!-- Competitor Price Table -->
      <div class="card" style="padding: 0; overflow: hidden;">
        <div style="padding: 16px 20px; border-bottom: 1px solid var(--border);">
          <span class="card-title">竞品价格对比</span>
        </div>
        <table class="data-table">
          <thead>
            <tr>
              <th>产品</th>
              <th class="text-right">当前价</th>
              <th class="text-right">建议价</th>
              <th class="text-center">信号</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="s in store.pricingSuggestions" :key="s.product">
              <td style="font-weight: 600;">{{ s.product }}</td>
              <td class="num text-right">${{ s.currentPrice.toFixed(2) }}</td>
              <td class="num text-right" style="color: var(--status-good);">${{ s.suggestedPrice.toFixed(2) }}</td>
              <td class="text-center">
                <AlertBadge v-if="s.severity === 'critical'" type="critical" dot>高价</AlertBadge>
                <AlertBadge v-else-if="s.severity === 'warning'" type="warning" dot>偏高</AlertBadge>
                <AlertBadge v-else type="good" dot>合理</AlertBadge>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Suggestions + Profit Simulation -->
      <div style="display: flex; flex-direction: column; gap: 16px;">
        <!-- AI Suggestions -->
        <div class="card">
          <div class="card-header">
            <span class="card-title">智能定价建议</span>
            <AlertBadge type="info" dot>AI生成</AlertBadge>
          </div>
          <div style="display: flex; flex-direction: column; gap: 12px;">
            <div v-for="s in store.pricingSuggestions" :key="s.product"
                 :style="['padding: 12px; border: 1px solid var(--border); border-radius: 8px; border-left: 3px solid;',
                           s.severity === 'critical' ? 'border-left-color: var(--status-critical);' : '',
                           s.severity === 'warning' ? 'border-left-color: var(--status-warn);' : '',
                           s.severity === 'good' ? 'border-left-color: var(--status-good);' : '']">
              <div style="font-size: 13px; font-weight: 600; margin-bottom: 4px;">{{ s.product }}</div>
              <div style="font-size: 12px; color: var(--muted);">
                建议调低至 <strong style="color: var(--fg);">${{ s.suggestedPrice }}</strong>，{{ s.reason }}
              </div>
            </div>
          </div>
        </div>

        <!-- Profit Simulation -->
        <div class="card">
          <div class="card-header">
            <span class="card-title">利润模拟</span>
          </div>
          <div style="display: flex; flex-direction: column; gap: 8px;">
            <div style="display: flex; justify-content: space-between; font-size: 13px; padding: 6px 0; border-bottom: 1px solid var(--border);">
              <span style="color: var(--muted);">建议售价</span>
              <span class="num" style="font-weight: 700;">${{ store.profitSimulation.suggestedPrice.toFixed(2) }}</span>
            </div>
            <div style="display: flex; justify-content: space-between; font-size: 13px; padding: 6px 0; border-bottom: 1px solid var(--border);">
              <span style="color: var(--muted);">成本价</span>
              <span class="num">${{ store.profitSimulation.cost.toFixed(2) }}</span>
            </div>
            <div style="display: flex; justify-content: space-between; font-size: 13px; padding: 6px 0; border-bottom: 1px solid var(--border);">
              <span style="color: var(--muted);">平台佣金 15%</span>
              <span class="num">-${{ store.profitSimulation.platformFee.toFixed(2) }}</span>
            </div>
            <div style="display: flex; justify-content: space-between; font-size: 13px; padding: 6px 0; border-bottom: 1px solid var(--border);">
              <span style="color: var(--muted);">物流成本</span>
              <span class="num">-${{ store.profitSimulation.shipping.toFixed(2) }}</span>
            </div>
            <div style="display: flex; justify-content: space-between; font-size: 14px; padding: 8px 0; font-weight: 700;">
              <span>预估利润</span>
              <span class="num" style="color: var(--status-good);">
                ${{ store.profitSimulation.profit.toFixed(2) }} ({{ store.profitSimulation.profitRate }}%)
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useMarketStore } from '@/stores/market'
import TrendChart from '@/components/TrendChart.vue'
import AlertBadge from '@/components/AlertBadge.vue'

const store = useMarketStore()
const kpis = computed(() => store.pricingKpis)

onMounted(() => {
  store.fetchPricingData()
})
</script>

<style scoped>
.kpi-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}
.kpi-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 20px;
  position: relative;
  overflow: hidden;
}
.kpi-card::before {
  content: '';
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 3px;
  background: var(--accent);
}
.kpi-card.warn::before { background: var(--status-warn); }
.kpi-label {
  font-size: 12px;
  color: var(--muted);
  font-weight: 500;
  margin-bottom: 8px;
}
.kpi-value {
  font-family: ui-monospace, 'JetBrains Mono', monospace;
  font-size: 32px;
  font-weight: 700;
  letter-spacing: -0.03em;
  color: var(--fg);
  line-height: 1;
}
.kpi-value span { font-size: 16px; font-weight: 400; color: var(--muted); margin-left: 2px; }
.kpi-delta {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  font-size: 12px;
  font-weight: 600;
  margin-top: 8px;
  padding: 2px 8px;
  border-radius: 999px;
}
.kpi-delta.up {
  color: var(--status-good);
  background: color-mix(in oklch, var(--status-good) 12%, transparent);
}
.kpi-delta.down {
  color: var(--status-critical);
  background: color-mix(in oklch, var(--status-critical) 12%, transparent);
}
.kpi-delta.neutral {
  color: var(--muted);
  background: var(--bg);
}
.kpi-delta.warning {
  color: #92400e;
  background: color-mix(in oklch, var(--status-warn) 12%, transparent);
}

.chart-container {
  background: #ffffff;
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 24px;
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
  color: var(--fg);
}
.chart-subtitle {
  font-size: 12px;
  color: var(--muted);
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
  color: var(--muted);
}
.legend-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.card-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--muted);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  font-family: ui-monospace, 'JetBrains Mono', monospace;
}
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}
.data-table th, .data-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid var(--border);
}
.data-table th {
  color: var(--muted);
  font-weight: 600;
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  font-family: ui-monospace, 'JetBrains Mono', monospace;
  background: var(--bg);
}
.data-table tbody tr:hover {
  background: var(--accent-soft);
}
.data-table .num {
  font-family: ui-monospace, 'JetBrains Mono', monospace;
  font-variant-numeric: tabular-nums;
}
.data-table .text-right {
  text-align: right;
}
.data-table .text-center {
  text-align: center;
}
</style>
