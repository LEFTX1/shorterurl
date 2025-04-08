<template>
  <el-dialog
    v-model="dialogVisible"
    :title="'访问统计 - ' + (isGroupMode ? '分组统计' : linkInfo.fullShortUrl)"
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
          value-format="YYYY-MM-DD"
              :shortcuts="dateShortcuts"
              size="default"
            />
          </el-col>
          <el-col :span="4">
            <el-button type="primary" @click="fetchStatsData">查询统计</el-button>
            <el-button 
              v-if="!showLogsTable" 
              type="info" 
              @click="showLogsTable = true">
              查看访问日志
            </el-button>
            <el-button 
              v-else 
              type="info" 
              @click="showLogsTable = false">
              隐藏访问日志
            </el-button>
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
            <div class="card-value">{{ linkInfo.totalPv || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>
              <div class="card-header">
                <span>独立访客数(UV)</span>
        </div>
            </template>
            <div class="card-value">{{ linkInfo.totalUv || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <template #header>
              <div class="card-header">
                <span>独立IP数(UIP)</span>
        </div>
            </template>
            <div class="card-value">{{ linkInfo.totalUip || 0 }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 访问日志表格 (条件显示) -->
      <el-card v-if="showLogsTable" shadow="hover" class="logs-card">
        <template #header>
          <div class="card-header">
            <span>访问日志</span>
            <div>
              <el-pagination
                v-model:current-page="currentPage"
                v-model:page-size="pageSize"
                :total="totalLogs"
                :page-sizes="[10, 20, 50, 100]"
                layout="total, sizes, prev, pager, next, jumper"
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
              />
          </div>
          </div>
        </template>
        <el-table :data="accessLogs" border stripe height="300" v-loading="logsLoading">
          <el-table-column prop="ip" label="IP" width="140" />
          <el-table-column prop="browser" label="浏览器" width="100">
            <template #default="scope">
              <el-tag size="small">{{ scope.row.browser }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="os" label="系统" width="100">
            <template #default="scope">
              <el-tag size="small" type="success">{{ scope.row.os }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="device" label="设备" width="100">
            <template #default="scope">
              <el-tag size="small" type="warning">{{ scope.row.device }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="network" label="网络" width="100">
              <template #default="scope">
              <el-tag size="small" type="info">{{ scope.row.network }}</el-tag>
              </template>
            </el-table-column>
          <el-table-column prop="locale" label="地区" />
          <el-table-column prop="accessTime" label="访问时间" width="180" />
          </el-table>
      </el-card>

      <!-- 详细统计数据 - 两列布局 -->
      <el-row :gutter="20" class="stats-detail">
        <!-- 左侧列 -->
        <el-col :span="12">
          <el-card shadow="hover" class="stats-card">
            <template #header>
              <div class="card-header">
                <span>访问设备分布</span>
          </div>
            </template>
            <div ref="deviceChartRef" class="chart"></div>
          </el-card>

          <el-card shadow="hover" class="stats-card">
            <template #header>
              <div class="card-header">
                <span>浏览器分布</span>
          </div>
            </template>
            <div ref="browserChartRef" class="chart"></div>
          </el-card>

          <el-card shadow="hover" class="stats-card">
            <template #header>
              <div class="card-header">
                <span>操作系统分布</span>
            </div>
            </template>
            <div ref="osChartRef" class="chart"></div>
          </el-card>
        </el-col>

        <!-- 右侧列 -->
        <el-col :span="12">
          <el-card shadow="hover" class="stats-card">
            <template #header>
              <div class="card-header">
                <span>网络类型分布</span>
                </div>
            </template>
            <div ref="networkChartRef" class="chart"></div>
          </el-card>

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
import type { 
  DeviceStats, 
  BrowserStats, 
  OsStats, 
  NetworkStats, 
  LocaleStats, 
  AccessLog,
  ShortLinkStatsRespDTO
} from '../api/stats'
import type { PropType } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    required: true
  },
  linkInfo: {
    type: Object,
    required: true
  },
  statsData: {
    type: Object as PropType<ShortLinkStatsRespDTO | null>,
    default: null
  }
})

const emit = defineEmits(['update:visible'])
const router = useRouter()

const dialogVisible = ref(props.visible)
const statsData = ref<ShortLinkStatsRespDTO | null>(props.statsData)
const accessLogs = ref<AccessLog[]>([])

// 查询参数
const queryMode = ref('single')  // 'single' 或 'group'
const dateRange = ref([
  new Date().toISOString().split('T')[0],
  new Date().toISOString().split('T')[0]
])
const isGroupMode = ref(false)
const showLogsTable = ref(false)
const logsLoading = ref(false)

// 分页参数
const currentPage = ref(1)
const pageSize = ref(20)
const totalLogs = ref(0)

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
const deviceChartRef = ref<HTMLElement>()
const browserChartRef = ref<HTMLElement>()
const osChartRef = ref<HTMLElement>()
const networkChartRef = ref<HTMLElement>()
const localeChartRef = ref<HTMLElement>()

// 图表实例
let deviceChart: echarts.ECharts | null = null
let browserChart: echarts.ECharts | null = null
let osChart: echarts.ECharts | null = null
let networkChart: echarts.ECharts | null = null
let localeChart: echarts.ECharts | null = null

// 防抖定时器
let fetchDataTimer: ReturnType<typeof setTimeout> | null = null

// 处理查询模式变更
const handleModeChange = (mode: string) => {
  isGroupMode.value = mode === 'group'
  fetchStatsData()
}

// 处理分页大小变化
const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchAccessLogs()
}

// 处理页码变化
const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchAccessLogs()
}

// 监听对话框可见性
watch(() => props.visible, (val) => {
  dialogVisible.value = val
  if (val) {
    // 重置统计数据
    statsData.value = null
    
    // 设置日期范围：起始日期为链接创建日期，结束日期为今天
    if (!isGroupMode.value && props.linkInfo && props.linkInfo.createTime) {
      try {
        // 提取创建日期的纯日期部分 (YYYY-MM-DD)
        const createDate = props.linkInfo.createTime.split(' ')[0].trim()
        const endDate = new Date().toISOString().split('T')[0] // 今天
        console.log(`设置日期范围：${createDate} 至 ${endDate}`)
        
        // 设置日期范围但不立即触发查询
        dateRange.value = [createDate, endDate]
      } catch (error) {
        console.error('处理日期格式时出错:', error)
        // 出错时使用当天日期
        const today = new Date().toISOString().split('T')[0]
        dateRange.value = [today, today]
      }
    }
    
    // 这里不调用 fetchStatsData，等待 dateRange 的 watch 触发
  }
})

// 监听对话框关闭
watch(() => dialogVisible.value, (val) => {
  emit('update:visible', val)
})

// 监听查看日志表格状态
watch(() => showLogsTable.value, (val) => {
  if (val) {
    fetchAccessLogs()
  }
})

// 监听统计数据变化
watch(() => props.statsData, (newData) => {
  if (newData) {
    statsData.value = newData
    updateCharts(newData)
  } else {
    statsData.value = null
  }
}, { deep: true })

// 监听日期范围变化
watch([dateRange], () => {
  // 如果对话框可见，则重新获取数据
  if (dialogVisible.value) {
    // 防抖处理，避免短时间内多次请求
    if (fetchDataTimer) {
      clearTimeout(fetchDataTimer)
    }
    fetchDataTimer = setTimeout(() => {
      fetchStatsData()
    }, 200)
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
    const fullShortUrl = props.linkInfo.fullShortUrl
    const gid = props.linkInfo.gid
    const [startDate, endDate] = dateRange.value
    
    // 确保日期格式正确 (YYYY-MM-DD)
    const formatDate = (dateStr: string) => {
      // 如果日期包含时间部分，只保留日期部分
      if (dateStr.includes('T')) {
        return dateStr.split('T')[0]
      }
      return dateStr
    }
    
    const formattedStartDate = formatDate(startDate)
    const formattedEndDate = formatDate(endDate)
    
    console.log(`发送统计请求，参数: fullShortUrl=${fullShortUrl}, startDate=${formattedStartDate}, endDate=${formattedEndDate}`)

    // 根据查询模式选择不同的API
    let statsRes
    if (isGroupMode.value) {
      // 获取分组统计数据
      const response = await statsApi.getShortLinkGroupStats({
        gid,
        startDate: formattedStartDate,
        endDate: formattedEndDate
      })
      statsRes = response.data
    } else {
      // 获取单个链接统计数据
      const response = await statsApi.getShortLinkStats({
        fullShortUrl,
        gid,
        startDate: formattedStartDate,
        endDate: formattedEndDate
      })
      statsRes = response.data
    }

    // 更新数据
    if (statsRes) {
      statsData.value = statsRes
      // 重新初始化图表
      setTimeout(() => {
        initCharts()
      }, 100)
    }
  } catch (error) {
    console.error('获取统计数据失败', error)
    ElMessage.error('获取统计数据失败')
    statsData.value = null
  }
}

// 获取访问日志
const fetchAccessLogs = async () => {
  logsLoading.value = true
  try {
    const fullShortUrl = props.linkInfo.fullShortUrl
    const gid = props.linkInfo.gid
    const [startDate, endDate] = dateRange.value

    // 根据查询模式选择不同的API
    let logsRes
    if (isGroupMode.value) {
      // 获取分组访问日志
      const response = await statsApi.getShortLinkGroupAccessRecord({
        gid,
        startDate,
        endDate,
        current: currentPage.value,
        size: pageSize.value
      })
      logsRes = response.data
    } else {
      // 获取单个链接访问日志
      const response = await statsApi.getShortLinkAccessRecord({
        fullShortUrl,
        gid,
        startDate,
        endDate,
        current: currentPage.value,
        size: pageSize.value
      })
      logsRes = response.data
    }

    // 更新数据
    accessLogs.value = logsRes.records || []
    totalLogs.value = logsRes.total || 0
  } catch (error) {
    console.error('获取访问日志失败', error)
    ElMessage.error('获取访问日志失败')
  } finally {
    logsLoading.value = false
  }
}

// 初始化图表
const initCharts = () => {
  if (!statsData.value) return

  // 设备分布图表 - 环形图
  if (deviceChartRef.value && statsData.value?.deviceStats) {
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
        data: statsData.value.deviceStats.map(item => item.device)
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
          data: statsData.value.deviceStats.map(item => ({
            name: item.device,
            value: item.cnt
          }))
        }
      ]
    })
  }

  // 浏览器分布图表 - 柱状图
  if (browserChartRef.value && statsData.value?.browserStats) {
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
        data: statsData.value.browserStats.map(item => item.browser),
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
          data: statsData.value.browserStats.map(item => item.cnt),
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
  }

  // 操作系统分布图表 - 饼图
  if (osChartRef.value && statsData.value?.osStats) {
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
        data: statsData.value.osStats.map(item => item.os)
    },
    series: [
      {
        name: '操作系统',
        type: 'pie',
          radius: '70%',
          center: ['40%', '50%'],
          data: statsData.value.osStats.map(item => ({
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
  }

  // 网络类型分布图表 - 雷达图
  if (networkChartRef.value && statsData.value?.networkStats) {
    networkChart = echarts.init(networkChartRef.value)
    networkChart.setOption({
      tooltip: {},
      radar: {
        indicator: statsData.value.networkStats.map(item => ({
          name: item.network,
          max: Math.max(...statsData.value.networkStats.map(i => i.cnt)) * 1.2
        }))
      },
      series: [{
        name: '网络类型',
        type: 'radar',
        data: [{
          value: statsData.value.networkStats.map(item => item.cnt),
          name: '访问量',
          areaStyle: {
            color: 'rgba(24, 144, 255, 0.2)'
          }
        }]
      }]
    })
  }

  // 地区分布图表 - 条形图
  if (localeChartRef.value && statsData.value?.localeCnStats) {
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
        data: statsData.value.localeCnStats.map(item => item.locale),
        axisLabel: {
          interval: 0
        }
    },
    series: [
      {
        name: '访问量',
          type: 'bar',
          data: statsData.value.localeCnStats.map(item => item.cnt),
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
  }
}

// 更新图表数据
const updateCharts = (data: ShortLinkStatsRespDTO) => {
  if (!data) return
  
  // 更新设备分布图表
  if (deviceChart && data.deviceStats && Array.isArray(data.deviceStats)) {
    deviceChart.setOption({
      series: [{
        data: data.deviceStats.map(item => ({
          value: item.cnt,
          name: item.device
        }))
      }]
    })
  }

  // 更新浏览器分布图表
  if (browserChart && data.browserStats && Array.isArray(data.browserStats)) {
    browserChart.setOption({
      series: [{
        data: data.browserStats.map(item => ({
          value: item.cnt,
          name: item.browser
        }))
      }]
    })
  }

  // 更新操作系统分布图表
  if (osChart && data.osStats && Array.isArray(data.osStats)) {
    osChart.setOption({
      series: [{
        data: data.osStats.map(item => ({
          value: item.cnt,
          name: item.os
        }))
      }]
    })
  }

  // 更新网络类型分布图表
  if (networkChart && data.networkStats && Array.isArray(data.networkStats)) {
    networkChart.setOption({
      series: [{
        data: data.networkStats.map(item => ({
          value: item.cnt,
          name: item.network
        }))
      }]
    })
  }

  // 更新地区分布图表
  if (localeChart && data.localeCnStats && Array.isArray(data.localeCnStats)) {
    localeChart.setOption({
      series: [{
        data: data.localeCnStats.map(item => ({
          value: item.cnt,
          name: item.locale
        }))
      }]
    })
  }
}

// 监听窗口大小变化
const handleResize = () => {
  deviceChart?.resize()
  browserChart?.resize()
  osChart?.resize()
  networkChart?.resize()
  localeChart?.resize()
}

// 组件挂载时添加窗口大小变化监听
onMounted(() => {
  // 初始化时如果对话框可见，则获取数据
  if (dialogVisible.value) {
    fetchStatsData()
  }
  window.addEventListener('resize', handleResize)
})

// 组件卸载时移除窗口大小变化监听
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  // 销毁图表实例
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

.logs-card {
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