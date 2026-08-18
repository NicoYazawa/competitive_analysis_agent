<template>
  <div :class="['score-card', { hoverable }]">
    <div class="score-card-top">
      <span v-if="rank" :class="['score-rank', rankClass]">{{ rank }}</span>
      <AlertBadge v-if="badge" :type="badgeType" dot>{{ badge }}</AlertBadge>
    </div>
    <div class="score-name">{{ name }}</div>
    <div class="score-meta">
      <span v-if="rating">⭐ {{ rating }}</span>
      <span v-if="reviews">评论 {{ formatNumber(reviews) }}</span>
      <span v-if="price">均价 ${{ price }}</span>
    </div>
    <div v-if="score !== undefined" class="score-bar-wrap">
      <div class="score-bar-bg">
        <div class="score-bar-fill" :style="{ width: score + '%' }"></div>
      </div>
      <span class="score-num">{{ score }}</span>
    </div>
    <div v-if="description" class="score-description">{{ description }}</div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import AlertBadge from './AlertBadge.vue'

const props = defineProps<{
  name: string
  rank?: number
  score?: number
  rating?: number
  reviews?: number
  price?: number
  badge?: string
  badgeType?: 'critical' | 'warning' | 'good' | 'info'
  description?: string
  hoverable?: boolean
}>()

const rankClass = computed(() => {
  if (props.rank === 1) return 'gold'
  if (props.rank === 2) return 'silver'
  if (props.rank === 3) return 'bronze'
  return ''
})

const formatNumber = (num: number) => {
  if (num >= 1000) return (num / 1000).toFixed(1) + 'k'
  return num.toString()
}
</script>

<style scoped>
.score-card {
  background: #ffffff;
  border: 1px solid #e5e5e5;
  border-radius: 12px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  transition: border-color 0.15s ease;
}
.score-card.hoverable:hover {
  border-color: #2f6feb;
}
.score-card-top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
}
.score-rank {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: color-mix(in oklch, #2f6feb 12%, transparent);
  color: #2f6feb;
  font-size: 12px;
  font-weight: 700;
  display: grid;
  place-items: center;
}
.score-rank.gold {
  background: #fef3c7;
  color: #92400e;
}
.score-rank.silver {
  background: #f1f5f9;
  color: #475569;
}
.score-rank.bronze {
  background: #fff7ed;
  color: #9a3412;
}
.score-name {
  font-size: 14px;
  font-weight: 600;
  color: #111111;
  margin-top: 4px;
}
.score-meta {
  font-size: 12px;
  color: #6b6b6b;
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}
.score-bar-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}
.score-bar-bg {
  flex: 1;
  height: 6px;
  background: #e5e5e5;
  border-radius: 999px;
  overflow: hidden;
}
.score-bar-fill {
  height: 100%;
  border-radius: 999px;
  background: #2f6feb;
  transition: width 0.3s ease;
}
.score-num {
  font-family: ui-monospace, 'JetBrains Mono', monospace;
  font-size: 13px;
  font-weight: 700;
  color: #2f6feb;
  min-width: 32px;
  text-align: right;
}
.score-description {
  font-size: 12px;
  color: #6b6b6b;
  margin-top: 4px;
}
</style>
