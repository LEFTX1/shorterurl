<template>
  <el-dialog
    v-model="dialogVisible"
    title="分组访问统计"
    width="90%"
    class="stats-dialog"
  >
    <div class="stats-container">
      <!-- 查询参数设置 -->
      <el-card shadow="hover" class="query-card">
        <el-row :gutter="20" align="middle">
          <el-col :span="12">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              format="YYYY-MM-DD"
              :shortcuts="dateRangePickerOptions.shortcuts"
              :disabled-date="dateRangePickerOptions.disabledDate"
              size="default"
            />
          </el-col>
          <el-col :span="4">
            <el-button type="primary" @click="fetchStatsData">查询统计</el-button>
          </el-col>
        </el-row>
      </el-card>

      <!-- 基础统计信息 -->
      <el-row :gutter="20" class="stats-cards">
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>
              <div class="card-header">
                <span>总访问量(PV)</span>
              </div>
            </template>
            <div class="card-value">{{ statsData?.overallPvUvUipStats?.pv || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>
              <div class="card-header">
                <span>独立访客数(UV)</span>
              </div>
            </template>
            <div class="card-value">{{ statsData?.overallPvUvUipStats?.uv || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>
              <div class="card-header">
                <span>独立IP数(UIP)</span>
              </div>
            </template>
            <div class="card-value">{{ statsData?.overallPvUvUipStats?.uip || 0 }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 详细统计数据 - 两列布局 -->
      <el-row :gutter="20" class="stats-detail">
        <!-- 左侧列 -->
        <el-col :span="12">
          <!-- 每日访问趋势 -->
          <el-card shadow="hover" class="stats-card">
            <template #header>
              <div class="card-header">
                <span>每日访问趋势</span>
              </div>
            </template>
            <div ref="dailyTrendChartRef" class="chart"></div>
          </el-card>

          <!-- 设备分布 -->
          <el-card shadow="hover" class="stats-card">
            <template #header>
              <div class="card-header">
                <span>设备分布</span>
              </div>
            </template>
            <div ref="deviceChartRef" class="chart"></div>
          </el-card>

          <!-- 浏览器分布 -->
          <el-card shadow="hover" class="stats-card">
            <template #header>
              <div class="card-header">
                <span>浏览器分布</span>
              </div>
            </template>
            <div ref="browserChartRef" class="chart"></div>
          </el-card>
        </el-col>

        <!-- 右侧列 -->
        <el-col :span="12">
          <!-- 操作系统分布 -->
          <el-card shadow="hover" class="stats-card">
            <template #header>
              <div class="card-header">
                <span>操作系统分布</span>
              </div>
            </template>
            <div ref="osChartRef" class="chart"></div>
          </el-card>

          <!-- 网络类型分布 -->
          <el-card shadow="hover" class="stats-card">
            <template #header>
              <div class="card-header">
                <span>网络类型分布</span>
              </div>
            </template>
            <div ref="networkChartRef" class="chart"></div>
          </el-card>

          <!-- 地区分布 -->
          <el-card shadow="hover" class="stats-card">
            <template #header>
              <div class="card-header">
                <span>地区分布</span>
              </div>
            </template>
            <div ref="localeChartRef" class="chart"></div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, defineProps, defineEmits, watch, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import statsApi from '../api/stats'
import { useUserStore } from '../store/user'

// 定义响应数据类型
interface PvUvUipStats {
  date: string
  pv: number
  uv: number
  uip: number
}

interface OverallPvUvUipStats {
  date: string
  pv: number
  uv: number
  uip: number
}

interface LocaleStats {
  locale: string
  cnt: number
  ratio: number
}

interface BrowserStats {
  browser: string
  cnt: number
  ratio: number
}

interface OsStats {
  os: string
  cnt: number
  ratio: number
}

interface DeviceStats {
  device: string
  cnt: number
  ratio: number
}

interface NetworkStats {
  network: string
  cnt: number
  ratio: number
}

interface GroupStatsRespDTO {
  pvUvUipStatsList: PvUvUipStats[]
  overallPvUvUipStats: OverallPvUvUipStats
  localeCnStats: LocaleStats[]
  hourStats: number[]
  topIpStats: { ip: string; cnt: number; ratio: number }[]
  weekdayStats: number[]
  browserStats: BrowserStats[]
  osStats: OsStats[]
  uvTypeStats: { uvType: string; cnt: number; ratio: number }[]
  deviceStats: DeviceStats[]
  networkStats: NetworkStats[]
}

const props = defineProps({
  visible: {
    type: Boolean,
    required: true
  },
  gid: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['update:visible'])
const router = useRouter()

const dialogVisible = ref(props.visible)
const statsData = ref<GroupStatsRespDTO | null>(null)

// 日期范围选择器配置
const dateRangePickerOptions = {
  shortcuts: [
    {
      text: '最近7天',
      value: () => {
        const end = new Date()
        const start = new Date()
        start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
        return [start, end]
      }
    },
    {
      text: '最近30天',
      value: () => {
        const end = new Date()
        const start = new Date()
        start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
        return [start, end]
      }
    },
    {
      text: '最近90天',
      value: () => {
        const end = new Date()
        const start = new Date()
        start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
        return [start, end]
      }
    }
  ],
  disabledDate: (time: Date) => {
    return time.getTime() > Date.now()
  }
}

// 日期范围
const dateRange = ref<[Date, Date]>([new Date(), new Date()])

// 日期快捷选项
const dateShortcuts = [
  {
    text: '今天',
    value: () => {
      const today = new Date().toISOString().split('T')[0]
      return [today, today]
    }
  },
  {
    text: '最近一周',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
      return [start.toISOString().split('T')[0], end.toISOString().split('T')[0]]
    }
  },
  {
    text: '最近一个月',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setMonth(start.getMonth() - 1)
      return [start.toISOString().split('T')[0], end.toISOString().split('T')[0]]
    }
  }
]

// 图表引用
const dailyTrendChartRef = ref<HTMLElement>()
const deviceChartRef = ref<HTMLElement>()
const browserChartRef = ref<HTMLElement>()
const osChartRef = ref<HTMLElement>()
const networkChartRef = ref<HTMLElement>()
const localeChartRef = ref<HTMLElement>()

// 图表实例
let dailyTrendChart: echarts.ECharts | null = null
let deviceChart: echarts.ECharts | null = null
let browserChart: echarts.ECharts | null = null
let osChart: echarts.ECharts | null = null
let networkChart: echarts.ECharts | null = null
let localeChart: echarts.ECharts | null = null

// 监听对话框可见性
watch(() => props.visible, (val) => {
  dialogVisible.value = val
})

// 监听对话框关闭
watch(() => dialogVisible.value, (val) => {
  emit('update:visible', val)
})

// 监听对话框显示状态
watch(dialogVisible, (newVal) => {
  if (newVal) {
    // 设置默认日期范围为一个非常早的日期到今天
    const end = new Date()
    const start = new Date('2020-01-01') // 设置一个固定的非常早的开始日期
    dateRange.value = [start, end]
  }
})

// 监听日期范围变化
watch(dateRange, (newVal) => {
  if (newVal && newVal.length === 2) {
    fetchStatsData()
  }
}, { deep: true })

// 获取统计数据
const fetchStatsData = async () => {
  const userStore = useUserStore()
  
  // 检查登录状态
  if (!userStore.checkLogin()) {
    ElMessage.error('请先登录')
    dialogVisible.value = false
    router.push('/login')
    return
  }

  try {
    const [startDate, endDate] = dateRange.value

    // 格式化日期为 YYYY-MM-DD
    const formatDate = (date: Date) => {
      return date.toISOString().split('T')[0]
    }

    console.log('开始获取分组统计数据，参数:', {
      gid: props.gid,
      startDate: formatDate(startDate),
      endDate: formatDate(endDate)
    })

    // 获取分组统计数据
    const response = await statsApi.getShortLinkGroupStats({
      gid: props.gid,
      startDate: formatDate(startDate),
      endDate: formatDate(endDate)
    })

    console.log('获取分组统计数据响应:', response)
    console.log('响应数据类型:', typeof response)
    console.log('响应数据是否为空:', !response)
    console.log('响应数据内容:', JSON.stringify(response, null, 2))

    if (response && response.pvUvUipStatsList) {
      console.log('设置统计数据前:', statsData.value)
      statsData.value = response
      console.log('设置统计数据后:', statsData.value)
      
      // 重新初始化图表
      setTimeout(() => {
        console.log('开始初始化图表')
        initCharts()
      }, 100)
    } else {
      console.error('响应数据格式错误')
      ElMessage.error('获取统计数据失败：响应数据格式错误')
      statsData.value = null
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
    ElMessage.error('获取统计数据失败')
    statsData.value = null
  }
}

// 初始化图表
const initCharts = () => {
  console.log('开始初始化图表，当前统计数据:', statsData.value)
  
  if (!statsData.value) {
    console.error('统计数据为空，无法初始化图表')
    return
  }

  const data = statsData.value
  console.log('图表数据:', data)

  // 每日访问趋势图表
  if (dailyTrendChartRef.value && data.pvUvUipStatsList) {
    console.log('初始化每日访问趋势图表')
    dailyTrendChart = echarts.init(dailyTrendChartRef.value)
    dailyTrendChart.setOption({
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross'
        }
      },
      legend: {
        data: ['PV', 'UV', 'UIP']
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: data.pvUvUipStatsList.map(item => item.date)
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: 'PV',
          type: 'line',
          data: data.pvUvUipStatsList.map(item => item.pv)
        },
        {
          name: 'UV',
          type: 'line',
          data: data.pvUvUipStatsList.map(item => item.uv)
        },
        {
          name: 'UIP',
          type: 'line',
          data: data.pvUvUipStatsList.map(item => item.uip)
        }
      ]
    })
  } else {
    console.error('每日访问趋势图表初始化失败:', {
      hasRef: !!dailyTrendChartRef.value,
      hasData: !!data.pvUvUipStatsList
    })
  }

  // 设备分布图表
  if (deviceChartRef.value && data.deviceStats) {
    console.log('初始化设备分布图表')
    deviceChart = echarts.init(deviceChartRef.value)
    deviceChart.setOption({
      tooltip: {
        trigger: 'item',
        formatter: '{b}: {c} ({d}%)'
      },
      legend: {
        orient: 'vertical',
        right: 10,
        top: 'center',
        data: data.deviceStats.map(item => item.device)
      },
      series: [
        {
          name: '设备类型',
          type: 'pie',
          radius: ['40%', '70%'],
          center: ['40%', '50%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: {
            show: true,
            position: 'outside',
            formatter: '{b}: {c}'
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 14,
              fontWeight: 'bold'
            }
          },
          data: data.deviceStats.map(item => ({
            name: item.device,
            value: item.cnt
          }))
        }
      ]
    })
  } else {
    console.error('设备分布图表初始化失败:', {
      hasRef: !!deviceChartRef.value,
      hasData: !!data.deviceStats
    })
  }

  // 浏览器分布图表
  if (browserChartRef.value && data.browserStats) {
    console.log('初始化浏览器分布图表')
    browserChart = echarts.init(browserChartRef.value)
    browserChart.setOption({
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: data.browserStats.map(item => item.browser),
        axisLabel: {
          interval: 0,
          rotate: 30
        }
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: '访问量',
          type: 'bar',
          data: data.browserStats.map(item => item.cnt),
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#83bff6' },
              { offset: 0.5, color: '#188df0' },
              { offset: 1, color: '#188df0' }
            ])
          }
        }
      ]
    })
  } else {
    console.error('浏览器分布图表初始化失败:', {
      hasRef: !!browserChartRef.value,
      hasData: !!data.browserStats
    })
  }

  // 操作系统分布图表
  if (osChartRef.value && data.osStats) {
    console.log('初始化操作系统分布图表')
    osChart = echarts.init(osChartRef.value)
    osChart.setOption({
      tooltip: {
        trigger: 'item',
        formatter: '{b}: {c} ({d}%)'
      },
      legend: {
        orient: 'vertical',
        right: 10,
        top: 'center',
        data: data.osStats.map(item => item.os)
      },
      series: [
        {
          name: '操作系统',
          type: 'pie',
          radius: '70%',
          center: ['40%', '50%'],
          data: data.osStats.map(item => ({
            name: item.os,
            value: item.cnt
          })),
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)'
            }
          }
        }
      ]
    })
  } else {
    console.error('操作系统分布图表初始化失败:', {
      hasRef: !!osChartRef.value,
      hasData: !!data.osStats
    })
  }

  // 网络类型分布图表
  if (networkChartRef.value && data.networkStats) {
    console.log('初始化网络类型分布图表')
    networkChart = echarts.init(networkChartRef.value)
    networkChart.setOption({
      tooltip: {},
      radar: {
        indicator: data.networkStats.map(item => ({
          name: item.network,
          max: Math.max(...data.networkStats.map(i => i.cnt)) * 1.2
        }))
      },
      series: [{
        name: '网络类型',
        type: 'radar',
        data: [{
          value: data.networkStats.map(item => item.cnt),
          name: '访问量',
          areaStyle: {
            color: 'rgba(24, 144, 255, 0.2)'
          }
        }]
      }]
    })
  } else {
    console.error('网络类型分布图表初始化失败:', {
      hasRef: !!networkChartRef.value,
      hasData: !!data.networkStats
    })
  }

  // 地区分布图表
  if (localeChartRef.value && data.localeCnStats) {
    console.log('初始化地区分布图表')
    localeChart = echarts.init(localeChartRef.value)
    localeChart.setOption({
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        }
      },
      grid: {
        top: '3%',
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'value'
      },
      yAxis: {
        type: 'category',
        data: data.localeCnStats.map(item => item.locale),
        axisLabel: {
          interval: 0
        }
      },
      series: [
        {
          name: '访问量',
          type: 'bar',
          data: data.localeCnStats.map(item => item.cnt),
          itemStyle: {
            color: new echarts.graphic.LinearGradient(1, 0, 0, 0, [
              { offset: 0, color: '#ffd85c' },
              { offset: 0.5, color: '#ff9a45' },
              { offset: 1, color: '#ff7343' }
            ])
          }
        }
      ]
    })
  } else {
    console.error('地区分布图表初始化失败:', {
      hasRef: !!localeChartRef.value,
      hasData: !!data.localeCnStats
    })
  }
}

