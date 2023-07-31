import { createRouter, createWebHistory } from "vue-router"
import HomePage from '../views/HomePage.vue'
import ArticlePage from '../views/ArticlePage.vue'
import QueryPage from '../views/QueryPage.vue'
import UpdateUserInfoPage from '../views/UpdateUserInfoPage.vue'
import UserProfilePage from '../views/UserProfilePage.vue'
import PostArticlePage from "../views/PostArticlePage.vue";
import NotFound from "../views/NotFound.vue";
import UpdateArticlePage from "../views/UpdateArticlePage.vue";

const routes = [
    {
        path :'/', component: HomePage
    },
    {
        path: '/userProfile/:uid', component: UserProfilePage
    },
    {
        path: '/article/:articleID', component: ArticlePage
    },
    {
        path: '/query/:keyword', component: QueryPage
    },
    {
        path: '/updateUserInfo/:uid', component: UpdateUserInfoPage
    },
    {
        path: '/postArticle', component: PostArticlePage
    },
    {
        path: '/updateArticle/:articleID', component: UpdateArticlePage
    },
    { path: '/:pathMatch(.*)*', name: 'NotFound', component: NotFound }, // 是顺序匹配,所以这里可以
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router
