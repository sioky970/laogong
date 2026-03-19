<template>
  <div class="permit-page anti-copy">
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <div class="spinner"></div>
      <p>Loading...</p>
    </div>

    <!-- 错误提示 -->
    <div v-else-if="error" class="error-container">
      <div class="error-icon">⚠️</div>
      <p class="error-message">{{ error }}</p>
      <button class="retry-btn" @click="loadData">Reload</button>
    </div>

    <!-- 详情内容 -->
    <div v-else class="page-content">
      <!-- 顶部横幅 -->
      <div class="header-banner">
        <img src="/tou.png" alt="Ministry of Labour and Vocational Training" class="banner-img" />
      </div>

      <!-- Work Permit Card 卡片 -->
      <div class="card-container">
        <!-- 标题栏 -->
        <div class="card-title-bar">
          <div class="card-icon-circle">
            <svg class="card-icon-svg" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
            </svg>
          </div>
          <span class="card-title-text">Foreigner Work Permit Card</span>
        </div>
        <!-- 头像区域：上蓝下白 -->
        <div class="photo-area">
          <!-- 装饰性半圆弧 -->
          <div class="photo-arc"></div>
          <!-- 头像 -->
          <div class="photo-wrapper">
            <img v-if="photoUrl" :src="photoUrl" alt="Photo" class="photo" @error="handlePhotoError" />
            <div v-else class="photo-placeholder">No Photo</div>
          </div>
        </div>
        <!-- 个人信息区域 -->
        <div class="info-area">
          <h2 class="person-name">{{ displayName }}</h2>
          <p class="person-info">Position: {{ displayPosition }}</p>
          <p class="person-info">Gender: {{ displayGender }}</p>
        </div>
      </div>

      <!-- Foreigner Profile -->
      <div class="profile-section">
        <h3 class="section-title">Foreigner Profile</h3>
        <div class="info-list">
          <div class="info-row">
            <span class="label">Date of Birth:</span>
            <span class="value">{{ profile.dateOfBirth || 'N/A' }}</span>
          </div>
          <div class="info-row">
            <span class="label">Country of Origin:</span>
            <span class="value">{{ profile.countryOfOrigin || 'N/A' }}</span>
          </div>
          <div class="info-row">
            <span class="label">Passport No:</span>
            <span class="value">{{ profile.passportNo || 'N/A' }}</span>
          </div>
          <div class="info-row">
            <span class="label">Passport Issued Date:</span>
            <span class="value">{{ profile.passportIssuedDate || 'N/A' }}</span>
          </div>
          <div class="info-row">
            <span class="label">Passport Expired Date:</span>
            <span class="value">{{ profile.passportExpiredDate || 'N/A' }}</span>
          </div>
          <div class="info-row">
            <span class="label">Education Background:</span>
            <span class="value">{{ profile.educationBackground || 'N/A' }}</span>
          </div>
          <div class="info-row">
            <span class="label">Visa Entry Date:</span>
            <span class="value highlight">{{ profile.visaEntryDate || 'N/A' }}</span>
          </div>
          <div class="info-row">
            <span class="label">Latest Card Issued Date:</span>
            <span class="value">{{ profile.latestCardIssuedDate || 'N/A' }}</span>
          </div>
          <div class="info-row">
            <span class="label">Latest Card Expired Date:</span>
            <span class="value">{{ profile.latestCardExpiredDate || 'N/A' }}</span>
          </div>
        </div>
      </div>

      <!-- Working History -->
      <div class="history-section">
        <h3 class="section-title">Working History</h3>
        <div class="table-wrapper">
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
              <tr v-for="(row, index) in profile.workingHistory" :key="index">
                <td>{{ row.session }}</td>
                <td>{{ row.companyName }}</td>
                <td>{{ row.startWorkingDate }}</td>
                <td>{{ row.stopWorking || '-' }}</td>
              </tr>
              <tr v-if="!profile.workingHistory || profile.workingHistory.length === 0">
                <td colspan="4" class="no-data">No working history</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';

interface WorkingHistoryItem {
  session: string;
  companyName: string;
  startWorkingDate: string;
  stopWorking: string;
}

