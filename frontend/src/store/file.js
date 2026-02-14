import { defineStore } from 'pinia';
import API from '../api/axios';

export const useFileStore = defineStore('file', {
    state: () => ({
        currentKey: '',
        keyStack: [],
        contents: []
    }),

    actions: {
        async fetchContents() {
            try {
                const response = await API.get('/file/list', {
                    params: {
                        parentPath: this.currentKey
                    }
                });
                console.log('Fetched contents:', response?.data?.files);
                this.contents = response?.data?.files;
                return this.contents;
            } catch (error) {
                console.error('Error fetching contents:', error);
                throw error;
            }
        },

        async uploadFile(file) {
            const formData = new FormData();
            formData.append('uploadType', 'file');
            formData.append('file', file);
            formData.append('folderKey', this.currentKey);
            try {
                const response = await API.post('/file/upload', formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                    }
                });
                console.log(response.data);
                return response.data;
            } catch (error) {
                console.error('Error uploading file:', error);
                throw error;
            }
        },

        async uploadFolder(folderName) {
            try {
                const response = await API.post('/file/upload', { uploadType: 'folder', folderKey: this.currentKey, folderName }, {
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded'
                    }
                });
                console.log(response.data);
                return response.data;
            } catch (error) {
                console.error('Error creating folder:', error.AxiosError);
                throw error;
            }
        },

        async enterFolder(folderName) {
            this.keyStack.push(this.currentKey);
            this.currentKey += folderName + '/';
            console.log('Current key:', this.currentKey);
            await this.fetchContents(this.currentKey);
        },

        async goBack() {
            if (this.keyStack.length > 0) {
                this.currentKey = this.keyStack.pop();
                await this.fetchContents(this.currentKey);
            }
        },

        async deleteContent(fileName, type) {
            try {
                console.log('Deleting content:', fileName, type, this.currentKey);
                const response = await API.delete('/file/delete', {
                    data: { parentPath: this.currentKey, fileName, type }
                });
                console.log(response.data);
                return response.data;
            } catch (error) {
                console.error('Error deleting content:', error);
                throw error;
            }
        },

        async downloadFile(id) {
            try {
                const response = await API.get(`/file/download/${id}`);
                const url = response.data.downloadURL;
                window.location.href = url;
                return true;
            } catch (error) {
                console.error('Error downloading file:', error);
                throw error;
            }
        }
    }
});