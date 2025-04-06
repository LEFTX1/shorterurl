<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="$emit('update:visible', $event)"
    title="数据统计"
    width="80%"
    :before-close="handleClose"
    class="stats-dialog"
  >
    <div class="stats-header">
      <div class="link-info">
        <div class="url">{{ linkInfo.fullShortUrl }}</div>
        <div class="origin-url">{{ linkInfo.originUrl }}</div>
      </div>
      <div class="date-picker">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
          @change="handleDateChange"
        />
      </div>
    </div>

    <div v-loading="loading" class="stats-content">
      <div class="stats-summary">
        <div class="stats-card">
          <div class="label">总访问量(PV)</div>
          <div class="value">{{ statsData.overallPvUvUipStats?.pv || 0 }}</div>
        </div>
        <div class="stats-card">
          <div class="label">独立访客(UV)</div>
          <div class="value">{{ statsData.overallPvUvUipStats?.uv || 0 }}</div>
        </div>
        <div class="stats-card">
          <div class="label">IP数量(UIP)</div>
          <div class="value">{{ statsData.overallPvUvUipStats?.uip || 0 }}</div>
        </div>
      </div>

      <el-divider />

      <el-tabs v-model="activeTab" class="stats-tabs">
        <el-tab-pane label="访问趋势" name="trend">
          <div class="chart-container">
            <div id="trend-chart" class="chart"></div>
          </div>
        </el-tab-pane>
        <el-tab-pane label="地域分布" name="area">
          <div class="chart-container">
            <div id="area-chart" class="chart"></div>
            <div id="map-chart" class="chart"></div>
          </div>
        </el-tab-pane>
        <el-tab-pane label="设备分析" name="device">
          <div class="chart-container multiple-charts">
            <div id="browser-chart" class="chart"></div>
            <div id="os-chart" class="chart"></div>
            <div id="device-chart" class="chart"></div>
          </div>
        </el-tab-pane>
        <el-tab-pane label="访问详情" name="details">
          <div class="chart-container multiple-charts">
            <div id="hour-chart" class="chart"></div>
            <div id="weekday-chart" class="chart"></div>
          </div>
        </el-tab-pane>
        <el-tab-pane label="访问记录" name="records">
          <el-table :data="accessRecords" border style="width: 100%">
            <el-table-column prop="ip" label="IP地址" width="140" />
            <el-table-column prop="locale" label="地区" width="120" />
            <el-table-column prop="browser" label="浏览器" width="120" />
            <el-table-column prop="os" label="操作系统" width="120" />
            <el-table-column prop="device" label="设备" width="100" />
            <el-table-column prop="network" label="网络" width="100" />
            <el-table-column prop="accessTime" label="访问时间" min-width="180" />
          </el-table>
          <div class="records-pagination">
            <el-pagination
              v-model:current-page="recordsCurrentPage"
              v-model:page-size="recordsPageSize"
              :page-sizes="[10, 20, 50, 100]"
              layout="total, sizes, prev, pager, next, jumper"
              :total="recordsTotal"
              @size-change="handleRecordsSizeChange"
              @current-change="handleRecordsCurrentChange"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, nextTick } from 'vue';
import * as echarts from 'echarts';
import { stats as statsApi } from '@/api';
import type { ShortLinkRecord } from '@/api/link';
import type { ShortLinkStatsRespDTO, AccessRecord } from '@/api/stats';

// 定义必要的类型
interface ShortLinkRecord {
  fullShortUrl: string;
  originUrl: string;
  domain: string;
  gid: string;
  createTime: string;
  validDate: string;
  describe: string;
  validDateType?: number;
  totalPv: number;
  totalUv: number;
  totalUip: number;
}

// PV/UV/UIP统计
interface PvUvUipStats {
  date: string;
  pv: number;
  uv: number;
  uip: number;
}

// 地域统计
interface LocaleCnStat {
  locale: string;
  cnt: number;
  ratio: number;
}

// 浏览器统计
interface BrowserStat {
  browser: string;
  cnt: number;
  ratio: number;
}

// 操作系统统计
interface OsStat {
  os: string;
  cnt: number;
  ratio: number;
}

// 设备统计
interface DeviceStat {
  device: string;
  cnt: number;
  ratio: number;
}

