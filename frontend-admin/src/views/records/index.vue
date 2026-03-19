<template>
  <n-card title="记录管理">
    <template #header-extra>
      <n-button type="primary" @click="handleCreate">新增记录</n-button>
    </template>

    <n-data-table
      :columns="columns"
      :data="records"
      :loading="loading"
      :pagination="pagination"
    />

    <!-- 新增/编辑记录弹窗 -->
    <n-modal
      v-model:show="showFormModal"
      :title="isEditing ? '编辑记录' : '新增记录'"
      preset="card"
      style="width: 1200px"
    >
      <n-form :model="form" :rules="rules" ref="formRef">
        <!-- 头像上传 -->
        <n-form-item label="头像">
          <n-upload
            v-model:file-list="avatarFileList"
            :custom-request="handleAvatarUpload"
            @remove="handleAvatarRemove"
            list-type="image-card"
            :max="1"
            style="display: flex; justify-content: center;"
          >
            <n-button v-if="!form.avatar">上传头像</n-button>
          </n-upload>
        </n-form-item>
        
        <!-- 三列布局 -->
        <div style="display: flex; gap: 24px;">
          <!-- 第一列：基本信息 -->
          <div style="flex: 1;">
            <n-divider title-placement="left">基本信息</n-divider>
            <n-form-item label="姓名" path="name">
              <n-input v-model:value="form.name" placeholder="请输入姓名（英文大写）" />
            </n-form-item>
            <n-form-item label="出生日期">
              <n-date-picker v-model:value="form.date_of_birth" type="date" clearable style="width: 100%;" format="dd MMM yyyy" />
            </n-form-item>
            <n-form-item label="国籍">
              <n-input-group>
                <n-input v-model:value="form.country_of_origin" placeholder="示例：CHINA" />
                <n-button type="primary" ghost @click="form.country_of_origin = 'CHINA'">中国</n-button>
              </n-input-group>
            </n-form-item>
            <n-form-item label="护照号码">
              <n-input v-model:value="form.passport_no" placeholder="示例：EG5895277" />
            </n-form-item>
            <n-form-item label="教育背景">
              <n-input-group>
                <n-input v-model:value="form.education_background" placeholder="示例：Bachelor" />
                <n-popselect v-model:value="form.education_background" :options="educationOptions" trigger="click">
                  <n-button type="primary" ghost>选择</n-button>
                </n-popselect>
              </n-input-group>
            </n-form-item>
          </div>
          
          <!-- 第二列：护照与签证 -->
          <div style="flex: 1;">
            <n-divider title-placement="left">护照与签证</n-divider>
            <n-form-item label="护照签发日期">
              <n-date-picker v-model:value="form.passport_issued_date" type="date" clearable style="width: 100%;" format="dd MMM yyyy" />
            </n-form-item>
            <n-form-item label="护照过期日期">
              <n-date-picker v-model:value="form.passport_expired_date" type="date" clearable style="width: 100%;" format="dd MMM yyyy" />
            </n-form-item>
            <n-form-item label="签证入境日期">
              <n-date-picker v-model:value="form.visa_entry_date" type="date" clearable style="width: 100%;" format="dd MMM yyyy" />
            </n-form-item>
            <n-form-item label="卡片签发日期">
              <n-date-picker v-model:value="form.card_issued_date" type="date" clearable style="width: 100%;" format="dd MMM yyyy" />
            </n-form-item>
            <n-form-item label="卡片过期日期">
              <n-date-picker v-model:value="form.card_expired_date" type="date" clearable style="width: 100%;" format="dd MMM yyyy" />
            </n-form-item>
          </div>
          
          <!-- 第三列：工作经历与职位 -->
          <div style="flex: 1;">
            <n-divider title-placement="left">工作经历</n-divider>
            <n-form-item label="工作年度">
              <n-input v-model:value="form.working_session" placeholder="示例：2026" />
            </n-form-item>
            <n-form-item label="公司名称">
              <n-input v-model:value="form.company_name" placeholder="示例：SELF EMPLOYED" />
            </n-form-item>
            <n-form-item label="开始工作日期">
              <n-date-picker v-model:value="form.start_working_date" type="date" clearable style="width: 100%;" format="dd MMM yyyy" />
            </n-form-item>
            <n-form-item label="停止工作日期">
              <n-date-picker v-model:value="form.stop_working_date" type="date" clearable style="width: 100%;" format="dd MMM yyyy" />
            </n-form-item>
            
            <n-divider title-placement="left">职位信息</n-divider>
            <n-form-item label="职位" path="position">
              <n-input v-model:value="form.position" placeholder="请输入职位" />
            </n-form-item>
            <n-form-item label="性别" path="gender">
              <n-input-group>
                <n-select v-model:value="form.gender" placeholder="请选择性别" :options="genderOptions" style="flex: 1;" />
                <n-button type="primary" ghost @click="form.gender = 'Male'">男</n-button>
                <n-button type="primary" ghost @click="form.gender = 'Female'">女</n-button>
              </n-input-group>
            </n-form-item>
          </div>
        </div>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showFormModal = false">取消</n-button>
          <n-button type="primary" :loading="submitting" @click="handleSubmit">保存</n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 二维码弹窗 -->
    <n-modal
      v-model:show="showQrcodeModal"
      title="二维码"
      preset="card"
      style="width: 400px"
    >
      <div class="flex flex-col items-center gap-4">
        <img
          :src="qrcodeImageUrl"
          width="200"
          height="200"
          alt="QR Code"
        />
        <div class="text-center">
          <p class="mb-2 text-sm text-gray-500">公共访问链接</p>
          <n-input-group>
            <n-input v-model:value="qrcodeLink" readonly />
            <n-button type="primary" @click="copyLink">复制链接</n-button>
          </n-input-group>
          <n-button class="mt-4" type="info" @click="downloadQRCode">下载二维码</n-button>
        </div>
      </div>
    </n-modal>
  </n-card>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, h } from 'vue';
