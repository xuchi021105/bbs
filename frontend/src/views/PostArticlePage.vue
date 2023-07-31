<script setup>
import {onMounted, reactive, ref} from "vue";
import {useLoginStore} from "../stores/login.js";
import {ElMessage} from "element-plus";
import http from "../utils/request.js";
import {useRoute, useRouter} from "vue-router";

// form要用reactive
const article = reactive({
  title:'',
  content:'',
})

const loginStore = useLoginStore()

const router = useRouter()
const route = useRoute()

async function postArticle(){
  try{
    const response = await http.post("/article/auth/article", article, {
      headers: loginStore.headers
    })
    await router.push("/userProfile/".concat(loginStore.getCurrentUser().uid.toString()))
    ElMessage("发布成功")
  }catch (e){
    ElMessage("出错")
    console.log(e)
  }

}

const editorHeight = ref((window.innerHeight - 175).toString().concat('px'))

onMounted(() =>{
  window.addEventListener('resize', () =>{
    editorHeight.value = (window.innerHeight - 175).toString().concat('px')
  })
})

const posted = ref(false)

function prePostArticle(){

  if (article.title.length == 0){
    ElMessage("标题不能为空")
    return
  }else if(article.content.length == 0){
    ElMessage("正文不能为空")
    return
  }
  posted.value = true
}

</script>

<template>
  <el-dialog v-model="posted" title="post article">
    <span>
      您确认要发布文章吗?
    </span>
    <template #footer>
      <el-button type="danger" @click="posted=false" style="border-radius: 25px">
        没有,我再想想
      </el-button>
      <el-button type="primary" @click="posted=false;postArticle()" style="border-radius: 25px">
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
      <el-menu-item @click="prePostArticle" style="background-color: #409eff; border-radius: 15px;width: 75px">
        发布
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