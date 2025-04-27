import { createRouter, createWebHistory } from 'vue-router';
import SignIn from '../views/SignIn.vue';
import SignUp from '../views/SignUp.vue';
import Home from '../views/Home.vue';

const routes = [
    { path: '/user/login', component: SignIn},
    { path: '/user/signup', component: SignUp },
    { path: '/', component: Home }
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

export default router;