import { NButton, NSpace, NSelect, NPopselect, useMessage, useDialog, type UploadCustomRequestOptions, type UploadFileInfo } from 'naive-ui';
import { fetchGetRecords, fetchCreateRecord, fetchDeleteRecord, fetchUploadImage, fetchUpdateRecord } from '@/service/api';
import { getServiceBaseURL } from '@/utils/service';

// 获取 API baseURL
const isHttpProxy = import.meta.env.DEV && import.meta.env.VITE_HTTP_PROXY === 'Y';
const { baseURL: apiBaseUrl } = getServiceBaseURL(import.meta.env, isHttpProxy);

const message = useMessage();
const dialog = useDialog();

const records = ref<Api.Record.RecordItem[]>([]);
const loading = ref(false);
const showFormModal = ref(false);
const submitting = ref(false);
const formRef = ref<any>(null);
const fileList = ref<UploadFileInfo[]>([]);
const avatarFileList = ref<UploadFileInfo[]>([]);
const editingRecordId = ref<number | null>(null);
const isEditing = computed(() => editingRecordId.value !== null);

// 二维码相关
const showQrcodeModal = ref(false);
const qrcodeLink = ref('');
const qrcodeImageUrl = ref('');

const form = ref<Api.Record.CreateRecordParams & { avatar?: string }>({
  name: '',
  position: '',
  gender: '',
  title: '',
  content: '',
  images: [],
  avatar: '',
  // Foreigner Profile (日期字段使用时间戳)
  date_of_birth: null as number | null,
  country_of_origin: '',
  passport_no: '',
  passport_issued_date: null as number | null,
  passport_expired_date: null as number | null,
  education_background: '',
  visa_entry_date: null as number | null,
  card_issued_date: null as number | null,
  card_expired_date: null as number | null,
  // Working History
  working_session: '',
  company_name: '',
  start_working_date: null as number | null,
  stop_working_date: null as number | null
});

// 日期格式化函数：时间戳转字符串 (dd MMM yyyy)
function formatTimestamp(timestamp: number | null): string {
  if (!timestamp) return '';
  const date = new Date(timestamp);
  const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
  const day = String(date.getDate()).padStart(2, '0');
  const month = months[date.getMonth()];
  const year = date.getFullYear();
  return `${day} ${month} ${year}`;
}

// 字符串日期转时间戳
function parseDateToTimestamp(dateStr: string): number | null {
  if (!dateStr) return null;
  const months: Record<string, number> = {
    'Jan': 0, 'Feb': 1, 'Mar': 2, 'Apr': 3, 'May': 4, 'Jun': 5,
    'Jul': 6, 'Aug': 7, 'Sep': 8, 'Oct': 9, 'Nov': 10, 'Dec': 11
  };
  const parts = dateStr.split(' ');
  if (parts.length === 3) {
    const day = parseInt(parts[0]);
    const month = months[parts[1]];
    const year = parseInt(parts[2]);
    if (!isNaN(day) && !isNaN(month) && !isNaN(year)) {
      return new Date(year, month, day).getTime();
    }
  }
  // 尝试其他格式
  const date = new Date(dateStr);
  return isNaN(date.getTime()) ? null : date.getTime();
}

