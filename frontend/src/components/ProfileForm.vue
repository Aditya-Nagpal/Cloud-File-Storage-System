<template>
  <div class="container mt-4">
    <h3>User Profile</h3>

    <div v-if="user" class="card p-4 mt-3">
      <!-- Display Picture Section -->
      <div class="mb-3 text-center">
        <img
          :src="displayPictureUrl"
          alt="Display Picture"
          class="rounded-circle mb-2"
          width="100"
          height="100"
        />

        <div v-if="!isChangingDP">
          <button
            class="btn btn-sm btn-outline-primary me-2"
            @click="triggerDPChange"
          >
            {{ user?.displayPicture ? 'Change DP' : 'Upload DP' }}
          </button>
          <button
            v-if="user?.displayPicture"
            class="btn btn-sm btn-outline-danger"
            @click="removeDP"
          >
            Remove DP
          </button>
        </div>

        <div v-else>
          <input type="file" @change="handleFileChange" class="form-control mt-2" />
          <button class="btn btn-sm btn-success mt-2 me-2" :disabled="!dpInput" @click="uploadDP">Save</button>
          <button class="btn btn-sm btn-secondary mt-2" @click="cancelDPChange">Cancel</button>
        </div>
      </div>

      <!-- User Info Section -->
      <div class="mb-3">
        <label>Email</label>
        <input type="email" class="form-control" :value="user.email" disabled />
      </div>

      <div class="mb-3">
        <label>Name</label>
        <input
          type="text"
          class="form-control"
          v-model="form.name"
          :disabled="!isEditing"
        />
      </div>

      <div class="mb-3">
        <label>Age</label>
        <input
          type="number"
          class="form-control"
          v-model="form.age"
          :disabled="!isEditing"
        />
      </div>

      <div>
        <button
          class="btn btn-primary me-2"
          v-if="isEditing"
          @click="updateProfile"
        >
          Save Changes
        </button>
        <button
          class="btn btn-secondary"
          v-if="isEditing"
          @click="cancelEdit"
        >
          Cancel
        </button>
        <button
          class="btn btn-outline-primary"
          v-else
          @click="isEditing = true"
        >
          Edit Info
        </button>
      </div>
    </div>

    <div v-else>Loading...</div>
  </div>
</template>


<script setup>
import { onMounted, computed, reactive, ref } from 'vue';
import { useUserStore } from '../store/user';
import defaultImage from '../images/person.png';
import { toast } from 'vue3-toastify';

const userStore = useUserStore();
// const user = userStore.user;
let user = computed(() => userStore.user);
console.log('user: ', user.value);
const isEditing = ref(false);
const isChangingDP = ref(false);
const dpInput = ref(null);

const form = reactive({
  name: null,
  age: null,
});

const displayPictureUrl = computed(() => {
  return user?.value?.displayPicture || defaultImage;
});

onMounted(async () => {
    if(user){
      form.name = user?.value?.name;
      form.age = user?.value?.age;
    }
});

const triggerDPChange = () => {
  isChangingDP.value = true;
};

const cancelDPChange = () => {
  dpInput.value = null;
  isChangingDP.value = false;
};

const handleFileChange = (e) => {
  dpInput.value = e.target.files?.[0];
};

const uploadDP = async () => {
  const file = dpInput.value;
  if(!file){
    console.log('We are here 2')
    toast.info('No file selected');
    return;
  }

  const formData = new FormData();
  formData.append('displayPicture', file);

  try {
    await userStore.updateDisplayPicture(formData);
    toast.success('DP updated successfully');
    isChangingDP.value = false;
  } catch (error) {
    console.error('DP upload failed', error);
    toast.error('Failed to update display picture');
  }
};

const removeDP = async () => {
  try {
    await userStore.updateDisplayPicture({
      'displayPicture': null
    });
    toast.success('Display picture removed');
  } catch (err) {
    console.error('Remove DP failed', err);
    toast.error('Failed to remove display picture');
  }
};

const updateProfile = async () => {
  try {
    await userStore.updateUserProfile({
      name: form.name,
      age: form.age
    });
    toast.success('Profile updated');
    isEditing.value = false;
  } catch (err) {
    console.error('Update failed', err);
    toast.error('Failed to update profile');
  }
};

const cancelEdit = () => {
  if (user) {
    form.name = user.name;
    form.age = user.age;
  }
  isEditing.value = false;
};
</script>