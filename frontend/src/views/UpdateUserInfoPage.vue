<script setup>

import {useRoute, useRouter} from "vue-router";
import {onMounted, reactive, ref} from "vue";
import http from "../utils/request.js";
import {ElMessage, genFileId } from "element-plus";
import {useLoginStore} from "../stores/login.js";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";

const router = useRouter()
const route = useRoute()

const user = reactive({nickname:"", signature:""})

onMounted(() =>{
  http.post("/user/altAuth/info", {
    uid: Number(route.params.uid),
  }).then((response) =>{
    user.nickname = response.data.data.nickname
    user.signature = response.data.data.signature

  }).catch((e) =>{
    ElMessage("出错")
    console.log(e)
  })
})

const loginStore = useLoginStore()
async function updateUserInfo(){
  try{
    const response = await http.put("/user/auth/info", user
    ,{
      headers: loginStore.headers,
    })
    router.push("/userProfile/".concat(String(route.params.uid)))
    ElMessage("修改成功")
  }catch (e){
    ElMessage("出错")
    console.log(e)
  }
}

// 实质上是一个控制el-upload的有着各种方法的对象
const upload = ref()

const handleExceed= (files) => {
  upload.value.clearFiles()
  const file = files[0]
  file.uid = genFileId()
  upload.value.handleStart(file)

  // console.log(upload.value)
}

const showUploadButtonFlag = ref(false)

const submitUpload = () => {
  try{
    upload.value.submit()
    router.push("/userProfile/".concat(String(route.params.uid)))
    ElMessage("上传成功")
  }catch (e) {
    ElMessage("出错")
    console.log(e)
  }

}

function handleRemove(uploadFile, uploadFiles){
  showUploadButtonFlag.value = false
}

// const user = {
//   nickname : '',
//   signature: '',
// }

// 这个代码用于客户端验证文件类型,因为这个会带rawFile
// const beforeAvatarUpload = (rawFile) => {
//   if (rawFile.type !== 'image/jpeg') {
//     ElMessage.error('Avatar picture must be image format!')
//     return false
//   } else if (rawFile.size / 1024 / 1024 > 2) {
//     ElMessage.error('Avatar picture size can not exceed 2MB!')
//     return false
//   }
//   return true
// }

</script>


<template>
  <el-image src="/images/banner/updateBanner.jpg" style="width: 100%"/>

  <el-card style="margin: 10px;border-radius: 24px">

    <span>
      修改你的个人信息
    </span>
    <el-divider/>

    <el-tabs tab-position="left">

      <el-tab-pane label="UserInfo">
        <div style="width: max-content;margin: auto">

          <el-card style="width: 380px;border-radius: 25px;margin: 15px">
            <el-form :model="user">
              <el-form-item label="nickname">
                <div style="width: 250px">
                  <el-input v-model="user.nickname" placeholder="please" />
                </div>
              </el-form-item>
              <el-form-item label="signature">
                <div style="width: 250px">
                  <el-input v-model="user.signature" placeholder="content" type="textarea" autosize resize="none"/>
                </div>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="updateUserInfo" style="border-radius: 25px;margin-left: 130px">
                  修改
                </el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </div>
      </el-tab-pane>

      <el-tab-pane label="Avatar">
        <!--      这里传文件,做文件处理-->
        <!--      不能设置headers-->
        <!--      <form action="127.0.0.1:8080/user/auth/avatarPicture" method="post" enctype="multipart/form-data" >-->
        <!--        <input type="file" name="f1">-->
        <!--        <input type="submit" value="上传">-->
        <!--      </form>-->

<!--        list-type用picture-card也挺好的-->
        <el-upload
            ref="upload"
            action="http://127.0.0.1:8080/user/auth/avatarPicture" // 这里应该也用axios的变量来表示,不然每次都要改这里,但是懒,就算了 其中127.0.0.1:8080表示用来测试的本地的后端ip地址
            :headers="loginStore.headers"
            :limit="1"
            :on-exceed="handleExceed"
            :auto-upload="false"
            :on-error="() => ElMessage('上传失败')"
            :on-remove="handleRemove"
            :on-change="() => showUploadButtonFlag = true"
            accept="image/*"
            list-type="picture"
            drag
        >
<!--          <template #trigger>-->
<!--            <el-button type="primary">select file</el-button>-->
<!--          </template>-->

          <font-awesome-icon :icon="['fas', 'upload']" size="2xl" style=""/>
          <br>
          <br>
          <span>
            drop file here or <em style="color: #409eff">click to upload</em>
          </span>

          <template #tip>
            <el-text tag="i" >
              jpg/png files with a size less than 2MB
            </el-text>
          </template>

        </el-upload>

        <div v-if="showUploadButtonFlag" style="margin: auto; width: max-content">
          <el-button type="success" @click="submitUpload" style="border-radius: 25px;">
            upload to server
          </el-button>
        </div>

<!--        <el-button @click="console.log()"/>-->
      </el-tab-pane>

    </el-tabs>

  </el-card>

</template>

