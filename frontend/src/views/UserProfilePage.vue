<script setup>
import {computed, onMounted, ref} from "vue";
import {Star} from "@element-plus/icons-vue";
import {useLoginStore} from "../stores/login.js";
import {useRoute, useRouter} from "vue-router";
import {ElMessage} from "element-plus";
import http from "../utils/request.js";
import LoginDialog from "../components/LoginDialog.vue";

const tab = ref("articles")

const loginStore = useLoginStore()

const router = useRouter()
const route = useRoute()

const user = ref({})

const tableData = computed(() =>{
  return [user.value]
})

const articles = ref([])
const staredArticles = ref([])

// 复用页面 setup函数不会再执行
onMounted(() =>{

  // console.log(loginStore.headers)

  http.post('/user/altAuth/info', {uid: Number(route.params.uid)},{
    headers: loginStore.headers,
  }).then((response) =>{
    // 这里还要错误处理的,但是没时间了,就先这样

    if(response.data.data == null){

    }else {
      user.value = response.data.data
    }
  }
  ).catch(
      (e) => {
        ElMessage('出错')
        console.log(e)
      }
  )

  http.post("/article/altAuth/articles", {uid:Number(route.params.uid)}, {
    headers: loginStore.headers
  }).then((response) =>{
    // 要对空情况做处理
    if(response.data.data == null){

    }else {
      articles.value = response.data.data.articles
      console.log(response.data.data.articles)
    }
  }).catch((error) =>{
    ElMessage('出错')
    console.log(error)
  })

  http.post("/article/altAuth/userStaredArticles", {uid:Number(route.params.uid)}, {
    headers: loginStore.headers
  }).then((response) =>{
    if (response.data.data == null){

    }else {
      staredArticles.value = response.data.data.stared_articles
      console.log(response.data.data.stared_articles)
    }
  }).catch((error) =>{
    ElMessage('出错')
    console.log(error)
  })

  http.post("/user/followers", {
    uid:Number(route.params.uid)
  }).then((response) =>{

    if (response.data.data == null){

    }else {
      followers.value = response.data.data
      console.log('followers', response.data.data)
    }
  }).catch((error) =>{
    ElMessage('出错')
    console.log(error)
  })

  http.post("/user/followMePersonList", {
    uid:Number(route.params.uid)
  }).then((response) =>{
    if (response.data.data == null){

    }else {
      fans.value = response.data.data
      console.log('fans', response.data.data)
    }
  }).catch((error) =>{
    ElMessage('出错')
    console.log(error)
  })
})

// const articles = [
//   {
//     articleId: 0,
//     nickname: 'xuchi',
//     uid: 0,
//     avatarUrl: 'https://static.runoob.com/images/demo/demo1.jpg',
//     title:"title",
//     releaseTime : 'now',
//     stared: false
//   },
//   {
//     articleId: 0,
//     nickname: 'xuchi',
//     uid: 0,
//     avatarUrl: 'https://static.runoob.com/images/demo/demo1.jpg',
//     title:"title",
//     releaseTime : 'now',
//     stared: true,
//   },
//   {
//     articleId: 0,
//     nickname: 'xuchi',
//     uid: 0,
//     avatarUrl: 'https://static.runoob.com/images/demo/demo1.jpg',
//     title:"title",
//     releaseTime : 'now',
//     stared: true,
//   },
// ]
//
// const staredArticles = [
//   {
//     articleId: 0,
//     nickname: 'xuchi',
//     uid: 0,
//     avatarUrl: 'https://static.runoob.com/images/demo/demo1.jpg',
//     title:"title",
//     signature: "signature",
//     releaseTime : 'now',
//   },
//   {
//     articleId: 0,
//     nickname: 'xuchi',
//     uid: 0,
//     avatarUrl: 'https://static.runoob.com/images/demo/demo1.jpg',
//     title:"title",
//     signature: "signature",
//     releaseTime : 'now',
//   },
//   {
//     articleId: 0,
//     nickname: 'xuchi',
//     uid: 0,
//     avatarUrl: 'https://static.runoob.com/images/demo/demo1.jpg',
//     title:"title",
//     signature: "signature",
//     releaseTime : 'now',
//   },
// ]

