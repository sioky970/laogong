<template>
  <n-layout has-sider style="height: 100vh">
    <n-layout-sider
      bordered
      collapse-mode="width"
      :collapsed-width="64"
      :width="240"
      :collapsed="collapsed"
      show-trigger
      @collapse="collapsed = true"
      @expand="collapsed = false"
    >
      <n-menu
        :collapsed="collapsed"
        :collapsed-width="64"
        :collapsed-icon-size="22"
        :options="menuOptions"
        :value="activeKey"
        @update:value="handleMenuSelect"
      />
    </n-layout-sider>
    <n-layout>
      <n-layout-header bordered style="padding: 16px; display: flex; justify-content: space-between; align-items: center">
        <h2>记录管理</h2>
        <n-space>
          <span>{{ username }}</span>
          <n-button @click="logout" size="small">退出</n-button>
        </n-space>
      </n-layout-header>
      <n-layout-content style="padding: 24px">
        <n-space vertical>
          <n-space>
            <n-button type="primary" @click="showCreateModal = true">新增记录</n-button>
          </n-space>
          
          <n-data-table
            :columns="columns"
            :data="records"
            :loading="loading"
            :pagination="pagination"
          />
        </n-space>
      </n-layout-content>
    </n-layout>
  </n-layout>

  <!-- 创建/编辑模态框 -->
  <n-modal
    v-model:show="showCreateModal"
    title="新增记录"
    preset="card"
    style="width: 600px"
  >
    <n-form :model="form" :rules="rules" ref="formRef">
      <n-form-item label="标题" path="title">
        <n-input v-model:value="form.title" placeholder="请输入标题" />
      </n-form-item>
      <n-form-item label="内容" path="content">
        <n-input
          v-model:value="form.content"
          type="textarea"
          :rows="5"
          placeholder="请输入内容"
        />
      </n-form-item>
      <n-form-item label="图片">
        <n-upload
          multiple
          :custom-request="handleUpload"
          :file-list="fileList"
          @remove="handleRemove"
          list-type="image-card"
        >
          <n-button>上传图片</n-button>
        </n-upload>
      </n-form-item>
    </n-form>
    <template #footer>
      <n-space justify="end">
        <n-button @click="showCreateModal = false">取消</n-button>
        <n-button type="primary" :loading="submitting" @click="handleSubmit">保存</n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<script>
import { ref, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NSpace, useMessage, useDialog } from 'naive-ui'
import { getProfile } from '../api/auth.js'
import { getRecords, createRecord, deleteRecord, uploadImage } from '../api/record.js'

export default {
  name: 'Records',
  setup() {
    const router = useRouter()
    const message = useMessage()
    const dialog = useDialog()
    const collapsed = ref(false)
    const activeKey = ref('records')
    const username = ref('')
    const records = ref([])
    const loading = ref(false)
    const showCreateModal = ref(false)
    const submitting = ref(false)
    const formRef = ref(null)
    const fileList = ref([])

    const form = ref({
      title: '',
      content: '',
      images: []
    })

    const rules = {
      title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
      content: [{ required: true, message: '请输入内容', trigger: 'blur' }]
    }

    const menuOptions = [
      { label: '仪表盘', key: 'dashboard' },
      { label: '记录管理', key: 'records' }
    ]

    const columns = [
      { title: 'ID', key: 'id', width: 80 },
      { title: '标题', key: 'title' },
      { title: '创建时间', key: 'created_at' },
      {
        title: '操作',
        key: 'actions',
        width: 150,
        render(row) {
          return h(NSpace, null, {
            default: () => [
              h(NButton, { size: 'small', onClick: () => viewImage(row) }, { default: () => '查看图片' }),
              h(NButton, { size: 'small', type: 'error', onClick: () => handleDelete(row) }, { default: () => '删除' })
            ]
          })
        }
      }
    ]

    const pagination = ref({
      page: 1,
      pageSize: 10,
      showSizePicker: true,
      pageSizes: [10, 20, 50],
      onChange: (page) => {
        pagination.value.page = page
      },
      onUpdatePageSize: (pageSize) => {
        pagination.value.pageSize = pageSize
        pagination.value.page = 1
      }
    })

    const fetchRecords = async () => {
      loading.value = true
      try {
        records.value = await getRecords()
      } catch (error) {
        message.error('获取记录失败')
      } finally {
        loading.value = false
      }
    }

    const handleMenuSelect = (key) => {
      if (key === 'dashboard') {
        router.push('/')
      }
    }

    const handleUpload = async ({ file, onFinish, onError }) => {
      try {
        const res = await uploadImage(file.file)
        form.value.images.push(res.url)
        onFinish()
      } catch (error) {
        message.error('上传失败')
        onError()
      }
    }

    const handleRemove = ({ file }) => {
      const index = form.value.images.indexOf(file.url)
      if (index > -1) {
        form.value.images.splice(index, 1)
      }
    }

    const handleSubmit = async () => {
      try {
        await formRef.value?.validate()
        submitting.value = true
        await createRecord(form.value)
        message.success('创建成功')
        showCreateModal.value = false
        form.value = { title: '', content: '', images: [] }
        fileList.value = []
        fetchRecords()
      } catch (error) {
        message.error(error.response?.data?.error || '创建失败')
      } finally {
        submitting.value = false
      }
    }

    const viewImage = (row) => {
      window.open(`/api/records/${row.id}/image`, '_blank')
    }

    const handleDelete = (row) => {
      dialog.warning({
        title: '确认删除',
        content: `确定要删除记录 "${row.title}" 吗？`,
        positiveText: '确定',
        negativeText: '取消',
        onPositiveClick: async () => {
          try {
            await deleteRecord(row.id)
            message.success('删除成功')
            fetchRecords()
          } catch (error) {
            message.error('删除失败')
          }
        }
      })
    }

    const logout = () => {
      localStorage.removeItem('token')
      router.push('/login')
    }

    onMounted(async () => {
      try {
        const profile = await getProfile()
        username.value = profile.username
      } catch (error) {
        console.error(error)
      }
      fetchRecords()
    })

    return {
      collapsed,
      activeKey,
      menuOptions,
      username,
      records,
      loading,
      columns,
      pagination,
      showCreateModal,
      form,
      formRef,
      rules,
      submitting,
      fileList,
      handleMenuSelect,
      handleUpload,
      handleRemove,
      handleSubmit,
      logout
    }
  }
}
</script>
