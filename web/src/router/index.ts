import { createRouter, createWebHistory } from 'vue-router';
import HomePage from '../pages/HomePage.vue';
import StatusPage from '../pages/StatusPage.vue';
import ConnectionPage from '../pages/ConnectionPage.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: HomePage },
    { path: '/status', name: 'status', component: StatusPage },
    { path: '/connection', name: 'connection', component: ConnectionPage },
  ],
});

export default router;
