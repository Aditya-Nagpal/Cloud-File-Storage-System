import { defineStore } from 'pinia';
import axios from 'axios';

const BASE_URL = import.meta.env.VITE_API_BASE_URL;
const SIGN_IN_API = BASE_URL + '/auth/login';
const SIGN_UP_API = BASE_URL + '/auth/register';
const LOGOUT_API = BASE_URL + '/auth/logout';
const PROTECT_API = BASE_URL+'/auth/protected';

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
        return {user, accessToken};
      } catch (error) {
        throw error;
      }
    },

    async checkAuth() {
      if(!this.accessToken){throw new Error;}
      try {
        console.log(PROTECT_API)
        const response = await axios.get(PROTECT_API, {
          headers: {
            Authorization: `Bearer ${this.accessToken}`
          }
        })
        console.log(response.data)
        return {accessToken: this.accessToken, user: this.user};
      } catch (error) {
        console.error('Authentication check failed:', error);
        await this.logout();
        throw error;
      }
    },

    async logout() {
      try {
        await axios.post(LOGOUT_API, {}, {
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

    setAccessToken(token) {
      this.accessToken = token;
      localStorage.setItem('accessToken', token);
    }
  },
  getters: {
    isAuthenticated(state) {
      return !!state.accessToken;
    },
  }
});