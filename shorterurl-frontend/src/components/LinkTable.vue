<template>
  <div class="link-table">
    <el-table :data="links" style="width: 100%" v-loading="loading" border>
      <el-table-column prop="fullShortUrl" label="短链网址" :min-width="getColWidth('shortUrl') || 180">
        <template #default="scope">
          <div class="link-info">
            <el-tooltip content="点击复制链接" placement="top">
              <div class="url-wrapper short-url-cell" @click="copyLink(scope.row.fullShortUrl)">
                {{ scope.row.fullShortUrl }}
                <el-icon><Document /></el-icon>
              </div>
            </el-tooltip>
            <div class="origin-url" :title="scope.row.originUrl">
              {{ scope.row.originUrl }}
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="describe" label="短链接标题" min-width="150">
        <template #default="scope">
          <el-tooltip :content="scope.row.describe" placement="top" :disabled="!scope.row.describe">
            <span class="link-title">{{ scope.row.describe || '无标题' }}</span>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column prop="createTime" label="创建时间" width="160">
        <template #default="scope">
          {{ formatDate(scope.row.createTime) }}
        </template>
      </el-table-column>
      <el-table-column prop="validDate" label="有效期" width="160">
        <template #default="scope">
          {{ scope.row.validDateType === 0 ? '永久有效' : scope.row.validDate }}
        </template>
      </el-table-column>
      <el-table-column label="访问数据" width="220">
        <template #default="scope">
          <div class="visit-data">
            <div class="data-item">
              <span class="label">总访问：</span>
              <span class="value">{{ scope.row.totalPv || 0 }}</span>
            </div>
            <div class="data-item">
              <span class="label">独立访客：</span>
              <span class="value">{{ scope.row.totalUv || 0 }}</span>
            </div>
            <div class="data-item">
              <span class="label">IP数：</span>
              <span class="value">{{ scope.row.totalUip || 0 }}</span>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="scope">
          <div class="operate-buttons">
            <el-tooltip content="编辑" placement="top">
              <el-button type="primary" :icon="Edit" circle size="small" @click="handleEdit(scope.row)" />
            </el-tooltip>
            <el-tooltip content="统计数据" placement="top">
              <el-button type="success" :icon="DataLine" circle size="small" @click="handleStats(scope.row)" />
            </el-tooltip>
            <el-tooltip :content="isRecycleBin ? '恢复' : '删除'" placement="top">
              <el-button :type="isRecycleBin ? 'success' : 'danger'" :icon="isRecycleBin ? 'Refresh' : 'Delete'" 
                circle size="small" @click="handleRecycleBin(scope.row)" />
            </el-tooltip>
          </div>
        </template>
      </el-table-column>
    </el-table>
    
    <div class="pagination-container">
      <el-pagination
        :current-page="currentPage"
        :page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        @update:current-page="$emit('update:currentPage', $event)"
        @update:page-size="$emit('update:pageSize', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { Edit, Document, DataLine, Delete, Refresh } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';

// 短链接记录接口
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

const props = defineProps({
  links: {
    type: Array as () => ShortLinkRecord[],
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  isRecycleBin: {
    type: Boolean,
    default: false
  },
  total: {
    type: Number,
    default: 0
  },
  currentPage: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 10
  },
  dateFormat: {
    type: String,
    default: 'YYYY-MM-DD HH:mm:ss'
  },
  colWidths: {
    type: Object,
    default: () => ({})
  }
});

const emit = defineEmits([
  'edit', 
  'stats', 
  'recycleBin',
  'update:currentPage',
  'update:pageSize',
  'refresh'
]);

// 复制链接到剪贴板
const copyLink = (link: string) => {
  navigator.clipboard.writeText(link).then(() => {
    ElMessage.success('链接已复制到剪贴板');
  }).catch(() => {
    ElMessage.error('复制失败，请手动复制');
  });
};

// 编辑短链接
const handleEdit = (row: ShortLinkRecord) => {
  emit('edit', row);
};

// 查看统计数据
const handleStats = (row: ShortLinkRecord) => {
  emit('stats', row);
};

// 回收站操作
const handleRecycleBin = (row: ShortLinkRecord) => {
  emit('recycleBin', row);
};

// 分页大小变化
const handleSizeChange = (val: number) => {
  emit('update:pageSize', val);
  emit('refresh');
};

// 页码变化
const handleCurrentChange = (val: number) => {
  emit('update:currentPage', val);
  emit('refresh');
};

// 日期格式化
const formatDate = (dateString: string) => {
  if (!dateString) return '-';
  const date = new Date(dateString);
  
  // 简单格式化，根据dateFormat属性配置
  // 这里实现一个简单版本，实际可以使用dayjs或其他库
  const formatMap: { [key: string]: string } = {
    'YYYY': date.getFullYear().toString(),
    'MM': (date.getMonth() + 1).toString().padStart(2, '0'),
    'DD': date.getDate().toString().padStart(2, '0'),
    'HH': date.getHours().toString().padStart(2, '0'),
    'mm': date.getMinutes().toString().padStart(2, '0'),
    'ss': date.getSeconds().toString().padStart(2, '0')
  };
  
  let result = props.dateFormat;
  Object.entries(formatMap).forEach(([pattern, value]) => {
    result = result.replace(pattern, value);
  });
  
  return result;
};

// 获取列宽度
const getColWidth = (colName: string) => {
  if (props.colWidths && typeof props.colWidths === 'object') {
    try {
      // 处理字符串形式的JSON
      const widthsObj = typeof props.colWidths === 'string' 
        ? JSON.parse(props.colWidths) 
        : props.colWidths;
        
      return widthsObj[colName];
    } catch (e) {
      console.error('解析colWidths失败', e);
      return null;
    }
  }
  return null;
};
</script>

<style scoped>
.link-table {
  width: 100%;
}

.link-info {
  display: flex;
  flex-direction: column;
}

.url-wrapper {
  display: flex;
  align-items: center;
  gap: 5px;
  font-weight: bold;
  color: #409eff;
  cursor: pointer;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.link-title {
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.origin-url {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.visit-data {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.data-item {
  font-size: 12px;
  display: flex;
}

.data-item .label {
  color: #666;
  min-width: 70px;
}

.data-item .value {
  font-weight: bold;
  color: #333;
}

.operate-buttons {
  display: flex;
  gap: 5px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style> 