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
            <el-button type="info" size="small" @click="showSortGroupsDialog">
              <el-icon><sort /></el-icon>
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
        
        <div 
          class="recycle-bin" 
          :class="{ active: isRecycleBin }" 
          @click="toggleRecycleBin"
        >
          <el-icon><delete /></el-icon>
          <span>回收站</span>
        </div>
      </div>
      
      <div class="content">
        <div class="content-header">
          <div class="breadcrumb">
            <el-breadcrumb separator="/">
              <el-breadcrumb-item :to="{ path: '/link' }">项目首页</el-breadcrumb-item>
              <el-breadcrumb-item>{{ getCurrentGroupName() }}</el-breadcrumb-item>
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
              <span class="header-title">{{ getCurrentGroupName() }}</span>
              <div class="mode-tag-container">
                <transition name="mode-fade" mode="out-in">
                  <el-tag v-if="isRecycleBin" key="recycle" type="danger" effect="dark" class="mode-tag">回收站模式</el-tag>
                  <el-tag v-else key="normal" type="success" effect="plain" class="mode-tag">常规模式</el-tag>
                </transition>
              </div>
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
          
          <!-- 回收站状态说明 -->
          <div v-if="isRecycleBin" class="recycle-bin-tip">
            <el-alert
              :title="`回收站模式 - 分组: ${getCurrentGroupName()}`"
              type="info"
              description="此视图显示状态为 enable_status=1(未启用/回收站) 且 del_flag=0(未永久删除) 的短链接，可以选择恢复或永久删除"
              show-icon
              :closable="false"
              style="margin-bottom: 15px;"
            />
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
            @recycle-bin="handleDelete"
            @recover="handleRecover"
            @remove="handleRemove"
            @refresh="fetchData"
            date-format="YYYY-MM-DD HH:mm"
            :col-widths="{ shortUrl: '200px' }"
          />
          
          <div class="empty-data" v-if="linkList.length === 0 && !loading">
            <el-empty :description="isRecycleBin ? 
              `分组 '${getCurrentGroupName()}' 的回收站中暂无短链接` : 
              `分组 '${getCurrentGroupName()}' 中暂无短链接`" />
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
        <el-form-item label="原始链接" prop="originUrl">
          <el-input v-model="editForm.originUrl" placeholder="请输入原始链接" disabled />
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
          <el-input v-model="editForm.describe" placeholder="短链接标题" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitEdit" :loading="submitting">保存</el-button>
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
    
    <!-- 新增：分组排序对话框 -->
    <el-dialog
      v-model="sortGroupsDialogVisible"
      title="分组排序"
      width="500px"
    >
      <el-alert
        type="info"
        title="提示"
        description="拖拽分组调整排序位置，点击保存按钮提交更改"
        :closable="false"
        show-icon
        style="margin-bottom: 20px;"
      />
      
      <draggable
        v-model="sortableGroups"
        ghost-class="ghost"
        handle=".drag-handle"
        item-key="gid"
        :animation="150"
      >
        <template #item="{element}">
          <div class="sort-group-item">
            <el-icon class="drag-handle"><operation /></el-icon>
            <span>{{ element.name }}</span>
            <span class="sort-order">#{{ element.sortOrder }}</span>
          </div>
        </template>
      </draggable>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="sortGroupsDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveGroupsSort" :loading="submitting">保存</el-button>
        </span>
      </template>
    </el-dialog>
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
  Operation,
  Sort
} from '@element-plus/icons-vue';
import { useUserStore } from '../../store/user';
import LinkTable from '@/components/LinkTable.vue';
import StatsDialog from '@/components/StatsDialog.vue';
import draggable from 'vuedraggable';

// 引入真实的API
import groupApi from '../../api/group';
import linkApi from '../../api/link';
import recycleApi from '../../api/recycle';

// 添加接口类型导入
import type { RecycleBinPageReq } from '../../api/recycle';
import type { ShortLinkGroupResp, SortGroup } from '../../api/group';

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
const sortGroupsDialogVisible = ref(false);
const sortableGroups = ref<ShortLinkGroupResp[]>([]);

// 短链列表
const linkList = ref<ShortLinkRecord[]>([]);
const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);
const searchKeyword = ref('');
const originUrl = ref('');