interface ProfileData {
  name: string;
  position: string;
  gender: string;
  dateOfBirth: string;
  countryOfOrigin: string;
  passportNo: string;
  passportIssuedDate: string;
  passportExpiredDate: string;
  educationBackground: string;
  visaEntryDate: string;
  latestCardIssuedDate: string;
  latestCardExpiredDate: string;
  workingHistory: WorkingHistoryItem[];
}

const route = useRoute();
const loading = ref(true);
const error = ref<string | null>(null);
const photoUrl = ref<string>('');

// 用于显示的字段（优先使用 API 顶层字段，回退到 content JSON）
const displayName = ref('');
const displayPosition = ref('');
const displayGender = ref('');

const profile = reactive<ProfileData>({
  name: '',
  position: '',
  gender: '',
  dateOfBirth: '',
  countryOfOrigin: '',
  passportNo: '',
  passportIssuedDate: '',
  passportExpiredDate: '',
  educationBackground: '',
  visaEntryDate: '',
  latestCardIssuedDate: '',
  latestCardExpiredDate: '',
  workingHistory: []
});

// 获取 API base URL
function getApiBaseUrl(): string {
  const isDev = import.meta.env.DEV;
  const isHttpProxy = isDev && import.meta.env.VITE_HTTP_PROXY === 'Y';

  if (isHttpProxy) {
    return '/proxy-default';
  }

  return import.meta.env.VITE_SERVICE_BASE_URL || '';
}

async function loadData() {
  const id = route.params.id as string;
  if (!id) {
    error.value = 'Invalid record ID';
    loading.value = false;
    return;
  }

  loading.value = true;
  error.value = null;

  const baseURL = getApiBaseUrl();

  try {
    const response = await fetch(`${baseURL}/api/records/${id}`);
    if (!response.ok) {
      if (response.status === 404) {
        throw new Error('Record not found');
      }
      throw new Error('Failed to fetch record');
    }
    const data = await response.json();

    // 解析 content JSON
    if (data.content) {
      try {
        const contentData = typeof data.content === 'string' ? JSON.parse(data.content) : data.content;
        Object.assign(profile, contentData);
      } catch {
        // 如果解析失败，尝试将 content 作为名字使用
        profile.name = data.content;
      }
    }

    // 设置头像 URL
    if (data.images && Array.isArray(data.images) && data.images.length > 0) {
      const imagePath = data.images[0];
      // 如果图片路径是相对路径，拼接 baseURL
      if (imagePath.startsWith('http')) {
        photoUrl.value = imagePath;
      } else {
        photoUrl.value = baseURL + imagePath;
      }
    }

    // 如果没有姓名，使用 title
    if (!profile.name && data.title) {
      profile.name = data.title;
    }

    // 设置显示字段（优先使用 API 顶层字段，回退到 content JSON）
    displayName.value = data.name || profile.name || 'N/A';
    displayPosition.value = data.position || profile.position || 'N/A';
    displayGender.value = data.gender || profile.gender || 'N/A';
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load, please try again';
  } finally {
    loading.value = false;
  }
}

function handlePhotoError() {
  photoUrl.value = '';
}

// 防复制事件处理
function preventAction(e: Event) {
  e.preventDefault();
}

