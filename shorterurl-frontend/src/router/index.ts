import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    // 根据用户登录状态重定向到不同页面
    redirect: () => {
      const token = localStorage.getItem('token')
      return token ? '/link' : '/login'
    }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/login/index.vue'),
    // 已登录用户访问登录页时重定向到短链接页面
    beforeEnter: (to, from, next) => {
      const token = localStorage.getItem('token')
      if (token) {
        next({ path: '/link' })
      } else {
        next()
      }
    }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/register/index.vue')
  },
  {
    path: '/link',
    name: 'ShortLink',
    component: () => import('../views/link/index.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/404',
    name: '404',
    component: () => import('../views/error/404.vue')
  },
  {
    path: '/:catchAll(.*)',
    redirect: '/404'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 导航守卫
router.beforeEach((to, from, next) => {
  // 检查路由是否需要认证
  if (to.matched.some(record => record.meta.requiresAuth)) {
    // 检查本地存储中是否有 token
    const token = localStorage.getItem('token')
    if (!token) {
      // 如果没有 token，重定向到登录页
      next({ name: 'Login' })
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router 