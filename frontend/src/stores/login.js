import {defineStore} from "pinia";
import {ref} from "vue";

export const useLoginStore = defineStore('login',() =>{
    const isLogin = ref(false)
    const currentUser = ref({
        token: "",
        uid:0,
        nickname : "",
        avatarUrl : "",
        signature: "",
    })

    const headers = ref({
    })

    function setCurrentUser(userInfo){
        currentUser.value.token = userInfo.token
        currentUser.value.uid = userInfo.uid
        currentUser.value.nickname = userInfo.nickname
        currentUser.value.avatarUrl = userInfo.avatarUrl
        currentUser.value.signature = userInfo.signature
        // 应该给header单独分个函数出来的
        headers.value.Authorization = userInfo.token
    }

    function getCurrentUser(){
        return currentUser.value
    }

    function getHeaders(){
        return headers.value
    }

    function exitLogin(){
        isLogin.value = false
        currentUser.value.token = ""
        currentUser.value.uid = 0
        currentUser.value.nickname = ""
        currentUser.value.avatarUrl = ""
        currentUser.value.signature = ""
        headers.value.Authorization = ""
    }

    return { isLogin ,currentUser ,setCurrentUser, getCurrentUser , headers, exitLogin}


},{
    persist: {
        enabled: true,
    },
})