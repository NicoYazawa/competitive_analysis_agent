<template>
  <span :class="['alert-badge', type]">
    <span v-if="dot" class="alert-dot"></span>
    <slot />
  </span>
</template>

<script setup lang="ts">
defineProps<{
  type?: 'critical' | 'warning' | 'good' | 'info'
  dot?: boolean
}>()
</script>

<style scoped>
.alert-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 10px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 700;
  font-family: ui-monospace, 'JetBrains Mono', monospace;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}
.alert-badge.critical {
  background: color-mix(in oklch, #d03b3b 14%, transparent);
  color: #d03b3b;
}
.alert-badge.warning {
  background: color-mix(in oklch, #fab219 14%, transparent);
  color: #92400e;
}
.alert-badge.good {
  background: color-mix(in oklch, #0ca30c 14%, transparent);
  color: #0ca30c;
}
.alert-badge.info {
  background: color-mix(in oklch, #2f6feb 12%, transparent);
  color: #2f6feb;
}
.alert-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
}
</style>
