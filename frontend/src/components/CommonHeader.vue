<script setup>
import {Avatar, Edit, Flag, Fold, Location, Search, Star} from "@element-plus/icons-vue";
import {ref} from "vue";
import {useRoute, useRouter} from "vue-router";
import LoginDialog from "./LoginDialog.vue";
import {useLoginStore} from "../stores/login.js";
import {ElMessage} from "element-plus";

// const user = {
//   avatarUrl : "https://static.runoob.com/images/demo/demo1.jpg"
// }

const router = useRouter()

function goToHomePage(){
  router.push("/")
}

function goToPostArticlePage(){
  if (loginStore.isLogin){
    router.push("/postArticle")
  }
  else{
    toggleLoginDialog()
  }
}

const queryKeyword = ref("")

const loginStore = useLoginStore()

const showLoginDialog = ref(false)

function toggleLoginDialog(){
  if (showLoginDialog.value){
    showLoginDialog.value = false
    setTimeout(() =>{ // 好玩又离谱的地方
      showLoginDialog.value = true
    }, 1)
  }
  else{
    showLoginDialog.value = true
  }
}



</script>

<template>


<!--  <div class="header-container">
    <el-row>
      <el-col :span="6">
        <el-input
            v-model="queryKeyword"
            placeholder="don't repeat yourself"
            :prefix-icon="Search"
            style="width: 200px"
        />
        <el-button @click="router.push('/query/'.concat(queryKeyword))" v-show="queryKeyword!=''">查询</el-button>
      </el-col>

      <el-col :span="6">
        <el-button type="primary" :icon="Edit" @click="goToPostArticlePage">写文章</el-button>
      </el-col>

      <el-col :span="6">
        <el-button type="primary" :icon="Flag" @click="goToHomePage">首页</el-button>
      </el-col>

      <el-col :span="6">
        <div v-if="loginStore.isLogin">
&lt;!&ndash;          这里需要改,改头像的话得刷新,因为其实url是没有变化的,不能触发响应式状态&ndash;&gt;
          <el-avatar  :src="loginStore.getCurrentUser().avatarUrl"></el-avatar>
          <el-button @click="router.push('/userProfile/'.concat(loginStore.getCurrentUser().uid))">去我的主页</el-button>
        </div>
        <el-button v-else @click="toggleLoginDialog">登录</el-button>
      </el-col>

      <login-dialog :show-login-dialog="showLoginDialog"></login-dialog>

    </el-row>

&lt;!&ndash;    如果登录了,显示头像;如果没有登录,显示注册的按钮&ndash;&gt;

  </div>-->

  <login-dialog :show-login-dialog="showLoginDialog" @close-login-dialog="showLoginDialog=false"></login-dialog>

  <el-menu
      mode="horizontal"
      :ellipsis="false"
  >
    <el-menu-item v-if="loginStore.isLogin" @click="router.push('/userProfile/'.concat(loginStore.getCurrentUser().uid))">
      <el-popover
          placement="bottom"
          width="200"
          trigger="hover">
        <template #reference>
          <el-avatar  :src="loginStore.getCurrentUser().avatarUrl" size="large"></el-avatar>
        </template>
        <div style="text-align: center">
          <el-text type="success" size="large" tag="b">{{ loginStore.getCurrentUser().nickname}}</el-text><br>
          <el-text type="warning" tag="i">{{ loginStore.getCurrentUser().signature}}</el-text><br>
          <el-divider/>
          <el-button class="transparent-button" icon="Flag" @click="goToHomePage">首页</el-button>
          <el-button class="transparent-button" icon="Edit" @click="goToPostArticlePage">写文章</el-button>
          <el-divider/>
          <el-button class="transparent-button" icon="Remove" @click="loginStore.exitLogin();ElMessage('退出登录成功')">退出登录</el-button>
        </div>
      </el-popover>
    </el-menu-item>

    <el-menu-item v-else @click="toggleLoginDialog">
      <el-icon><user/></el-icon>
      <span>登录</span>
    </el-menu-item>

    <div class="flex-grow" />

    <el-menu-item>
      <div class="search-container">
      <el-input
          class="search-input"
          v-model="queryKeyword"
          placeholder="don't repeat yourself"
          :suffix-icon="Search"
          @keyup.enter="router.push('/query/'.concat(queryKeyword))"
          @clear="queryKeyword=''"
          clearable
      />
      <el-button class="search-button" @click="router.push('/query/'.concat(queryKeyword))" v-show="queryKeyword!=''" style="vertical-align: 0px">查询</el-button>
      </div>
    </el-menu-item>

    <el-menu-item @click="goToHomePage">
      <el-icon><flag/></el-icon>
      <span>首页</span>
    </el-menu-item>
    <el-menu-item @click="goToPostArticlePage">
      <el-icon><edit/></el-icon>
      <span>写文章</span>
    </el-menu-item>

    <el-menu-item @click="">
      <el-link href="https://github.com/xuchi021105/bbs">
        <font-awesome-icon :icon="['fab', 'github']" style="height: 30px;margin-top: 14px" @click=""/>
      </el-link>
    </el-menu-item>
<!--    <el-menu-item>-->

<!--    </el-menu-item>-->

  </el-menu>


</template>

<style scoped>

.flex-grow {
  flex-grow: 1;
}

.search-container {
  align-items: center;
}

/*其中的--el-input-border-radius是el-input__wrapper中的用var定义的变量,可以通过这种方式覆盖*/
.search-input {
  width: 250px;
  height: 45px;
  padding-bottom: 8px;
  margin-right: 5px;
  --el-input-border-radius: 23px;
  --el-input-bg-color: #f5f7fa;
}

.search-button {
  padding: 10px 20px;
  border-radius: 20px;
  background-color: #409eff;
  color: #fff;
  transition: background-color 0.3s ease;
}

.search-button:hover {
  background-color: #66b1ff;
}

.transparent-button {
  background-color: transparent;
  color: #409eff;
  border-color: transparent;
}
</style>