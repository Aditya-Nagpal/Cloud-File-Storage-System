<template>
  <nav class="navbar navbar-light bg-white border-bottom px-4 d-flex justify-content-between">
    <h4 class="mt-2 fw-bold">FastFiles</h4>

    <div class="d-flex align-items-center gap-3">
      <div class="dropdown">
        <button
          class="btn btn-primary dropdown-toggle"
          type="button"
          data-bs-toggle="dropdown"
        >
          Upload
        </button>
        <ul class="dropdown-menu">
          <li><a class="dropdown-item" href="#" @click.prevent="showFileModal = true">Upload File</a></li>
          <li><a class="dropdown-item" href="#" @click.prevent="showFolderModal = true">Upload Folder</a></li>
        </ul>

        <FileModal :show="showFileModal" @close="showFileModal = false" />
        <FolderModal :show="showFolderModal" @close="showFolderModal = false" />

      </div>

      <div class="dropdown">
        <img
          :src="displayPictureUrl"
          class="rounded-circle"
          style="width: 40px; height: 40px; cursor: pointer;"
          data-bs-toggle="dropdown"
        />
        <ul class="dropdown-menu dropdown-menu-end">
          <li><a class="dropdown-item" href="/user/profile">View Profile</a></li>
          <li><a class="dropdown-item" href="#" @click="handleLogout">Logout</a></li>
        </ul>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router'
import FolderModal from './FolderModal.vue';
import FileModal from './FileModal.vue';
import defaultImage from '../images/person.png';
import { useAuthStore } from '../store/auth.js';
import { useUserStore } from '../store/user.js';

const auth = useAuthStore();
const userStore = useUserStore();
const user = computed(() => userStore.user);
const router = useRouter();
const showFileModal = ref(false);
const showFolderModal = ref(false);

const displayPictureUrl = computed(() => {
  return user?.value?.display_picture || defaultImage;
});

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