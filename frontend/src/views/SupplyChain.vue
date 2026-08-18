<template>
  <div class="supply-chain fade-in">
    <div class="page-header">
      <h1>供应链预警</h1>
      <p>原材料/物流/政策风险信号早知道</p>
    </div>

    <!-- KPI Row -->
    <div class="kpi-grid" style="margin-bottom: 32px;">
      <div class="kpi-card danger">
        <div class="kpi-label">高风险信号</div>
        <div class="kpi-value">{{ kpis.highRisk }}</div>
        <div class="kpi-delta down">需立即处理</div>
      </div>
      <div class="kpi-card warn">
        <div class="kpi-label">中风险信号</div>
        <div class="kpi-value">{{ kpis.mediumRisk }}</div>
        <div class="kpi-delta neutral">本周关注</div>
      </div>
      <div class="kpi-card">
        <div class="kpi-label">已化解风险</div>
        <div class="kpi-value">{{ kpis.resolved }}</div>
        <div class="kpi-delta up">本月</div>
      </div>
      <div class="kpi-card">
        <div class="kpi-label">监控品类</div>
        <div class="kpi-value">{{ kpis.categories }}</div>
        <div class="kpi-delta neutral">全部品类</div>
      </div>
    </div>

    <!-- Supply Alert Table -->
    <div class="card" style="padding: 0; overflow: hidden; margin-bottom: 32px;">
      <div style="padding: 16px 20px; border-bottom: 1px solid var(--border); display: flex; align-items: center; justify-content: space-between;">
        <span class="card-title">预警列表</span>
        <div style="display: flex; gap: 8px;">
          <select v-model="riskFilter" class="filter-select" style="font-size: 12px; padding: 6px 10px;">
            <option value="">全部风险等级</option>
            <option>严重</option>
            <option>警告</option>
            <option>低风险</option>
          </select>
          <select v-model="typeFilter" class="filter-select" style="font-size: 12px; padding: 6px 10px;">
            <option value="">全部类型</option>
            <option>原材料</option>
            <option>物流</option>
            <option>政策</option>
            <option>竞品产能</option>
          </select>
        </div>
      </div>
      <table class="data-table">
        <thead>
          <tr>
            <th>风险等级</th>
            <th>品类</th>
            <th>信号类型</th>
            <th>详情</th>
            <th>影响市场</th>
            <th>发现时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="alert in filteredAlerts" :key="alert.id"
              :class="{ 'row-highlight-critical': alert.riskLevel === 'critical', 'row-highlight-warning': alert.riskLevel === 'warning' }">
            <td>
              <AlertBadge :type="alert.riskLevel" dot>
                {{ alert.riskLevel === 'critical' ? '严重' : alert.riskLevel === 'warning' ? '警告' : '低风险' }}
              </AlertBadge>
            </td>
            <td style="font-weight: 600;">{{ alert.category }}</td>
            <td>{{ alert.signalType }}</td>
            <td style="font-size: 13px;">{{ alert.detail }}</td>
            <td>{{ alert.market }}</td>
            <td class="num">{{ alert.time }}</td>
            <td>
              <button class="action-btn">{{ alert.riskLevel === 'info' ? '已处理' : '分析' }}</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Risk Radar Cards -->
    <div class="grid-3">
      <!-- Raw Materials -->
      <div class="card">
        <div class="card-header">
          <span class="card-title">原材料风险</span>
        </div>
        <div style="display: flex; flex-direction: column; gap: 12px;">
          <div v-for="item in store.riskRadar.rawMaterials" :key="item.name">
            <div style="display: flex; justify-content: space-between; font-size: 13px; margin-bottom: 6px;">
              <span>{{ item.name }}</span>
              <span class="num" :style="{ color: item.level === 'critical' ? 'var(--status-critical)' : item.level === 'warning' ? 'var(--status-warn)' : 'var(--status-good)', fontWeight: '700' }">
                {{ item.change || item.status }}
              </span>
            </div>
            <div style="height: 6px; background: var(--border); border-radius: 999px;">
              <div :style="{ width: item.percent + '%', height: '6px', background: item.level === 'critical' ? 'var(--status-critical)' : item.level === 'warning' ? 'var(--status-warn)' : 'var(--status-good)', borderRadius: '999px' }"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Logistics -->
      <div class="card">
        <div class="card-header">
          <span class="card-title">物流风险</span>
        </div>
        <div style="display: flex; flex-direction: column; gap: 12px;">
          <div v-for="item in store.riskRadar.logistics" :key="item.name">
            <div style="display: flex; justify-content: space-between; font-size: 13px; margin-bottom: 6px;">
              <span>{{ item.name }}</span>
              <span class="num" :style="{ color: item.level === 'critical' ? 'var(--status-critical)' : item.level === 'warning' ? 'var(--status-warn)' : 'var(--status-good)', fontWeight: '700' }">
                {{ item.status }}
              </span>
            </div>
            <div style="height: 6px; background: var(--border); border-radius: 999px;">
              <div :style="{ width: item.percent + '%', height: '6px', background: item.level === 'critical' ? 'var(--status-critical)' : item.level === 'warning' ? 'var(--status-warn)' : 'var(--status-good)', borderRadius: '999px' }"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Policy -->
      <div class="card">
        <div class="card-header">
          <span class="card-title">政策风险</span>
        </div>
        <div style="display: flex; flex-direction: column; gap: 12px;">
          <div v-for="item in store.riskRadar.policy" :key="item.name">
            <div style="display: flex; justify-content: space-between; font-size: 13px; margin-bottom: 6px;">
              <span>{{ item.name }}</span>
              <span class="num" :style="{ color: item.level === 'warning' ? 'var(--status-warn)' : 'var(--status-good)', fontWeight: '700' }">
                {{ item.cost || item.status }}
              </span>
            </div>
            <div style="height: 6px; background: var(--border); border-radius: 999px;">
              <div :style="{ width: item.percent + '%', height: '6px', background: item.level === 'warning' ? 'var(--status-warn)' : 'var(--status-good)', borderRadius: '999px' }"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useMarketStore } from '@/stores/market'
import AlertBadge from '@/components/AlertBadge.vue'

const store = useMarketStore()
const kpis = computed(() => store.supplyKpis)

onMounted(() => {
  store.fetchSupplyAlerts()
})

const riskFilter = ref('')
const typeFilter = ref('')

const filteredAlerts = computed(() => {
  let result = [...store.supplyAlerts]

  if (riskFilter.value) {
    const levelMap: Record<string, string> = { '严重': 'critical', '警告': 'warning', '低风险': 'info' }
    result = result.filter(a => a.riskLevel === levelMap[riskFilter.value])
  }

  if (typeFilter.value) {
    result = result.filter(a => a.signalType.includes(typeFilter.value))
  }

  return result
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
}
.kpi-value {
  font-family: ui-monospace, 'JetBrains Mono', monospace;
  font-size: 32px;
  font-weight: 700;
  letter-spacing: -0.03em;
  color: var(--fg);
  line-height: 1;
}
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

.filter-select {
  padding: 6px 10px;
  border: 1px solid var(--border);
  border-radius: 8px;
  background: var(--surface);
  font: inherit;
  font-size: 12px;
  color: var(--fg);
  cursor: pointer;
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
