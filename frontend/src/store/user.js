import { defineStore } from 'pinia';
import API from '../api/axios';

const FETCH_USER_API = `/user/profile`;
const UPDATE_USER_API = `/user/update`;``

export const useUserStore = defineStore('user', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('user')) || null,
  }),

  actions: {
    async fetchUserProfile() {
      try {
        console.log('Fetching user profile...');
        const response = await API.get(FETCH_USER_API);
        this.user = response.data.profile;
        console.log('User profile fetched:', this.user);
        localStorage.setItem('user', JSON.stringify(this.user));
        return this.user;
      } catch (error) {
        console.error('Failed to fetch user:', error);
        throw error;
      }
    },

    async updateUserProfile(payload) {
      try {
        const response = await API.put(UPDATE_USER_API, payload);
        this.user = response.data;
        localStorage.setItem('user', JSON.stringify(this.user));
        return this.user;
      } catch (error) {
        console.error('Failed to update user:', error);
        throw error;
      }
    }
  }
});