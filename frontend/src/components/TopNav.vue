<template>
  <header class="topnav">
    <div class="container topnav-inner">
      <span class="logo">Compete<span>AI</span></span>
      <nav class="main-tabs">
        <router-link v-for="tab in tabs" :key="tab.path"
                     :to="tab.path"
                     :class="['tab-btn', { active: isActive(tab.path) }]">
          {{ tab.label }}
        </router-link>
      </nav>
      <div class="topnav-right">
        <button class="notif-btn" title="通知">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
            <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/>
            <path d="M13.73 21a2 2 0 0 1-3.46 0"/>
          </svg>
          <span class="notif-dot"></span>
        </button>
        <div class="avatar">N</div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()

const tabs = [
  { label: 'Dashboard', path: '/' },
  { label: '竞品情报', path: '/competitors' },
  { label: '选品分析', path: '/products' },
  { label: '定价策略', path: '/pricing' },
  { label: '供应链预警', path: '/supply-chain' }
]

const isActive = (path: string) => {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}
</script>

<style scoped>
.topnav {
  position: sticky;
  top: 0;
  z-index: 100;
  background: color-mix(in oklch, #fafafa 94%, transparent);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid #e5e5e5;
}
.container {
  max-width: 1280px;
  margin-inline: auto;
  padding-inline: 24px;
}
.topnav-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-block: 14px;
}
.logo {
  font-family: 'Inter', -apple-system, system-ui, sans-serif;
  font-size: 17px;
  font-weight: 700;
  letter-spacing: -0.01em;
  color: #111111;
}
.logo span {
  color: #2f6feb;
}
.main-tabs {
  display: flex;
  gap: 4px;
}
.tab-btn {
  padding: 8px 16px;
  border-radius: 8px;
  border: 1px solid transparent;
  background: transparent;
  color: #6b6b6b;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease;
  white-space: nowrap;
  text-decoration: none;
}
.tab-btn:hover {
  color: #111111;
  background: color-mix(in oklch, #2f6feb 12%, transparent);
}
.tab-btn.active {
  background: #2f6feb;
  color: #ffffff;
  border-color: #2f6feb;
}
.topnav-right {
  display: flex;
  align-items: center;
  gap: 12px;
}
.notif-btn {
  position: relative;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  border: 1px solid #e5e5e5;
  background: #ffffff;
  display: grid;
  place-items: center;
  cursor: pointer;
  transition: border-color 0.15s;
}
.notif-btn:hover {
  border-color: #2f6feb;
}
.notif-dot {
  position: absolute;
  top: 6px;
  right: 7px;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #d03b3b;
}
.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: color-mix(in oklch, #2f6feb 12%, transparent);
  display: grid;
  place-items: center;
  font-size: 13px;
  font-weight: 700;
  color: #2f6feb;
}
</style>
