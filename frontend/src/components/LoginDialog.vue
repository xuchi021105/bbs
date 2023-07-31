<script setup>
import {reactive, ref} from "vue";
import http from "../utils/request.js";
import {ElMessage, ElMessageBox} from "element-plus";
import {useLoginStore} from "../stores/login.js";

defineProps(["showLoginDialog"])
const emit = defineEmits(["closeLoginDialog"])

const tab = ref("login")

const loginStore = useLoginStore()

// form中的不能用ref,要用reactive
// 这是文档里面说的,我也不知道为什么
const loginInfo = reactive({
  uid: "", // 这个得是Number类型
  password:"",
})
const registerInfo = reactive({
  nickname : "",
  password: "",
})
const confirmPassword = ref("")


async function login(){
  try{

    const response = await http.post("/login", {
      uid: parseInt(loginInfo.uid),

      password: loginInfo.password,
    })

    loginInfo.password = ""

    if(response.data.code != 0){
      ElMessage(response.data.msg)
    }
    else{
      loginStore.isLogin = true
      loginStore.setCurrentUser(response.data.data)

      emit("closeLoginDialog")

      ElMessage("登录成功")

    }
  }catch (e){
    ElMessage("请求有错")
    console.log(e)
  }

}

function register(){
  http.post("/register", registerInfo).then((response) =>{

    registerInfo.nickname = ""
    registerInfo.password = ""
    confirmPassword.value = ""

    if(response.data.msg == 'success'){
      tab.value = "login"
      ElMessageBox.alert("您的用户id为: ".concat(String(response.data.data.uid)), "注册成功")
    }

  }).catch((e) =>{
    ElMessage("请求有错")
    console.log(e)
  })

}

const numberRe = new RegExp("^[1-9][0-9]*$")

</script>

<template>
<!--  props不要双向绑定,而v-model是双向绑定(用函数实现的),所以这里不能用v-model,而是要单向绑定-->
  <el-dialog  :model-value="showLoginDialog" title="login or register" style="border-radius: 25px" @closed="$emit('closeLoginDialog')">

    <el-tabs v-model="tab">
      <el-tab-pane label="登录" name="login">
        <el-card style="border-radius: 25px;margin: 10px">
          <el-form :model="loginInfo" label-width="100px">
            <el-form-item label="用户ID">
              <el-input v-model="loginInfo.uid" placeholder="uid" style="--el-input-border-radius: 25px"/>
              <el-text v-if="!(loginInfo.uid.length == 0 || numberRe.test(loginInfo.uid)) " type="danger">请输入合法数字</el-text>
            </el-form-item>
            <el-form-item label="密码">
              <el-input v-model="loginInfo.password" type="password" placeholder="password" show-password style="--el-input-border-radius: 25px"/>
            </el-form-item>
          </el-form>
        </el-card>

        <div style="width: max-content;margin: auto">
          <el-button type="primary" @click="login" style="border-radius: 25px;" :disabled="loginInfo.password.length == 0 || !(loginInfo.uid.length == 0 || numberRe.test(loginInfo.uid))">
            登录
          </el-button>
        </div>
      </el-tab-pane>
      <el-tab-pane label="注册" name="register">
        <el-card style="border-radius: 25px;margin: 10px">
          <el-form :model="registerInfo" label-width="100px">
            <el-form-item label="昵称">
              <el-input v-model="registerInfo.nickname" placeholder="nickname" style="--el-input-border-radius: 25px"/>
            </el-form-item>
            <el-form-item label="密码">
              <el-input show-password v-model="registerInfo.password" placeholder="password" style="--el-input-border-radius: 25px"/>
            </el-form-item>
            <el-form-item label="确认密码">
              <el-input show-password v-model="confirmPassword" placeholder="confirm password" style="--el-input-border-radius: 25px"/>

              <el-text v-show="(registerInfo.password != confirmPassword)" type="danger">和之前密码不一致</el-text>

            </el-form-item>
          </el-form>
        </el-card>
        <div style="width: max-content;margin: auto">
          <el-button type="primary" :disabled="registerInfo.nickname.length == 0 || registerInfo.password.length == 0 || confirmPassword.length == 0 || registerInfo.password != confirmPassword" @click="register" style="border-radius: 25px">
            注册
          </el-button>
        </div>
      </el-tab-pane>
    </el-tabs>
  </el-dialog>
</template>