// 使用计算属性从userStore获取isRecycleBin状态
const isRecycleBin = computed({
  get: () => userStore.isRecycleBinMode,
  set: (value) => userStore.setViewMode(value ? 'recycle' : 'normal')
});

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
  originGid: '',
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
    // 按照 sortOrder 升序排列分组
    groups.value = (res.data || []).sort((a, b) => a.sortOrder - b.sortOrder);
    
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
        groups.value = (newRes.data || []).sort((a, b) => a.sortOrder - b.sortOrder);
        
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
  loading.value = true;
  linkList.value = []; // 先清空当前列表
  total.value = 0;
  
  try {
    if (isRecycleBin.value) {
      // 查询回收站 - 获取enable_status=1且del_flag=0的链接
      const params: RecycleBinPageReq = {
        current: currentPage.value,
        size: pageSize.value
      };
      
      // 必须提供gid参数，即使未选择分组也要使用默认分组ID
      if (currentGroup.value) {
        params.gid = currentGroup.value;
      } else if (groups.value.length > 0) {
        // 如果没有选中的分组但有分组数据，使用第一个分组的ID
        params.gid = groups.value[0].gid;
        currentGroup.value = params.gid; // 更新当前选中的分组
        console.log('未选择分组，自动使用默认分组:', params.gid);
      } else {
        console.error('没有可用的分组，无法查询回收站数据');
        ElMessage.error('没有可用的分组，请先创建分组');
        loading.value = false;
        return;
      }
      
      console.log(`查询分组 "${getCurrentGroupName()}" 的回收站数据 (enable_status=1/未启用/回收站, del_flag=0/未永久删除)，参数:`, params);
      try {
        const res = await recycleApi.pageRecycleBin(params);
        console.log('回收站查询结果:', res.data);
        
        if (res.data) {
          linkList.value = res.data.records || [];
          total.value = res.data.total || 0;
        } else {
          ElMessage.warning('获取回收站数据失败');
        }
      } catch (recycleError) {
        console.error('回收站查询失败:', recycleError);
        ElMessage.error('回收站查询失败，请检查网络连接或联系管理员');
        linkList.value = [];
        total.value = 0;
      }
    } else {
      // 查询正常短链 - 获取enable_status=0且del_flag=0的链接
      const params = {
        gid: currentGroup.value,
        current: currentPage.value,
        size: pageSize.value
      };
      
      console.log(`查询分组 "${getCurrentGroupName()}" 的正常短链数据 (enable_status=0/启用/正常, del_flag=0/未永久删除)，参数:`, params);
      const res = await linkApi.pageShortLink(params);
      console.log('正常短链查询结果:', res.data);
      
      if (res.data) {
        linkList.value = res.data.records || [];
        total.value = res.data.total || 0;
      } else {
        ElMessage.warning('获取短链列表失败');
      }
    }
  } catch (error) {
    console.error('获取数据失败', error);
    ElMessage.error('获取数据失败，请稍后重试');
  } finally {
    loading.value = false;
  }
};

// 选择分组
const selectGroup = (gid: string) => {
  if (gid === currentGroup.value) {
    return; // 如果点击的是当前已选中的分组，不做任何操作
  }
  
  currentGroup.value = gid;
  currentPage.value = 1; // 切换分组时重置页码
  
  // 在当前模式下切换分组（回收站或常规模式）
  const modeText = isRecycleBin.value ? '回收站' : '常规';
  ElMessage.info({
    message: `已切换到分组 "${getCurrentGroupName()}" 的${modeText}视图`,
    duration: 2000
  });
  
  fetchData();
};

// 切换到回收站
const toggleRecycleBin = () => {
  // 更新store中的状态，通过计算属性自动同步到本地
  isRecycleBin.value = !isRecycleBin.value;
  userStore.setViewMode(isRecycleBin.value ? 'recycle' : 'normal');
  
  currentPage.value = 1; // 重置页码
  
  // 确保有选中的分组
  if (!currentGroup.value && groups.value.length > 0) {
    currentGroup.value = groups.value[0].gid;
    console.log('切换回收站模式时自动选择默认分组:', currentGroup.value);
  }
  
  if (isRecycleBin.value) {
    ElMessage.info({
      message: `已切换到回收站模式，显示分组 "${getCurrentGroupName()}" 中enable_status=1(未启用/回收站)且del_flag=0(未永久删除)的链接`,
      duration: 3000
    });
  } else {
    ElMessage.info({
      message: `已退出回收站模式，显示分组 "${getCurrentGroupName()}" 中enable_status=0(启用/正常)且del_flag=0(未永久删除)的链接`,
      duration: 3000
    });
  }
  
  // 延迟一点执行，给用户视觉反馈
  setTimeout(() => {
    fetchData();
  }, 300);
};

