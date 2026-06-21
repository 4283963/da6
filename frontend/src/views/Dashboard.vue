<template>
  <div class="page-container">
    <div v-if="deviceStatus.power_saving" class="night-mode-bar glass-card moon-glow">
      <div class="night-mode-content">
        <div class="moon-icon-wrapper">
          <el-icon class="moon-icon"><Moon /></el-icon>
        </div>
        <div class="night-info">
          <div class="night-title">深夜省电模式已激活 🌙</div>
          <div class="night-desc">
            当前时间 00:00-05:00 · 水温 {{ deviceStatus.current_temp?.toFixed(1) }}°C 在安全范围(22-28°C)，
            灯光和气泵功率自动减半运行
          </div>
        </div>
        <el-tag type="warning" effect="dark" size="large" class="night-tag">
          POWER SAVING
        </el-tag>
      </div>
    </div>

    <div class="stats-grid">
      <div class="stat-card glass-card" style="background: linear-gradient(135deg, rgba(56, 189, 248, 0.2), rgba(59, 130, 246, 0.1);">
        <div class="stat-header">
          <el-icon class="stat-icon" style="color: #38bdf8;"><Thermometer /></el-icon>
          <span class="stat-label">当前水温</span>
        </div>
        <div class="stat-value">{{ sensorStats.current_temperature || '--' }}°C</div>
        <div class="stat-sub">平均: {{ sensorStats.avg_temperature || '--' }}°C</div>
      </div>

      <div class="stat-card glass-card" style="background: linear-gradient(135deg, rgba(251, 191, 36, 0.2), rgba(245, 158, 11, 0.1);">
        <div class="stat-header">
          <el-icon class="stat-icon" style="color: #fbbf24;"><Lightning /></el-icon>
          <span class="stat-label">灯光瓦数</span>
        </div>
        <div class="stat-value">{{ sensorStats.current_light_wattage || 0 }}W</div>
        <div class="stat-sub">平均: {{ sensorStats.avg_light_wattage || 0 }}W</div>
      </div>

      <div class="stat-card glass-card" :class="{ 'power-saving-active': lightStatus.power_saving }"
           :style="lightStatus.power_saving 
             ? 'background: linear-gradient(135deg, rgba(251, 191, 36, 0.15), rgba(168, 85, 247, 0.1);'
             : 'background: linear-gradient(135deg, rgba(34, 197, 94, 0.2), rgba(16, 185, 129, 0.1);'">
        <div class="stat-header">
          <div class="stat-icon-wrap">
            <el-icon class="stat-icon" :style="{ color: lightStatus.power_saving ? '#fbbf24' : '#22c55e' }">
              <Sunny v-if="!lightStatus.power_saving" />
              <Moon v-else />
            </el-icon>
          </div>
          <span class="stat-label">灯光状态</span>
          <el-tag v-if="lightStatus.power_saving" size="small" type="warning" effect="dark" class="mini-moon-tag">
            <el-icon><Moon /></el-icon>
            省电
          </el-tag>
        </div>
        <div class="stat-value">
          <span v-if="lightStatus.is_on" :style="{ color: lightStatus.power_saving ? '#fbbf24' : '#4ade80' }">
            {{ lightStatus.power_saving ? '省电运行中' : '开启中' }}
          </span>
          <span v-else style="color: #6b7280;">已关闭</span>
        </div>
        <div class="stat-sub">
          <span>亮度: {{ lightStatus.brightness || 0 }}%</span>
          <span v-if="lightStatus.power_saving && lightStatus.original_brightness" class="saving-hint">
             (原 {{ lightStatus.original_brightness }}% )
          </span>
        </div>
      </div>

      <div class="stat-card glass-card" :class="{ 'power-saving-active': deviceStatus.power_saving && deviceStatus.current_pump_level > 0 }"
           :style="deviceStatus.power_saving && deviceStatus.current_pump_level > 0
             ? 'background: linear-gradient(135deg, rgba(251, 191, 36, 0.15), rgba(139, 92, 246, 0.15);'
             : 'background: linear-gradient(135deg, rgba(139, 92, 246, 0.2), rgba(124, 58, 237, 0.1);'">
        <div class="stat-header">
          <div class="stat-icon-wrap">
            <el-icon class="stat-icon" :style="{ color: deviceStatus.power_saving ? '#fbbf24' : '#8b5cf6' }">
              <Wind v-if="!deviceStatus.power_saving" />
              <Moon v-else />
            </el-icon>
          </div>
          <span class="stat-label">气泵档位</span>
          <el-tag v-if="deviceStatus.power_saving && deviceStatus.current_pump_level > 0" size="small" type="warning" effect="dark" class="mini-moon-tag">
            <el-icon><Moon /></el-icon>
            省电
          </el-tag>
        </div>
        <div class="stat-value">
          <span v-if="deviceStatus.current_pump_level > 0" :style="{ color: deviceStatus.power_saving ? '#fbbf24' : '#a78bfa' }">
            {{ deviceStatus.power_saving ? '省电 ' : '' }}{{ deviceStatus.current_pump_level }}档
          </span>
          <span v-else style="color: #6b7280;">关闭</span>
        </div>
        <div class="stat-sub">开启: {{ deviceStatus.pumps_on || 0 }}/{{ deviceStatus.total_pumps || 0 }}</div>
      </div>
    </div>

    <div class="dashboard-grid">
      <div class="glass-card panel">
        <h2 class="section-title">
          <el-icon><Timer /></el-icon>
          光照排程
        </h2>
        <div class="schedule-info">
          <div v-if="lightStatus.schedule_name" class="active-schedule">
            <el-icon style="color: #4ade80;"><VideoPlay /></el-icon>
            <span>当前执行: {{ lightStatus.schedule_name }}</span>
          </div>
          <div class="next-action">
            <span>下一步: {{ lightStatus.next_action || '--' }}</span>
            <span class="next-time">{{ lightStatus.next_time || '--' }}</span>
          </div>
        </div>
        <el-table :data="schedules.slice(0, 3)" class="data-table" style="margin-top: 16px;">
          <el-table-column prop="name" label="名称" />
          <el-table-column prop="start_time" label="开灯时间" width="120" />
          <el-table-column prop="end_time" label="关灯时间" width="120" />
          <el-table-column prop="brightness" label="亮度" width="100">
            <template #default="{ row }">
              <el-tag :type="row.brightness > 60 ? 'warning' : 'success'" size="small">
                {{ row.brightness }}%
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="enabled" label="状态" width="80">
            <template #default="{ row }">
              <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
                {{ row.enabled ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div class="glass-card panel">
        <h2 class="section-title">
          <el-icon><Cpu /></el-icon>
          智能匹配计算
        </h2>
        <div class="match-calculator">
          <div class="input-group">
            <label>灯光瓦数 (W)</label>
            <el-slider v-model="matchInput.light_wattage" :min="0" :max="100" show-input />
          </div>
          <div class="input-group">
            <label>水温 (°C)</label>
            <el-slider v-model="matchInput.temperature" :min="18" :max="32" :step="0.1" show-input />
          </div>
          <el-button type="primary" @click="calculateMatch" :loading="calculating" style="width: 100%;">
            计算气泵档位
          </el-button>
          <div v-if="matchResult" class="match-result" :class="{ 'power-saving-result': matchResult.power_saving }">
            <div class="result-level">
              <div style="display: flex; align-items: center; gap: 8px;">
                <span class="level-label">推荐档位</span>
                <el-tag v-if="matchResult.power_saving" type="warning" effect="dark" size="small">
                  <el-icon style="margin-right: 2px;"><Moon /></el-icon>
                  深夜省电
                </el-tag>
              </div>
              <div style="display: flex; align-items: center; gap: 10px;">
                <span v-if="matchResult.power_saving && matchResult.original_pump_level" 
                      class="original-level" title="原始档位">
                  <s>{{ matchResult.original_pump_level }}档</s>
                </span>
                <span class="level-value" :style="{ color: matchResult.power_saving ? '#fbbf24' : '#38bdf8' }">
                  {{ matchResult.pump_level }}档
                </span>
              </div>
            </div>
            <div class="result-detail">
              <p><strong>匹配模式:</strong> {{ matchResult.description }}</p>
              <p><strong>计算公式:</strong> {{ matchResult.formula }}</p>
              <p><strong>说明:</strong> {{ matchResult.reason }}</p>
            </div>
          </div>
        </div>
      </div>

      <div class="glass-card panel full-width">
        <h2 class="section-title">
          <el-icon><TrendCharts /></el-icon>
          环境数据趋势
        </h2>
        <div ref="chartRef" class="chart-container"></div>
      </div>

      <div class="glass-card panel full-width">
        <h2 class="section-title">
          <el-icon><Monitor /></el-icon>
          设备状态概览
        </h2>
        <div class="device-grid">
          <div v-for="device in allDevices" :key="device.id" class="device-card"
               :class="{ 'device-power-saving': deviceStatus.power_saving && device.status && !device.manual_mode }">
            <div class="device-header">
              <div class="device-info">
                <span class="status-dot" :class="{ on: device.status, off: !device.status }"></span>
                <span class="device-name">{{ device.device_name }}</span>
                <el-tag size="small" :type="device.device_type === 'light' ? 'warning' : 'primary'">
                  {{ device.device_type === 'light' ? '灯光' : '气泵' }}
                </el-tag>
                <el-icon v-if="deviceStatus.power_saving && device.status && !device.manual_mode" 
                         class="device-moon-icon" title="深夜省电模式">
                  <Moon />
                </el-icon>
              </div>
              <div style="display: flex; gap: 6px; align-items: center;">
                <el-tag size="small" :type="device.manual_mode ? 'danger' : 'success'">
                  {{ device.manual_mode ? '手动' : '自动' }}
                </el-tag>
              </div>
            </div>
            <div class="device-status">
              <span>
                状态: {{ device.status ? '开启' : '关闭' }}
                <span v-if="deviceStatus.power_saving && device.status && !device.manual_mode" class="saving-hint-inline">
                  (省电模式)
                </span>
              </span>
              <span v-if="device.current_value !== null">
                {{ device.device_type === 'light' ? '亮度' : '档位' }}: {{ device.current_value }}{{ device.device_type === 'light' ? '%' : '档' }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { lightingApi, oxygenApi, deviceApi, sensorApi } from '@/utils/api'

const chartRef = ref(null)
let chartInstance = null
let refreshTimer = null

const lightStatus = ref({ is_on: false, brightness: 0, night_mode: false, power_saving: false, original_brightness: 0 })
const schedules = ref([])
const sensorStats = ref({})
const deviceStatus = ref({})
const allDevices = ref([])

const calculating = ref(false)
const matchInput = reactive({ light_wattage: 50, temperature: 25 })
const matchResult = ref(null)

const calculateMatch = async () => {
  calculating.value = true
  try {
    const res = await oxygenApi.calculateMatch(matchInput)
    matchResult.value = res
  } catch (e) {
    ElMessage.error('计算失败')
  } finally {
    calculating.value = false
  }
}

const fetchData = async () => {
  try {
    const [lightRes, scheduleRes, sensorRes, deviceRes, devicesRes] = await Promise.all([
      lightingApi.getStatus(),
      lightingApi.getSchedules(),
      sensorApi.getStats(24),
      deviceApi.getDashboardStatus(),
      deviceApi.getDevices()
    ])
    lightStatus.value = lightRes || { is_on: false, brightness: 0 }
    schedules.value = scheduleRes || []
    sensorStats.value = sensorRes || {}
    deviceStatus.value = deviceRes || {}
    allDevices.value = devicesRes || []
    updateChart()
  } catch (e) {
    console.error('Fetch data failed:', e)
  }
}

const initChart = () => {
  if (!chartRef.value) return
  chartInstance = echarts.init(chartRef.value, 'dark')
  updateChart()
}

const updateChart = async () => {
  if (!chartInstance) return
  try {
    const data = await sensorApi.getData(50, 24)
    const times = data.map(d => d.recorded_at?.substring(11, 16) || '').reverse()
    const temps = data.map(d => d.temperature).reverse()
    const lights = data.map(d => d.light_wattage).reverse()

    chartInstance.setOption({
      backgroundColor: 'transparent',
      tooltip: {
        trigger: 'axis',
        backgroundColor: 'rgba(12, 20, 69, 0.9)',
        borderColor: 'rgba(255,255,255,0.1)'
      },
      legend: {
        data: ['水温', '灯光瓦数'],
        textStyle: { color: 'rgba(255,255,255,0.8)' }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: times,
        axisLine: { lineStyle: { color: 'rgba(255,255,255,0.2)' } },
        axisLabel: { color: 'rgba(255,255,255,0.6)' }
      },
      yAxis: [
        {
          type: 'value',
          name: '°C',
          axisLine: { lineStyle: { color: 'rgba(255,255,255,0.2)' } },
          axisLabel: { color: 'rgba(255,255,255,0.6)' },
          splitLine: { lineStyle: { color: 'rgba(255,255,255,0.05)' } }
        },
        {
          type: 'value',
          name: 'W',
          axisLine: { lineStyle: { color: 'rgba(255,255,255,0.2)' } },
          axisLabel: { color: 'rgba(255,255,255,0.6)' },
          splitLine: { lineStyle: { color: 'rgba(255,255,255,0.05)' } }
        }
      ],
      series: [
        {
          name: '水温',
          type: 'line',
          smooth: true,
          data: temps,
          yAxisIndex: 0,
          itemStyle: { color: '#38bdf8' },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(56, 189, 248, 0.3)' },
              { offset: 1, color: 'rgba(56, 189, 248, 0)' }
            ])
          }
        },
        {
          name: '灯光瓦数',
          type: 'line',
          smooth: true,
          data: lights,
          yAxisIndex: 1,
          itemStyle: { color: '#fbbf24' },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(251, 191, 36, 0.3)' },
              { offset: 1, color: 'rgba(251, 191, 36, 0)' }
            ])
          }
        }
      ]
    })
  } catch (e) {
    console.error('Chart update failed:', e)
  }
}

