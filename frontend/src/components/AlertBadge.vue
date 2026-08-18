<template>
  <div :class="['alert-badge', type]">
    {{ count }}
  </div>
</template>

<script setup lang="ts">
defineProps<{
  count: number
  type: 'warning' | 'danger' | 'info' | 'success'
}>()
</script>

<style scoped>
.alert-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: bold;
}
.warning {
  background: #fff3cd;
  color: #856404;
}
.danger {
  background: #f8d7da;
  color: #721c24;
}
.info {
  background: #d1ecf1;
  color: #0c5460;
}
.success {
  background: #d4edda;
  color: #155724;
}
</style>
