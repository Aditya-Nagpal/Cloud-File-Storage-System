import { createRouter, createWebHistory } from 'vue-router';
import SignIn from '../views/SignIn.vue';
import SignUp from '../views/SignUp.vue';

const routes = [
  {
    path: '/user/signin',
    name: 'SignIn',
    component: SignIn,
  },
  {
    path: '/user/signup',
    name: 'SignUp',
    component: SignUp,
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

export default router;