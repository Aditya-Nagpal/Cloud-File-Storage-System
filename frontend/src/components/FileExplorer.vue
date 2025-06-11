<template>
  <div class="p-4">
    <button
      v-if="canGoBack"
      class="btn btn-outline-secondary mb-3"
      @click="goBack"
    >
      ← Back
    </button>

    <h5>{{ fileStore.currentKey ? fileStore.currentKey : '/' }}</h5>
    <ul class="list-group">
      <li
        v-for="item in fileStore.contents"
        :key="item.filename"
        class="list-group-item d-flex justify-content-between"
        @click="handleItemClick(item)"
        style="cursor: pointer;"
      >
        <span>
          <i :class="item.type === 'folder' ? 'bi bi-folder-fill' : 'bi bi-file-earmark'" class="me-2"></i>
          {{ item.filename }}
        </span>
        <small>{{ item.type }}</small>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useFileStore } from '../store/file.js';

const fileStore = useFileStore();

onMounted(async () => {
  await fileStore.fetchContents();
});

const canGoBack = computed(() => fileStore.keyStack.length > 0);

const currentFolderName = computed(() => {
  const parts = fileStore.currentKey.split('/').filter(Boolean)
  return parts.length ? parts[parts.length - 1] : 'root'
});

const handleItemClick = (item) => {
  if (item.type === 'folder') {
    fileStore.enterFolder(item.filename);
    console.log(fileStore.currentKey);
  }
};

const goBack = () => {
  fileStore.goBack();
  console.log(fileStore.currentKey);
};
</script>
