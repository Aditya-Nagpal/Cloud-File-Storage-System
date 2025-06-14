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
          <li><a class="dropdown-item" href="#" @click.prevent="showModal = true">Upload Folder</a></li>
        </ul>

        <FolderModal :show="showModal" @close="showModal = false" />

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
import { ref } from 'vue';
import { useRouter } from 'vue-router'
import FolderModal from './FolderModal.vue';
import { useAuthStore } from '../store/auth.js';
import { useFileStore } from '../store/file.js';

const auth = useAuthStore();
const fileStore = useFileStore();
const router = useRouter();
const fileInput = ref(null);
const showModal = ref(false);

const triggerFileUpload = () => {
  console.log('triggerFileUpload called');
  fileInput.value.click();
};

const handleFileSelected = async (event) => {
  const file = event.target.files[0];
  if (!file) return;

  try {
    console.log('we are here', fileStore.currentKey)
    await fileStore.uploadFile(file);
    await fileStore.fetchContents();
  } catch (error) {
    console.error('Upload failed:', error);
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