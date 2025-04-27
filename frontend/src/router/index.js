import { createRouter, createWebHistory } from 'vue-router';
import SignIn from '../views/SignIn.vue';
import SignUp from '../views/SignUp.vue';
import Home from '../views/Home.vue';
import { useAuthStore } from '../store/auth';

const routes = [
    { path: '/user/signin', component: SignIn},
    { path: '/user/signup', component: SignUp },
    { path: '/', redirect: '/user/signin', component: Home }
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

router.beforeEach((to, from, next) => {
    const authStore = useAuthStore();
    const publicPages = ['/user/signin', '/user/signup'];
    const authRequired = !publicPages.includes(to.path);
    const loggedIn = authStore.isAuthenticated;
    if (authRequired && !loggedIn) {
      return next('/user/signin');
    }
    next();
});

export default router;