// const user = {
//   uid : 0,
//   avatarUrl: "https://static.runoob.com/images/demo/demo1.jpg",
//   nickname: 'xuchi',
//   signature: 'signature',
//   createAtTime: 'now',
//   followersNumber: 10,
//   fansNumber: 10,
//   articlesNumber:10,
//   staredArticlesNumber:10,
//   isMyself:false,
//   followed:false
// }

const followers = ref([])
const fans = ref([])

const showLoginDialog = ref(false)
async function followUser(){
  try{

    if(loginStore.isLogin == false){
      showLoginDialog.value = true
      return
    }

    const response = await http.post("/user/auth/followUid",{
      uid:Number(route.params.uid),
    },{
      headers : loginStore.headers
    })

    user.value.followed = true
    ElMessage("关注成功")

  }catch (e){
    ElMessage("出错")
    console.log(e)
  }

}

function cancelFollowUser(){

  if(loginStore.isLogin == false){
    showLoginDialog.value = true
    return
  }

  http.put("/user/auth/followUid",{
    uid:Number(route.params.uid),
  },{
    headers : loginStore.headers
  }).then((response) =>{

    user.value.followed = false
    ElMessage("取消关注成功")

  }).catch((e) =>{
    console.log(e)
  })
}

</script>

<template>
    <login-dialog :show-login-dialog="showLoginDialog" @close-login-dialog="showLoginDialog=false"/>
    <div style="padding: 5px;margin: 5px">
      <el-card shadow="never" style="border-color: white">
        <el-row>
          <el-avatar :size="100" :src="user.avatar_url" style="margin-right: 15px"/>
          <div>
            <el-row>
              <el-text tag="b" style="font-size: 44px; margin-top: 15px">{{user.nickname}}</el-text>
            </el-row>
            <el-row>
              <el-text>{{user.signature}}</el-text>
            </el-row>
          </div>
          <div style="margin-top: 35px;margin-left: 10px">
            <el-button v-if="user.is_myself" type="info" @click="router.push('/updateUserInfo/'.concat(route.params.uid))" style="border-radius: 25px;">修改信息</el-button>
            <el-button v-else-if="!user.followed" type="primary" @click="followUser" style="border-radius: 25px;">关注</el-button>
            <el-button v-else @click="cancelFollowUser" style="border-radius: 25px;">已关注</el-button>
          </div>
        </el-row>
        <el-text>uid: {{ user.uid}}</el-text>
        <br/>
        <el-text>注册时间: {{user.signInTime}}</el-text>
      </el-card>

      <el-table :data="tableData" style="width: 100%">
        <el-table-column prop="followers_number" label="关注数" />
        <el-table-column prop="my_followers_number" label="粉丝数" />
        <el-table-column prop="articles_number" label="文章数" />
        <el-table-column prop="stared_articles_number" label="收藏文章数" />
      </el-table>


      <el-tabs v-model="tab" class="demo-tabs">
        <el-tab-pane label="文章" name="articles">

          <el-empty v-if="articles.length==0" description="没有写文章"/>
          <el-card v-for="article in articles" style="border-radius: 20px;padding: 5px; margin: 15px; border-color: lightgrey;background-color: #fffff7;" shadow="hover">
            <el-link @click="router.push('/article/'.concat(article.article_id))" >
              <h1 style="font-size: 19px">
                title : {{ article.title }}
              </h1>
            </el-link>
<!--            <el-text>{{ article.title}}</el-text>-->
            <br>
            <el-text tag="i" style="vertical-align: 5px">{{article.release_time}}</el-text>
<!--            <el-button @click="router.push('/article/'.concat(article.article_id))">查看文章</el-button>-->

