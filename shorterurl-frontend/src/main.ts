import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import './style.css'

// 导入Element Plus图标
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

// 创建应用实例
const app = createApp(App)

// 注册所有Element Plus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 设置axios默认值
axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL || ''
axios.defaults.timeout = 15000

// 添加请求拦截器，自动添加token和username到请求头
axios.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    const username = localStorage.getItem('username')
    
    if (token && username) {
      // 确保headers对象存在
      config.headers = config.headers || {}
      // 使用标准格式的header名称（小写）
      config.headers['token'] = token
      config.headers['username'] = username
      
      // 打印请求信息，方便调试
      console.log(`请求添加头信息: username=${username}, token=${token.substring(0, 8)}...`)
    }
    
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 使用插件
app.use(createPinia())
app.use(router)
app.use(ElementPlus)

// 挂载应用
app.mount('#app')
