import { useAuthStore } from './auth';
import { useUserStore } from './user';
import { useFileStore } from './file';

export const resetAllStores = () => {
  useAuthStore().$reset();
  useUserStore().$reset();
  useFileStore().$reset();
};

export const resetLocalStorage = () => {
  localStorage.removeItem('accessToken');
  localStorage.removeItem('user');
};