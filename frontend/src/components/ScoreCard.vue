<template>
  <div class="score-card">
    <div class="icon">{{ icon }}</div>
    <div class="content">
      <div class="title">{{ title }}</div>
      <div class="value">{{ value }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  value: number
  icon: string
}>()
</script>

<style scoped>
.score-card {
  background: white;
  border-radius: 8px;
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}
.icon {
  font-size: 32px;
}
.title {
  font-size: 12px;
  color: #666;
  text-transform: uppercase;
}
.value {
  font-size: 24px;
  font-weight: bold;
  color: #2c5282;
}
</style>
