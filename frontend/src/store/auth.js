import { defineStore } from 'pinia';
import axios from 'axios';
import API from '../api/axios';

const BASE_URL = import.meta.env.VITE_API_BASE_URL;
const SIGN_IN_API = BASE_URL + '/auth/login';
const SIGN_UP_API = BASE_URL + '/auth/register';
const LOGOUT_API = BASE_URL + '/auth/logout';
const PROTECT_API = BASE_URL+'/auth/protected';
const REFRESH_API = BASE_URL + '/auth/refresh';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('user')) || null,
    accessToken: localStorage.getItem('accessToken') || null
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
        const { user, accessToken } = response.data;
        this.user = user;
        this.accessToken = accessToken;

        // save to local storage
        localStorage.setItem('accessToken', this.accessToken);
        localStorage.setItem('user', JSON.stringify(this.user));
        return { user, accessToken };
      } catch (error) {
        throw error;
      }
    },

    async checkAuth() {
      try {
        console.log('Checking authentication...', this.accessToken);
        const response = await API.get('/auth/protected');
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
        await axios.get(LOGOUT_API, {
          withCredentials: true
        })
      } catch (error) {
        console.error('Logout failed:', error);
        throw error;
      }

      this.user = null;
      this.accessToken = null;

      // remove from local storage
      localStorage.removeItem('accessToken');
      localStorage.removeItem('user');
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