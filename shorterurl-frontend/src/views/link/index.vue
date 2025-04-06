<template>
  <div class="link-management">
    <div class="header">
      <div class="left-section">
        <div class="logo">
          <img src="@/assets/logo.svg" alt="Logo" />
          <span class="title">短链</span>
        </div>
        <el-dropdown class="project-select" trigger="click">
          <span class="project-name">
            {{ currentProject || '我的项目' }}
            <el-icon class="el-icon--right"><arrow-down /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>我的项目</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
      
      <div class="search-section">
        <el-input
          v-model="originUrl"
          placeholder="请输入 http:// 或 https:// 开头的链接或应用链接链接"
          clearable
          class="search-input"
        >
          <template #append>
            <el-button type="primary" @click="showCreateDialog">创建短链</el-button>
          </template>
        </el-input>
        <el-button class="create-button" @click="showBatchCreateDialog">
          <el-icon><plus /></el-icon>
          批量创建
        </el-button>
      </div>
      
      <div class="user-section">
        <el-dropdown trigger="click">
          <el-avatar :size="32" class="user-avatar">{{ userInitials }}</el-avatar>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="showUserInfoDialog">用户信息</el-dropdown-item>
              <el-dropdown-item @click="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
    
    <div class="main-content">
      <div class="sidebar">
        <div class="group-section">
          <div class="section-header">
            <span>短链分组</span>
            <el-button type="primary" size="small" circle @click="showCreateGroupDialog">
              <el-icon><plus /></el-icon>
            </el-button>
          </div>
          
          <div class="group-list">
            <div 
              v-for="group in groups" 
              :key="group.gid" 
              class="group-item"
              :class="{ active: currentGroup === group.gid }"
              @click="selectGroup(group.gid)"
            >
              <span class="group-name">{{ group.name }}</span>
              <span class="link-count">{{ group.shortLinkCount }}</span>
            </div>
          </div>
        </div>
        
        <div class="recycle-bin" @click="toggleRecycleBin">
          <el-icon><delete /></el-icon>
          <span>回收站</span>
        </div>
      </div>
      
      <div class="content">
        <div class="content-header">
          <div class="breadcrumb">
            <el-breadcrumb separator="/">
              <el-breadcrumb-item :to="{ path: '/link' }">项目首页</el-breadcrumb-item>
              <el-breadcrumb-item>{{ isRecycleBin ? '回收站' : '短链管理' }}</el-breadcrumb-item>
            </el-breadcrumb>
          </div>
          
          <div class="header-actions" v-if="!isRecycleBin">
            <el-button type="primary" @click="showCreateDialog">创建短链</el-button>
          </div>
        </div>
        
        <div class="table-container">
          <div class="table-header">
            <div class="left">
              <span class="header-title">{{ isRecycleBin ? '回收站' : '短链列表' }}</span>
            </div>
            <div class="right">
              <el-input
                v-model="searchKeyword"
                placeholder="搜索短链/原始链接"
                class="search-input"
                clearable
                @input="handleSearch"
              >
                <template #prefix>
                  <el-icon><search /></el-icon>
                </template>
              </el-input>
            </div>
          </div>
          
          <link-table
            :links="linkList"
            :loading="loading"
            :is-recycle-bin="isRecycleBin"
            :total="total"
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            @edit="handleEdit"
            @stats="handleStats"
            @recycle-bin="handleRecycleBin"
            @refresh="fetchData"
            date-format="YYYY-MM-DD HH:mm"
            :col-widths="{ shortUrl: '200px' }"
          />
          
          <div class="empty-data" v-if="linkList.length === 0 && !loading">
            <el-empty description="暂无数据" />
          </div>
        </div>
      </div>
    </div>
    
    <!-- 创建短链对话框 -->
    <el-dialog
      v-model="createDialogVisible"
      title="创建短链"
      width="500px"
    >
      <el-form 
        ref="createFormRef"
        :model="createForm"
        :rules="createRules"
        label-width="100px"
      >
        <el-form-item label="原始链接" prop="originUrl">
          <el-input v-model="createForm.originUrl" placeholder="请输入原始链接" @blur="getUrlTitle" />
        </el-form-item>
        <el-form-item label="分组" prop="gid">
          <el-select v-model="createForm.gid" placeholder="请选择分组">
            <el-option
              v-for="group in groups"
              :key="group.gid"
              :label="group.name"
              :value="group.gid"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="有效期" prop="validDateType">
          <el-radio-group v-model="createForm.validDateType">
            <el-radio :label="0">永久有效</el-radio>
            <el-radio :label="1">自定义</el-radio>
          </el-radio-group>
          <el-date-picker
            v-if="createForm.validDateType === 1"
            v-model="createForm.validDate"
            type="date"
            placeholder="选择有效期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 100%; margin-top: 10px;"
          />
        </el-form-item>
        <el-form-item label="短链接标题" prop="describe">
          <el-input v-model="createForm.describe" placeholder="短链接标题将自动获取，也可手动修改" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="createDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="createShortLink" :loading="submitting">创建</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 编辑短链对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      title="编辑短链"
      width="500px"
    >
      <el-form 
        ref="editFormRef"
        :model="editForm"
        :rules="createRules"
        label-width="100px"
      >
        <el-form-item label="短链接">
          <el-input v-model="editForm.fullShortUrl" disabled />
        </el-form-item>
        <el-form-item label="原始链接" prop="originUrl">
          <el-input v-model="editForm.originUrl" placeholder="请输入原始链接" @blur="getEditUrlTitle" />
        </el-form-item>
        <el-form-item label="分组" prop="gid">
          <el-select v-model="editForm.gid" placeholder="请选择分组">
            <el-option
              v-for="group in groups"
              :key="group.gid"
              :label="group.name"
              :value="group.gid"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="有效期" prop="validDateType">
          <el-radio-group v-model="editForm.validDateType">
            <el-radio :label="0">永久有效</el-radio>
            <el-radio :label="1">自定义</el-radio>
          </el-radio-group>
          <el-date-picker
            v-if="editForm.validDateType === 1"
            v-model="editForm.validDate"
            type="date"
            placeholder="选择有效期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 100%; margin-top: 10px;"
          />
        </el-form-item>
        <el-form-item label="短链接标题" prop="describe">
          <el-input v-model="editForm.describe" placeholder="短链接标题将自动获取，也可手动修改" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="updateShortLink" :loading="submitting">保存</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 创建分组对话框 -->
    <el-dialog
      v-model="createGroupDialogVisible"
      title="创建分组"
      width="400px"
    >
      <el-form 
        ref="groupFormRef"
        :model="groupForm"
        :rules="groupRules"
        label-width="80px"
      >
        <el-form-item label="分组名称" prop="name">
          <el-input v-model="groupForm.name" placeholder="请输入分组名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="createGroupDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="createGroup" :loading="submitting">创建</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 批量创建对话框 -->
    <el-dialog
      v-model="batchCreateDialogVisible"
      title="批量创建短链"
      width="600px"
    >
      <el-form 
        ref="batchCreateFormRef"
        :model="batchCreateForm"
        :rules="batchCreateRules"
        label-width="100px"
      >
        <el-form-item label="原始链接" prop="originUrls">
          <el-input 
            v-model="batchCreateForm.originUrls" 
            type="textarea" 
            rows="5"
            placeholder="每行输入一个原始链接"
          />
        </el-form-item>
        <el-form-item label="分组" prop="gid">
          <el-select v-model="batchCreateForm.gid" placeholder="请选择分组">
            <el-option
              v-for="group in groups"
              :key="group.gid"
              :label="group.name"
              :value="group.gid"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="有效期" prop="validDateType">
          <el-radio-group v-model="batchCreateForm.validDateType">
            <el-radio :label="0">永久有效</el-radio>
            <el-radio :label="1">自定义</el-radio>
          </el-radio-group>
          <el-date-picker
            v-if="batchCreateForm.validDateType === 1"
            v-model="batchCreateForm.validDate"
            type="date"
            placeholder="选择有效期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 100%; margin-top: 10px;"
          />
        </el-form-item>
        <el-form-item label="批量标题" prop="describe">
          <el-input 
            v-model="batchCreateForm.describe" 
            placeholder="为所有链接设置相同的标题，也可留空后自动获取"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="batchCreateDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="batchCreateShortLink" :loading="submitting">创建</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 统计数据对话框 -->
    <stats-dialog
      v-model:visible="statsDialogVisible"
      :link-info="currentLinkInfo"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { 
  Plus, 
  Search, 
  Delete, 
  ArrowDown, 
  Document, 
  Edit, 
  DataLine,
} from '@element-plus/icons-vue';
import { useUserStore } from '../../store/user';
import LinkTable from '@/components/LinkTable.vue';
import StatsDialog from '@/components/StatsDialog.vue';