// 显示分组排序对话框
const showSortGroupsDialog = () => {
  // 直接使用当前排序好的分组数据，不需要再次排序
  sortableGroups.value = [...groups.value];
  sortGroupsDialogVisible.value = true;
};

// 保存分组排序
const saveGroupsSort = async () => {
  submitting.value = true;
  try {
    // 根据新的顺序更新sortOrder
    const sortGroups: SortGroup[] = sortableGroups.value.map((group, index) => ({
      gid: group.gid,
      sortOrder: index + 1 // 顺序从1开始
    }));
    
    const response = await groupApi.sortGroups({ groups: sortGroups });
    
    if (response.data && response.data.success) {
      ElMessage.success('分组排序保存成功');
      sortGroupsDialogVisible.value = false;
      
      // 重新获取分组列表以更新顺序
      await fetchGroups();
    } else {
      ElMessage.error('分组排序保存失败');
    }
  } catch (error) {
    console.error('保存分组排序失败', error);
    ElMessage.error('保存分组排序失败，请稍后重试');
  } finally {
    submitting.value = false;
  }
};

// 获取当前选中分组的名称
const getCurrentGroupName = () => {
  if (!currentGroup.value) return '所有分组';
  const group = groups.value.find(g => g.gid === currentGroup.value);
  return group ? group.name : '未知分组';
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
const handleEdit = (link: ShortLinkRecord) => {
  // 设置编辑表单数据
  editForm.value = {
    fullShortUrl: link.fullShortUrl,
    originUrl: link.originUrl,
    originGid: link.gid,
    gid: link.gid,
    validDateType: link.validDateType || 0,
    validDate: link.validDate,
    describe: link.describe
  };
  
  // 显示编辑对话框
  editDialogVisible.value = true;
};

// 处理查看统计
const handleStats = (row: ShortLinkRecord) => {
  currentLinkInfo.value = row;
  statsDialogVisible.value = true;
};

// 删除短链（保存到回收站）
const handleDelete = (link: ShortLinkRecord) => {
  ElMessageBox.confirm(`确定要将短链接 ${link.fullShortUrl} 移动到回收站吗?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const data = {
        gid: link.gid,
        fullShortUrl: link.fullShortUrl
      };
      
      console.log('正在将短链接移动到回收站 (设置enable_status=1/未启用/回收站):', data);
      
      // 调用保存到回收站API
      const res = await recycleApi.saveToRecycleBin(data);
      
      console.log('移动到回收站响应:', res.data);
      
      if (res.data && res.data.success) {
        ElMessage.success('已成功移动到回收站');
        // 手动延迟一下再刷新数据，确保后端数据已更新
        setTimeout(() => {
          fetchData(); // 刷新数据
        }, 300);
      } else if (res.data && res.data.code === "0") {
        // 如果返回了特定code，可能表示短链接已在回收站
        ElMessage.info('该短链接已在回收站中');
        fetchData(); // 刷新数据
      } else {
        ElMessage.error('移动到回收站失败: ' + (res.data?.message || '未知错误'));
      }
    } catch (error: any) {
      console.error('移动到回收站失败', error);
      // 显示详细错误信息
      const errorMsg = error.response?.data?.message || error.message || '请稍后重试';
      ElMessage.error('移动到回收站失败: ' + errorMsg);
    }
  }).catch(() => {
    // 用户取消操作
  });
};

// 从回收站恢复短链
const handleRecover = (link: ShortLinkRecord) => {
  ElMessageBox.confirm(`确定要将短链接 ${link.fullShortUrl} 从回收站恢复吗?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'info'
  }).then(async () => {
    try {
      const data = {
        gid: link.gid,
        fullShortUrl: link.fullShortUrl
      };
      
      console.log('正在恢复短链接 (设置enable_status=0/启用/正常):', data);
      
      // 调用从回收站恢复API
      const res = await recycleApi.recoverFromRecycleBin(data);
      
      console.log('从回收站恢复响应:', res.data);
      
      if (res.data && res.data.success) {
        ElMessage.success('短链接已恢复');
        // 手动延迟一下再刷新数据，确保后端数据已更新
        setTimeout(() => {
          fetchData(); // 刷新数据
        }, 300);
      } else {
        ElMessage.error('恢复失败: ' + (res.data?.message || '未知错误'));
      }
    } catch (error: any) {
      console.error('恢复短链接失败', error);
      // 显示详细错误信息
      const errorMsg = error.response?.data?.message || error.message || '请稍后重试';
      ElMessage.error('恢复短链接失败: ' + errorMsg);
    }
  }).catch(() => {
    // 用户取消操作
  });
};