onMounted(async () => {
  await fetchData()
  await nextTick()
  initChart()
  refreshTimer = setInterval(fetchData, 10000)
  window.addEventListener('resize', () => chartInstance?.resize())
})

onUnmounted(() => {
  if (refreshTimer) clearInterval(refreshTimer)
  if (chartInstance) chartInstance.dispose()
  window.removeEventListener('resize', () => chartInstance?.resize())
})
</script>

<style scoped lang="scss">
.night-mode-bar {
  margin-bottom: 24px;
  overflow: hidden;
  position: relative;

  &::before {
    content: '';
    position: absolute;
    top: -50%;
    right: -10%;
    width: 200px;
    height: 200px;
    background: radial-gradient(circle, rgba(251, 191, 36, 0.15), transparent 70%);
    border-radius: 50%;
  }
}

.moon-glow {
  background: linear-gradient(135deg, rgba(30, 27, 75, 0.6), rgba(251, 191, 36, 0.1), rgba(168, 85, 247, 0.15)) !important;
  border: 1px solid rgba(251, 191, 36, 0.3) !important;
}

.night-mode-content {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 8px 4px;
  position: relative;
  z-index: 1;
}

.moon-icon-wrapper {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(251, 191, 36, 0.2), rgba(168, 85, 247, 0.2));
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  animation: moonPulse 2s ease-in-out infinite;
}

