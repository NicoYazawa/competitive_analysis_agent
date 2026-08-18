<template>
  <div class="dashboard fade-in">
    <div class="page-header">
      <h1>监控总览</h1>
      <p>实时追踪竞品动态，辅助出海选品与定价决策</p>
    </div>

    <!-- KPI Grid -->
    <div class="kpi-grid" style="margin-bottom: 32px;">
      <div class="kpi-card">
        <div class="kpi-label">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="8" r="4"/><path d="M6 20v-2a4 4 0 0 1 4-4h4a4 4 0 0 1 4 4v2"/>
          </svg>
          监控竞品数
        </div>
        <div class="kpi-value">{{ kpis.competitorsCount }}<span>家</span></div>
        <div class="kpi-delta up">↑ {{ kpis.competitorsDelta }}%</div>
      </div>
      <div class="kpi-card warn">
        <div class="kpi-label">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="22 7 13.5 15.5 8.5 10.5 2 17"/><polyline points="16 7 22 7 22 13"/>
          </svg>
          本月价格变动
        </div>
        <div class="kpi-value">{{ kpis.priceChanges.toLocaleString() }}<span>次</span></div>
        <div class="kpi-delta down">↓ {{ Math.abs(kpis.priceChangesDelta) }}%</div>
      </div>
      <div class="kpi-card danger">
        <div class="kpi-label">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
            <line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/>
          </svg>
          活跃预警数
        </div>
        <div class="kpi-value">{{ kpis.activeAlerts }}</div>
        <div class="kpi-delta down">↑ {{ kpis.activeAlertsNew }} 新增</div>
      </div>
      <div class="kpi-card">
        <div class="kpi-label">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="18" height="18" rx="2"/><path d="M3 9h18M9 21V9"/>
          </svg>
          覆盖市场
        </div>
        <div class="kpi-value">{{ kpis.markets }}<span>个</span></div>
        <div class="kpi-delta neutral">{{ kpis.marketNames }}</div>
      </div>
    </div>

    <!-- Main Chart + Sidebar -->
    <div class="grid-2" style="margin-bottom: 32px;">
      <TrendChart
        title="价格与评分趋势"
        subtitle="近30天 · 全部监控品类"
        :legend="[
          { label: '平均价格指数', color: '#2a78d6' },
          { label: '品类评分均值', color: '#eb6834' }
        ]"
        :xData="store.trend30Days.xData"
        :series="[
          { name: '平均价格指数', data: store.trend30Days.priceIndex, color: '#2a78d6' },
          { name: '品类评分均值', data: store.trend30Days.scoreAvg, color: '#eb6834' }
        ]"
      />

      <div style="display: flex; flex-direction: column; gap: 16px;">
        <!-- Market Overview -->
        <div class="card">
          <div class="card-header">
            <span class="card-title">市场概览</span>
            <AlertBadge type="info" dot>实时</AlertBadge>
          </div>
          <div style="display: flex; flex-direction: column; gap: 12px;">
            <div v-for="market in store.marketOverview" :key="market.name"
                 style="display: flex; justify-content: space-between; align-items: center; padding: 10px 0; border-bottom: 1px solid var(--border);">
              <div>
                <div style="font-size: 14px; font-weight: 600;">{{ market.name }}</div>
                <div style="font-size: 12px; color: var(--muted);">竞争激烈度：{{ market.competition }}</div>
              </div>
              <div style="text-align: right;">
                <div class="num" :style="{ fontSize: '16px', fontWeight: '700', color: market.growth.startsWith('+') && parseFloat(market.growth) > 15 ? 'var(--status-good)' : 'var(--accent)' }">
                  {{ market.growth }}
                </div>
                <div style="font-size: 11px; color: var(--muted);">月增长率</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Top Categories -->
        <div class="card">
          <div class="card-header">
            <span class="card-title">TOP 5 热门品类</span>
          </div>
          <div style="display: flex; flex-direction: column; gap: 10px;">
            <div v-for="(cat, idx) in store.hotCategories" :key="cat.name"
                 style="display: flex; align-items: center; gap: 10px;">
              <span :class="['score-rank', idx === 0 ? 'gold' : idx === 1 ? 'silver' : idx === 2 ? 'bronze' : '']">
                {{ idx + 1 }}
              </span>
              <div style="flex: 1;">
                <div style="font-size: 13px; font-weight: 600;">{{ cat.name }}</div>
                <div style="height: 4px; background: var(--border); border-radius: 999px; margin-top: 4px;">
                  <div :style="{ width: cat.score + '%', height: '4px', background: 'var(--accent)', borderRadius: '999px' }"></div>
                </div>
              </div>
              <span class="num" style="font-size: 13px; color: var(--muted);">{{ cat.score }}%</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Latest Alerts -->
    <div class="card">
      <div class="card-header">
        <span class="card-title">最新预警</span>
        <span class="filter-chip active">全部 {{ kpis.activeAlerts }}</span>
      </div>
      <table class="data-table">
        <thead>
          <tr>
            <th>竞品</th>
            <th>信号类型</th>
            <th>时间</th>
            <th>风险等级</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="alert in store.latestAlerts" :key="alert.id"
              :class="{ 'row-highlight-critical': alert.riskLevel === 'critical', 'row-highlight-warning': alert.riskLevel === 'warning' }">
            <td style="font-weight: 600;">{{ alert.competitor }}</td>
            <td>{{ alert.signalType }}</td>
            <td class="num">{{ alert.time }}</td>
            <td>
              <AlertBadge :type="alert.riskLevel" dot>{{ alert.riskText }}</AlertBadge>
            </td>
            <td>
              <button class="action-btn">查看</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useMarketStore } from '@/stores/market'
import TrendChart from '@/components/TrendChart.vue'
import AlertBadge from '@/components/AlertBadge.vue'

const store = useMarketStore()
const kpis = computed(() => store.kpis)

onMounted(() => {
  store.fetchDashboard()
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
.kpi-card.danger::before { background: var(--status-critical); }
.kpi-label {
  font-size: 12px;
  color: var(--muted);
  font-weight: 500;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
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

.score-rank {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: color-mix(in oklch, var(--accent) 12%, transparent);
  color: var(--accent);
  font-size: 12px;
  font-weight: 700;
  display: grid;
  place-items: center;
  flex-shrink: 0;
}
.score-rank.gold { background: #fef3c7; color: #92400e; }
.score-rank.silver { background: #f1f5f9; color: #475569; }
.score-rank.bronze { background: #fff7ed; color: #9a3412; }

.filter-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border: 1px solid var(--border);
  border-radius: 999px;
  font-size: 11px;
  color: var(--muted);
  background: var(--surface);
  cursor: pointer;
}
.filter-chip.active {
  border-color: var(--accent);
  color: var(--accent);
  background: var(--accent-soft);
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
.row-highlight-warning {
  background: color-mix(in oklch, var(--status-warn) 8%, transparent) !important;
}
.row-highlight-critical {
  background: color-mix(in oklch, var(--status-critical) 8%, transparent) !important;
}
.action-btn {
  background: var(--accent-soft);
  color: var(--accent);
  border: none;
  padding: 4px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
}
</style>
