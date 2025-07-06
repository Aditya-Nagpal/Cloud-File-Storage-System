import { useAuthStore } from './auth';
import { useUserStore } from './user';
import { useFileStore } from './file';

export const resetAllStores = () => {
  useAuthStore().$reset();
  useFileStore().$reset();
  useUserStore().$reset();
};

export const resetLocalStorage = () => {
  localStorage.removeItem('accessToken');
  localStorage.removeItem('user');
};