<template>
  <div class="p-4">
    <button
      v-if="canGoBack"
      class="btn btn-outline-secondary mb-3"
      @click="goBack"
    >
      ← Back
    </button>

    <h5>{{ currentFolder.name }}</h5>
    <ul class="list-group">
      <li
        v-for="item in currentFolder.children"
        :key="item.name"
        class="list-group-item d-flex justify-content-between"
        @click="handleItemClick(item)"
        style="cursor: pointer;"
      >
        <span>
          <i :class="item.type === 'folder' ? 'bi bi-folder-fill' : 'bi bi-file-earmark'" class="me-2"></i>
          {{ item.name }}
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

// Sample dummy file tree
const root = {
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

const canGoBack = computed(() => fileStore.keyStack.length > 0);

const currentFolderName = computed(() => {
  const parts = fileStore.currentKey.split('/').filter(Boolean)
  return parts.length ? parts[parts.length - 1] : 'root'
})

// const handleItemClick = (item) => {
//   if (item.type === 'folder') {
//     fileStore.enterFolder(item.name);
//     console.log(fileStore.currentKey);
//   }
// };

const goBack = () => {
  fileStore.goBack();
  console.log(fileStore.currentKey);
};
</script>
