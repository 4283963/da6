<template>
  <div class="page-container">
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

      <div class="stat-card glass-card" style="background: linear-gradient(135deg, rgba(34, 197, 94, 0.2), rgba(16, 185, 129, 0.1);">
        <div class="stat-header">
          <el-icon class="stat-icon" style="color: #22c55e;"><Sunny /></el-icon>
          <span class="stat-label">灯光状态</span>
        </div>
        <div class="stat-value">
          <span v-if="lightStatus.is_on" style="color: #4ade80;">开启中</span>
          <span v-else style="color: #6b7280;">已关闭</span>
        </div>
        <div class="stat-sub">亮度: {{ lightStatus.brightness || 0 }}%</div>
      </div>

      <div class="stat-card glass-card" style="background: linear-gradient(135deg, rgba(139, 92, 246, 0.2), rgba(124, 58, 237, 0.1);">
        <div class="stat-header">
          <el-icon class="stat-icon" style="color: #8b5cf6;"><Wind /></el-icon>
          <span class="stat-label">气泵档位</span>
        </div>
        <div class="stat-value">
          <span v-if="deviceStatus.current_pump_level > 0">{{ deviceStatus.current_pump_level }}档</span>
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
          <div v-if="matchResult" class="match-result">
            <div class="result-level">
              <span class="level-label">推荐档位</span>
              <span class="level-value">{{ matchResult.pump_level }}档</span>
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
          <div v-for="device in allDevices" :key="device.id" class="device-card">
            <div class="device-header">
              <div class="device-info">
                <span class="status-dot" :class="{ on: device.status, off: !device.status }"></span>
                <span class="device-name">{{ device.device_name }}</span>
                <el-tag size="small" :type="device.device_type === 'light' ? 'warning' : 'primary'">
                  {{ device.device_type === 'light' ? '灯光' : '气泵' }}
                </el-tag>
              </div>
              <el-tag size="small" :type="device.manual_mode ? 'danger' : 'success'">
                {{ device.manual_mode ? '手动' : '自动' }}
              </el-tag>
            </div>
            <div class="device-status">
              <span>状态: {{ device.status ? '开启' : '关闭' }}</span>
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

const lightStatus = ref({ is_on: false, brightness: 0 })
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
