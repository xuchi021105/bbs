<script setup>
import LoginDialog from "../components/LoginDialog.vue";
import {useRouter} from "vue-router";
import http from "../utils/request.js";
import {ElMessage} from "element-plus";
import {onMounted, ref} from "vue";
import Carousel from "../components/Carousel.vue";

const router = useRouter()

const articles = ref([])

// 放在onMounted里面和外面效果差不多
onMounted(() =>{
// 从这里回来开始写,GET都要改为POST,我的评价是很无语...改吧...
  http.post("/query/articleByKeyword",  {keyword:""})
      .then(function (response){
        articles.value = response.data.data
      }).catch(function (error){
    ElMessage("出错")
    console.log(error)
  })
})


//
// const articles = [
//   {
//     articleId: 0,
//     nickname: 'xuchi',
//     uid: 0,
//     avatarUrl: 'https://static.runoob.com/images/demo/demo1.jpg',
//     title:"title",
//   },
//   {
//     articleId: 0,
//     nickname: 'xuchi',
//     uid: 0,
//     avatarUrl: 'https://static.runoob.com/images/demo/demo1.jpg',
//     title:"title",
//   },
//   {
//     articleId: 0,
//     nickname: 'xuchi',
//     uid: 0,
//     avatarUrl: 'https://static.runoob.com/images/demo/demo1.jpg',
//     title:"title",
//   },
// ]


</script>

<template>
    <div style="padding: 25px">
      <el-row :gutter="20">
        <el-col :span="10">
          <carousel/>
        </el-col>
        <el-col :span="14">
          <el-card style="border-radius: 25px;height: 98%;">
            <span style="font-size: 25px;">
              官网超链接
            </span>
            <el-divider/>
<!--            这里最好用for循环遍历一个js数组或者对象,但是懒,就算了,直接填链接得了-->
            <el-row>
              <el-col :span="8">
                <a href="https://www.vim.org/">vim</a>
              </el-col>
              <el-col :span="8">
                <a href="https://vuejs.org/">vue</a>
              </el-col>
              <el-col :span="8">
                <a href="https://go.dev/">go</a>
              </el-col>
            </el-row>
            <el-row>
              <el-col :span="8">
                <a href="https://vitejs.dev/">vite</a>
              </el-col>
              <el-col :span="8">
                <a href="https://element-plus.org/">element-plus</a>
              </el-col>
              <el-col :span="8">
                <a href="https://gin-gonic.com/">gin</a>
              </el-col>
            </el-row>
            <el-row>
              <el-col :span="8">
                <a href="https://pinia.vuejs.org/">pinia</a>
              </el-col>
              <el-col :span="8">
                <a href="https://router.vuejs.org/">vue-router</a>
              </el-col>
              <el-col :span="8">
                <a href="https://gorm.io/">gorm</a>
              </el-col>
            </el-row>
            <el-row>
              <el-col :span="8">
                <a href="https://animate.style/">animate.css</a>
              </el-col>
              <el-col :span="8">
                <a href="https://fontawesome.com/">fontawesome</a>
              </el-col>
              <el-col :span="8">
                <a href="https://jwt.io/">jwt</a>
              </el-col>
            </el-row>
            <el-row>
              <el-col :span="8">
                <a href="https://axios-http.com/">axios</a>
              </el-col>
              <el-col :span="8">
                <a href="https://github.com/code-farmer-i/vue-markdown-editor">vue-markdown-editor</a>
              </el-col>
              <el-col :span="8">
                <a href="https://github.com/aliyun/aliyun-oss-go-sdk">aliyun-oss-go-sdk</a>
              </el-col>
            </el-row>
            <el-row>
              <el-col :span="8">
                <a href="https://www.jetbrains.com/webstorm/">jetbrains/webstorm</a>
              </el-col>
              <el-col :span="8">
                <a href="https://github.com/surmon-china/vue-awesome-swiper">vue-awesome-swiper</a>
              </el-col>
              <el-col :span="8">
                <a href="https://code.visualstudio.com/">vscode</a>
              </el-col>
            </el-row>
            <el-divider></el-divider>
            <el-row>
              <el-col :span="8">
                <a href="http://github.com">github</a>
              </el-col>
              <el-col :span="8">
                <a href="https://www.kali.org/">kali</a>
              </el-col>
              <el-col :span="8">
                <a href="https://kotlinlang.org/">kotlin</a>
              </el-col>
            </el-row>
            <el-row>
              <el-col :span="8">
                <a href=""></a>
              </el-col>
              <el-col :span="8">
                <a href=""></a>
              </el-col>
              <el-col :span="8">
                <a href=""></a>
              </el-col>
            </el-row>
          </el-card>
        </el-col>
      </el-row>
      <div>
  <!--      这里应该抽象出ArticleCard来,但是懒,就算了-->

  <!--      可以用@click.native而不是click(因为click事件在这个组件上被废弃了),但是这样也有问题,子组件的点击事件也会触发这个事件-->
        <el-card v-for="article in articles" style="border-radius: 20px;padding: 5px; margin: 15px; border-color: lightgrey;background-color: #faf7f7;" shadow="hover">
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
  <!--        <el-button @click="router.push('/article/'.concat(article.article_id))">查看文章</el-button>-->
  <!--        本来想法是点击头像就可以了,但是好像没有这个接口,我也懒,就以后再做吧-->
  <!--        <el-button @click="router.push('/userProfile/'.concat(article.uid))">查看用户</el-button>-->
        </el-card>
      </div>
    </div>
</template>
<style scoped>
</style>