import axios from 'axios'
import { ElMessage } from 'element-plus'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000
})

api.interceptors.response.use(
  response => {
    const res = response.data
    if (res.code !== 0) {
      ElMessage.error(res.message || '请求失败')
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    return res.data
  },
  error => {
    ElMessage.error(error.message || '网络错误')
    return Promise.reject(error)
  }
)

export const lightingApi = {
  getSchedules: () => api.get('/lighting/schedules'),
  getSchedule: (id) => api.get(`/lighting/schedules/${id}`),
  createSchedule: (data) => api.post('/lighting/schedules', data),
  updateSchedule: (id, data) => api.put(`/lighting/schedules/${id}`, data),
  deleteSchedule: (id) => api.delete(`/lighting/schedules/${id}`),
  getStatus: () => api.get('/lighting/status')
}

export const oxygenApi = {
  getConfigs: () => api.get('/oxygen/configs'),
  getConfig: (id) => api.get(`/oxygen/configs/${id}`),
  createConfig: (data) => api.post('/oxygen/configs', data),
  updateConfig: (id, data) => api.put(`/oxygen/configs/${id}`, data),
  deleteConfig: (id) => api.delete(`/oxygen/configs/${id}`),
  calculateMatch: (data) => api.post('/oxygen/match', data),
  getMatrix: () => api.get('/oxygen/matrix')
}

export const deviceApi = {
  getDevices: (type) => api.get(`/device${type ? `?type=${type}` : ''}`),
  getDevice: (type, name) => api.get(`/device/${type}/${name}`),
  toggleDevice: (type, name, status) => api.put(`/device/${type}/${name}/toggle`, { status }),
  updateValue: (type, name, value) => api.put(`/device/${type}/${name}/value`, { current_value: value }),
  setManualMode: (type, name, manualMode) => api.put(`/device/${type}/${name}/manual`, { manual_mode: manualMode }),
  getDashboardStatus: () => api.get('/device/dashboard')
}

export const sensorApi = {
  getData: (limit, hours) => api.get(`/sensor/data?limit=${limit || 100}&hours=${hours || 24}`),
  getLatest: () => api.get('/sensor/data/latest'),
  getStats: (hours) => api.get(`/sensor/stats?hours=${hours || 24}`),
  createData: (data) => api.post('/sensor/data', data)
}

export default api
