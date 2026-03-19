import { request } from '../request';

/** Get all records */
export function fetchGetRecords() {
  return request<Api.Record.RecordItem[]>({
    url: '/api/records',
    method: 'get'
  });
}

/**
 * Create record
 *
 * @param params Record params
 */
export function fetchCreateRecord(params: Api.Record.CreateRecordParams) {
  return request<Api.Record.RecordItem>({
    url: '/api/admin/records',
    method: 'post',
    data: params
  });
}

/**
 * Delete record
 *
 * @param id Record id
 */
export function fetchDeleteRecord(id: number) {
  return request({
    url: `/api/admin/records/${id}`,
    method: 'delete'
  });
}

/**
 * Upload image
 *
 * @param file Image file
 */
export function fetchUploadImage(file: File) {
  const formData = new FormData();
  formData.append('file', file);
  return request<{ url: string }>({
    url: '/api/admin/upload',
    method: 'post',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' }
  });
}

/** Get public records (no auth required) */
export function fetchGetPublicRecords() {
  return request<Api.Record.RecordItem[]>({
    url: '/api/records',
    method: 'get'
  });
}

/**
 * Get single record by ID (public)
 *
 * @param id Record id
 */
export function fetchGetRecordById(id: number) {
  return request<Api.Record.RecordItem>({
    url: `/api/records/${id}`,
    method: 'get'
  });
}

/**
 * Update record (auth required)
 *
 * @param id Record id
 * @param params Record params
 */
export function fetchUpdateRecord(id: number, params: Api.Record.CreateRecordParams) {
  return request<Api.Record.RecordItem>({
    url: `/api/admin/records/${id}`,
    method: 'put',
    data: params
  });
}

/**
 * Get admin records with pagination (auth required)
 *
 * @param page Page number
 * @param pageSize Page size
 */
export function fetchGetAdminRecords(page: number = 1, pageSize: number = 20) {
  return request({
    url: '/api/admin/records',
    method: 'get',
    params: { page, pageSize }
  });
}