// 访问记录
interface AccessRecord {
  ip: string;
  browser: string;
  os: string;
  network: string;
  device: string;
  locale: string;
  accessTime: string;
}

// 统计响应
interface ShortLinkStatsRespDTO {
  pvUvUipStatsList: PvUvUipStats[];
  overallPvUvUipStats: PvUvUipStats;
  localeCnStats: LocaleCnStat[];
  hourStats: number[];
  topIpStats: any[];
  weekdayStats: number[];
  browserStats: BrowserStat[];
  osStats: OsStat[];
  uvTypeStats: any[];
  deviceStats: DeviceStat[];
  networkStats: any[];
}

// Mock API for demo purposes
const statsApi = {
  getShortLinkStats: async (params: any) => {
    // In a real app, this would be an actual API call
    console.log('Fetching stats for:', params);
    return {
      data: {
        pvUvUipStatsList: [],
        overallPvUvUipStats: { pv: 0, uv: 0, uip: 0, date: '' },
        localeCnStats: [],
        hourStats: Array(24).fill(0),
        weekdayStats: Array(7).fill(0),
        browserStats: [],
        osStats: [],
        deviceStats: [],
        topIpStats: [],
        uvTypeStats: [],
        networkStats: []
      } as ShortLinkStatsRespDTO
    };
  },
  getShortLinkAccessRecord: async (params: any) => {
    // In a real app, this would be an actual API call
    console.log('Fetching access records for:', params);
    return {
      data: {
        records: [],
        total: 0,
        size: 10,
        current: 1
      }
    };
  }
};

const props = defineProps<{
  visible: boolean;
  linkInfo: ShortLinkRecord;
}>();

const emit = defineEmits(['update:visible']);

// 状态
const dateRange = ref<[string, string]>(['', '']);
const loading = ref(false);
const activeTab = ref('trend');
const statsData = ref<ShortLinkStatsRespDTO>({} as ShortLinkStatsRespDTO);
const accessRecords = ref<AccessRecord[]>([]);
const recordsTotal = ref(0);
const recordsCurrentPage = ref(1);
const recordsPageSize = ref(10);

// 图表实例
let trendChart: echarts.ECharts | null = null;
let areaChart: echarts.ECharts | null = null;
let mapChart: echarts.ECharts | null = null;
let browserChart: echarts.ECharts | null = null;
let osChart: echarts.ECharts | null = null;
let deviceChart: echarts.ECharts | null = null;
let hourChart: echarts.ECharts | null = null;
let weekdayChart: echarts.ECharts | null = null;

// 关闭对话框
const handleClose = () => {
  emit('update:visible', false);
};

// 初始化日期范围
const initDateRange = () => {
  const end = new Date();
  const start = new Date();
  start.setDate(start.getDate() - 7);
  
  const formatDate = (date: Date) => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  };
  
  dateRange.value = [formatDate(start), formatDate(end)];
};

// 日期变化处理
const handleDateChange = () => {
  fetchStatsData();
};

// 获取统计数据
const fetchStatsData = async () => {
  if (!props.linkInfo || !props.linkInfo.fullShortUrl || !dateRange.value[0] || !dateRange.value[1]) return;
  
  loading.value = true;
  try {
    const res = await statsApi.getShortLinkStats({
      fullShortUrl: props.linkInfo.fullShortUrl,
      gid: props.linkInfo.gid,
      startDate: dateRange.value[0],
      endDate: dateRange.value[1],
      enableStatus: 0
    });
    
    statsData.value = res.data;
    
    nextTick(() => {
      initCharts();
    });
  } catch (error) {
    console.error('获取统计数据失败', error);
  } finally {
    loading.value = false;
  }
};

// 获取访问记录
const fetchAccessRecords = async () => {
  if (!props.linkInfo || !props.linkInfo.fullShortUrl || !dateRange.value[0] || !dateRange.value[1]) return;
  
  loading.value = true;
  try {
    const res = await statsApi.getShortLinkAccessRecord({
      fullShortUrl: props.linkInfo.fullShortUrl,
      gid: props.linkInfo.gid,
      startDate: dateRange.value[0],
      endDate: dateRange.value[1],
      enableStatus: 0,
      current: recordsCurrentPage.value,
      size: recordsPageSize.value
    });
    
    accessRecords.value = res.data.records;
    recordsTotal.value = res.data.total;
  } catch (error) {
    console.error('获取访问记录失败', error);
  } finally {
    loading.value = false;
  }
};

