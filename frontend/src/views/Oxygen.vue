<template>
  <div class="page-container">
    <div class="content-grid">
      <div class="glass-card panel main-panel">
        <div class="panel-header">
          <h2 class="section-title">
            <el-icon><Wind /></el-icon>
            溶氧量匹配配置
          </h2>
          <el-button type="primary" @click="openAddDialog">
            <el-icon><Plus /></el-icon>
            新增配置
          </el-button>
        </div>

        <el-table :data="configs" class="data-table" v-loading="loading">
          <el-table-column prop="id" label="ID" width="60" />
          <el-table-column label="灯光瓦数范围" width="140">
            <template #default="{ row }">
              <span style="color: #fbbf24;">{{ row.min_light_wattage }} - {{ row.max_light_wattage }}W</span>
            </template>
          </el-table-column>
          <el-table-column label="水温范围" width="140">
            <template #default="{ row }">
              <span style="color: #38bdf8;">{{ row.min_temperature }} - {{ row.max_temperature }}°C</span>
            </template>
          </el-table-column>
          <el-table-column prop="pump_level" label="气泵档位" width="120">
            <template #default="{ row }">
              <el-tag :type="getPumpLevelType(row.pump_level)" size="large">
                {{ row.pump_level }}档
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="description" label="说明" min-width="120" />
          <el-table-column label="操作" width="140" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" size="small" link @click="openEditDialog(row)">
                编辑
              </el-button>
              <el-button type="danger" size="small" link @click="deleteConfig(row)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div class="glass-card panel side-panel">
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
            <el-icon><MagicStick /></el-icon>
            计算匹配结果
          </el-button>

          <div v-if="matchResult" class="match-result">
            <div class="result-level">
              <span class="level-label">推荐档位</span>
              <span class="level-value">{{ matchResult.pump_level }}档</span>
            </div>
            <div class="result-detail">
              <p><strong>模式:</strong> {{ matchResult.description }}</p>
              <p><strong>公式:</strong> {{ matchResult.formula }}</p>
              <p><strong>说明:</strong> {{ matchResult.reason }}</p>
            </div>
          </div>
        </div>

        <h2 class="section-title" style="margin-top: 24px;">
          <el-icon><Grid /></el-icon>
          匹配矩阵
        </h2>
        <div class="match-matrix">
          <div class="matrix-header">
            <div class="matrix-corner"></div>
            <div v-for="temp in tempRanges" :key="temp" class="matrix-header-cell">
              {{ temp }}
            </div>
          </div>
          <div v-for="light in lightRanges" :key="light" class="matrix-row">
            <div class="matrix-row-header">{{ light }}</div>
            <div v-for="temp in tempRanges" :key="temp" class="matrix-cell" :class="getMatrixClass(light, temp)">
              {{ getMatrixValue(light, temp) }}
            </div>
          </div>
        </div>
        <div class="matrix-legend">
          <div><span class="legend-dot" style="background: #22c55e;"></span> 1-2档 (低)</div>
          <div><span class="legend-dot" style="background: #eab308;"></span> 3档 (中)</div>
          <div><span class="legend-dot" style="background: #f97316;"></span> 4-5档 (高)</div>
        </div>
      </div>
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑配置' : '新增配置'"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="最小瓦数" prop="min_light_wattage">
          <el-input-number v-model="formData.min_light_wattage" :min="0" :max="100" />
        </el-form-item>
        <el-form-item label="最大瓦数" prop="max_light_wattage">
          <el-input-number v-model="formData.max_light_wattage" :min="0" :max="100" />
        </el-form-item>
        <el-form-item label="最低水温" prop="min_temperature">
          <el-input-number v-model="formData.min_temperature" :min="18" :max="32" :step="0.1" :precision="1" />
        </el-form-item>
        <el-form-item label="最高水温" prop="max_temperature">
          <el-input-number v-model="formData.max_temperature" :min="18" :max="32" :step="0.1" :precision="1" />
        </el-form-item>
        <el-form-item label="气泵档位" prop="pump_level">
          <el-select v-model="formData.pump_level" style="width: 100%;">
            <el-option v-for="i in 5" :key="i" :label="`${i}档`" :value="i" />
          </el-select>
        </el-form-item>
        <el-form-item label="说明" prop="description">
          <el-input v-model="formData.description" placeholder="请输入说明文字" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">
          {{ isEdit ? '保存修改' : '创建配置' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { oxygenApi } from '@/utils/api'

const loading = ref(false)
const calculating = ref(false)
const submitting = ref(false)
const configs = ref([])
const matrixData = ref({})
const lightRanges = ref([])
const tempRanges = ref([])

const dialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref(null)
const formRef = ref(null)

const formData = reactive({
  min_light_wattage: 0,
  max_light_wattage: 30,
  min_temperature: 18,
  max_temperature: 24,
  pump_level: 1,
  description: ''
})

const formRules = {
  min_light_wattage: [{ required: true, message: '请输入最小瓦数', trigger: 'blur' }],
  max_light_wattage: [{ required: true, message: '请输入最大瓦数', trigger: 'blur' }],
  min_temperature: [{ required: true, message: '请输入最低水温', trigger: 'blur' }],
  max_temperature: [{ required: true, message: '请输入最高水温', trigger: 'blur' }],
  pump_level: [{ required: true, message: '请选择气泵档位', trigger: 'change' }]
}

const matchInput = reactive({ light_wattage: 50, temperature: 25 })
const matchResult = ref(null)

const getPumpLevelType = (level) => {
  if (level <= 2) return 'success'
  if (level <= 3) return 'warning'
  return 'danger'
}

const getMatrixClass = (light, temp) => {
  const val = matrixData.value[light]?.[temp] || 0
  if (val <= 2) return 'level-low'
  if (val <= 3) return 'level-medium'
  return 'level-high'
}

const getMatrixValue = (light, temp) => {
  const val = matrixData.value[light]?.[temp] || 0
  return val ? `${val}档` : '-'
}

const fetchData = async () => {
  loading.value = true
  try {
    const [configsRes, matrixRes] = await Promise.all([
      oxygenApi.getConfigs(),
      oxygenApi.getMatrix()
    ])
    configs.value = configsRes || []
    matrixData.value = matrixRes.matrix || {}
    lightRanges.value = matrixRes.light_ranges || []
    tempRanges.value = matrixRes.temp_ranges || []
  } catch (e) {
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

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

const openAddDialog = () => {
  isEdit.value = false
  editingId.value = null
  Object.assign(formData, {
    min_light_wattage: 0,
    max_light_wattage: 30,
    min_temperature: 18,
    max_temperature: 24,
    pump_level: 1,
    description: ''
  })
  dialogVisible.value = true
}

const openEditDialog = (row) => {
  isEdit.value = true
  editingId.value = row.id
  Object.assign(formData, { ...row })
  dialogVisible.value = true
}

const submitForm = async () => {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value) {
      await oxygenApi.updateConfig(editingId.value, formData)
      ElMessage.success('修改成功')
    } else {
      await oxygenApi.createConfig(formData)
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

const deleteConfig = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除配置「${row.description}吗？`,
      '确认删除',
      { type: 'warning' }
    )
    await oxygenApi.deleteConfig(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  fetchData()
  calculateMatch()
})
</script>

<style scoped lang="scss">
.content-grid {
  display: grid;
  grid-template-columns: 1fr 420px;
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

.match-matrix {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  overflow: hidden;
  margin-top: 16px;
}

.matrix-header,
.matrix-row {
  display: flex;
}

.matrix-corner {
  width: 80px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.05);
}

.matrix-header-cell {
  flex: 1;
  padding: 12px 8px;
  text-align: center;
  font-size: 12px;
  font-weight: 500;
  background: rgba(255, 255, 255, 0.05);
  border-left: 1px solid rgba(255, 255, 255, 0.1);
}

.matrix-row-header {
  width: 80px;
  padding: 16px 12px;
  font-size: 12px;
  font-weight: 500;
  background: rgba(255, 255, 255, 0.03);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.matrix-cell {
  flex: 1;
  padding: 16px 8px;
  text-align: center;
  font-weight: 600;
  font-size: 14px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  border-left: 1px solid rgba(255, 255, 255, 0.05);
  transition: all 0.3s ease;
}

.matrix-cell:hover {
  transform: scale(1.05);
}

.matrix-cell.level-low {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
}

.matrix-cell.level-medium {
  background: rgba(234, 179, 8, 0.2);
  color: #eab308;
}

.matrix-cell.level-high {
  background: rgba(249, 115, 22, 0.2);
  color: #f97316;
}

.matrix-legend {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-top: 16px;
  font-size: 12px;
  opacity: 0.8;
}

.matrix-legend > div {
  display: flex;
  align-items: center;
  gap: 6px;
}

.legend-dot {
  width: 12px;
  height: 12px;
  border-radius: 3px;
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

:deep(.el-dialog__title),
:deep(.el-form-item__label) {
  color: #fff;
}

:deep(.el-input__wrapper),
:deep(.el-input-number),
:deep(.el-select__wrapper) {
  background: rgba(255, 255, 255, 0.05);
  box-shadow: none;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

:deep(.el-input__inner),
:deep(.el-input-number__decrease),
:deep(.el-input-number__increase) {
  color: #fff;
}

:deep(.el-select-dropdown) {
  background: rgba(12, 20, 69, 0.95);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

:deep(.el-select-dropdown__item) {
  color: rgba(255, 255, 255, 0.9);
}

:deep(.el-select-dropdown__item:hover) {
  background: rgba(255, 255, 255, 0.1);
}

:deep(.el-select-dropdown__item.selected) {
  color: #38bdf8;
  font-weight: 600;
}
</style>
