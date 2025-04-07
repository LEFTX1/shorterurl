<template>
  <div class="ip-location-viewer">
    <div ref="mapContainer" class="map-container"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, onUnmounted } from 'vue';
import type { IPLocationResponse } from '../api/location';

// 导入高德地图，使用动态导入避免全局污染
// 注意：这里使用了Vue3的特性，动态导入脚本
const AMap = ref<any>(null);

// 定义props
const props = defineProps({
  locationData: {
    type: Object as () => IPLocationResponse | null,
    required: true
  }
});

// 地图容器ref
const mapContainer = ref<HTMLDivElement | null>(null);
// 地图实例
const map = ref<any>(null);

// 加载高德地图
const loadAMap = async () => {
  if (window.AMap) {
    AMap.value = window.AMap;
    return;
  }

  return new Promise((resolve, reject) => {
    const script = document.createElement('script');
    script.type = 'text/javascript';
    script.async = true;
    script.src = `https://webapi.amap.com/maps?v=2.0&key=9891e494403818e3fc79fb61fcf06b84&callback=initAMap`;
    
    // 定义回调函数
    window.initAMap = () => {
      AMap.value = window.AMap;
      resolve(window.AMap);
    };
    
    script.onerror = reject;
    document.head.appendChild(script);
  });
};

// 初始化地图
const initMap = async () => {
  if (!mapContainer.value) return;
  
  try {
    // 确保AMap已加载
    if (!AMap.value) {
      await loadAMap();
    }
    
    // 创建地图实例
    map.value = new AMap.value.Map(mapContainer.value, {
      zoom: 12, // 初始缩放级别
      resizeEnable: true,
      viewMode: '3D' // 3D模式
    });
    
    // 如果有位置数据，则显示
    updateMapLocation();
  } catch (error) {
    console.error('初始化地图失败：', error);
  }
};

// 更新地图位置
const updateMapLocation = () => {
  if (!map.value || !props.locationData) return;
  
  try {
    // 清除所有标记
    map.value.clearMap();
    
    const { rectangle } = props.locationData;
    
    if (rectangle) {
      // 解析矩形区域坐标
      const rectangleCoords = rectangle.split(';')
        .map(point => point.split(',').map(Number));
      
      const southWest = rectangleCoords[0];
      const northEast = rectangleCoords[1];
      
      // 创建矩形范围
      const rectangleBounds = new AMap.value.Bounds(
        new AMap.value.LngLat(southWest[0], southWest[1]),
        new AMap.value.LngLat(northEast[0], northEast[1])
      );
      
      // 调整地图视野以包含矩形区域
      map.value.setBounds(rectangleBounds);
      
      // 添加矩形覆盖物
      const rectangleOverlay = new AMap.value.Rectangle({
        bounds: rectangleBounds,
        strokeColor: '#409EFF',
        strokeWeight: 2,
        strokeOpacity: 0.5,
        fillColor: '#409EFF',
        fillOpacity: 0.2,
        zIndex: 50
      });
      
      rectangleOverlay.setMap(map.value);
      
      // 计算中心点
      const center = [
        (southWest[0] + northEast[0]) / 2,
        (southWest[1] + northEast[1]) / 2
      ];
      
      // 添加标记
      const marker = new AMap.value.Marker({
        position: new AMap.value.LngLat(center[0], center[1]),
        title: `${props.locationData.province} ${props.locationData.city}`,
        icon: new AMap.value.Icon({
          size: new AMap.value.Size(40, 40),
          image: 'https://webapi.amap.com/theme/v1.3/markers/n/mark_r.png',
          imageSize: new AMap.value.Size(40, 40)
        }),
        offset: new AMap.value.Pixel(-20, -40)
      });
      
      map.value.add(marker);
      
      // 添加信息窗体
      const infoWindow = new AMap.value.InfoWindow({
        content: `
          <div class="info-window">
            <div class="info-title">${props.locationData.province} ${props.locationData.city}</div>
            <div class="info-content">
              <div>区域编码: ${props.locationData.adcode}</div>
            </div>
          </div>
        `,
        offset: new AMap.value.Pixel(0, -30)
      });
      
      marker.on('mouseover', () => {
        infoWindow.open(map.value, marker.getPosition());
      });
      
      // 自动打开信息窗体
      infoWindow.open(map.value, marker.getPosition());
    }
  } catch (error) {
    console.error('更新地图位置失败：', error);
  }
};

// 监听位置数据变化
watch(() => props.locationData, (newVal) => {
  if (newVal && map.value) {
    updateMapLocation();
  }
}, { deep: true });

// 组件挂载后初始化地图
onMounted(() => {
  initMap();
});

// 组件卸载前清理地图
onUnmounted(() => {
  if (map.value) {
    map.value.destroy();
    map.value = null;
  }
});
</script>

<style scoped>
.ip-location-viewer {
  width: 100%;
  height: 100%;
}

.map-container {
  width: 100%;
  height: 400px;
  border-radius: 4px;
  overflow: hidden;
}

/* 信息窗体样式 */
:deep(.amap-info-content) {
  padding: 10px;
}

:deep(.info-window) {
  padding: 6px 0;
}

:deep(.info-title) {
  font-weight: bold;
  font-size: 14px;
  margin-bottom: 8px;
  color: #333;
}

:deep(.info-content) {
  font-size: 12px;
  color: #666;
}
</style>

<script lang="ts">
// 为TypeScript增加高德地图类型声明
declare global {
  interface Window {
    AMap: any;
    initAMap: () => void;
  }
}
</script> 