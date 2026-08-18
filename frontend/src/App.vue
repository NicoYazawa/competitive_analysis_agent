<template>
  <div id="app">
    <TopNav />
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import TopNav from '@/components/TopNav.vue'
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

:root {
  --bg: #fafafa;
  --surface: #ffffff;
  --fg: #111111;
  --muted: #6b6b6b;
  --border: #e5e5e5;
  --accent: #2f6feb;
  --accent-soft: color-mix(in oklch, #2f6feb 12%, transparent);
  --status-good: #0ca30c;
  --status-warn: #fab219;
  --status-serious: #ec835a;
  --status-critical: #d03b3b;
  --series-1: #2a78d6;
  --series-2: #eb6834;
  --series-3: #1baf7a;
}

body {
  font-family: 'Inter', -apple-system, system-ui, sans-serif;
  font-size: 16px;
  line-height: 1.5;
  background: var(--bg);
  color: var(--fg);
  -webkit-font-smoothing: antialiased;
}

#app {
  min-height: 100vh;
  background-color: var(--bg);
}

.main-content {
  max-width: 1280px;
  margin-inline: auto;
  padding-inline: 24px;
  padding-block: 32px 48px;
}

.page-header {
  padding-block: 0 8px;
}
.page-header h1 {
  font-family: 'Inter', -apple-system, system-ui, sans-serif;
  font-size: 28px;
  font-weight: 700;
  letter-spacing: -0.02em;
}
.page-header p {
  color: var(--muted);
  font-size: 14px;
  margin-top: 4px;
}

.grid-2 {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 32px;
}
.grid-3 {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}
.grid-4 {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

@media (max-width: 1024px) {
  .grid-2 {
    grid-template-columns: 1fr;
  }
}
@media (max-width: 640px) {
  .grid-3, .grid-4 {
    grid-template-columns: 1fr;
  }
}

/* Animations */
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(6px); }
  to { opacity: 1; transform: translateY(0); }
}
.fade-in {
  animation: fadeIn 0.25s ease both;
}

.card {
  background: var(--surface);
  border: 1px solid var(--border);
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
  color: var(--muted);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  font-family: ui-monospace, 'JetBrains Mono', monospace;
}
</style>