<!--            <el-button  v-if="article.stared" :icon="Star" type="primary">已收藏</el-button>-->
<!--            <el-button  v-else :icon="Star" >未收藏</el-button>-->

            <Transition mode="out-in"
                        enter-active-class="animate__animated animate__bounceIn"
                        leave-active-class="animate__animated animate__bounceOut"
                        style="margin-left: 10px;"
            >
              <font-awesome-icon v-if="!article.stared" :icon="['far', 'star']" style="" size="2xs"/>
              <font-awesome-icon v-else :icon="['fas', 'star']" style="color: rgb(200,204,39);" size="2xs"/>
            </Transition>

          </el-card>
        </el-tab-pane>

        <el-tab-pane label="收藏的文章" name="staredArticles">

          <el-empty v-if="staredArticles.length==0" description="没有收藏文章"/>
          <el-card v-for="staredArticle in staredArticles" style="border-radius: 20px;padding: 5px; margin: 15px; border-color: lightgrey;background-color: #fff7ff;" shadow="hover">
            <el-avatar :src="staredArticle.avatar_url" @click.native="router.push('/userProfile/'.concat(staredArticle.uid))"/>

            <el-link @click="router.push('/userProfile/'.concat(staredArticle.uid))" style="vertical-align: 15px;margin-left: 5px">
              <span style="font-size: 20px">
                {{ staredArticle.nickname }}
              </span>
            </el-link>
            <el-divider style="margin-top: 2px;margin-bottom: 10px"/>

<!--            <el-text>{{staredArticle.nickname}}</el-text>-->
<!--            <br>-->
<!--            <el-text>{{staredArticle.title}}</el-text>-->

            <el-link @click="router.push('/article/'.concat(staredArticle.article_id))" >
              <span style="font-size: 19px">
                title : {{ staredArticle.title }}
              </span>
            </el-link>

            <br>
            <el-text>{{staredArticle.release_time}}</el-text>

<!--            <el-button @click="router.push('/article/'.concat(staredArticle.article_id))">查看文章</el-button>-->
<!--            <el-button @click="router.push('/userProfile/'.concat(staredArticle.uid))">查看用户</el-button>-->
          </el-card>

        </el-tab-pane>
        <el-tab-pane label="关注" name="followers">

          <el-empty v-if="followers.length==0" description="没有关注的人"/>

          <el-card v-for="user in followers" style="border-radius: 20px;padding: 5px; margin: 15px; border-color: lightgrey;background-color: #f7ffff;" shadow="hover">
            <el-row>
              <el-avatar :size="50" :src="user.avatar_url" @click.native="router.push('/userProfile/'.concat(user.uid))"/>
              <div>
                <el-row style="margin-top: 5px">
                  <el-link @click="router.push('/userProfile/'.concat(user.uid))" style="margin-left: 5px">
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
<!--            <el-text>{{user.nickname}}</el-text>-->
<!--            <el-button @click="router.push('/userProfile/'.concat(user.uid))">查看用户</el-button>-->
          </el-card>

        </el-tab-pane>
        <el-tab-pane label="粉丝" name="fans" >

          <el-empty v-if="fans.length==0" description="暂时没有关注ta的人"/>

          <el-card v-for="user in fans" style="border-radius: 20px;padding: 5px; margin: 15px; border-color: lightgrey;background-color: #f7f7f7;" shadow="hover">
            <el-row>
              <el-avatar :size="50" :src="user.avatar_url" @click.native="router.push('/userProfile/'.concat(user.uid))"/>
              <div>
                <el-row style="margin-top: 5px">
                  <el-link @click="router.push('/userProfile/'.concat(user.uid))" style="margin-left: 5px">
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
<!--            <el-avatar :size="50" :src="user.avatar_url" @click.native="router.push('/userProfile/'.concat(user.uid))"></el-avatar>-->
<!--            <el-text>{{user.nickname}}</el-text>-->
<!--            <br>-->
<!--            <el-text>{{user.signature}}</el-text>-->
<!--            <el-button @click="router.push('/userProfile/'.concat(user.uid))">查看用户</el-button>-->
          </el-card>

        </el-tab-pane>
      </el-tabs>

    </div>
</template>

<style>
.demo-tabs > .el-tabs__content {
  padding: 10px;
  color: #6b778c;
  font-size: 32px;
  font-weight: 600;
}
</style>