// 引入真实的API
import groupApi from '../../api/group';
import linkApi from '../../api/link';
import recycleApi from '../../api/recycle';

// 添加接口类型导入
import type { RecycleBinPageReq } from '../../api/recycle';

// 定义必要的类型接口
interface ShortLinkRecord {
  id: number;
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

interface ShortLinkGroupResp {
  gid: string;
  name: string;
  sortOrder: number;
  shortLinkCount: number;
}

// 路由实例
const router = useRouter();

// 用户存储
const userStore = useUserStore();

// 数据加载状态
const loading = ref(false);
const submitting = ref(false);

// 用户信息
const userInitials = computed(() => {
  const realname = userStore.getRealname;
  return realname && realname.length > 0 ? realname.substring(0, 1) : 'U';
});

// 项目信息
const currentProject = ref('我的项目');

// 分组管理
const groups = ref<ShortLinkGroupResp[]>([]);
const currentGroup = ref('');

// 短链列表
const linkList = ref<ShortLinkRecord[]>([]);
const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);
const searchKeyword = ref('');
const originUrl = ref('');
const isRecycleBin = ref(false);

// 分组表单
const createGroupDialogVisible = ref(false);
const groupForm = ref({
  name: ''
});
const groupRules = {
  name: [
    { required: true, message: '请输入分组名称', trigger: 'blur' },
    { min: 1, max: 20, message: '长度在 1 到 20 个字符', trigger: 'blur' }
  ]
};

