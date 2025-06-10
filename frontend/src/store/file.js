import { defineStore } from 'pinia';
import API from '../api/axios';

export const useFileStore = defineStore('file', {
    state: () => ({
        currentKey: '',
        keyStack: [],
        contents: []
    }),

    actions: {
        async fetchContents(key = '') {
            try {
                return {
                name: 'root',
                    children: [
                        { name: 'Documents', type: 'folder', children: [
                        { name: 'Resume.pdf', type: 'file' },
                        { name: 'Project', type: 'folder', children: [
                            { name: 'code.js', type: 'file' }
                        ]}
                        ]},
                        { name: 'image.png', type: 'file' }
                    ]
                };
            } catch (error) {
                console.error('Error fetching contents:', error);
                throw error;
            }
        },

        async uploadFile(file, folderKey) {
            const formData = new FormData();
            formData.append('uploadType', 'file');
            formData.append('file', file);
            formData.append('folderKey', folderKey);
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

        async uploadFolder(folderKey, folderName) {
            try {
                const response = await API.post('/file/upload', { uploadType: 'folder', folderKey, folderName }, {
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
        },

        async goBack() {
            if (this.keyStack.length > 0) {
                this.currentKey = this.keyStack.pop();
            }
        }
    }
});