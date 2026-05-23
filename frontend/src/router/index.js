import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import AuthView from '../views/Auth.vue'
import ProfileView from '../views/Profile.vue'
import NumAuthGuard from '../router/NunAuthGuard'
import PostDetails from '@/components/post/PostDetails.vue'
import Search from '@/components/search/Search.vue'
import Notification from '@/components/notification/Notification.vue'
import Chat from '@/components/chat/Chat.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView
  },
  {
    path: '/Auth',
    name: 'Auth',
    component: AuthView,
    beforeEnter: [NumAuthGuard]
  },
  {
    path: '/Search',
    name: 'Search',
    component: Search
  },
  {
    path: '/PostDetails/:id',
    name: "PostDetails",
    component: PostDetails
  },
  {
    path: '/Profile/:id',
    name: 'Profile',
    component: ProfileView
  },
  {
    path: '/Notification',
    name: 'Notification',
    component: Notification
  },
  {
    path: '/Chat',
    name: 'Chat',
    component: Chat
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
