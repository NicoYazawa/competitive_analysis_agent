<template>
  <div class="supply-chain">
    <h1>Supply Chain Alerts</h1>
    <div class="alert-list">
      <div class="loading">No alerts</div>
    </div>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.supply-chain {
  padding: 24px;
}
.alert-list {
  background: white;
  border-radius: 8px;
  margin-top: 16px;
}
.loading {
  padding: 40px;
  text-align: center;
  color: #666;
}
</style>