// 初始化图表
const initCharts = () => {
  initTrendChart();
  initAreaChart();
  initMapChart();
  initBrowserChart();
  initOsChart();
  initDeviceChart();
  initHourChart();
  initWeekdayChart();
};

// 初始化趋势图表
const initTrendChart = () => {
  const chartDom = document.getElementById('trend-chart');
  if (!chartDom) return;
  
  trendChart = echarts.init(chartDom);
  
  const option = {
    title: {
      text: '访问趋势'
    },
    tooltip: {
      trigger: 'axis'
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
      boundaryGap: false,
      data: statsData.value.pvUvUipStatsList?.map((item: PvUvUipStats) => item.date) || []
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: 'PV',
        type: 'line',
        data: statsData.value.pvUvUipStatsList?.map((item: PvUvUipStats) => item.pv) || []
      },
      {
        name: 'UV',
        type: 'line',
        data: statsData.value.pvUvUipStatsList?.map((item: PvUvUipStats) => item.uv) || []
      },
      {
        name: 'UIP',
        type: 'line',
        data: statsData.value.pvUvUipStatsList?.map((item: PvUvUipStats) => item.uip) || []
      }
    ]
  };
  
  trendChart.setOption(option);
};

// 初始化地域饼图
const initAreaChart = () => {
  const chartDom = document.getElementById('area-chart');
  if (!chartDom) return;
  
  areaChart = echarts.init(chartDom);
  
  const option = {
    title: {
      text: '地域分布'
    },
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 10,
      data: statsData.value.localeCnStats?.map((item: LocaleCnStat) => item.locale) || []
    },
    series: [
      {
        name: '访问来源',
        type: 'pie',
        radius: ['50%', '70%'],
        avoidLabelOverlap: false,
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 16,
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: statsData.value.localeCnStats?.map((item: LocaleCnStat) => ({
          value: item.cnt,
          name: item.locale
        })) || []
      }
    ]
  };
  
  areaChart.setOption(option);
};

// 初始化中国地图
const initMapChart = () => {
  const chartDom = document.getElementById('map-chart');
  if (!chartDom) return;
  
  mapChart = echarts.init(chartDom);
  
  // 地图配置
  // 注意：实际项目中需要加载中国地图数据
  const option = {
    title: {
      text: '访问地图'
    },
    tooltip: {
      trigger: 'item'
    },
    visualMap: {
      min: 0,
      max: statsData.value.localeCnStats?.[0]?.cnt || 10,
      left: 'left',
      top: 'bottom',
      text: ['高', '低'],
      calculable: true
    },
    series: [
      {
        name: '访问量',
        type: 'map',
        map: 'china',
        roam: true,
        label: {
          show: true
        },
        data: statsData.value.localeCnStats?.map((item: LocaleCnStat) => ({
          name: item.locale,
          value: item.cnt
        })) || []
      }
    ]
  };
  
  mapChart.setOption(option);
};

// 初始化浏览器图表
const initBrowserChart = () => {
  const chartDom = document.getElementById('browser-chart');
  if (!chartDom) return;
  
  browserChart = echarts.init(chartDom);
  
  const option = {
    title: {
      text: '浏览器分布'
    },
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'horizontal',
      bottom: 0,
      data: statsData.value.browserStats?.map((item: BrowserStat) => item.browser) || []
    },
    series: [
      {
        name: '浏览器',
        type: 'pie',
        radius: '60%',
        data: statsData.value.browserStats?.map((item: BrowserStat) => ({
          value: item.cnt,
          name: item.browser
        })) || []
      }
    ]
  };
  
  browserChart.setOption(option);
};

// 初始化操作系统图表
const initOsChart = () => {
  const chartDom = document.getElementById('os-chart');
  if (!chartDom) return;
  
  osChart = echarts.init(chartDom);
  
  const option = {
    title: {
      text: '操作系统分布'
    },
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'horizontal',
      bottom: 0,
      data: statsData.value.osStats?.map((item: OsStat) => item.os) || []
    },
    series: [
      {
        name: '操作系统',
        type: 'pie',
        radius: '60%',
        data: statsData.value.osStats?.map((item: OsStat) => ({
          value: item.cnt,
          name: item.os
        })) || []
      }
    ]
  };
  
  osChart.setOption(option);
};