.moon-icon {
  font-size: 36px;
  color: #fbbf24;
  filter: drop-shadow(0 0 10px rgba(251, 191, 36, 0.5));
}

@keyframes moonPulse {
  0%, 100% {
    box-shadow: 0 0 20px rgba(251, 191, 36, 0.3);
    transform: scale(1);
  }
  50% {
    box-shadow: 0 0 40px rgba(251, 191, 36, 0.5);
    transform: scale(1.05);
  }
}

.night-info {
  flex: 1;
}

.night-title {
  font-size: 20px;
  font-weight: 700;
  color: #fbbf24;
  margin-bottom: 6px;
  letter-spacing: 0.5px;
}

.night-desc {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.75);
  line-height: 1.5;
}

.night-tag {
  font-weight: 700;
  letter-spacing: 1px;
  padding: 8px 16px;
  flex-shrink: 0;
}

.mini-moon-tag {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 2px;

  .el-icon {
    font-size: 12px;
  }
}

.stat-icon-wrap {
  display: flex;
  align-items: center;
}

.saving-hint {
  color: #fbbf24;
  margin-left: 6px;
  font-weight: 500;
  font-size: 12px;
}

.saving-hint-inline {
  color: #fbbf24;
  font-size: 12px;
}

.power-saving-active {
  border: 1px solid rgba(251, 191, 36, 0.35) !important;
  box-shadow: 0 0 20px rgba(251, 191, 36, 0.08), inset 0 0 30px rgba(251, 191, 36, 0.03);
  transition: all 0.3s ease;
}