onMounted(() => {
  loadData();

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
.permit-page {
  min-height: 100vh;
  background: #f0f2f5;
  display: flex;
  flex-direction: column;
  align-items: center;
}

/* 防复制样式 */
.anti-copy {
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
  -webkit-touch-callout: none;
}

/* 加载状态 */
.loading-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #666;
  width: 100%;
}

.spinner {
  width: 50px;
  height: 50px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #2e6eb5;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 错误状态 */
.error-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
  width: 100%;
}

.error-icon {
  font-size: 60px;
  margin-bottom: 20px;
}

.error-message {
  font-size: 18px;
  color: #e74c3c;
  margin-bottom: 30px;
}

.retry-btn {
  padding: 12px 30px;
  background: #2e6eb5;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.retry-btn:hover {
  background: #245a94;
}

/* 页面内容 */
.page-content {
  width: 100%;
  max-width: 540px;
  background: #f0f2f5;
}

/* 顶部横幅 */
.header-banner {
  width: 100%;
  background: white;
}

.banner-img {
  width: 100%;
  display: block;
}

/* Work Permit Card */
.card-container {
  margin: 15px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

/* 标题栏 */
.card-title-bar {
  background: linear-gradient(135deg, #3B82F6 0%, #60A5FA 100%);
  color: white;
  padding: 14px 18px;
  font-size: 16px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 12px;
}

.card-icon-circle {
  width: 28px;
  height: 28px;
  background: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.card-icon-svg {
  width: 18px;
  height: 18px;
  color: #3B82F6;
}

.card-title-text {
  font-weight: 600;
  letter-spacing: 0.3px;
}

/* 头像区域：上蓝下白 */
.photo-area {
  position: relative;
  background: linear-gradient(to bottom, #4A9BE8 50%, #ffffff 50%);
  padding: 25px 0 35px;
  text-align: center;
}

/* 装饰性半圆弧 */
.photo-arc {
  position: absolute;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  width: 170px;
  height: 85px;
  border: 3px solid rgba(255, 255, 255, 0.4);
  border-bottom: none;
  border-radius: 85px 85px 0 0;
}

.photo-wrapper {
  position: relative;
  z-index: 1;
  display: inline-block;
}

.photo {
  width: 140px;
  height: 170px;
  object-fit: cover;
  border: 3px solid white;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  background: linear-gradient(to bottom, #4A9BE8 50%, #ffffff 50%);
}

.photo-placeholder {
  width: 140px;
  height: 170px;
  background: #e0e0e0;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 14px;
  border: 3px solid white;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 个人信息区域 */
.info-area {
  background: white;
  padding: 20px 20px 28px;
  text-align: center;
}

.person-name {
  font-size: 22px;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0 0 12px 0;
  letter-spacing: 0.5px;
}

.person-info {
  font-size: 15px;
  color: #666;
  margin: 6px 0;
}

/* Foreigner Profile */
.profile-section {
  margin: 15px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.section-title {
  color: #2e6eb5;
  font-size: 18px;
  font-weight: 700;
  text-align: center;
  margin: 0 0 20px 0;
  padding-bottom: 10px;
  border-bottom: 2px solid #2e6eb5;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-row {
  display: flex;
  font-size: 14px;
  line-height: 1.5;
}

.label {
  flex: 0 0 180px;
  color: #555;
  font-weight: 500;
}

.value {
  flex: 1;
  color: #333;
  font-weight: 600;
}

.value.highlight {
  color: #e74c3c;
}

/* Working History */
.history-section {
  margin: 15px;
  margin-bottom: 30px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.table-wrapper {
  overflow-x: auto;
}

.history-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.history-table th,
.history-table td {
  border: 1px solid #ddd;
  padding: 10px 8px;
  text-align: center;
}

.history-table th {
  background: #f5f5f5;
  font-weight: 600;
  color: #333;
}

.history-table td {
  color: #555;
}

.no-data {
  color: #999;
  font-style: italic;
}

/* 移动端适配 */
@media (max-width: 540px) {
  .card-container,
  .profile-section,
  .history-section {
    margin: 10px;
  }

  .card-title-bar {
    padding: 12px 15px;
    font-size: 15px;
  }

  .card-icon-circle {
    width: 24px;
    height: 24px;
  }

  .card-icon-svg {
    width: 16px;
    height: 16px;
  }

  .photo-area {
    padding: 20px 0 30px;
  }

  .photo-arc {
    width: 150px;
    height: 75px;
    border-radius: 75px 75px 0 0;
    top: 15px;
  }

  .photo {
    width: 120px;
    height: 145px;
  }

  .photo-placeholder {
    width: 120px;
    height: 145px;
  }

  .info-area {
    padding: 15px 15px 24px;
  }

  .person-name {
    font-size: 20px;
  }

  .person-info {
    font-size: 14px;
  }

  .section-title {
    font-size: 16px;
  }

  .label {
    flex: 0 0 140px;
    font-size: 13px;
  }

  .value {
    font-size: 13px;
  }

  .history-table {
    font-size: 12px;
  }

  .history-table th,
  .history-table td {
    padding: 8px 5px;
  }
}
</style>
