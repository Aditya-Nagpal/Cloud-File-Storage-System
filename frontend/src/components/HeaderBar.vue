<template>
  <nav class="navbar navbar-light bg-white border-bottom px-4 d-flex justify-content-between">
    <h4 class="mt-2 fw-bold">FastFiles</h4>

    <div class="d-flex align-items-center gap-3">
      <!-- Upload dropdown -->
      <div class="dropdown">
        <button
          class="btn btn-primary dropdown-toggle"
          type="button"
          data-bs-toggle="dropdown"
        >
          Upload
        </button>
        <ul class="dropdown-menu">
          <li><a class="dropdown-item" href="#" @click.prevent="triggerFileUpload">Upload File</a></li>
          <li><a class="dropdown-item" href="#" @click.prevent="handleFolderUpload">Upload Folder</a></li>
        </ul>
        <!-- Hidden input for file selection -->
        <input
          type="file"
          ref="fileInput"
          class="d-none"
          @change="handleFileSelected"
        />
      </div>

      <!-- Profile dropdown -->
      <div class="dropdown">
        <img
          src="../images/person.png"
          class="rounded-circle"
          style="width: 40px; height: 40px; cursor: pointer;"
          data-bs-toggle="dropdown"
        />
        <ul class="dropdown-menu dropdown-menu-end">
          <li><a class="dropdown-item" href="#">View Profile</a></li>
          <li><a class="dropdown-item" href="#" @click="handleLogout">Logout</a></li>
        </ul>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useAuthStore } from '../store/auth.js';
import { useFileStore } from '../store/file.js';
import { useRouter, useRoute } from 'vue-router'

const auth = useAuthStore();
const fileStore = useFileStore();
const router = useRouter();
const fileInput = ref(null);

const triggerFileUpload = () => {
  fileInput.value.click();
};

const handleFileSelected = async (event) => {
  const file = event.target.files[0];
  if (!file) return;

  try {
    console.log(fileStore.currentKey)
    // await fileStore.uploadFile(file, fileStore.currentKey);
    await fileStore.uploadFile(file, "user/");

    // await fileStore.fetchContents(fileStore.currentKey);
  } catch (error) {
    console.error('Upload failed:', error);
  }
};

const handleFolderUpload = async () => {
  try {
    console.log(fileStore.currentKey)
    const folderKey = "user/";
    await fileStore.uploadFolder(folderKey, "data");
  } catch (error) {
    console.error('Error uploading folder:', error);
  }
};

const handleLogout = async () => {
  try {
    await auth.logout();
    router.push('/user/login');
  } catch (error) {
    console.error('logout failed: ', error);
    return;
  }
};

</script>