// 创建短链表单
const createDialogVisible = ref(false);
const createForm = ref({
  originUrl: '',
  gid: '',
  validDateType: 0,
  validDate: '',
  describe: ''
});

// 批量创建短链表单
const batchCreateDialogVisible = ref(false);
const batchCreateForm = ref({
  originUrls: '',
  gid: '',
  validDateType: 0,
  validDate: '',
  describe: ''
});

// 编辑短链表单
const editDialogVisible = ref(false);
const editForm = ref({
  fullShortUrl: '',
  originUrl: '',
  gid: '',
  validDateType: 0,
  validDate: '',
  describe: ''
});

// 统计数据对话框
const statsDialogVisible = ref(false);
const currentLinkInfo = ref<ShortLinkRecord>({} as ShortLinkRecord);

// 表单验证规则
const createRules = {
  originUrl: [
    { required: true, message: '请输入原始链接', trigger: 'blur' },
    { pattern: /^https?:\/\//, message: '链接必须以 http:// 或 https:// 开头', trigger: 'blur' }
  ],
  gid: [
    { required: true, message: '请选择分组', trigger: 'change' }
  ]
};

const batchCreateRules = {
  originUrls: [
    { required: true, message: '请输入原始链接', trigger: 'blur' }
  ],
  gid: [
    { required: true, message: '请选择分组', trigger: 'change' }
  ]
};

// 获取分组列表
const fetchGroups = async () => {
  try {
    loading.value = true;
    const res = await groupApi.listGroups();
    groups.value = res.data || [];
    
    // 如果没有选择分组但有分组数据，默认选择第一个（通常是默认分组）
    if (groups.value.length > 0 && !currentGroup.value) {
      currentGroup.value = groups.value[0].gid;
      // 获取该分组下的短链接
      fetchData();
    } else if (groups.value.length === 0) {
      // 如果没有分组，可能是系统问题，创建一个默认分组
      try {
        await groupApi.createGroup({ name: '默认分组' });
        // 创建成功后重新获取分组列表
        const newRes = await groupApi.listGroups();
        groups.value = newRes.data || [];
        
        if (groups.value.length > 0) {
          currentGroup.value = groups.value[0].gid;
          // 获取该分组下的短链接
          fetchData();
        }
      } catch (error) {
        console.error('创建默认分组失败', error);
      }
    }
  } catch (error) {
    console.error('获取分组列表失败', error);
  } finally {
    loading.value = false;
  }
};

// 获取短链列表
const fetchData = async () => {
  if (!currentGroup.value && !isRecycleBin.value) return;
  
  loading.value = true;
  try {
    if (isRecycleBin.value) {
      // 查询回收站
      const params: RecycleBinPageReq = {
        current: currentPage.value,
        size: pageSize.value
      };
      
      // 仅当选择了分组时才添加gidList参数
      if (currentGroup.value) {
        params.gidList = [currentGroup.value];
      }
      
      const res = await recycleApi.pageRecycleBin(params);
      linkList.value = res.data.records || [];
      total.value = res.data.total;
    } else {
      // 查询短链
      const res = await linkApi.pageShortLink({
        gid: currentGroup.value,
        current: currentPage.value,
        size: pageSize.value
      });
      linkList.value = res.data.records || [];
      total.value = res.data.total;
    }
  } catch (error) {
    console.error('获取短链列表失败', error);
  } finally {
    loading.value = false;
  }
};

// 选择分组
const selectGroup = (gid: string) => {
  currentGroup.value = gid;
  isRecycleBin.value = false;
  currentPage.value = 1;
  fetchData();
};

// 切换回收站
const toggleRecycleBin = () => {
  isRecycleBin.value = true;
  currentPage.value = 1;
  fetchData();
};

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1;
  fetchData();
};

