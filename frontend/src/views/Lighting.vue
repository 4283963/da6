<template>
  <div class="page-container">
    <div class="content-grid">
      <div class="glass-card panel main-panel">
        <div class="panel-header">
          <h2 class="section-title">
            <el-icon><Sunny /></el-icon>
            光照排程管理
          </h2>
          <el-button type="primary" @click="openAddDialog">
            <el-icon><Plus /></el-icon>
            新增排程
          </el-button>
        </div>

        <el-table :data="schedules" class="data-table" v-loading="loading">
          <el-table-column prop="name" label="排程名称" min-width="120" />
          <el-table-column prop="start_time" label="开灯时间" width="120">
            <template #default="{ row }">
              <el-icon style="color: #4ade80; margin-right: 4px;"><Sunrise /></el-icon>
              {{ row.start_time }}
            </template>
          </el-table-column>
          <el-table-column prop="end_time" label="关灯时间" width="120">
            <template #default="{ row }">
              <el-icon style="color: #f97316; margin-right: 4px;"><Sunset /></el-icon>
              {{ row.end_time }}
            </template>
          </el-table-column>
          <el-table-column prop="brightness" label="亮度" width="150">
            <template #default="{ row }">
              <div class="brightness-display">
                <el-slider
                  v-model="row.brightness"
                  :min="0"
                  :max="100"
                  :disabled="true"
                  style="width: 100px; margin-right: 12px;"
                />
                <span class="brightness-value">{{ row.brightness }}%</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="enabled" label="状态" width="100">
            <template #default="{ row }">
              <el-switch
                v-model="row.enabled"
                @change="toggleSchedule(row)"
                active-text="启用"
                inactive-text="禁用"
              />
            </template>
          </el-table-column>
          <el-table-column label="操作" width="140" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" size="small" link @click="openEditDialog(row)">
                编辑
              </el-button>
              <el-button type="danger" size="small" link @click="deleteSchedule(row)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div class="glass-card panel side-panel">
        <h2 class="section-title">
          <el-icon><Clock /></el-icon>
          当前状态
        </h2>
        <div class="status-display">
          <div class="status-item">
            <div class="status-label">灯光状态</div>
            <div class="status-value">
              <span class="status-dot" :class="{ on: lightStatus.is_on, off: !lightStatus.is_on }"></span>
              <span :class="{ 'text-green': lightStatus.is_on, 'text-gray': !lightStatus.is_on }">
                {{ lightStatus.is_on ? '已开启' : '已关闭' }}
              </span>
            </div>
          </div>
          <div class="status-item">
            <div class="status-label">当前亮度</div>
            <div class="status-value highlight">{{ lightStatus.brightness }}%</div>
          </div>
          <div class="status-item" v-if="lightStatus.schedule_name">
            <div class="status-label">执行排程</div>
            <div class="status-value highlight">{{ lightStatus.schedule_name }}</div>
          </div>
          <div class="status-item">
            <div class="status-label">下一步</div>
            <div class="status-value">
              <span class="next-action">{{ lightStatus.next_action || '--' }}</span>
              <span class="next-time">{{ lightStatus.next_time || '--' }}</span>
            </div>
          </div>
        </div>

        <div class="timeline-section">
          <h3 class="section-title" style="font-size: 16px; margin-top: 24px;">
            <el-icon><Calendar /></el-icon>
            今日排程时间线
          </h3>
          <div class="timeline">
            <div
              v-for="(schedule, index) in sortedSchedules"
              :key="schedule.id"
              class="timeline-item"
              :class="{ active: isScheduleActive(schedule) }"
            >
              <div class="timeline-dot"></div>
              <div class="timeline-content">
                <div class="timeline-time">
                  {{ schedule.start_time }} - {{ schedule.end_time }}
                </div>
                <div class="timeline-name">{{ schedule.name }}</div>
                <div class="timeline-brightness">
                  <el-tag size="small" :type="schedule.brightness > 60 ? 'warning' : 'success'">
                    亮度 {{ schedule.brightness }}%
                  </el-tag>
                  <el-tag v-if="!schedule.enabled" size="small" type="info">已禁用</el-tag>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑排程' : '新增排程'"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="排程名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入排程名称" />
        </el-form-item>
        <el-form-item label="开灯时间" prop="start_time">
          <el-time-picker
            v-model="formData.start_time"
            format="HH:mm:ss"
            value-format="HH:mm:ss"
            placeholder="选择开灯时间"
            style="width: 100%;"
          />
        </el-form-item>
        <el-form-item label="关灯时间" prop="end_time">
          <el-time-picker
            v-model="formData.end_time"
            format="HH:mm:ss"
            value-format="HH:mm:ss"
            placeholder="选择关灯时间"
            style="width: 100%;"
          />
        </el-form-item>
        <el-form-item label="亮度" prop="brightness">
          <el-slider v-model="formData.brightness" :min="0" :max="100" show-input />
        </el-form-item>
        <el-form-item label="启用状态" prop="enabled">
          <el-switch v-model="formData.enabled" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">
          {{ isEdit ? '保存修改' : '创建排程' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { lightingApi } from '@/utils/api'
import dayjs from 'dayjs'

const loading = ref(false)
const submitting = ref(false)
const schedules = ref([])
const lightStatus = ref({ is_on: false, brightness: 0 })

const dialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref(null)
const formRef = ref(null)
const formData = ref({
  name: '',
  start_time: '',
  end_time: '',
  brightness: 50,
  enabled: true
})

const formRules = {
  name: [{ required: true, message: '请输入排程名称', trigger: 'blur' }],
  start_time: [{ required: true, message: '请选择开灯时间', trigger: 'change' }],
  end_time: [{ required: true, message: '请选择关灯时间', trigger: 'change' }]
}

const sortedSchedules = computed(() => {
  return [...schedules.value].sort((a, b) => a.start_time.localeCompare(b.start_time))
})

const isScheduleActive = (schedule) => {
  if (!schedule.enabled) return false
  const now = dayjs().format('HH:mm:ss')
  return schedule.start_time <= now && now < schedule.end_time
}

const fetchData = async () => {
  loading.value = true
  try {
    const [schedulesRes, statusRes] = await Promise.all([
      lightingApi.getSchedules(),
      lightingApi.getStatus()
    ])
    schedules.value = schedulesRes || []
    lightStatus.value = statusRes || { is_on: false, brightness: 0 }
  } catch (e) {
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const openAddDialog = () => {
  isEdit.value = false
  editingId.value = null
  formData.value = {
    name: '',
    start_time: '',
    end_time: '',
    brightness: 50,
    enabled: true
  }
  dialogVisible.value = true
}

const openEditDialog = (row) => {
  isEdit.value = true
  editingId.value = row.id
  formData.value = { ...row }
  dialogVisible.value = true
}

const submitForm = async () => {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value) {
      await lightingApi.updateSchedule(editingId.value, formData.value)
      ElMessage.success('修改成功')
    } else {
      await lightingApi.createSchedule(formData.value)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (e) {
    ElMessage.error(isEdit.value ? '修改失败' : '创建失败')
  } finally {
    submitting.value = false
  }
}

const toggleSchedule = async (row) => {
  try {
    await lightingApi.updateSchedule(row.id, { enabled: row.enabled })
    ElMessage.success(row.enabled ? '已启用' : '已禁用')
    fetchData()
  } catch (e) {
    row.enabled = !row.enabled
    ElMessage.error('操作失败')
  }
}

const deleteSchedule = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除排程「${row.name}」吗？`,
      '确认删除',
      { type: 'warning' }
    )
    await lightingApi.deleteSchedule(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

onMounted(fetchData)
</script>

<style scoped lang="scss">
.content-grid {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 20px;
}

.panel {
  padding: 24px;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.brightness-display {
  display: flex;
  align-items: center;
}

.brightness-value {
  font-weight: 600;
  color: #fbbf24;
  min-width: 45px;
}

.status-display {
  background: rgba(255, 255, 255, 0.05);
  padding: 16px;
  border-radius: 12px;
}

.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.status-item:last-child {
  border-bottom: none;
}

.status-label {
  font-size: 14px;
  opacity: 0.7;
}

.status-value {
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-value.highlight {
  color: #38bdf8;
  font-size: 18px;
  font-weight: 600;
}

.text-green {
  color: #4ade80;
}

.text-gray {
  color: #6b7280;
}

.next-action {
  font-size: 13px;
  opacity: 0.8;
}

.next-time {
  font-family: 'Courier New', monospace;
  color: #38bdf8;
  margin-left: 8px;
}

.timeline {
  position: relative;
  padding-left: 24px;
}

.timeline::before {
  content: '';
  position: absolute;
  left: 4px;
  top: 8px;
  bottom: 8px;
  width: 2px;
  background: rgba(255, 255, 255, 0.1);
}

.timeline-item {
  position: relative;
  padding-bottom: 20px;
  opacity: 0.6;
  transition: opacity 0.3s ease;
}

.timeline-item.active {
  opacity: 1;
}

.timeline-dot {
  position: absolute;
  left: -20px;
  top: 6px;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #6b7280;
  border: 2px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;
}

.timeline-item.active .timeline-dot {
  background: #4ade80;
  box-shadow: 0 0 12px #4ade80;
}

.timeline-time {
  font-size: 14px;
  font-weight: 600;
  color: #38bdf8;
  margin-bottom: 4px;
}

.timeline-name {
  font-size: 14px;
  margin-bottom: 6px;
}

.timeline-brightness {
  display: flex;
  gap: 8px;
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

.data-table :deep(.el-table__empty-block) {
  background: transparent;
}

:deep(.el-dialog) {
  background: rgba(12, 20, 69, 0.95);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
}

:deep(.el-dialog__header) {
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

:deep(.el-dialog__title) {
  color: #fff;
}

:deep(.el-form-item__label) {
  color: rgba(255, 255, 255, 0.8);
}

:deep(.el-input__wrapper),
:deep(.el-slider__runway),
:deep(.el-picker-panel) {
  background: rgba(255, 255, 255, 0.05);
  box-shadow: none;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

:deep(.el-input__inner) {
  color: #fff;
}

:deep(.el-picker-panel__content),
:deep(.el-time-spinner__list) {
  background: rgba(12, 20, 69, 0.95);
}

:deep(.el-time-spinner__item:hover:not(.disabled)),
:deep(.el-picker-panel__icon-btn:hover) {
  background: rgba(255, 255, 255, 0.1);
}

:deep(.el-time-spinner__item.active:not(.disabled)) {
  color: #38bdf8;
}
</style>
