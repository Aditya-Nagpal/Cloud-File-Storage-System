<template>
    <div>
      <h1>Welcome to the Home Page</h1>
      <p>This is a simple home page.</p>
      <router-link to="/user/signup">Sign Up</router-link>
      <router-link to="/user/login">Sign In</router-link>
    </div>
  
    <form @submit.prevent="handleUpload" enctype="multipart/form-data">
      <input type="file" @change="handleFileChange" multiple accept="image/*" ref="fileInput">
      <button type="submit">Upload</button>
    </form>
  </template>
  
  <script setup>
  import { ref } from 'vue'
  import API from '../api/axios'
  
  const selectedFiles = ref([])
  
  const handleFileChange = (event) => {
    selectedFiles.value = Array.from(event.target.files)
  }
  
  const handleUpload = async () => {
    if (selectedFiles.value.length === 0) {
      alert('Please select at least one file.')
      return
    }
  
    const formData = new FormData()
    selectedFiles.value.forEach((file, index) => {
      formData.append('files', file) // Use 'files' or whatever field your backend expects
    })
  
    try {
      const response = await API.post('http://localhost:8000/file/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })
      console.log('Upload successful:', response.data)
    } catch (error) {
      console.error('Upload failed:', error)
    }
  }
  </script>