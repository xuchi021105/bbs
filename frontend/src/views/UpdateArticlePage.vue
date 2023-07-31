<script setup>
import {onMounted, reactive, ref} from "vue";
import {useRoute, useRouter} from "vue-router";
import {useLoginStore} from "../stores/login.js";
import http from "../utils/request.js";
import {ElMessage} from "element-plus";

const article = reactive({
  title:'',
  content:'',
})

const router = useRouter()
const route = useRoute()
const loginStore = useLoginStore()

const editorHeight = ref((window.innerHeight - 175).toString().concat('px'))


const posted = ref(false)

function preUpdateArticle(){

  if (article.title.length == 0){
    ElMessage("标题不能为空")
    return
  }else if(article.content.length == 0){
    ElMessage("正文不能为空")
    return
  }
  posted.value = true
}

onMounted(()=>{

  window.addEventListener('resize', () =>{
    editorHeight.value = (window.innerHeight - 175).toString().concat('px')
  })

  http.post("/article/altAuth/getArticle",{
    article_id : Number(route.params.articleID)
  },{
    headers:loginStore.headers,
  }).then((response) =>{
    article.title = response.data.data.title
    article.content = response.data.data.content
  }).catch((e) =>{
    ElMessage("出错")
    console.log(e)
  })
})

function updateArticle(){
  http.put("/article/auth/article",{
    title : article.title,
    content : article.content,
    article_id : Number(route.params.articleID),
  },{
    headers:loginStore.headers,
  }).then((response) =>{
    router.push("/article/".concat(route.params.articleID))
    ElMessage("修改文章成功")
  }).catch((e) =>{
    ElMessage("出错")
    console.log(e)
  })
}

</script>

<template>
<!--  <el-form :model="article">-->
<!--    <el-form-item label="title">-->
<!--      <el-input v-model="article.title" placeholder="please"/>-->
<!--    </el-form-item>-->
<!--    <el-form-item label="content">-->
<!--      <el-input type="textarea" v-model="article.content" placeholder="content" :rows="10"/>-->
<!--    </el-form-item>-->
<!--    <el-form-item>-->
<!--      <el-button type="primary" @click="updateArticle">-->
<!--        提交保存-->
<!--      </el-button>-->
<!--    </el-form-item>-->
<!--  </el-form>-->

  <el-dialog v-model="posted" title="update article">
    <span>
      您确认修改好了吗?
    </span>
    <template #footer>
      <el-button type="danger" @click="posted=false" style="border-radius: 25px">
        没有,我再想想
      </el-button>
      <el-button type="primary" @click="posted=false;updateArticle()" style="border-radius: 25px">
        确认
      </el-button>
    </template>
  </el-dialog>

  <div style="overflow: visible;">

    <el-input v-model="article.title" placeholder="请输入标题" style="height: 50px;font-size: 30px;"/>
    <v-md-editor  v-model="article.content" placeholder="请输入正文" :height="editorHeight"></v-md-editor>
    <div style="height: 5px"></div>
    <!--    <el-form :model="article" style="padding: 0;margin: 0">-->
    <!--      <el-form-item >-->
    <!--        <el-input v-model="article.title" placeholder="请输入标题" style="height: 50px;font-size: 30px;"/>-->
    <!--      </el-form-item>-->
    <!--      <el-form-item >-->
    <!--        &lt;!&ndash;      <el-input type="textarea" v-model="article.content" placeholder="content" :rows="10"/>&ndash;&gt;-->
    <!--        <v-md-editor  v-model="article.content" placeholder="请输入正文" :height="editorHeight"></v-md-editor>-->
    <!--      </el-form-item>-->
    <!--    </el-form>-->
    <el-menu style="position: sticky;bottom: 0;" mode="horizontal" :ellipsis="false" text-color="white">
      <div class="flex-grow"/>
      <el-menu-item @click="preUpdateArticle" style="background-color: #409eff; border-radius: 15px;width: 75px">
        修改
      </el-menu-item>
      <el-menu-item>

      </el-menu-item>
      <el-menu-item>

      </el-menu-item>
      <el-menu-item>

      </el-menu-item>
      <el-menu-item>

      </el-menu-item>
    </el-menu>

  </div>

</template>

<style scoped>

.flex-grow {
  flex-grow: 2;
}

</style>