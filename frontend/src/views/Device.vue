<template>
  <div class="page-container">
    <div class="tabs-wrapper">
      <el-tabs v-model="activeTab" class="custom-tabs">
        <el-tab-pane label="灯光控制" name="light">
          <div class="devices-grid">
            <div
              v-for="device in lightDevices"
              :key="device.id"
              class="device-control-card glass-card"
            >
              <div class="card-header">
                <div class="device-info">
                  <el-icon class="device-icon" style="color: #fbbf24;"><Lightbulb /></el-icon>
                  <div>
                    <h3 class="device-name">{{ device.device_name }}</h3>
                    <span class="device-type">灯光设备</span>
                  </div>
                </div>
                <el-tag :type="device.manual_mode ? 'danger' : 'success'" size="small">
                  {{ device.manual_mode ? '手动模式' : '自动模式' }}
                </el-tag>
              </div>

              <div class="card-body">
                <div class="status-row">
                  <div class="status-item">
                    <span class="status-label">电源状态</span>
                    <div class="status-value">
                      <span class="status-dot" :class="{ on: device.status, off: !device.status }"></span>
                      <span :class="{ 'text-green': device.status, 'text-gray': !device.status }">
                        {{ device.status ? '已开启' : '已关闭' }}
                      </span>
                    </div>
                  </div>
                  <div class="status-item">
                    <span class="status-label">当前亮度</span>
                    <span class="status-value highlight">{{ device.current_value || 0 }}%</span>
                  </div>
                </div>

                <div class="control-section">
                  <label class="control-label">亮度调节</label>
                  <div class="slider-row">
                    <el-slider
                      v-model="device.temp_brightness"
                      :min="0"
                      :max="100"
                      :disabled="!device.manual_mode"
                      @change="(val) => updateBrightness(device, val)"
                    />
                    <span class="slider-value">{{ device.temp_brightness || 0 }}%</span>
                  </div>
                </div>

                <div class="button-group">
                  <el-button
                    :type="device.status ? 'danger' : 'success'"
                    :disabled="!device.manual_mode"
                    @click="toggleDevice(device)"
                    style="flex: 1;"
                  >
                    <el-icon><component :is="device.status ? 'SwitchOff' : 'SwitchOn'" /></el-icon>
                    {{ device.status ? '关闭' : '开启' }}
                  </el-button>
                  <el-button
                    :type="device.manual_mode ? 'warning' : 'primary'"
                    @click="toggleManualMode(device)"
                    style="flex: 1;"
                  >
                    <el-icon><Setting /></el-icon>
                    {{ device.manual_mode ? '切自动' : '切手动' }}
                  </el-button>
                </div>
              </div>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane label="气泵控制" name="pump">
          <div class="devices-grid">
            <div
              v-for="device in pumpDevices"
              :key="device.id"
              class="device-control-card glass-card"
            >
              <div class="card-header">
                <div class="device-info">
                  <el-icon class="device-icon" style="color: #8b5cf6;"><Wind /></el-icon>
                  <div>
                    <h3 class="device-name">{{ device.device_name }}</h3>
                    <span class="device-type">气泵设备</span>
                  </div>
                </div>
                <el-tag :type="device.manual_mode ? 'danger' : 'success'" size="small">
                  {{ device.manual_mode ? '手动模式' : '自动模式' }}
                </el-tag>
              </div>

              <div class="card-body">
                <div class="status-row">
                  <div class="status-item">
                    <span class="status-label">电源状态</span>
                    <div class="status-value">
                      <span class="status-dot" :class="{ on: device.status, off: !device.status }"></span>
                      <span :class="{ 'text-green': device.status, 'text-gray': !device.status }">
                        {{ device.status ? '已开启' : '已关闭' }}
                      </span>
                    </div>
                  </div>
                  <div class="status-item">
                    <span class="status-label">当前档位</span>
                    <span class="status-value highlight">{{ device.current_value || 0 }}档</span>
                  </div>
                </div>

                <div class="control-section">
                  <label class="control-label">档位调节</label>
                  <div class="level-buttons">
                    <el-button
                      v-for="level in 5"
                      :key="level"
                      :type="device.current_value === level ? 'primary' : 'default'"
                      :disabled="!device.manual_mode"
                      @click="updatePumpLevel(device, level)"
                      class="level-btn"
                    >
                      {{ level }}档
                    </el-button>
                  </div>
                </div>

                <div class="button-group">
                  <el-button
                    :type="device.status ? 'danger' : 'success'"
                    :disabled="!device.manual_mode"
                    @click="toggleDevice(device)"
                    style="flex: 1;"
                  >
                    <el-icon><component :is="device.status ? 'SwitchOff' : 'SwitchOn'" /></el-icon>
                    {{ device.status ? '关闭' : '开启' }}
                  </el-button>
                  <el-button
                    :type="device.manual_mode ? 'warning' : 'primary'"
                    @click="toggleManualMode(device)"
                    style="flex: 1;"
                  >
                    <el-icon><Setting /></el-icon>
                    {{ device.manual_mode ? '切自动' : '切手动' }}
                  </el-button>
                </div>
              </div>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <div class="glass-card quick-action-panel">
      <h2 class="section-title">
        <el-icon><Lightning /></el-icon>
        快捷操作
      </h2>
      <div class="quick-buttons">
        <el-button type="warning" size="large" @click="turnAllLightsOn">
          <el-icon><Sunny /></el-icon>
          全开灯光
        </el-button>
        <el-button type="info" size="large" @click="turnAllLightsOff">
          <el-icon><Moon /></el-icon>
          全关灯光
        </el-button>
        <el-button type="success" size="large" @click="turnAllPumpsOn">
          <el-icon><Wind /></el-icon>
          全开气泵
        </el-button>
        <el-button type="danger" size="large" @click="turnAllPumpsOff">
          <el-icon><Close /></el-icon>
          全关气泵
        </el-button>
        <el-button type="primary" size="large" @click="setAllAuto">
          <el-icon><Cpu /></el-icon>
          全部自动
        </el-button>
      </div>
    </div>

    <div class="glass-card status-panel">
      <h2 class="section-title">
        <el-icon><InfoFilled /></el-icon>
        模式说明
      </h2>
      <div class="info-content">
        <div class="info-item">
          <el-tag type="success">自动模式</el-tag>
          <p>系统根据光照排程和溶氧量匹配模型自动控制设备，适合日常运行。</p>
        </div>
        <div class="info-item">
          <el-tag type="danger">手动模式</el-tag>
          <p>用户手动控制设备开关和参数，适合调试或特殊情况使用。切换到自动模式后，设备会恢复系统控制。</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { deviceApi } from '@/utils/api'

