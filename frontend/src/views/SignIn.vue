<template>
  <form @submit.prevent="handleSignIn">
    <h2>Sign In</h2>
    <div>
      <input v-model="email" type="email" placeholder="Email" required />
      <input v-model="password" type="password" placeholder="Password" required />
      <button type="submit">Sign In</button>
      <a href="/user/signup">Don't have an account? Sign up</a>
    </div>
  </form>
</template>

<script setup>
import { ref } from 'vue';
import { useAuthStore } from '../store/auth';
import { useUserStore } from '../store/user';
import { useRouter } from 'vue-router';
import { toast } from 'vue3-toastify';

const auth = useAuthStore();
const user = useUserStore();
const router = useRouter();

const email = ref('');
const password = ref('');

const handleSignIn = async () => {
  try {
    await auth.signIn(email.value, password.value);
    toast.success('Login successful');
    resetForm();
    router.push('/'); // redirect after login
    try {
      console.log('Fetching user profile ...');
      await user.fetchUserProfile();
      console.log('User profile fetched:', user.user);
    } catch (error) {
      console.error('Error in fetchUserProfile:', error);
    }
  } catch (error) {
    console.error('Error in handleSignIn', error);
    toast.error(error.response.data.message);
    resetForm();
    return;
  }
}

const resetForm = () => {
  email.value = ''
  password.value = ''
}
</script>