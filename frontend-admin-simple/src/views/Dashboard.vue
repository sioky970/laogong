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
        <h2>劳工签证管理系统</h2>
        <n-space>
          <span>{{ username }}</span>
          <n-button @click="logout" size="small">退出</n-button>
        </n-space>
      </n-layout-header>
      <n-layout-content style="padding: 24px">
        <n-h3>欢迎使用劳工签证管理系统</n-h3>
        <n-p>请从左侧菜单选择功能</n-p>
        
        <n-grid :cols="3" :x-gap="16" :y-gap="16" style="margin-top: 24px">
          <n-grid-item>
            <n-card title="记录总数">
              <n-statistic :value="stats.total" />
            </n-card>
          </n-grid-item>
        </n-grid>
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<script>
import { ref, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { getProfile } from '../api/auth.js'
import { getRecords } from '../api/record.js'

export default {
  name: 'Dashboard',
  setup() {
    const router = useRouter()
    const collapsed = ref(false)
    const activeKey = ref('dashboard')
    const username = ref('')
    const stats = ref({ total: 0 })

    const menuOptions = [
      {
        label: '仪表盘',
        key: 'dashboard'
      },
      {
        label: '记录管理',
        key: 'records'
      }
    ]

    const handleMenuSelect = (key) => {
      if (key === 'records') {
        router.push('/records')
      }
    }

    const logout = () => {
      localStorage.removeItem('token')
      router.push('/login')
    }

    onMounted(async () => {
      try {
        const profile = await getProfile()
        username.value = profile.username
        const records = await getRecords()
        stats.value.total = records.length
      } catch (error) {
        console.error(error)
      }
    })

    return {
      collapsed,
      activeKey,
      menuOptions,
      username,
      stats,
      handleMenuSelect,
      logout
    }
  }
}
</script>
