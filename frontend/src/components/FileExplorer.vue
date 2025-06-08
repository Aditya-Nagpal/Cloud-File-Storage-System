<template>
  <div class="p-4">
    <button
      v-if="history.length"
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
import { ref } from 'vue'

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
}

const history = ref([])
const currentFolder = ref(root)

function handleItemClick(item) {
  if (item.type === 'folder') {
    history.value.push(currentFolder.value)
    currentFolder.value = item
  }
}

function goBack() {
  currentFolder.value = history.value.pop()
}
</script>
