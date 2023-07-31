<script setup>
import {onMounted, ref} from 'vue'
import {useRoute, useRouter} from "vue-router";
import {useLoginStore} from "../stores/login.js";
import http from "../utils/request.js";
import {Search} from "@element-plus/icons-vue";
import {ElMessage} from "element-plus";

const tab = ref("article")

const router = useRouter()
const route = useRoute()
const loginStore = useLoginStore()

const articles = ref([])

const users = ref([])

const userWithID = ref({
})

const showUserWithID = ref(false)

const showArticle = ref(false)

const showUser = ref(false)

onMounted(()=>{
  http.post("/query/articleByKeyword", {
    keyword : route.params.keyword,
  },).then((response) =>{
    if(response.data.data == null){
      showArticle.value = false
    }else{
      articles.value = response.data.data
      showArticle.value = true
    }
  })
})

function queryByKeyword(){
  switch (tab.value) {
    case 'article':{

      http.post("/query/articleByKeyword", {
        keyword : queryKeyword.value,
      },).then((response) =>{
        if(response.data.data == null){
          showArticle.value = false
        }else {
          articles.value = response.data.data
          showArticle.value = true
        }
      }).catch((e) =>{
        ElMessage("出错")
        console.log(e)
      })
      break
    }
    case 'user':{
      http.post("/query/userByKeyword", {
        keyword : queryKeyword.value,
      },).then((response) =>{
        if(response.data.data == null){
          showUser.value = false
        }else {
          users.value = response.data.data
          showUser.value = true
        }
      }).catch((e) =>{
        ElMessage("出错")
        console.log(e)
      })

      break
    }
    case  'userID':{
      http.post("/query/userByUid", {
        keyword : queryKeyword.value,
      },).then((response) =>{

        userWithID.value = response.data.data
        // console.log(response.data.msg)
        if(response.data.msg == 'success'){
          showUserWithID.value = true
        }else {
          showUserWithID.value = false
        }
      }).catch((e) =>{
        ElMessage("出错")
        console.log(e)
      })

      break
    }
  }
}

// const users = [
//   {
//     avatarUrl: "https://static.runoob.com/images/demo/demo1.jpg",
//     uid: 0,
//     nickname: 'xuchi',
//     signature: 'signature'
//   },
//   {
//     avatarUrl: "https://static.runoob.com/images/demo/demo1.jpg",
//     uid: 0,
//     nickname: 'xuchi',
//     signature: 'signature'
//   },
//   {
//     avatarUrl: "https://static.runoob.com/images/demo/demo1.jpg",
//     uid: 0,
//     nickname: 'xuchi',
//     signature: 'signature'
//   },
// ]
const queryKeyword = ref("")
</script>

<template>

<!--  <el-input-->
<!--      v-model="queryKeyword"-->
<!--      placeholder="don't repeat yourself"-->
<!--      :prefix-icon="Search"-->
<!--      style="width: 200px"-->
<!--  />-->
<!--  <el-button @click="queryByKeyword" v-show="queryKeyword!=''">查询</el-button>-->

  <el-row>
    <el-col :span="8" />
    <el-col :span="8">
      <div class="search-container">
        <el-input
            class="search-input"
            v-model="queryKeyword"
            placeholder="don't repeat yourself"
            :suffix-icon="Search"
            @keyup.enter="queryByKeyword"
            @clear="queryKeyword=''"
            clearable
        />
        <el-button class="search-button" @click="queryByKeyword" v-show="queryKeyword!=''" style="vertical-align: 4px">查询</el-button>
      </div>
    </el-col>

    <el-col :span="8" />
  </el-row>

  <div style="padding: 10px;">
    <el-tabs v-model="tab" class="demo-tabs">
      <el-tab-pane label="Article" name="article">

        <el-empty v-if="!showArticle" description="没有找到文章"/>

        <el-card v-else v-for="article in articles" style="border-radius: 20px;padding: 5px; margin: 15px; border-color: lightgrey;background-color: #f7f7f7;" shadow="hover">
          <el-row>
            <el-link @click="router.push('/article/'.concat(article.article_id))" >
              <span style="font-size: 19px">
                title : {{ article.title }}
              </span>
            </el-link>
          </el-row>
          <el-divider style="margin-top: 2px;margin-bottom: 10px"/>
          <div>
            <el-avatar :size="45" :src="article.avatar_url" @click.native="router.push('/userProfile/'.concat(article.uid))"/>
            <el-link @click="router.push('/userProfile/'.concat(article.uid))" style="vertical-align: 15px">
              <span >
                author : {{ article.nickname }}
              </span>
            </el-link>
          </div>
        </el-card>


