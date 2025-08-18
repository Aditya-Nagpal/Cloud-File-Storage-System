import { defineStore } from 'pinia';
import API from '../api/axios';

const FETCH_USER_API = `/user/profile`;
const UPDATE_USER_API = `/user/profile`;

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

    async updateDisplayPicture(payload) {
      try {
        const response = await API.patch(`${UPDATE_USER_API}?dp=true`, payload, {
          headers: {
            "Content-Type": "multipart/form-data"
          }
        });
        console.log('Display picture updated:', response.data);
        this.user.displayPicture = response.data.displayPicture;
        localStorage.setItem('user', JSON.stringify(this.user));
        return true;
      } catch (error) {
        console.error('Failed to update dp: ', error);
        throw error;
      }
    },

    async updateUserProfile(payload) {
      try {
        console.log('Updating user profile with payload:', payload);
        const response = await API.patch(`${UPDATE_USER_API}?dp=false`, payload);
        console.log('User profile updated:', response.data);
        this.user = { ...this.user, ...response.data.updatedUser }
        localStorage.setItem('user', JSON.stringify(this.user));
        return true;
      } catch (error) {
        console.error('Failed to update user:', error);
        throw error;
      }
    }
  }
});