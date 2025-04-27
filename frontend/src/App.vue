<script setup>
import { onMounted } from 'vue';
import { useAuthStore } from './store/auth';
import { useRouter, useRoute } from 'vue-router'

const auth = useAuthStore();
const router = useRouter();
const route = useRoute();

onMounted(async () => {
  await router.isReady();
  const publicPages = ['/user/login', '/user/signup'];
  try {
    await auth.checkAuth();
    if(publicPages.includes(route.path)){
      router.push('/');
      return;
    }
  } catch (error) {
    console.error('checkAuth failed: ', error);
    if(!publicPages.includes(route.path)) {
      router.push('/user/login');
      return;
    }
  }
});
</script>

<template>
  <div>
    <router-view />
  </div>
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>