const activeTab = ref('light')
const devices = ref([])
const loading = ref(false)

const lightDevices = computed(() => {
  return devices.value
    .filter(d => d.device_type === 'light')
    .map(d => ({ ...d, temp_brightness: d.current_value || 0 }))
})

const pumpDevices = computed(() => {
  return devices.value.filter(d => d.device_type === 'pump')
})

const fetchData = async () => {
  loading.value = true
  try {
    const data = await deviceApi.getDevices()
    devices.value = data || []
  } catch (e) {
    ElMessage.error('获取设备状态失败')
  } finally {
    loading.value = false
  }
}

const toggleDevice = async (device) => {
  if (!device.manual_mode) {
    ElMessage.warning('请先切换到手动模式')
    return
  }
  const newStatus = !device.status
  try {
    await deviceApi.toggleDevice(device.device_type, device.device_name, newStatus)
    device.status = newStatus
    if (!newStatus) {
      device.current_value = 0
      device.temp_brightness = 0
    }
    ElMessage.success(newStatus ? '已开启' : '已关闭')
  } catch (e) {
    ElMessage.error('操作失败')
  }
}

const updateBrightness = async (device, value) => {
  if (!device.manual_mode) {
    ElMessage.warning('请先切换到手动模式')
    return
  }
  try {
    await deviceApi.updateValue(device.device_type, device.device_name, value)
    device.current_value = value
    device.status = value > 0
    ElMessage.success(`亮度已调整为 ${value}%`)
  } catch (e) {
    ElMessage.error('调整亮度失败')
  }
}

const updatePumpLevel = async (device, level) => {
  if (!device.manual_mode) {
    ElMessage.warning('请先切换到手动模式')
    return
  }
  try {
    await deviceApi.updateValue(device.device_type, device.device_name, level)
    device.current_value = level
    device.status = true
    ElMessage.success(`档位已调整为 ${level}档`)
  } catch (e) {
    ElMessage.error('调整档位失败')
  }
}

const toggleManualMode = async (device) => {
  const newMode = !device.manual_mode
  const action = newMode ? '手动模式' : '自动模式'
  try {
    await ElMessageBox.confirm(
      `确定要切换到${action}吗？`,
      '切换模式',
      { type: 'warning' }
    )
    await deviceApi.setManualMode(device.device_type, device.device_name, newMode)
    device.manual_mode = newMode
    if (!newMode) {
      device.status = false
      device.current_value = 0
      device.temp_brightness = 0
    }
    ElMessage.success(`已切换到${action}`)
    fetchData()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('切换失败')
    }
  }
}

