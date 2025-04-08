/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

declare module '@amap/amap-jsapi-loader' {
  const AMapLoader: any
  export default AMapLoader
}

declare module '@/api/*' {
  const api: any
  export default api
}