// 初始化设备图表
const initDeviceChart = () => {
  const chartDom = document.getElementById('device-chart');
  if (!chartDom) return;
  
  deviceChart = echarts.init(chartDom);
  
  const option = {
    title: {
      text: '设备分布'
    },
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'horizontal',
      bottom: 0,
      data: statsData.value.deviceStats?.map((item: DeviceStat) => item.device) || []
    },
    series: [
      {
        name: '设备',
        type: 'pie',
        radius: '60%',
        data: statsData.value.deviceStats?.map((item: DeviceStat) => ({
          value: item.cnt,
          name: item.device
        })) || []
      }
    ]
  };
  
  deviceChart.setOption(option);
};

// 初始化小时分布图表
const initHourChart = () => {
  const chartDom = document.getElementById('hour-chart');
  if (!chartDom) return;
  
  hourChart = echarts.init(chartDom);
  
  const hours = Array.from({length: 24}, (_, i) => `${i}时`);
  
  const option = {
    title: {
      text: '时段分布'
    },
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: hours
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        data: statsData.value.hourStats || new Array(24).fill(0),
        type: 'bar'
      }
    ]
  };
  
  hourChart.setOption(option);
};

// 初始化星期分布图表
const initWeekdayChart = () => {
  const chartDom = document.getElementById('weekday-chart');
  if (!chartDom) return;
  
  weekdayChart = echarts.init(chartDom);
  
  const weekdays = ['星期一', '星期二', '星期三', '星期四', '星期五', '星期六', '星期日'];
  
  const option = {
    title: {
      text: '星期分布'
    },
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: weekdays
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        data: statsData.value.weekdayStats || new Array(7).fill(0),
        type: 'bar'
      }
    ]
  };
  
  weekdayChart.setOption(option);
};

// 访问记录分页处理
const handleRecordsSizeChange = (val: number) => {
  recordsPageSize.value = val;
  fetchAccessRecords();
};

const handleRecordsCurrentChange = (val: number) => {
  recordsCurrentPage.value = val;
  fetchAccessRecords();
};

// 监听tab切换
watch(() => activeTab.value, (newVal) => {
  if (newVal === 'records') {
    fetchAccessRecords();
  } else {
    nextTick(() => {
      // 重绘图表
      trendChart?.resize();
      areaChart?.resize();
      mapChart?.resize();
      browserChart?.resize();
      osChart?.resize();
      deviceChart?.resize();
      hourChart?.resize();
      weekdayChart?.resize();
    });
  }
});

// 监听对话框可见性
watch(() => props.visible, (newVal) => {
  if (newVal) {
    initDateRange();
    fetchStatsData();
  }
});

// 组件卸载前销毁图表实例
onMounted(() => {
  window.addEventListener('resize', () => {
    trendChart?.resize();
    areaChart?.resize();
    mapChart?.resize();
    browserChart?.resize();
    osChart?.resize();
    deviceChart?.resize();
    hourChart?.resize();
    weekdayChart?.resize();
  });
});
</script>

<style scoped>
.stats-dialog {
  font-family: Arial, sans-serif;
}

.stats-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.link-info {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.link-info .url {
  font-weight: bold;
  font-size: 16px;
  color: #409eff;
}

.link-info .origin-url {
  font-size: 14px;
  color: #666;
}

.stats-summary {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.stats-card {
  background-color: #f5f7fa;
  padding: 20px;
  border-radius: 4px;
  flex: 1;
  text-align: center;
  margin: 0 10px;
}

.stats-card:first-child {
  margin-left: 0;
}

.stats-card:last-child {
  margin-right: 0;
}

.stats-card .label {
  font-size: 14px;
  color: #666;
  margin-bottom: 5px;
}

.stats-card .value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.chart-container {
  width: 100%;
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
}

.multiple-charts .chart {
  width: 32%;
  height: 350px;
  margin-bottom: 20px;
}

.chart-container:not(.multiple-charts) .chart {
  width: 48%;
  height: 400px;
  margin-bottom: 20px;
}

.records-pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style> 