<!--        <el-card v-for="article in articles" >-->
<!--          <el-avatar :size="50" :src="article.avatar_url" />-->
<!--          <el-text>{{ article.nickname }}</el-text>-->
<!--          <br>-->
<!--          <el-text>{{ article.title }}</el-text>-->
<!--          <el-button @click="router.push('/article/'.concat(article.article_id))">查看文章</el-button>-->
<!--          &lt;!&ndash;        本来想法是点击头像就可以了,但是好像没有这个接口,我也懒,就以后再做吧&ndash;&gt;-->
<!--          <el-button @click="router.push('/userProfile/'.concat(article.uid))">查看用户</el-button>-->
<!--        </el-card>-->

      </el-tab-pane>

      <el-tab-pane label="User" name="user">

<!--        <el-card v-for="user in users" >-->
<!--          <el-avatar :size="50" :src="user.avatar_url"></el-avatar>-->
<!--          <el-text>{{user.nickname}}</el-text>-->
<!--          <br>-->
<!--          <el-text>{{user.signature}}</el-text>-->
<!--          <el-button @click="router.push('/userProfile/'.concat(user.uid))">查看用户</el-button>-->
<!--        </el-card>-->

        <el-empty v-if="!showUser" description="没有找到用户"/>
        <el-card v-else v-for="user in users" style="border-radius: 20px;padding: 5px; margin: 15px; border-color: lightgrey;background-color: #f7f7f7;" shadow="hover">
          <el-row>
            <el-avatar :size="50" :src="user.avatar_url" @click.native="router.push('/userProfile/'.concat(user.uid))"/>
            <div>
              <el-row style="margin-top: 5px">
                <el-link @click="router.push('/userProfile/'.concat(staredArticle.uid))" style="margin-left: 5px">
                  <span style="font-size: 20px">
                    {{ user.nickname }}
                  </span>
                </el-link>
              </el-row>
              <el-row>
                <span style="font-size: 13px;margin-left: 5px">{{user.signature}}</span>
              </el-row>
            </div>
          </el-row>
        </el-card>


      </el-tab-pane>

      <el-tab-pane label="UserID" name="userID">

        <el-card v-if="showUserWithID" style="border-radius: 20px;padding: 5px; margin: 15px; border-color: lightgrey;background-color: #f7f7f7;" shadow="hover">
          <el-row>
            <el-avatar :size="50" :src="userWithID.avatar_url" @click.native="router.push('/userProfile/'.concat(userWithID.uid))"/>
            <div>
              <el-row style="margin-top: 5px">
                <el-link @click="router.push('/userProfile/'.concat(userWithID.uid))" style="margin-left: 5px">
                  <span style="font-size: 20px">
                    {{ userWithID.nickname }}
                  </span>
                </el-link>
              </el-row>
              <el-row>
                <span style="font-size: 13px;margin-left: 5px">{{userWithID.signature}}</span>
              </el-row>
            </div>
          </el-row>
        </el-card>

        <el-empty v-else description="没有找到用户"/>

<!--        <el-card>-->
<!--          <el-avatar :size="50" :src="userWithID.avatar_url"></el-avatar>-->
<!--          <el-text>{{userWithID.nickname}}</el-text>-->
<!--          <br>-->
<!--          <el-text>{{userWithID.signature}}</el-text>-->
<!--        </el-card>-->

      </el-tab-pane>

    </el-tabs>
  </div>
</template>
<style scoped>
.demo-tabs > .el-tabs__content {
  padding: 32px;
  color: #6b778c;
  font-size: 32px;
  font-weight: 600;
}

.search-container {
  align-items: center;
}

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
 : width: auto; transition: background-color 0.3s ease;
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