// 永久删除短链
const handleRemove = (link: ShortLinkRecord) => {
  ElMessageBox.confirm(`确定要永久删除短链接 ${link.fullShortUrl} 吗? 此操作将设置del_flag=1(永久删除)，删除后将无法恢复!`, '警告', {
    confirmButtonText: '确定删除',
    cancelButtonText: '取消',
    type: 'error'
  }).then(async () => {
    try {
      const data = {
        gid: link.gid,
        fullShortUrl: link.fullShortUrl
      };
      
      console.log('正在永久删除短链接 (设置del_flag=1/永久删除):', data);
      
      // 调用从回收站删除API
      const res = await recycleApi.removeFromRecycleBin(data);
      
      console.log('永久删除短链接响应:', res.data);
      
      if (res.data && res.data.success) {
        ElMessage.success('短链接已永久删除');
        // 手动延迟一下再刷新数据，确保后端数据已更新
        setTimeout(() => {
          fetchData(); // 刷新数据
        }, 300);
      } else {
        ElMessage.error('删除失败: ' + (res.data?.message || '未知错误'));
      }
    } catch (error: any) {
      console.error('删除短链接失败', error);
      // 显示详细错误信息
      const errorMsg = error.response?.data?.message || error.message || '请稍后重试';
      ElMessage.error('删除短链接失败: ' + errorMsg);
    }
  }).catch(() => {
    // 用户取消操作
  });
};

// 提交编辑
const submitEdit = async () => {
  if (!editForm.value.fullShortUrl || !editForm.value.originUrl || !editForm.value.gid) {
    ElMessage.error('请填写完整信息');
    return;
  }
  
  submitting.value = true;
  
  try {
    const data = {
      fullShortUrl: editForm.value.fullShortUrl,
      originUrl: editForm.value.originUrl,
      originGid: editForm.value.originGid,
      gid: editForm.value.gid,
      validDateType: editForm.value.validDateType,
      validDate: editForm.value.validDateType === 1 ? editForm.value.validDate : undefined,
      describe: editForm.value.describe
    };
    
    console.log('提交编辑短链接数据:', data);
    
    // 调用更新短链API
    const res = await linkApi.updateShortLink(data);
    
    console.log('编辑短链接响应:', res.data);
    
    if (res.data && res.data.success) {
      ElMessage.success('短链接已更新');
      editDialogVisible.value = false;
      // 手动延迟一下再刷新数据，确保后端数据已更新
      setTimeout(() => {
        fetchData(); // 刷新数据
      }, 300);
    } else {
      ElMessage.error('更新失败');
    }
  } catch (error) {
    console.error('更新短链接失败', error);
    ElMessage.error('更新短链接失败，请稍后重试');
  } finally {
    submitting.value = false;
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
  font-weight: bold;
  border-left: 3px solid #409eff;
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

.recycle-bin.active {
  background-color: #fef0f0;
  color: #f56c6c;
  font-weight: bold;
  border-left: 3px solid #f56c6c;
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
  margin-right: 10px;
}

.mode-tag-container {
  display: inline-block;
  min-width: 90px;
  height: 22px;
  margin-left: 10px;
}

.table-header .left {
  display: flex;
  align-items: center;
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

/* 模式标签切换动画 */
.mode-fade-enter-active,
.mode-fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.mode-fade-enter-from,
.mode-fade-leave-to {
  opacity: 0;
  transform: translateY(5px);
}

.sort-group-item {
  display: flex;
  align-items: center;
  padding: 10px;
  margin-bottom: 8px;
  background-color: #f5f7fa;
  border-radius: 4px;
  cursor: move;
}

.sort-group-item .drag-handle {
  cursor: move;
  margin-right: 10px;
  color: #909399;
}

.sort-group-item .sort-order {
  margin-left: auto;
  color: #909399;
}

.ghost {
  opacity: 0.5;
  background: #c8ebfb;
}

/* 修复drag-handle样式 */
.drag-handle {
  cursor: move;
}
</style> 