// 获取URL标题
const getUrlTitle = async () => {
  if (!createForm.value.originUrl) return;
  
  try {
    const response = await linkApi.getUrlTitle(createForm.value.originUrl);
    if (response && response.data && response.data.data) {
      createForm.value.describe = response.data.data;
      ElMessage.success('成功获取网站标题');
    }
  } catch (error) {
    console.error('获取URL标题失败:', error);
    ElMessage.warning('获取URL标题失败，请手动输入');
  }
};

// 获取URL标题（编辑模式）
const getEditUrlTitle = async () => {
  if (!editForm.value.originUrl) return;
  
  try {
    const response = await linkApi.getUrlTitle(editForm.value.originUrl);
    if (response && response.data && response.data.data) {
      editForm.value.describe = response.data.data;
      ElMessage.success('成功获取网站标题');
    }
  } catch (error) {
    console.error('获取URL标题失败:', error);
    ElMessage.warning('获取URL标题失败，请手动输入');
  }
};

// 显示创建短链对话框
const showCreateDialog = () => {
  createForm.value = {
    originUrl: originUrl.value,
    gid: currentGroup.value,
    validDateType: 0,
    validDate: '',
    describe: ''
  };
  createDialogVisible.value = true;
  
  // 如果已有URL，自动获取标题
  if (originUrl.value) {
    getUrlTitle();
  }
};

// 显示批量创建对话框
const showBatchCreateDialog = () => {
  batchCreateForm.value = {
    originUrls: '',
    gid: currentGroup.value,
    validDateType: 0,
    validDate: '',
    describe: ''
  };
  batchCreateDialogVisible.value = true;
};

// 显示创建分组对话框
const showCreateGroupDialog = () => {
  groupForm.value = { name: '' };
  createGroupDialogVisible.value = true;
};

// 显示用户信息对话框
const showUserInfoDialog = () => {
  // 实现用户信息对话框
};

// 处理编辑短链
const handleEdit = (row: ShortLinkRecord) => {
  editForm.value = {
    fullShortUrl: row.fullShortUrl,
    originUrl: row.originUrl,
    gid: row.gid,
    validDateType: row.validDate ? 1 : 0,
    validDate: row.validDate,
    describe: row.describe
  };
  editDialogVisible.value = true;
};

// 处理查看统计
const handleStats = (row: ShortLinkRecord) => {
  currentLinkInfo.value = row;
  statsDialogVisible.value = true;
};

// 处理回收站操作
const handleRecycleBin = (row: ShortLinkRecord) => {
  if (isRecycleBin.value) {
    // 从回收站恢复
    ElMessageBox.confirm('确定要恢复该短链接吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      try {
        await recycleApi.recoverFromRecycleBin({
          gid: row.gid,
          fullShortUrl: row.fullShortUrl
        });
        ElMessage.success('恢复成功');
        fetchData();
      } catch (error) {
        console.error('恢复失败', error);
      }
    }).catch(() => {});
  } else {
    // 移动到回收站
    ElMessageBox.confirm('确定要删除该短链接吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      try {
        await recycleApi.saveToRecycleBin({
          gid: row.gid,
          fullShortUrl: row.fullShortUrl
        });
        ElMessage.success('删除成功');
        fetchData();
      } catch (error) {
        console.error('删除失败', error);
      }
    }).catch(() => {});
  }
};

// 创建短链
const createShortLink = async () => {
  try {
    submitting.value = true;
    const res = await linkApi.createShortLink(createForm.value);
    ElMessage.success('创建成功');
    createDialogVisible.value = false;
    originUrl.value = '';
    fetchData();
  } catch (error) {
    console.error('创建短链失败', error);
    ElMessage.error('创建短链失败');
  } finally {
    submitting.value = false;
  }
};