const genderOptions = [
  { label: '男', value: 'Male' },
  { label: '女', value: 'Female' }
];

const educationOptions = [
  { label: '学士', value: 'Bachelor' },
  { label: '硕士', value: 'Master' },
  { label: '博士', value: 'PhD' },
  { label: '大专', value: 'Associate' },
  { label: '高中', value: 'High School' }
];

const rules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  position: [{ required: true, message: '请输入职位', trigger: 'blur' }],
  gender: [{ required: true, message: '请选择性别', trigger: 'change' }]
};

const columns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '标题', key: 'title' },
  {
    title: '创建时间',
    key: 'created_at',
    render(row: Api.Record.RecordItem) {
      if (!row.created_at) return '';
      return new Date(row.created_at).toLocaleString('zh-CN');
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 300,
    render(row: Api.Record.RecordItem) {
      return h(NSpace, null, {
        default: () => [
          h(
            NButton,
            { size: 'small', type: 'info', onClick: () => handleEdit(row) },
            { default: () => '编辑' }
          ),
          h(
            NButton,
            { size: 'small', type: 'warning', onClick: () => handleShowQrcode(row) },
            { default: () => '二维码' }
          ),
          h(
            NButton,
            { size: 'small', type: 'error', onClick: () => handleDelete(row) },
            { default: () => '删除' }
          )
        ]
      });
    }
  }
];

const pagination = ref({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50]
});

async function fetchRecordList() {
  loading.value = true;
  try {
    const { data } = await fetchGetRecords();
    records.value = (data as any) || [];
  } catch {
    message.error('获取记录失败');
  } finally {
    loading.value = false;
  }
}

async function handleAvatarUpload({ file, onFinish, onError }: UploadCustomRequestOptions) {
  try {
    const { data } = await fetchUploadImage(file.file as File);
    const imageUrl = (data as any).url;
    console.log('Avatar upload success, image URL:', imageUrl);
    
    // 设置头像 URL
    form.value.avatar = imageUrl;
    // 同时更新 images 数组（将头像作为第一个图片）
    form.value.images = [imageUrl];
    console.log('Current avatar:', form.value.avatar);
    
    // 更新 file 对象的属性，让组件正确显示
    file.url = imageUrl.startsWith('http') ? imageUrl : `${apiBaseUrl}${imageUrl}`;
    file.status = 'finished';
    file.thumbnailUrl = file.url;
    
    onFinish();
  } catch (error) {
    console.error('Avatar upload failed:', error);
    file.status = 'error';
    message.error('头像上传失败');
    onError();
  }
}

function handleAvatarRemove() {
  form.value.avatar = '';
  form.value.images = [];
  avatarFileList.value = [];
  return true;
}

async function handleUpload({ file, onFinish, onError }: UploadCustomRequestOptions) {
  try {
    const { data } = await fetchUploadImage(file.file as File);
    const imageUrl = (data as any).url;
    console.log('Upload success, image URL:', imageUrl);
    
    // 添加到 images 数组
    form.value.images.push(imageUrl);
    console.log('Current images array:', form.value.images);
    
    // 更新 file 对象的属性，让组件正确显示
    file.url = imageUrl.startsWith('http') ? imageUrl : `${apiBaseUrl}${imageUrl}`;
    file.status = 'finished';
    file.thumbnailUrl = file.url;
    
    onFinish();
  } catch (error) {
    console.error('Upload failed:', error);
    file.status = 'error';
    message.error('上传失败');
    onError();
  }
}

function handleRemove({ file }: { file: UploadFileInfo }) {
  // 获取文件的 URL（可能是完整 URL 或相对路径）
  const fileUrl = file.url || '';
  // 提取路径部分（去掉 baseURL 前缀）
  let relativeUrl = fileUrl;
  if (fileUrl.startsWith('http')) {
    // 如果是完整 URL，提取路径部分
    try {
      const url = new URL(fileUrl);
      relativeUrl = url.pathname;
    } catch {
      relativeUrl = fileUrl;
    }
  } else if (fileUrl.startsWith(apiBaseUrl)) {
    // 如果以代理前缀开头，去掉前缀
    relativeUrl = fileUrl.replace(apiBaseUrl, '');
  }
  
  // 从 images 数组中移除
  const index = form.value.images.indexOf(relativeUrl);
  if (index > -1) {
    form.value.images.splice(index, 1);
  }
  
  // 从 fileList 中移除
  const fileIndex = fileList.value.findIndex(f => f.id === file.id || f.url === file.url);
  if (fileIndex > -1) {
    fileList.value.splice(fileIndex, 1);
  }
  
  return true;
}

