<template>
  <div class="card" style="padding:0;overflow:hidden;">
    <div v-if="title" class="card-header">
      <span class="card-title">{{ title }}</span>
      <slot name="header-actions"></slot>
    </div>
    <table class="data-table">
      <thead>
        <tr>
          <th v-for="col in columns" :key="col.key"
              :class="{ 'text-right': col.align === 'right', 'text-center': col.align === 'center' }">
            {{ col.label }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(row, idx) in data" :key="idx"
            :class="{ 'row-highlight-warning': row._highlight === 'warning', 'row-highlight-critical': row._highlight === 'critical' }">
          <td v-for="col in columns" :key="col.key"
              :class="{ 'text-right': col.align === 'right', 'text-center': col.align === 'center' }">
            <template v-if="col.slot">
              <slot :name="col.slot" :row="row" :value="row[col.key]"></slot>
            </template>
            <template v-else-if="col.format === 'currency'">
              <span class="num">${{ row[col.key] }}</span>
            </template>
            <template v-else-if="col.format === 'number'">
              <span class="num">{{ formatNumber(row[col.key]) }}</span>
            </template>
            <template v-else>
              {{ row[col.key] }}
            </template>
          </td>
        </tr>
      </tbody>
    </table>
    <div v-if="pagination" class="pagination">
      <span class="pagination-info">显示 {{ startIndex }}-{{ endIndex }} / 共 {{ total }} 条</span>
      <div class="pagination-buttons">
        <button class="pagination-btn" :disabled="currentPage === 1" @click="currentPage--">上一页</button>
        <button v-for="p in visiblePages" :key="p"
                :class="['pagination-btn', { active: p === currentPage, ellipsis: p === '...' }]"
                :disabled="p === '...'" @click="typeof p === 'number' && (currentPage = p)">
          {{ p }}
        </button>
        <button class="pagination-btn" :disabled="currentPage === totalPages" @click="currentPage++">下一页</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

export interface TableColumn {
  key: string
  label: string
  align?: 'left' | 'center' | 'right'
  format?: 'currency' | 'number'
  slot?: string
}

const props = withDefaults(defineProps<{
  columns: TableColumn[]
  data: Record<string, any>[]
  title?: string
  pagination?: boolean
  pageSize?: number
  total?: number
}>(), {
  pagination: false,
  pageSize: 10,
  total: 0
})

const currentPage = ref(1)

const totalPages = computed(() => Math.ceil(props.total / props.pageSize) || 1)
const startIndex = computed(() => (currentPage.value - 1) * props.pageSize + 1)
const endIndex = computed(() => Math.min(currentPage.value * props.pageSize, props.total))

const visiblePages = computed(() => {
  const pages: (number | string)[] = []
  const total = totalPages.value
  const current = currentPage.value

  if (total <= 5) {
    for (let i = 1; i <= total; i++) pages.push(i)
  } else {
    if (current <= 3) {
      pages.push(1, 2, 3, 4, '...')
    } else if (current >= total - 2) {
      pages.push('...', total - 3, total - 2, total - 1, total)
    } else {
      pages.push('...', current - 1, current, current + 1, '...')
    }
  }
  return pages
})

const formatNumber = (num: number) => {
  if (!num) return '0'
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}
</script>

<style scoped>
.card {
  background: #ffffff;
  border: 1px solid #e5e5e5;
  border-radius: 16px;
  padding: 24px;
}
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
.card-title {
  font-size: 13px;
  font-weight: 600;
  color: #6b6b6b;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  font-family: ui-monospace, 'JetBrains Mono', monospace;
}
.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}
.data-table th, .data-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #e5e5e5;
}
.data-table th {
  color: #6b6b6b;
  font-weight: 600;
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  font-family: ui-monospace, 'JetBrains Mono', monospace;
  background: #fafafa;
  position: sticky;
  top: 0;
}
.data-table tbody tr:hover {
  background: color-mix(in oklch, #2f6feb 8%, transparent);
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
.row-highlight-warning {
  background: color-mix(in oklch, #fab219 8%, transparent) !important;
}
.row-highlight-critical {
  background: color-mix(in oklch, #d03b3b 8%, transparent) !important;
}
.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-top: 1px solid #e5e5e5;
}
.pagination-info {
  font-size: 13px;
  color: #6b6b6b;
}
.pagination-buttons {
  display: flex;
  gap: 4px;
}
.pagination-btn {
  padding: 6px 12px;
  border: 1px solid #e5e5e5;
  border-radius: 8px;
  background: #ffffff;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.15s;
}
.pagination-btn:hover:not(:disabled) {
  border-color: #2f6feb;
  color: #2f6feb;
}
.pagination-btn.active {
  background: #2f6feb;
  color: #ffffff;
  border-color: #2f6feb;
}
.pagination-btn:disabled {
  color: #6b6b6b;
  cursor: not-allowed;
}
.pagination-btn.ellipsis {
  border: none;
  background: transparent;
  cursor: default;
}
</style>
