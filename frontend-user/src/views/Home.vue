<template>
  <div class="home">
    <div class="intro">
      <h2>欢迎访问劳工签证信息查询系统</h2>
      <p>本系统提供劳工签证相关信息查询服务，所有信息以图片形式展示</p>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>

    <!-- 错误提示 -->
    <div v-else-if="error" class="error">
      <p>{{ error }}</p>
      <button @click="fetchRecords">重新加载</button>
    </div>

    <!-- 记录列表 -->
    <div v-else-if="records.length > 0" class="records">
      <div 
        v-for="record in records" 
        :key="record.id" 
        class="record-card"
        @click="showDetail(record)"
      >
        <div class="record-header">
          <h3>{{ record.title }}</h3>
          <span class="record-date">{{ formatDate(record.created_at) }}</span>
        </div>
        <div class="record-preview">
          <img 
            :src="getImageUrl(record.id)" 
            :alt="record.title"
            loading="lazy"
            @error="handleImageError"
          >
        </div>
        <div class="record-footer">
          <span class="view-btn">点击查看详情</span>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="empty">
      <p>暂无记录</p>
    </div>

    <!-- 详情弹窗 -->
    <div v-if="selectedRecord" class="modal" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ selectedRecord.title }}</h3>
          <button class="close-btn" @click="closeModal">&times;</button>
        </div>
        <div class="modal-body">
          <img 
            :src="getImageUrl(selectedRecord.id)" 
            :alt="selectedRecord.title"
            class="detail-image"
          >
          <p class="tip">提示：您可以截图保存此信息</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { getRecords } from '../api/record.js'

export default {
  name: 'Home',
  data() {
    return {
      records: [],
      loading: true,
      error: null,
      selectedRecord: null
    }
  },
  mounted() {
    this.fetchRecords()
  },
  methods: {
    async fetchRecords() {
      this.loading = true
      this.error = null
      try {
        const data = await getRecords()
        this.records = data || []
      } catch (err) {
        this.error = '加载失败，请稍后重试'
        console.error('Failed to fetch records:', err)
      } finally {
        this.loading = false
      }
    },
    getImageUrl(id) {
      return `/api/records/${id}/image`
    },
    formatDate(dateStr) {
      if (!dateStr) return ''
      const date = new Date(dateStr)
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      })
    },
    showDetail(record) {
      this.selectedRecord = record
      document.body.style.overflow = 'hidden'
    },
    closeModal() {
      this.selectedRecord = null
      document.body.style.overflow = ''
    },
    handleImageError(e) {
      e.target.src = 'data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" width="200" height="150"><rect fill="%23f0f0f0" width="200" height="150"/><text fill="%23999" x="50%" y="50%" text-anchor="middle">图片加载失败</text></svg>'
    }
  }
}
</script>

<style scoped>
.home {
  padding: 20px 0;
}

.intro {
  text-align: center;
  margin-bottom: 30px;
  padding: 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
}

.intro h2 {
  color: #333;
  margin-bottom: 10px;
  font-size: 24px;
}

.intro p {
  color: #666;
  font-size: 14px;
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
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error {
  text-align: center;
  padding: 60px;
  color: #e74c3c;
}

.error button {
  margin-top: 15px;
  padding: 10px 25px;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
}

.error button:hover {
  background: #5a6fd6;
}

.empty {
  text-align: center;
  padding: 80px;
  color: #999;
  background: white;
  border-radius: 12px;
}

.records {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.record-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;
}

.record-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0,0,0,0.15);
}

.record-header {
  padding: 15px 20px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.record-header h3 {
  font-size: 16px;
  color: #333;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 200px;
}

.record-date {
  font-size: 12px;
  color: #999;
}

.record-preview {
  height: 200px;
  overflow: hidden;
  background: #f8f8f8;
  display: flex;
  align-items: center;
  justify-content: center;
}

.record-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.record-card:hover .record-preview img {
  transform: scale(1.05);
}

.record-footer {
  padding: 12px 20px;
  text-align: center;
  background: #fafafa;
}

.view-btn {
  color: #667eea;
  font-size: 14px;
  font-weight: 500;
}

/* 弹窗样式 */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  background: white;
  border-radius: 12px;
  max-width: 900px;
  width: 100%;
  max-height: 90vh;
  overflow: hidden;
  animation: modalIn 0.3s ease;
}

@keyframes modalIn {
  from {
    opacity: 0;
    transform: scale(0.9);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.modal-header {
  padding: 20px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 28px;
  color: #999;
  cursor: pointer;
  line-height: 1;
  padding: 0;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.3s;
}

.close-btn:hover {
  background: #f0f0f0;
  color: #333;
}

.modal-body {
  padding: 20px;
  overflow-y: auto;
  max-height: calc(90vh - 80px);
  text-align: center;
}

.detail-image {
  max-width: 100%;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.15);
}

.tip {
  margin-top: 15px;
  color: #999;
  font-size: 13px;
}

@media (max-width: 768px) {
  .records {
    grid-template-columns: 1fr;
  }
  
  .modal-content {
    max-height: 95vh;
  }
}
</style>
