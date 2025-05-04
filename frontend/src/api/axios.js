import axios from 'axios';
import { useAuthStore } from '../store/auth';

const BASE_URL = import.meta.env.VITE_API_BASE_URL;
const REFRESH_API = BASE_URL+'/auth/refresh';

const API = axios.create({
  baseURL: BASE_URL,
  withCredentials: true // utmost important for refresh token
});

// Flag to avoid multiple refresh attempts
let isRefreshing = false;
let failedQueue = [];

const processQueue = (error, token = null) => {
    failedQueue.forEach(prom => {
        if (error) {
            prom.reject(error);
        } else {
            prom.resolve(token);
        }
    });
    failedQueue = [];
};

// Request interceptor - attact access token
API.interceptors.request.use((config) => {
    const authStore = useAuthStore();
    const token = authStore.accessToken;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor - handle 401 by try refresh
API.interceptors.response.use(
    response => response,
    async (error) => {
        const authStore = useAuthStore();
        const originalRequest = error.config;

        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true;
            if(isRefreshing) {
                return new Promise(function(resolve, reject) {
                    failedQueue.push({resolve, reject});
                }).then(token => {
                    originalRequest.headers.Authorization = `Bearer ${token}`;
                    return API(originalRequest);
                }).catch(err => {
                    return Promise.reject(err);
                });
            }
            isRefreshing = true;
            try {
                const refreshResponse = await axios.get(REFRESH_API, {
                    withCredentials: true
                });

                const newAccessToken = refreshResponse.data.accessToken;
                authStore.setAccessToken(newAccessToken); // update pinia + localstorage

                // Retry the original request with the new token
                originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
                processQueue(null, newAccessToken);

                return API(originalRequest);
            } catch (refreshError) {
                processQueue(refreshError, null);
                // await authStore.logout();
                throw refreshError;
            } finally {
                isRefreshing = false;
            }
        }
        return Promise.reject(error);
    }
);

export default API;