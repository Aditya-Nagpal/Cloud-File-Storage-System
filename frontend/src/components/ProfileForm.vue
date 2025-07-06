<template>
  <div class="container mt-4">
    <h3>User Profile</h3>

    <div v-if="user" class="card p-4 mt-3">
      <div class="mb-3 text-center">
        <img
          :src="displayPictureUrl"
          alt="Display Picture"
          class="rounded-circle mb-2"
          width="100"
          height="100"
        />
        <div>
          <input type="file" @change="handleFileChange" class="form-control mt-2" />
          <button class="btn btn-sm btn-outline-danger mt-1" @click="removeDP" v-if="user.displayPicture">
            Remove Display Picture
          </button>
        </div>
      </div>

      <div class="mb-3">
        <label>Email</label>
        <input type="email" class="form-control" :value="user.email" disabled />
      </div>

      <div class="mb-3">
        <label>Name</label>
        <input type="text" class="form-control" v-model="form.name" />
      </div>

      <div class="mb-3">
        <label>Age</label>
        <input type="number" class="form-control" v-model="form.age" />
      </div>

      <button class="btn btn-primary" @click="updateProfile">Save Changes</button>
    </div>

    <div v-else>Loading...</div>
  </div>
</template>

<script setup>
import { onMounted, computed, reactive } from 'vue';
import { useUserStore } from '../store/user';
import defaultImage from '../images/person.png';
import { toast } from 'vue3-toastify';

const userStore = useUserStore();
const user = userStore.user;

const form = reactive({
  name: '',
  age: '',
  displayPictureFile: null
});

const displayPictureUrl = computed(() => {
  return user.displayPicture || defaultImage;
});
console.log('displayPictureUrl: ', displayPictureUrl);

onMounted(async () => {
    const profile = user;
    console.log('User profile:', profile);
    form.name = profile.name;
    form.age = profile.age;
});

const handleFileChange = (e) => {
  form.displayPictureFile = e.target.files[0];
};

const updateProfile = async () => {
  try {
    await userStore.updateUserProfile(form);
    toast.success('Profile updated successfully');
  } catch (err) {
    console.error('Update failed', err);
    alert('Failed to update profile');
  }
};

const removeDP = async () => {
  try {
    await userStore.removeDisplayPicture();
  } catch (err) {
    console.error('Failed to remove DP', err);
  }
};
</script>