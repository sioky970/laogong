import request from './request.js'

export function getRecords() {
  return request.get('/admin/records')
}

export function createRecord(data) {
  return request.post('/admin/records', data)
}

export function updateRecord(id, data) {
  return request.put(`/admin/records/${id}`, data)
}

export function deleteRecord(id) {
  return request.delete(`/admin/records/${id}`)
}

export function uploadImage(file) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/admin/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
