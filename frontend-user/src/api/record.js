import request from './request.js'

// 获取所有记录
export function getRecords() {
  return request.get('/records')
}

// 获取单条记录
export function getRecordById(id) {
  return request.get(`/records/${id}`)
}

// 获取记录渲染后的图片URL
export function getRecordImageUrl(id) {
  return `/api/records/${id}/image`
}
