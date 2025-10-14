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
        const response = await API.get(FETCH_USER_API);
        this.user = response.data.profile;
        localStorage.setItem('user', JSON.stringify(this.user));
        return this.user;
      } catch (error) {
        console.error('Failed to fetch user:', error);
        throw error;
      }
    },

    async updateUserProfile(payload) {
      try {
        const response = await API.patch(`${UPDATE_USER_API}?dp=false`, payload);
        const updatedUser = {};
        for(let [key, value] of Object.entries(payload)) {
          if (["",null,undefined,0].includes(value)) continue;
          updatedUser[key] = value;
        }
        this.user = { ...this.user, ...updatedUser };
        localStorage.setItem('user', JSON.stringify(this.user));
        return true;
      } catch (error) {
        console.error('Failed to update user:', error);
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
        this.user.displayPicture = response.data.displayPicture;
        localStorage.setItem('user', JSON.stringify(this.user));
        return true;
      } catch (error) {
        console.error('Failed to update dp: ', error);
        throw error;
      }
    },

    async removeDisplayPicture() {
      try {
        const response = await API.patch(`${UPDATE_USER_API}?removeDp=true`);
        this.user.displayPicture = null;
        localStorage.setItem('user', JSON.stringify(this.user));
        return true;
      } catch (error) {
        console.error('Failed to remove dp: ', error);
        throw error;
      }
    }
  }
});