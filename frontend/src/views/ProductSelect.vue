<template>
  <div class="product-select">
    <h1>Product Selection</h1>
    <div class="filters">
      <input v-model="searchQuery" placeholder="Search products..." class="search-input" />
    </div>
    <div class="product-list">
      <div v-if="loading" class="loading">Loading...</div>
      <div v-else-if="products.length === 0" class="empty">No products found</div>
      <div v-else v-for="product in products" :key="product.id" class="product-item">
        {{ product.name }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMarketStore } from '@/stores/market'

const store = useMarketStore()
const searchQuery = ref('')
const products = ref([])
const loading = ref(false)

onMounted(() => {
  store.fetchProducts()
})
</script>

<style scoped>
.product-select {
  padding: 24px;
}
.filters {
  margin-bottom: 16px;
}
.search-input {
  padding: 8px 16px;
  border: 1px solid #ddd;
  border-radius: 4px;
  width: 100%;
  max-width: 400px;
}
.product-list {
  background: white;
  border-radius: 8px;
  padding: 16px;
}
.loading, .empty {
  padding: 40px;
  text-align: center;
  color: #666;
}
</style>