.original-level {
  color: rgba(255, 255, 255, 0.4);
  font-size: 18px;
}

.power-saving-result {
  background: rgba(251, 191, 36, 0.08) !important;
  border-color: rgba(251, 191, 36, 0.3) !important;
}

.device-power-saving {
  border-color: rgba(251, 191, 36, 0.25) !important;
  background: linear-gradient(135deg, rgba(251, 191, 36, 0.05), rgba(255, 255, 255, 0.05)) !important;
}

.device-moon-icon {
  color: #fbbf24;
  font-size: 18px;
  margin-left: 4px;
  animation: twinkle 2s ease-in-out infinite;
  filter: drop-shadow(0 0 4px rgba(251, 191, 36, 0.5));
}

@keyframes twinkle {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.6; }
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.stat-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stat-icon {
  font-size: 24px;
}

.stat-sub {
  font-size: 13px;
  opacity: 0.7;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.panel {
  padding: 24px;
}

.panel.full-width {
  grid-column: 1 / -1;
}

.schedule-info {
  background: rgba(255, 255, 255, 0.05);
  padding: 16px;
  border-radius: 12px;
  margin-top: 8px;
}

.active-schedule {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #4ade80;
  font-weight: 500;
  margin-bottom: 8px;
}

.next-action {
  display: flex;
  justify-content: space-between;
  opacity: 0.8;
}

.next-time {
  font-family: 'Courier New', monospace;
  color: #38bdf8;
}

.data-table {
  background: transparent !important;
}

.data-table :deep(.el-table) {
  background: transparent;
}

.data-table :deep(.el-table__body tr:hover > td) {
  background-color: rgba(255, 255, 255, 0.05);
}

.data-table :deep(.el-table th),
.data-table :deep(.el-table td) {
  background: transparent;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.9);
}

.input-group {
  margin-bottom: 20px;
}

.input-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  opacity: 0.9;
}

.match-result {
  margin-top: 20px;
  padding: 20px;
  background: rgba(56, 189, 248, 0.1);
  border: 1px solid rgba(56, 189, 248, 0.3);
  border-radius: 12px;
}

.result-level {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.level-label {
  font-size: 14px;
  opacity: 0.8;
}

.level-value {
  font-size: 36px;
  font-weight: 700;
  color: #38bdf8;
}

.result-detail p {
  margin: 8px 0;
  font-size: 13px;
  opacity: 0.85;
  line-height: 1.6;
}

.chart-container {
  height: 350px;
  margin-top: 16px;
}

.device-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 16px;
  margin-top: 16px;
}

.device-card {
  background: rgba(255, 255, 255, 0.05);
  padding: 16px;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.device-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.device-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.device-name {
  font-weight: 500;
}

.device-status {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
  opacity: 0.8;
}
</style>
