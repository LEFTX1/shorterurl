import { createPinia } from 'pinia';
import { useUserStore } from './user';

// 创建pinia实例
const pinia = createPinia();

export {
  useUserStore
};

export default pinia; 