// 批量创建短链接
const batchCreateShortLink = async () => {
  try {
    submitting.value = true;
    
    // 分割原始URL
    const urls = batchCreateForm.value.originUrls
      .split('\n')
      .map(url => url.trim())
      .filter(url => url);
    
    if (urls.length === 0) {
      ElMessage.warning('请输入至少一个有效的URL');
      submitting.value = false;
      return;
    }
    
    // 准备API请求参数
    const requestData = {
      originUrls: batchCreateForm.value.originUrls,
      gid: batchCreateForm.value.gid,
      validDateType: batchCreateForm.value.validDateType,
      describe: batchCreateForm.value.describe
    } as any; // 使用类型断言
    
    if (batchCreateForm.value.validDateType === 1 && batchCreateForm.value.validDate) {
      requestData.validDate = batchCreateForm.value.validDate;
    }
    
    const response = await linkApi.batchCreateShortLink(requestData);
    
    ElMessage.success(`成功创建${response.data.total || 0}个短链接`);
    batchCreateDialogVisible.value = false;
    fetchData(); // 刷新数据
  } catch (error) {
    console.error('批量创建短链接失败', error);
    ElMessage.error('批量创建短链接失败');
  } finally {
    submitting.value = false;
  }
};

// 更新短链
const updateShortLink = async () => {
  try {
    submitting.value = true;
    await linkApi.updateShortLink(editForm.value);
    ElMessage.success('更新成功');
    editDialogVisible.value = false;
    fetchData();
  } catch (error) {
    console.error('更新短链失败', error);
  } finally {
    submitting.value = false;
  }
};

// 创建分组
const createGroup = async () => {
  try {
    submitting.value = true;
    await groupApi.createGroup({ name: groupForm.value.name });
    ElMessage.success('创建成功');
    createGroupDialogVisible.value = false;
    await fetchGroups();
  } catch (error) {
    console.error('创建分组失败', error);
  } finally {
    submitting.value = false;
  }
};

// 退出登录
const logout = () => {
  ElMessageBox.confirm('确定要退出登录吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    userStore.logout();
  }).catch(() => {});
};

// 在组件的onMounted钩子中添加初始化逻辑
onMounted(async () => {
  // 检查用户是否已登录
  if (userStore.isLogin) {
    // 首先获取分组列表
    await fetchGroups();
  }
});

// 监听分组变化，加载对应的短链列表
watch(currentGroup, () => {
  if (currentGroup.value) {
    fetchData();
  }
});
</script>

<style scoped>
.link-management {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f5f7fa;
}

.header {
  height: 60px;
  background-color: #fff;
  border-bottom: 1px solid #ebeef5;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.left-section {
  display: flex;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  margin-right: 20px;
}

.logo img {
  width: 32px;
  height: 32px;
  margin-right: 8px;
}

.logo .title {
  font-size: 18px;
  font-weight: bold;
  color: #333;
}

.project-select {
  cursor: pointer;
}

.project-name {
  display: flex;
  align-items: center;
  gap: 5px;
}

.search-section {
  flex: 1;
  max-width: 600px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.search-input {
  width: 100%;
}

.user-section {
  margin-left: 20px;
}

.user-avatar {
  cursor: pointer;
  background-color: #409eff;
  color: #fff;
}

.main-content {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.sidebar {
  width: 240px;
  background-color: #fff;
  border-right: 1px solid #ebeef5;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.group-section {
  flex: 1;
  padding: 20px 0;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  margin-bottom: 15px;
  font-weight: bold;
}

.group-list {
  max-height: calc(100vh - 200px);
  overflow-y: auto;
}

.group-item {
  padding: 10px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  transition: background-color 0.3s;
}

.group-item:hover {
  background-color: #f5f7fa;
}

.group-item.active {
  background-color: #ecf5ff;
  color: #409eff;
}

.link-count {
  background-color: #f0f0f0;
  border-radius: 10px;
  padding: 2px 8px;
  font-size: 12px;
  color: #666;
}

.recycle-bin {
  padding: 15px 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  border-top: 1px solid #ebeef5;
  transition: background-color 0.3s;
}

.recycle-bin:hover {
  background-color: #f5f7fa;
}

.content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.table-container {
  background-color: #fff;
  border-radius: 4px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.05);
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-title {
  font-size: 16px;
  font-weight: bold;
}

.table-header .search-input {
  width: 250px;
}

.empty-data {
  padding: 30px 0;
}

/* 优化表格中短链接列宽度 */
:deep(.short-url-cell) {
  max-width: 200px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style> 