// 监听窗口大小变化
const handleResize = () => {
  dailyTrendChart?.resize()
  deviceChart?.resize()
  browserChart?.resize()
  osChart?.resize()
  networkChart?.resize()
  localeChart?.resize()
}

// 组件挂载时添加窗口大小变化监听
onMounted(() => {
  if (dialogVisible.value) {
    fetchStatsData()
  }
  window.addEventListener('resize', handleResize)
})

// 组件卸载时移除窗口大小变化监听
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  // 销毁图表实例
  dailyTrendChart?.dispose()
  deviceChart?.dispose()
  browserChart?.dispose()
  osChart?.dispose()
  networkChart?.dispose()
  localeChart?.dispose()
})
</script>

<style scoped>
.stats-container {
  padding: 20px;
}

.query-card {
  margin-bottom: 20px;
}

.stats-cards {
  margin-bottom: 20px;
}

.stats-detail {
  margin-top: 20px;
}

.stats-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-value {
  font-size: 24px;
  font-weight: bold;
  text-align: center;
  color: #409EFF;
}

.chart {
  height: 300px;
  width: 100%;
}

:deep(.el-card__header) {
  padding: 10px 20px;
  border-bottom: 1px solid #EBEEF5;
  box-sizing: border-box;
}

:deep(.el-card__body) {
  padding: 20px;
}
</style> 