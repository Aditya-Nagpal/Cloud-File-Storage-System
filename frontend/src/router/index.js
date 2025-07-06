import { createRouter, createWebHistory } from 'vue-router';
import SignIn from '../views/SignIn.vue';
import SignUp from '../views/SignUp.vue';
import Profile from '../views/Profile.vue';
import Home from '../views/Home.vue';
import { useAuthStore } from '../store/auth';

const routes = [
    { path: '/user/login', component: SignIn},
    { path: '/user/signup', component: SignUp },
    { path: '/user/profile', component: Profile },
    { path: '/', component: Home }
];

const publicPages = ['/user/login', '/user/signup'];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

router.beforeEach(async (to) => {
  const auth = useAuthStore();

  const isPublic = publicPages.includes(to.path);
  const accessToken = auth.accessToken;

  if(!accessToken && !isPublic){
    return '/user/login';
  }

  if(accessToken){
    try {
      await auth.checkAuth();
      if(isPublic){
        return '/';
      }
    } catch (error) {
      return '/user/login';
    }
  }
});

export default router;