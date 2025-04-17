import { createRouter, createWebHistory } from 'vue-router';
import SignIn from '../views/SignIn.vue';
import SignUp from '../views/SignUp.vue';

const routes = [
    { path: '/user/signin', component: SignIn},
    { path: '/user/signup', component: SignUp },
    { path: '/', redirect: '/user/signin' }
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

export default router;