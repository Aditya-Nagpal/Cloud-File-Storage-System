<template>
  <div>
    <h1>Welcome to the Home Page</h1>
    <p>This is a simple home page.</p>
    <router-link to="/user/signup">Sign Up</router-link>
    <router-link to="/user/login">Sign In</router-link>
  </div>

  <form @submit.prevent="handleUpload" enctype="multipart/form-data">
    <input type="file" @change="handleFileChange" accept=".jpg,.jpeg,.png,.pdf,.doc,.docx" ref="fileInput">
    <button type="submit">Upload</button>
  </form>
</template>
  
<script setup>
import { ref } from 'vue'
import API from '../api/axios'
import { toast } from 'vue3-toastify';

const selectedFiles = ref(null)
const fileInput = ref(null)

const handleFileChange = (event) => {
  const file = event.target.files[0]
  if(file){
    selectedFiles.value = file;
  }
}

const handleUpload = async () => {
  if (!selectedFiles.value) {
    toast.warn('Please select at least one file.')
    return
  }

  const formData = new FormData()
  formData.append('file', selectedFiles.value)
  
  for (let [key, value] of formData.entries()) {
    console.log(`${key}:`, value)
  }

  try {
    const response = await API.post('http://localhost:8000/file/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
    })
    resetFile();
    console.log('Upload successful:', response.data)
  } catch (error) {
    console.error('Upload failed:', error)
  }
}

const resetFile = () => {
  fileInput.value.value = '';
  selectedFiles.value = null;
}
</script>