function handleCreate() {
  editingRecordId.value = null;
  form.value = {
    name: '',
    position: '',
    gender: '',
    title: '',
    content: '',
    images: [],
    avatar: '',
    // Foreigner Profile (日期字段使用时间戳)
    date_of_birth: null,
    country_of_origin: '',
    passport_no: '',
    passport_issued_date: null,
    passport_expired_date: null,
    education_background: '',
    visa_entry_date: null,
    card_issued_date: null,
    card_expired_date: null,
    // Working History
    working_session: '',
    company_name: '',
    start_working_date: null,
    stop_working_date: null
  };
  fileList.value = [];
  avatarFileList.value = [];
  showFormModal.value = true;
}

function handleEdit(row: Api.Record.RecordItem) {
  editingRecordId.value = row.id;
  form.value = {
    name: row.name || '',
    position: row.position || '',
    gender: row.gender || '',
    title: row.title,
    content: row.content,
    images: row.images ? [...row.images] : [],
    avatar: (row.images && row.images.length > 0) ? row.images[0] : '',
    // Foreigner Profile (转换日期字符串为时间戳)
    date_of_birth: parseDateToTimestamp(row.date_of_birth || ''),
    country_of_origin: row.country_of_origin || '',
    passport_no: row.passport_no || '',
    passport_issued_date: parseDateToTimestamp(row.passport_issued_date || ''),
    passport_expired_date: parseDateToTimestamp(row.passport_expired_date || ''),
    education_background: row.education_background || '',
    visa_entry_date: parseDateToTimestamp(row.visa_entry_date || ''),
    card_issued_date: parseDateToTimestamp(row.card_issued_date || ''),
    card_expired_date: parseDateToTimestamp(row.card_expired_date || ''),
    // Working History
    working_session: row.working_session || '',
    company_name: row.company_name || '',
    start_working_date: parseDateToTimestamp(row.start_working_date || ''),
    stop_working_date: parseDateToTimestamp(row.stop_working_date || '')
  };
  // 构建头像 fileList
  if (row.images && row.images.length > 0) {
    avatarFileList.value = [{
      id: 'avatar-0',
      name: 'avatar.jpg',
      status: 'finished' as const,
      url: row.images[0].startsWith('http') ? row.images[0] : `${apiBaseUrl}${row.images[0]}`
    }];
  } else {
    avatarFileList.value = [];
  }
  // 构建 fileList 用于显示已有图片（保留兼容性）
  fileList.value = (row.images || []).map((url, index) => ({
    id: `existing-${index}`,
    name: `图片${index + 1}`,
    status: 'finished' as const,
    url: url.startsWith('http') ? url : `${apiBaseUrl}${url}`
  }));
  showFormModal.value = true;
}

async function handleSubmit() {
  try {
    await formRef.value?.validate();
    submitting.value = true;
    
    console.log('Submitting form with images:', form.value.images);
    
    // 准备提交数据，将时间戳转换为字符串格式
    const submitData = {
      name: form.value.name,
      position: form.value.position,
      gender: form.value.gender,
      title: form.value.title,
      content: form.value.content,
      images: form.value.images || [],
      // Foreigner Profile (转换时间戳为字符串)
      date_of_birth: formatTimestamp(form.value.date_of_birth),
      country_of_origin: form.value.country_of_origin,
      passport_no: form.value.passport_no,
      passport_issued_date: formatTimestamp(form.value.passport_issued_date),
      passport_expired_date: formatTimestamp(form.value.passport_expired_date),
      education_background: form.value.education_background,
      visa_entry_date: formatTimestamp(form.value.visa_entry_date),
      card_issued_date: formatTimestamp(form.value.card_issued_date),
      card_expired_date: formatTimestamp(form.value.card_expired_date),
      // Working History
      working_session: form.value.working_session,
      company_name: form.value.company_name,
      start_working_date: formatTimestamp(form.value.start_working_date),
      stop_working_date: formatTimestamp(form.value.stop_working_date)
    };
    
    console.log('Submit data:', submitData);
    
    if (isEditing.value && editingRecordId.value) {
      await fetchUpdateRecord(editingRecordId.value, submitData);
      message.success('更新成功');
    } else {
      await fetchCreateRecord(submitData);
      message.success('创建成功');
    }
    
    showFormModal.value = false;
    editingRecordId.value = null;
    form.value = {
      name: '',
      position: '',
      gender: '',
      title: '',
      content: '',
      images: [],
      avatar: '',
      // Foreigner Profile (日期字段使用时间戳)
      date_of_birth: null,
      country_of_origin: '',
      passport_no: '',
      passport_issued_date: null,
      passport_expired_date: null,
      education_background: '',
      visa_entry_date: null,
      card_issued_date: null,
      card_expired_date: null,
      // Working History
      working_session: '',
      company_name: '',
      start_working_date: null,
      stop_working_date: null
    };
    fileList.value = [];
    avatarFileList.value = [];
    fetchRecordList();
  } catch (error: any) {
    console.error('Submit error:', error);
    if (error?.response?.data?.error) {
      message.error(error.response.data.error);
    }
  } finally {
    submitting.value = false;
  }
}

