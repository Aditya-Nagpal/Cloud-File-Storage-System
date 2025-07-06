import { defineStore } from 'pinia';
import { resetAllStores, resetLocalStorage } from './reset';
import axios from 'axios';
import API from '../api/axios';

const BASE_URL = import.meta.env.VITE_API_BASE_URL;
const SIGN_IN_API = BASE_URL + '/auth/login';
const SIGN_UP_API = BASE_URL + '/auth/register';
const PROTECT_API = '/auth/protected';
const LOGOUT_API = '/auth/logout';
const REFRESH_API = BASE_URL + '/auth/refresh';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    accessToken: localStorage.getItem('accessToken') || null,
  }),

  actions: {
    async signUp(payload) {
      try {
        const response = await axios.post(SIGN_UP_API, payload);
        return response.data;
      } catch (error) {
        throw error;
      }
    },

    async signIn(email, password) {
      try {
        const response = await axios.post(SIGN_IN_API, {
          email,
          password
        }, {
          withCredentials: true
        });
        const { accessToken } = response.data;
        this.accessToken = accessToken;

        // save to local storage
        localStorage.setItem('accessToken', this.accessToken);
        return { accessToken };
      } catch (error) {
        throw error;
      }
    },

    async checkAuth() {
      try {
        console.log('Checking authentication...', this.accessToken);
        const response = await API.get(PROTECT_API);
        console.log('Authentication check response:', response.data);
        return response.data;
      } catch (error) {
        console.error('Authentication check failed:', error);
        await this.logout();
        throw error;
      }
    },

    async logout() {
      try {
        await API.get(LOGOUT_API)
      } catch (error) {
        console.error('Logout failed:', error);
        throw error;
      }

      resetLocalStorage();
      resetAllStores();
    },

    async refreshAccessToken() {
      try {
        const res = await axios.get(REFRESH_API, {
          withCredentials: true // sends HttpOnly cookie
        });
        const { accessToken } = res.data;

        this.setAccessToken(accessToken);
        return accessToken;
      } catch (error) {
        throw error;
      }
    },

    setAccessToken(accessToken) {
      this.accessToken = accessToken;
      localStorage.setItem('accessToken', accessToken);
    }
  },
  getters: {
    isAuthenticated(state) {
      return !!state.accessToken;
    },
  }
});