<template>
  <div class="p-4">
    <button
      v-if="canGoBack"
      class="btn btn-outline-secondary mb-3"
      @click="fileStore.goBack"
    >
      ← Back
    </button>

    <h5>{{ fileStore.currentKey ? fileStore.currentKey : '/' }}</h5>
    <ul class="list-group">
      <li
        v-for="item in fileStore.contents"
        :key="item.filename"
        class="list-group-item d-flex justify-content-between align-items-center"
        style="cursor: pointer;"
      >
        <div
          @click="item.type === 'folder' ? fileStore.enterFolder(item.filename) : null"
          style="cursor: pointer;"
        >
          <i :class="item.type === 'folder' ? 'bi bi-folder-fill' : 'bi bi-file-earmark'" class="me-2"></i>
          {{ item.filename }}
        </div>
        
        <div class="dropdown">
          <button
            class="btn btn-sm btn-light dropdown-toggle"
            type="button"
            data-bs-toggle="dropdown"
            aria-expanded="false"
            @click.stop
          >
            ⋮
          </button>
          <ul class="dropdown-menu dropdown-menu-end">
            <li v-if="canDownload(item.type)"><a class="dropdown-item" href="#" @click.prevent="downloadItem(item)">Download</a></li>
            <li><a class="dropdown-item" href="#" @click.prevent="deleteItem(item)">Delete</a></li>
            <li><a class="dropdown-item" href="#" @click.prevent="showInfo(item)">File Information</a></li>
          </ul>

          <FileInfoModal :show="showInfoModal" :item="selectedItem" @close="handleInfoModalClose" />

          <ConfirmDeleteModal
            :show="showDeleteModal"
            :item="selectedItemToDelete"
            @close="handleDeleteModalClose"
          />

        </div>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useFileStore } from '../store/file.js';
import FileInfoModal from './FileInfoModal.vue';
import ConfirmDeleteModal from './ConfirmDeleteModal.vue';
import { toast } from 'vue3-toastify';

const fileStore = useFileStore();

const showInfoModal = ref(false);
const selectedItem = ref(null);

const showDeleteModal = ref(false);
const selectedItemToDelete = ref(null);

onMounted(async () => {
    await fileStore.fetchContents();
});

const canGoBack = computed(() => fileStore.keyStack.length > 0);
const canDownload = (type) => type === 'file';

const handleInfoModalClose = () => {
    console.log('Info modal closed');
    showInfoModal.value = false;
    selectedItem.value = null;
};

const handleDeleteModalClose = () => {
    console.log('Delete modal closed');
    showDeleteModal.value = false;
    selectedItemToDelete.value = null;
};

const showInfo = (item) => {
    selectedItem.value = item
    showInfoModal.value = true;
};

const deleteItem = (item) => {
    selectedItemToDelete.value = item;
    showDeleteModal.value = true;
};

const downloadItem = async (item) => {
    try {
        await fileStore.downloadFile(item.id);
    } catch (error) {
        toast.error('Could not download file. Please try again later.');
    }
};

</script>
