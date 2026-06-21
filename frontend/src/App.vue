<template>
  <div class="app-layout">
    <header class="app-header">
      <div class="header-content">
        <div class="logo-section">
          <el-icon class="logo-icon"><Fish /></el-icon>
          <h1 class="app-title">水族箱智能控制系统</h1>
        </div>
        <div class="header-right">
          <div class="current-time">{{ currentTime }}</div>
        </div>
      </div>
      <nav class="nav-tabs">
        <router-link
          v-for="route in navRoutes"
          :key="route.path"
          :to="route.path"
          class="nav-tab"
          :class="{ active: $route.path === route.path }"
        >
          <el-icon><component :is="route.icon" /></el-icon>
          <span>{{ route.name }}</span>
        </router-link>
      </nav>
    </header>
    <main class="app-main">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { Fish, DataLine, Sunny, Wind, Setting } from '@element-plus/icons-vue'
import dayjs from 'dayjs'

const route = useRoute()
const currentTime = ref('')

const navRoutes = [
  { path: '/', name: '仪表盘', icon: DataLine },
  { path: '/lighting', name: '光照排程', icon: Sunny },
  { path: '/oxygen', name: '溶氧量控制', icon: Wind },
  { path: '/device', name: '设备控制', icon: Setting }
]

let timer = null

const updateTime = () => {
  currentTime.value = dayjs().format('YYYY-MM-DD HH:mm:ss')
}

onMounted(() => {
  updateTime()
  timer = setInterval(updateTime, 1000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped lang="scss">
.app-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  background: rgba(12, 20, 69, 0.9);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}

.logo-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  font-size: 32px;
  color: #38bdf8;
}

.app-title {
  font-size: 24px;
  font-weight: 700;
  background: linear-gradient(135deg, #38bdf8 0%, #818cf8 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.current-time {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.8);
  font-family: 'Courier New', monospace;
}

.nav-tabs {
  display: flex;
  gap: 4px;
  padding: 0 24px;
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.nav-tab {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  color: rgba(255, 255, 255, 0.6);
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  border-bottom: 2px solid transparent;
  transition: all 0.3s ease;
}

.nav-tab:hover {
  color: rgba(255, 255, 255, 0.9);
  background: rgba(255, 255, 255, 0.05);
}

.nav-tab.active {
  color: #38bdf8;
  border-bottom-color: #38bdf8;
  background: rgba(56, 189, 248, 0.1);
}

.app-main {
  flex: 1;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
