<template>
  <div class="portal anti-copy">
    <div class="header-banner">
      <img src="/tou.png" alt="Ministry of Labour and Vocational Training" class="banner-img" />
    </div>

    <main class="portal-main">
      <!-- 错误提示 -->
      <div v-if="error" class="error">
        <p>{{ error }}</p>
      </div>

      <!-- 记录列表 -->
      <div v-else-if="records.length > 0" class="records">
        <div
          v-for="record in records"
          :key="record.id"
          class="record-card"
        >
          <div class="card-header">
            <h3>Foreigner Work Permit Card</h3>
          </div>
          <div class="card-body">
            <!-- 基本信息区域：头像 + 姓名 -->
            <div class="basic-info">
              <div class="avatar-row">
                <div class="avatar-wrapper">
                  <img
                    v-if="record.images && record.images.length > 0"
                    :src="getFullImageUrl(record.images[0])"
                    :alt="record.name"
                    class="avatar"
                    loading="lazy"
                    @error="handleImageError"
                  >
                  <div v-else class="avatar-placeholder">No Photo</div>
                </div>
              </div>
              <div class="basic-details">
                <h4 class="worker-name">{{ (record.name || 'UNKNOWN').toUpperCase() }}</h4>
                <p class="worker-info"><span class="label">Position:</span> {{ record.position || '-' }}</p>
                <p class="worker-info"><span class="label">Gender:</span> {{ record.gender || '-' }}</p>
              </div>
            </div>

            <!-- Foreigner Profile -->
            <div class="profile-section">
              <div class="profile-title">Foreigner Profile</div>
              <div class="profile-table">
                <div class="profile-line">
                  <span class="profile-label">Date of Birth:</span>
                  <span class="profile-value">{{ record.date_of_birth || '-' }}</span>
                </div>
                <div class="profile-line">
                  <span class="profile-label">Country of Origin:</span>
                  <span class="profile-value">{{ record.country_of_origin || '-' }}</span>
                </div>
                <div class="profile-line">
                  <span class="profile-label">Passport No:</span>
                  <span class="profile-value">{{ record.passport_no || '-' }}</span>
                </div>
                <div class="profile-line">
                  <span class="profile-label">Passport Issued Date:</span>
                  <span class="profile-value">{{ record.passport_issued_date || '-' }}</span>
                </div>
                <div class="profile-line">
                  <span class="profile-label">Passport Expired Date:</span>
                  <span class="profile-value">{{ record.passport_expired_date || '-' }}</span>
                </div>
                <div class="profile-line">
                  <span class="profile-label">Education Background:</span>
                  <span class="profile-value">{{ record.education_background || '-' }}</span>
                </div>
                <div class="profile-line">
                  <span class="profile-label">Visa Entry Date:</span>
                  <span class="profile-value highlight">{{ record.visa_entry_date || '-' }}</span>
                </div>
                <div class="profile-line">
                  <span class="profile-label">Latest Card Issued Date:</span>
                  <span class="profile-value">{{ record.card_issued_date || '-' }}</span>
                </div>
                <div class="profile-line">
                  <span class="profile-label">Latest Card Expired Date:</span>
                  <span class="profile-value">{{ record.card_expired_date || '-' }}</span>
                </div>
              </div>
            </div>

            <!-- Working History -->
            <div class="profile-section">
              <div class="profile-title">Working History</div>
              <table class="history-table">
                <thead>
                  <tr>
                    <th>Session</th>
                    <th>Company Name</th>
                    <th>Start Working Date</th>
                    <th>Stop Working</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>{{ record.working_session || '-' }}</td>
                    <td>{{ record.company_name || '-' }}</td>
                    <td>{{ record.start_working_date || '-' }}</td>
                    <td>{{ record.stop_working_date || '-' }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="!loading && records.length === 0" class="empty">
        <p>No records found</p>
      </div>
    </main>

    <!-- 法律警告文本 -->
    <div class="legal-warning">
      <p class="warning-text">Anyone who counterfeits the information in QR Code of this foreign work permit and foreign employment card will be convicted by law.</p>
      <p class="warning-text khmer">អ្នកណាដែលក្លែងក្លាយព័ត៌មាននៅក្នុង QR Code នៃលិខិតអនុញ្ញាតធ្វើការជាបរទេស និងកាតជាបរទេសនេះ នឹងត្រូវផ្តន្ទាទោសតាមច្បាប់។</p>
    </div>

    <footer class="portal-footer">
      <p>Copyright 2025 Ministry of Labour and Vocational Training, Kingdom of Cambodia.</p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { fetchGetPublicRecords } from '@/service/api';
import { getServiceBaseURL } from '@/utils/service';

// 获取 API baseURL
const isHttpProxy = import.meta.env.DEV && import.meta.env.VITE_HTTP_PROXY === 'Y';
const { baseURL: apiBaseUrl } = getServiceBaseURL(import.meta.env, isHttpProxy);

interface RecordItem {
  id: number;
  name?: string;
  position?: string;
  gender?: string;
  title: string;
  created_at: string;
  image_url?: string;
  images?: string[];
  // Foreigner Profile
  date_of_birth?: string;
  country_of_origin?: string;
  passport_no?: string;
  passport_issued_date?: string;
  passport_expired_date?: string;
  education_background?: string;
  visa_entry_date?: string;
  card_issued_date?: string;
  card_expired_date?: string;
  // Working History
  working_session?: string;
  company_name?: string;
  start_working_date?: string;
  stop_working_date?: string;
}

const records = ref<RecordItem[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);

async function fetchRecords() {
  loading.value = true;
  error.value = null;
  try {
    const { data } = await fetchGetPublicRecords();
    records.value = (data as any) || [];
  } catch (err) {
    error.value = 'Failed to load, please try again later';
  } finally {
    loading.value = false;
  }
}

function getFullImageUrl(imageUrl?: string) {
  if (!imageUrl) return '';
  // 如果已经是完整URL，直接返回
  if (imageUrl.startsWith('http')) return imageUrl;
  // 对于静态文件路径（/uploads/），需要通过代理访问
  // 在开发模式下，apiBaseUrl 是 /proxy-default
  return `${apiBaseUrl}${imageUrl}`;
}

function handleImageError(e: Event) {
  (e.target as HTMLImageElement).src =
    'data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" width="200" height="150"><rect fill="%23f0f0f0" width="200" height="150"/><text fill="%23999" x="50%" y="50%" text-anchor="middle">Image load failed</text></svg>';
}

// 防复制事件处理
function preventAction(e: Event) {
  e.preventDefault();
}

onMounted(() => {
  fetchRecords();

  // 添加防复制事件监听
  document.addEventListener('contextmenu', preventAction);
  document.addEventListener('dragstart', preventAction);
  document.addEventListener('selectstart', preventAction);
});

onUnmounted(() => {
  // 移除防复制事件监听
  document.removeEventListener('contextmenu', preventAction);
  document.removeEventListener('dragstart', preventAction);
  document.removeEventListener('selectstart', preventAction);
});
</script>

<style scoped>
/* 防复制样式 */
.anti-copy {
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
  -webkit-touch-callout: none;
}

.portal {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
}

.header-banner {
  width: 100%;
}

.banner-img {
  width: 100%;
  display: block;
}

.portal-main {
  flex: 1;
  padding: 30px 20px;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
}

.portal-footer {
  background: #fff;
  color: #333;
  text-align: center;
  padding: 15px;
  font-size: 14px;
  border-top: 1px solid #eee;
}

.loading {
  text-align: center;
  padding: 60px;
  color: #666;
}

.spinner {
  width: 50px;
  height: 50px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 20px;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.error {
  text-align: center;
  padding: 60px;
  color: #e74c3c;
}



.empty {
  text-align: center;
  padding: 80px;
  color: #999;
  background: white;
  border-radius: 12px;
}

.records {
  display: flex;
  flex-direction: column;
  gap: 24px;
  max-width: 600px;
  margin: 0 auto;
}

.record-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  background: linear-gradient(135deg, #1e88e5 0%, #1565c0 100%);
  padding: 12px 20px;
  text-align: center;
}

.card-header h3 {
  font-size: 16px;
  color: white;
  margin: 0;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.card-body {
  padding: 20px;
}

/* 基本信息区域 */
.basic-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding-bottom: 16px;
  border-bottom: 1px solid #eee;
  margin-bottom: 16px;
}

.avatar-row {
  width: 100%;
  display: flex;
  justify-content: center;
}

.basic-details {
  width: 100%;
  text-align: center;
}

.avatar-wrapper {
  width: 80px;
  height: 100px;
  flex-shrink: 0;
  overflow: hidden;
  border-radius: 4px;
}

.avatar {
  width: 100%;
  height: 100%;
  object-fit: fill;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f0f0;
  color: #999;
  font-size: 11px;
}

.worker-name {
  font-size: 18px;
  font-weight: 700;
  color: #333;
  margin: 0 0 6px 0;
  letter-spacing: 1px;
  line-height: 1.2;
}

.worker-info {
  font-size: 13px;
  color: #666;
  margin: 3px 0;
  line-height: 1.3;
}

.worker-info .label {
  color: #999;
}

/* Profile Section */
.profile-section {
  margin-bottom: 20px;
  background: #fff;
  border-radius: 8px;
  padding: 10px 4px;
}

.profile-section:last-child {
  margin-bottom: 0;
}

.profile-title {
  font-size: 14px;
  font-weight: 600;
  color: #1565c0;
  text-align: center;
  margin-bottom: 16px;
}

/* Profile Table - 标签值两端对齐 */
.profile-table {
  display: flex;
  flex-direction: column;
}

.profile-line {
  display: flex;
  align-items: center;
  padding: 5px 0;
  gap: 8px;
}

.profile-line:last-child {
  border-bottom: none;
}

.profile-line .profile-label {
  font-size: 11px;
  color: #666;
  font-weight: 600;
  flex: 0 0 50%;
  text-align: left;
}

.profile-line .profile-value {
  font-size: 11px;
  color: #333;
  font-weight: 700;
  flex: 0 0 50%;
  text-align: left;
}

.profile-line .profile-value.highlight {
  color: #c62828;
}

/* Working History Table */
.history-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.history-table th,
.history-table td {
  padding: 10px 8px;
  text-align: left;
  border: 1px solid #e0e0e0;
}

.history-table th {
  background: #f5f5f5;
  font-weight: 600;
  color: #333;
}

.history-table td {
  color: #333;
}

@media (max-width: 768px) {
  .records {
    max-width: 100%;
  }
  
  .profile-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
  
  .profile-value {
    text-align: left;
  }
  
  .history-table {
    font-size: 11px;
  }
  
  .history-table th,
  .history-table td {
    padding: 8px 4px;
  }
}

/* 法律警告文本 */
.legal-warning {
  background: #fce4ec;
  border-top: 3px solid #f48fb1;
  padding: 16px 20px;
  text-align: center;
}

.warning-text {
  font-size: 12px;
  color: #880e4f;
  font-weight: 600;
  margin: 0 0 8px 0;
  line-height: 1.5;
}

.warning-text.khmer {
  font-size: 11px;
  margin-bottom: 0;
}

@media (max-width: 768px) {
  .legal-warning {
    padding: 12px 16px;
  }
  
  .warning-text {
    font-size: 11px;
  }
  
  .warning-text.khmer {
    font-size: 10px;
  }
}
</style>
