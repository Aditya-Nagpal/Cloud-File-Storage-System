import { defineStore } from 'pinia';
import axios from 'axios';

const BASE_URL = import.meta.env.API_BASE_URL;
const SIGN_IN_API = '';
// const SIGN_UP_API = `${BASE_URL}/auth/register`;
const SIGN_UP_API = ''

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: null
  }),

  actions: {
    async signIn(email, password) {
      try {
        const response = await axios.post(SIGN_IN_API, {
          email,
          password
        });
        const { user, token } = response.data;
        this.user = user;
        this.token = token;

        // save to local storage
        localStorage.setItem('token', this.token);
        localStorage.setItem('user', JSON.stringify(this.user));
      } catch (error) {
        console.error('Login failed:', error);
        throw error;
      }
    },

    async signUp(payload) {
        try {
            const response = await axios.post(SIGN_UP_API, payload);
            const { user, token } = response.data;
            this.user = user;
            this.token = token;

            // save to local storage
            localStorage.setItem('token', this.token);
            localStorage.setItem('user', JSON.stringify(this.user));
        } catch (error) {
          console.error('Registration failed:', error);
          throw error;
        }
    },

    logout() {
      this.user = null;
      this.token = null;

      // remove from local storage
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    }
  },

  getters: {
    isAuthenticated() {
      return !!this.token;
    },
    getUser() {
      return this.user;
    }
  }
});