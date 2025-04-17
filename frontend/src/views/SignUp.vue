<template>
  <form @submit.prevent="handleSignUp">
    <h2>Sign Up</h2>
    <input v-model="name" type="text" placeholder="Name" required />
    <input v-model="email" type="email" placeholder="Email" required />
    <input v-model="age" type="number" placeholder="Age" required />
    <input v-model="password" type="password" placeholder="Password" required />
    <input v-model="confirmPassword" type="password" placeholder="Confirm Password" required />
    <button type="submit">Sign Up</button>
    <a href="/user/signin">Already have an account? Sign in</a>
  </form>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../store/auth'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()

const name = ref('')
const email = ref('')
const age = ref('')
const password = ref('')
const confirmPassword = ref('')

const handleSignUp = async () => {
  if (password.value !== confirmPassword.value) {
    alert('Passwords do not match')
    return
  }

  try {
    await auth.signUp({
      name: name.value,
      email: email.value,
      age: Number(age.value),
      password: password.value,
    })
    router.push('/') // redirect after sign-up
  } catch (err) {
    alert('Sign up failed')
  }
}
</script>