const turnAllLightsOn = async () => {
  try {
    await ElMessageBox.confirm('确定要开启所有灯光吗？这将切换所有灯光到手动模式。', '确认操作', { type: 'warning' })
    for (const device of lightDevices.value) {
      if (!device.manual_mode) {
        await deviceApi.setManualMode('light', device.device_name, true)
      }
      await deviceApi.updateValue('light', device.device_name, 80)
    }
    ElMessage.success('所有灯光已开启')
    fetchData()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const turnAllLightsOff = async () => {
  try {
    await ElMessageBox.confirm('确定要关闭所有灯光吗？', '确认操作', { type: 'warning' })
    for (const device of lightDevices.value) {
      if (!device.manual_mode) {
        await deviceApi.setManualMode('light', device.device_name, true)
      }
      await deviceApi.toggleDevice('light', device.device_name, false)
    }
    ElMessage.success('所有灯光已关闭')
    fetchData()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const turnAllPumpsOn = async () => {
  try {
    await ElMessageBox.confirm('确定要开启所有气泵吗？这将切换所有气泵到手动模式。', '确认操作', { type: 'warning' })
    for (const device of pumpDevices.value) {
      if (!device.manual_mode) {
        await deviceApi.setManualMode('pump', device.device_name, true)
      }
      await deviceApi.updateValue('pump', device.device_name, 3)
    }
    ElMessage.success('所有气泵已开启')
    fetchData()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const turnAllPumpsOff = async () => {
  try {
    await ElMessageBox.confirm('确定要关闭所有气泵吗？', '确认操作', { type: 'warning' })
    for (const device of pumpDevices.value) {
      if (!device.manual_mode) {
        await deviceApi.setManualMode('pump', device.device_name, true)
      }
      await deviceApi.toggleDevice('pump', device.device_name, false)
    }
    ElMessage.success('所有气泵已关闭')
    fetchData()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const setAllAuto = async () => {
  try {
    await ElMessageBox.confirm('确定要将所有设备切换到自动模式吗？', '确认操作', { type: 'warning' })
    for (const device of devices.value) {
      if (device.manual_mode) {
        await deviceApi.setManualMode(device.device_type, device.device_name, false)
      }
    }
    ElMessage.success('所有设备已切换到自动模式')
    fetchData()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

onMounted(fetchData)
</script>

<style scoped lang="scss">
.tabs-wrapper {
  margin-bottom: 20px;
}

:deep(.el-tabs__item) {
  color: rgba(255, 255, 255, 0.6);
  font-size: 16px;
  height: 48px;
  line-height: 48px;
}

:deep(.el-tabs__item.is-active) {
  color: #38bdf8;
}

:deep(.el-tabs__active-bar) {
  background-color: #38bdf8;
}

:deep(.el-tabs__nav-wrap::after) {
  background-color: rgba(255, 255, 255, 0.1);
}

.devices-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 20px;
  margin-top: 8px;
}

.device-control-card {
  padding: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.device-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.device-icon {
  font-size: 36px;
}

.device-name {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
}

.device-type {
  font-size: 12px;
  opacity: 0.6;
}

.status-row {
  display: flex;
  justify-content: space-between;
  background: rgba(255, 255, 255, 0.05);
  padding: 16px;
  border-radius: 12px;
  margin-bottom: 20px;
}

.status-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.status-label {
  font-size: 13px;
  opacity: 0.7;
}

.status-value {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 500;
}

.status-value.highlight {
  font-size: 24px;
  font-weight: 700;
  color: #38bdf8;
}

.text-green {
  color: #4ade80;
}

.text-gray {
  color: #6b7280;
}

.control-section {
  margin-bottom: 20px;
}

.control-label {
  display: block;
  font-size: 14px;
  margin-bottom: 12px;
  opacity: 0.9;
}

.slider-row {
  display: flex;
  align-items: center;
  gap: 16px;
}

.slider-row :deep(.el-slider) {
  flex: 1;
}

.slider-value {
  min-width: 50px;
  text-align: right;
  font-weight: 600;
  color: #fbbf24;
}

.level-buttons {
  display: flex;
  gap: 8px;
}

.level-btn {
  flex: 1;
}

.button-group {
  display: flex;
  gap: 12px;
}

.quick-action-panel,
.status-panel {
  padding: 24px;
  margin-bottom: 20px;
}

.quick-buttons {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 12px;
}

.info-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 12px;
  background: rgba(255, 255, 255, 0.03);
  padding: 20px;
  border-radius: 12px;
}

.info-item p {
  margin: 0;
  font-size: 14px;
  opacity: 0.8;
  line-height: 1.6;
}
</style>
