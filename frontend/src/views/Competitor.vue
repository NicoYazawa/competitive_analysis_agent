<template>
  <div class="competitor">
    <h1>Competitor Intelligence</h1>
    <AlertBadge :count="0" type="info" />
    <div class="competitor-list">
      <div v-for="comp in competitors" :key="comp.id" class="competitor-item">
        <span>{{ comp.name }}</span>
        <span class="price">${{ comp.currentPrice }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMarketStore } from '@/stores/market'
import AlertBadge from '@/components/AlertBadge.vue'

const store = useMarketStore()
const competitors = ref([])

onMounted(() => {
  store.fetchCompetitors()
})
</script>

<style scoped>
.competitor {
  padding: 24px;
}
.competitor-list {
  background: white;
  border-radius: 8px;
  margin-top: 16px;
}
.competitor-item {
  padding: 12px 16px;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
}
.price {
  font-weight: bold;
  color: #2c5282;
}
</style>