function handleShowQrcode(row: Api.Record.RecordItem) {
  qrcodeLink.value = `${window.location.origin}/portal/${row.id}`;
  // 使用 QR Server API 生成二维码
  const encodedUrl = encodeURIComponent(qrcodeLink.value);
  qrcodeImageUrl.value = `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodedUrl}`;
  showQrcodeModal.value = true;
}

async function copyLink() {
  try {
    await navigator.clipboard.writeText(qrcodeLink.value);
    message.success('链接已复制到剪贴板');
  } catch {
    message.error('复制失败');
  }
}

async function downloadQRCode() {
  try {
    // 获取二维码图片
    const response = await fetch(qrcodeImageUrl.value);
    const blob = await response.blob();
    
    // 创建图片对象
    const img = new Image();
    img.crossOrigin = 'anonymous';
    
    await new Promise<void>((resolve, reject) => {
      img.onload = () => resolve();
      img.onerror = () => reject(new Error('图片加载失败'));
      img.src = URL.createObjectURL(blob);
    });
    
    // 使用 Canvas 处理图片，将白色背景转为透明
    const canvas = document.createElement('canvas');
    canvas.width = img.width;
    canvas.height = img.height;
    const ctx = canvas.getContext('2d');
    
    if (!ctx) {
      message.error('处理图片失败');
      return;
    }
    
    // 绘制原图
    ctx.drawImage(img, 0, 0);
    
    // 获取图片数据
    const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
    const data = imageData.data;
    
    // 将白色背景转为透明
    for (let i = 0; i < data.length; i += 4) {
      const r = data[i];
      const g = data[i + 1];
      const b = data[i + 2];
      
      // 如果接近白色，设为透明
      if (r > 240 && g > 240 && b > 240) {
        data[i + 3] = 0; // 设置 alpha 为 0（透明）
      }
    }
    
    // 重新绘制处理后的图片
    ctx.putImageData(imageData, 0, 0);
    
    // 转换为 PNG 并下载
    canvas.toBlob((newBlob) => {
      if (!newBlob) {
        message.error('处理图片失败');
        return;
      }
      
      const url = window.URL.createObjectURL(newBlob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `qrcode-${Date.now()}.png`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      message.success('二维码已下载（透明背景）');
    }, 'image/png');
    
  } catch {
    message.error('下载失败');
  }
}

async function viewImage(row: Api.Record.RecordItem) {
  // 如果记录有图片列表，直接打开第一张图片
  if (row.images && row.images.length > 0) {
    const imageUrl = row.images[0];
    const fullUrl = imageUrl.startsWith('http') ? imageUrl : `${apiBaseUrl}${imageUrl}`;
    window.open(fullUrl, '_blank');
    return;
  }
  
  // 否则请求 API 获取图片 URL
  try {
    const response = await fetch(`${apiBaseUrl}/api/records/${row.id}/image`);
    if (!response.ok) {
      message.error('获取图片失败');
      return;
    }
    const data = await response.json();
    if (data.image_url) {
      const fullUrl = data.image_url.startsWith('http') ? data.image_url : `${apiBaseUrl}${data.image_url}`;
      window.open(fullUrl, '_blank');
    } else {
      message.warning('该记录没有图片');
    }
  } catch {
    message.error('获取图片失败');
  }
}

function handleDelete(row: Api.Record.RecordItem) {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除记录 "${row.title}" 吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await fetchDeleteRecord(row.id);
        message.success('删除成功');
        fetchRecordList();
      } catch {
        message.error('删除失败');
      }
    }
  });
}

onMounted(() => {
  fetchRecordList();
});
</script>
