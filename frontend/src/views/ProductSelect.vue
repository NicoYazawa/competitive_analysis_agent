<template>
  <div class="product-select fade-in">
    <div class="page-header">
      <h1>选品分析</h1>
      <p>AI推荐产品排行，辅助选品决策</p>
    </div>

    <!-- Score Ranking Cards -->
    <div class="section-head">
      <h2>产品评分排名</h2>
      <div style="display: flex; gap: 8px;">
        <select v-model="categoryFilter" class="filter-select" style="font-size: 13px;">
          <option value="">全部品类</option>
          <option>智能手表</option>
          <option>无线耳机</option>
          <option>便携储能</option>
          <option>充电器</option>
        </select>
      </div>
    </div>

    <div class="grid-3" style="margin-bottom: 40px;">
      <ScoreCard
        v-for="(product, idx) in filteredProducts"
        :key="product.id"
        :name="product.name"
        :rank="idx + 1"
        :score="product.score"
        :rating="product.rating"
        :reviews="product.reviews"
        :price="product.price"
        :badge="product.badge"
        :badgeType="product.badgeType"
        :description="product.description"
        hoverable
      />
    </div>

    <!-- Price Comparison Table -->
    <div class="card">
      <div class="card-header">
        <span class="card-title">价格对比表</span>
        <span style="font-size: 12px; color: var(--muted);">数据更新于 2026-08-18 15:00</span>
      </div>
      <table class="data-table">
        <thead>
          <tr>
            <th>产品</th>
            <th class="text-center">品类</th>
            <th class="text-right">最低价</th>
            <th class="text-right">最高价</th>
            <th class="text-right">均价</th>
            <th class="text-center">价格趋势</th>
            <th class="text-right">建议售价</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in store.priceComparison" :key="item.id">
            <td style="font-weight: 600;">{{ item.name }}</td>
            <td class="text-center">{{ item.category }}</td>
            <td class="num text-right">${{ item.low }}</td>
            <td class="num text-right">${{ item.high }}</td>
            <td class="num text-right">${{ item.avg }}</td>
            <td class="text-center">
              <svg width="80" height="24" viewBox="0 0 80 24">
                <polyline
                  v-if="item.trend === 'up'"
                  points="0,18 13,16 26,14 39,12 52,10 65,11 78,8"
                  fill="none"
                  stroke="#1baf7a"
                  stroke-width="1.5"
                />
                <polyline
                  v-else-if="item.trend === 'down'"
                  points="0,8 13,10 26,9 39,14 52,16 65,15 78,17"
                  fill="none"
                  stroke="#eb6834"
                  stroke-width="1.5"
                />
                <polyline
                  v-else
                  points="0,12 13,11 26,13 39,12 52,14 65,13 78,14"
                  fill="none"
                  stroke="#2a78d6"
                  stroke-width="1.5"
                />
              </svg>
            </td>
            <td class="num text-right" style="color: var(--accent); font-weight: 700;">${{ item.suggestedPrice }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useMarketStore } from '@/stores/market'
import ScoreCard from '@/components/ScoreCard.vue'

const store = useMarketStore()
const categoryFilter = ref('')

onMounted(() => {
  store.fetchProductScores()
})

const filteredProducts = computed(() => {
  if (!categoryFilter.value) return store.productScores
  return store.productScores.filter(p => p.category === categoryFilter.value)
})
</script>

<style scoped>
.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
  gap: 12px;
}
.section-head h2 {
  font-size: 16px;
  font-weight: 600;
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
