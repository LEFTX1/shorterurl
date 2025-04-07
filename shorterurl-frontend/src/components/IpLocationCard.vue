<template>
  <div class="ip-location-card">
    <div v-if="loading" class="loading-overlay">
      <el-skeleton :rows="2" animated />
    </div>
    <div v-else-if="error" class="error-message">
      <el-alert
        title="获取IP位置信息失败"
        type="error"
        :description="error"
        show-icon
        :closable="false"
      />
    </div>
    <div v-else class="location-content">
      <div class="ip-info">
        <span class="label">IP地址:</span>
        <span class="value">{{ ipAddress }}</span>
      </div>
      <div class="location-info">
        <span class="label">位置:</span>
        <span class="value">{{ locationText }}</span>
      </div>
      <div v-if="adcode" class="adcode-info">
        <span class="label">区域代码:</span>
        <span class="value">{{ adcode }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import locationApi from '../api/location';
import type { IPLocationResponse } from '../api/location';

// 定义props
const props = defineProps({
  ipAddress: {
    type: String,
    required: true
  },
  autoLoad: {
    type: Boolean,
    default: true
  }
});

// 定义事件
const emit = defineEmits(['location-loaded', 'location-error']);

// 状态变量
const loading = ref(false);
const error = ref('');
const locationData = ref<IPLocationResponse | null>(null);

// 计算属性：位置文本
const locationText = computed(() => {
  if (!locationData.value) return '未知位置';
  
  const { province, city } = locationData.value;
  
  if (province === city) {
    return province || '未知位置';
  } else {
    return `${province} ${city}`.trim() || '未知位置';
  }
});

// 计算属性：区域代码
const adcode = computed(() => locationData.value?.adcode || '');

// 加载IP地理位置信息
const loadIpLocation = async () => {
  if (!props.ipAddress) return;
  
  loading.value = true;
  error.value = '';
  
  try {
    // 首先尝试使用后端API获取IP位置信息
    const response = await locationApi.getServerIpLocation(props.ipAddress);
    locationData.value = response.data;
    
    // 如果状态不为成功，尝试使用前端直接调用高德API
    if (response.data.status !== '1') {
      try {
        const amapResponse = await locationApi.getIpLocation(props.ipAddress);
        locationData.value = amapResponse.data;
        
        if (amapResponse.data.status !== '1') {
          error.value = amapResponse.data.info || '获取地理位置失败';
          emit('location-error', error.value);
        } else {
          emit('location-loaded', locationData.value);
        }
      } catch (amapErr: any) {
        error.value = `后端查询失败: ${response.data.info || '未知错误'}, 前端查询失败: ${amapErr.message || '网络请求失败'}`;
        emit('location-error', error.value);
      }
    } else {
      emit('location-loaded', locationData.value);
    }
  } catch (err: any) {
    // 后端API调用失败，尝试直接使用高德API
    try {
      const amapResponse = await locationApi.getIpLocation(props.ipAddress);
      locationData.value = amapResponse.data;
      
      if (amapResponse.data.status !== '1') {
        error.value = amapResponse.data.info || '获取地理位置失败';
        emit('location-error', error.value);
      } else {
        emit('location-loaded', locationData.value);
      }
    } catch (amapErr: any) {
      error.value = `后端查询失败: ${err.message || '网络请求失败'}, 前端查询失败: ${amapErr.message || '网络请求失败'}`;
      emit('location-error', error.value);
    }
  } finally {
    loading.value = false;
  }
};

// 监听IP地址变化重新加载
watch(() => props.ipAddress, (newIp) => {
  if (newIp) {
    loadIpLocation();
  }
});

// 组件挂载后自动加载（如果设置了autoLoad）
onMounted(() => {
  if (props.autoLoad && props.ipAddress) {
    loadIpLocation();
  }
});

// 暴露方法给父组件
defineExpose({
  loadIpLocation
});
</script>

<style scoped>
.ip-location-card {
  position: relative;
  padding: 15px;
  border-radius: 8px;
  background-color: #f8f9fa;
  margin-bottom: 15px;
  min-height: 80px;
}

.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(255, 255, 255, 0.7);
  padding: 15px;
  border-radius: 8px;
}

.error-message {
  padding: 5px 0;
}

.ip-info, .location-info, .adcode-info {
  margin-bottom: 8px;
}

.label {
  font-weight: bold;
  color: #606266;
  margin-right: 8px;
}

.value {
  color: #303133;
}
</style> 