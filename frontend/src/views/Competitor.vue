<template>
  <div class="competitor fade-in">
    <div class="page-header">
      <h1>竞品情报</h1>
      <p>全量竞品库，支持按品类/价格/评分多维筛选</p>
    </div>

    <!-- Filter Bar -->
    <div class="filter-bar">
      <div class="filter-search">
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="var(--muted)" stroke-width="2">
          <circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/>
        </svg>
        <input v-model="searchQuery" type="text" placeholder="搜索竞品名称…" />
      </div>
      <select v-model="categoryFilter" class="filter-select">
        <option value="">全部品类</option>
        <option>智能手表</option>
        <option>无线耳机</option>
        <option>便携储能</option>
        <option>充电器</option>
      </select>
      <select v-model="marketFilter" class="filter-select">
        <option value="">全部市场</option>
        <option>北美</option>
        <option>欧洲</option>
        <option>东南亚</option>
      </select>
      <select v-model="sortBy" class="filter-select">
        <option value="rating">评分排序</option>
        <option value="price-asc">价格从低到高</option>
        <option value="reviews">评论数最多</option>
        <option value="newest">最新上架</option>
      </select>
      <div style="margin-left: auto; display: flex; gap: 8px; align-items: center;">
        <span :class="['filter-chip', { active: statusFilter === '' }]" @click="statusFilter = ''">全部</span>
        <span :class="['filter-chip', { active: statusFilter === 'alert' }]" @click="statusFilter = 'alert'">有预警</span>
        <span :class="['filter-chip', { active: statusFilter === 'new' }]" @click="statusFilter = 'new'">新品</span>
      </div>
    </div>

    <!-- Competitor Table -->
    <div class="card" style="padding: 0; overflow: hidden;">
      <table class="data-table">
        <thead>
          <tr>
            <th>竞品名称</th>
            <th>品类</th>
            <th class="text-center">评分</th>
            <th class="text-right">价格区间</th>
            <th class="text-right">评论数</th>
            <th class="text-center">状态</th>
            <th class="text-right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="comp in filteredCompetitors" :key="comp.id">
            <td>
              <div style="font-weight: 600; font-size: 14px;">{{ comp.name }}</div>
              <div style="font-size: 11px; color: var(--muted); font-family: ui-monospace, monospace;">
                {{ comp.platform }} · {{ comp.region }}
              </div>
            </td>
            <td>{{ comp.category }}</td>
            <td class="text-center">
              <div style="display: flex; align-items: center; justify-content: center; gap: 6px;">
                <span style="font-family: ui-monospace; font-weight: 700; color: var(--status-good);">{{ comp.rating }}</span>
                <svg width="12" height="12" viewBox="0 0 24 24" fill="var(--status-warn)" stroke="none">
                  <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
                </svg>
              </div>
            </td>
            <td class="num text-right">${{ comp.currentPrice.toFixed(2) }}</td>
            <td class="num text-right">{{ comp.reviewCount.toLocaleString() }}</td>
            <td class="text-center">
              <AlertBadge :type="comp.status" dot>{{ comp.statusText }}</AlertBadge>
            </td>
            <td class="text-right">
              <button class="action-btn">详情</button>
            </td>
          </tr>
        </tbody>
      </table>
      <!-- Pagination -->
      <div class="pagination">
        <span class="pagination-info">显示 1-{{ filteredCompetitors.length }} / 共 {{ filteredCompetitors.length }} 条</span>
        <div class="pagination-buttons">
          <button class="pagination-btn" disabled>上一页</button>
          <button class="pagination-btn active">1</button>
          <button class="pagination-btn">2</button>
          <button class="pagination-btn">3</button>
          <button class="pagination-btn">…</button>
          <button class="pagination-btn">下一页</button>
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

onMounted(() => {
  store.fetchCompetitors()
})

const searchQuery = ref('')
const categoryFilter = ref('')
const marketFilter = ref('')
const sortBy = ref('rating')
const statusFilter = ref('')

const filteredCompetitors = computed(() => {
  let result = [...store.competitors]

  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(c => c.name.toLowerCase().includes(q))
  }

  if (categoryFilter.value) {
    result = result.filter(c => c.category === categoryFilter.value)
  }

  if (statusFilter.value === 'alert') {
    result = result.filter(c => c.status === 'critical' || c.status === 'warning')
  } else if (statusFilter.value === 'new') {
    result = result.filter(c => c.status === 'info')
  }

  return result
})
</script>

<style scoped>
.filter-bar {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  align-items: center;
  margin-bottom: 20px;
}
.filter-search {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 14px;
  border: 1px solid var(--border);
  border-radius: 8px;
  background: var(--surface);
  flex: 1;
  min-width: 200px;
  max-width: 320px;
}
.filter-search input {
  border: none;
  background: transparent;
  outline: none;
  font: inherit;
  font-size: 14px;
  color: var(--fg);
  width: 100%;
}
.filter-search input::placeholder {
  color: var(--muted);
}
.filter-select {
  padding: 8px 14px;
  border: 1px solid var(--border);
  border-radius: 8px;
  background: var(--surface);
  font: inherit;
  font-size: 14px;
  color: var(--fg);
  cursor: pointer;
}
.filter-select:focus {
  outline: 2px solid var(--accent-soft);
  border-color: var(--accent);
}
.filter-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border: 1px solid var(--border);
  border-radius: 999px;
  font-size: 13px;
  color: var(--muted);
  background: var(--surface);
  cursor: pointer;
  transition: all 0.15s;
}
.filter-chip:hover,
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
  position: sticky;
  top: 0;
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

.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-top: 1px solid var(--border);
}
.pagination-info {
  font-size: 13px;
  color: var(--muted);
}
.pagination-buttons {
  display: flex;
  gap: 4px;
}
.pagination-btn {
  padding: 6px 12px;
  border: 1px solid var(--border);
  border-radius: 8px;
  background: var(--surface);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.15s;
}
.pagination-btn:hover:not(:disabled) {
  border-color: var(--accent);
  color: var(--accent);
}
.pagination-btn.active {
  background: var(--accent);
  color: #fff;
  border-color: var(--accent);
}
.pagination-btn:disabled {
  color: var(--muted);
  cursor: not-allowed;
}
</style>
