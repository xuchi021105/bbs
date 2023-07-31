<script setup>
import {onMounted, ref} from "vue";
import {Star, StarFilled} from "@element-plus/icons-vue";
import http from "../utils/request.js";
import {useRoute, useRouter} from "vue-router";
import {ElMessage} from "element-plus";
import {useLoginStore} from "../stores/login.js";
import {FontAwesomeIcon} from "@fortawesome/vue-fontawesome";
import LoginDialog from "../components/LoginDialog.vue";

const COMMENTTYPE = 'comment'
const REPLYTYPE = 'reply'

const router = useRouter()
const route = useRoute()
const loginStore = useLoginStore()

const articleInfo = ref({})
const comments = ref([])

onMounted(() =>{
  http.post("/article/altAuth/getArticle", {article_id : Number(route.params.articleID)},{
    headers: loginStore.headers,
  }).then((response) =>{
    articleInfo.value = response.data.data
    console.log(response.data.data)
  }).catch((e) =>{
    ElMessage("出错")
    console.log(e)
  })

  http.post('/comment/getCommentArticle', {article_id: Number(route.params.articleID)},{
    headers: loginStore.headers,
  }).then((response) =>{
    if(response.data.data != null){
      for(let index = 0; index < response.data.data.length; index++) {
        response.data.data[index].offset = index
      }
      comments.value = response.data.data
      console.log(response.data.data)
    }else{

    }
  }).catch((e) =>{
    ElMessage("获取评论出错")
    console.log(e)
  })

})

// const commentDialog = ref(false)

const commentInput = ref('')

const replyDialog = ref(false)

const replyInput = ref('')

const commentArticle = () =>{
  console.log({article_id: Number(route.params.articleID),content:commentInput.value})

  if(loginStore.isLogin == false){
    showLoginDialog.value = true
    return
  }

  http.post("/comment/auth/commentArticle", {article_id: Number(route.params.articleID),content:commentInput.value},{
    headers: loginStore.headers,
  }).then((response)=>{
    ElMessage("评论成功")
    comments.value.push(
        {
          avatar_url : loginStore.getCurrentUser().avatarUrl,
          comment_id : response.data.data.comment_id,
          content : commentInput.value,
          nickname : loginStore.getCurrentUser().nickname,
          release_time : response.data.data.release_time,
          replys: [],
          uid : loginStore.getCurrentUser().uid,
        }
    )
    commentInput.value = ''
  }).catch((e) =>{
    console.log(e)
    ElMessage("出错")
  })
}

const commentID = ref(0)

const replyType = ref("")

const toUid = ref(0)

const toUserNickname = ref('')

const commentOffset = ref('0')
const replyComment = () =>{

  if(loginStore.isLogin == false){
    showLoginDialog.value = true
    return
  }

  http.post("/reply/auth/replyComment", {
    comment_id :commentID.value,
    content: replyInput.value,
    reply_type:replyType.value, // reply_type优先级更高
    to_uid:toUid.value,
  },{
    headers: loginStore.headers,
  }).then((response) =>{
    ElMessage("回复成功")
    // console.log(comments.value[commentOffset])
    comments.value[commentOffset.value].replys.push(
        {
          content: replyInput.value,
          from_avatar_url: loginStore.getCurrentUser().avatarUrl,
          from_uid: loginStore.getCurrentUser().uid,
          from_user_nickname: loginStore.getCurrentUser().nickname,
          release_time: response.data.data.release_time,
          reply_id: response.data.data.reply_id,
          reply_type: replyType.value,
          to_uid: toUid.value,
          to_user_nickname: toUserNickname.value,
        }
    )
  }).catch((e) =>{
    ElMessage("出错")
    console.log(e)
  })
}

function starArticle(){

  if(loginStore.isLogin == false){
    showLoginDialog.value = true
    return
  }

  http.post("/article/auth/star",{
    article_id : Number(route.params.articleID)
  },{
    headers : loginStore.headers
  }).then((response) =>{
    ElMessage("收藏成功")
    articleInfo.value.stared = true
  }).catch((e) =>{
    ElMessage("出错")
    console.log(e)
  })

}

function cancelStaredArticle(){

  if(loginStore.isLogin == false){
    showLoginDialog.value = true
    return
  }

  http.put("/article/auth/star",{
    article_id : Number(route.params.articleID)
  },{
    headers : loginStore.headers
  }).then((response) =>{
    ElMessage("取消收藏")
    articleInfo.value.stared = false
  }).catch((e) =>{
    ElMessage("出错")
    console.log(e)
  })
}


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

    articleInfo.value.followed = true
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

    articleInfo.value.followed = false
    ElMessage("取消关注成功")

  }).catch((e) =>{
    console.log(e)
  })
}

// const articleInfo = {
//   isMyArticle: true,
//   title:"title",
//   content : "content",
//   stared: "stared",
//   uid: 100,
//   nickname: 'xuchi',
//   signature: 'signature',
//   avatarUrl: "https://static.runoob.com/images/demo/demo1.jpg",
//   followed: false,
//   releaseTime: 'now'
// }
//

// const comments = [
//   {
//     uid:0,
//     commentID:0,
//     content:"content",
//     avatarUrl:"https://static.runoob.com/images/demo/demo1.jpg",
//     nickname:"xuchi",
//     releaseTime:"now",
//     replys:[
//       {
//         replyType:"comment",
//         fromUID:0,
//         content:'content',
//         fromAvatarUrl:"https://static.runoob.com/images/demo/demo1.jpg",
//         replyID:1,
//         fromUserNickname:'xuchi',
//         releaseTime:"now",
//       },
//       {
//         replyType:"reply",
//         fromUID:0,
//         content:'content',
//         fromAvatarUrl:"https://static.runoob.com/images/demo/demo1.jpg",
//         replyID:1,
//         fromUserNickname:'xuchi',
//         releaseTime:"now",
//         // 后面这两个是reply特有的属性
//         toUID:1,
//         toUserNickname:'xuchi',
//       },
//     ],
//   },
//   {
//     uid:0,
//     commentID:0,
//     content:"content",
//     avatarUrl:"https://static.runoob.com/images/demo/demo1.jpg",
//     nickname:"xuchi",
//     releaseTime:"now",
//     replys:[
//       {
//         replyType:"comment",
//         fromUID:0,
//         content:'content',
//         fromAvatarUrl:"https://static.runoob.com/images/demo/demo1.jpg",
//         replyID:1,
//         fromUserNickname:'xuchi',
//         releaseTime:"now",
//       },
//       {
//         replyType:"reply",
//         fromUID:0,
//         content:'content',
//         fromAvatarUrl:"https://static.runoob.com/images/demo/demo1.jpg",
//         replyID:1,
//         fromUserNickname:'xuchi',
//         releaseTime:"now",
//         // 后面这两个是reply特有的属性
//         toUID:1,
//         toUserNickname:'xuchi',
//       },
//     ],
//   },
// ]

// const markdownText = ref('# header')

// const stared = ref(false)

// const test = ref('')

const showLoginDialog = ref(false)

</script>

<template>
    <login-dialog :show-login-dialog="showLoginDialog" @close-login-dialog="showLoginDialog=false"></login-dialog>
    <div>
      <el-row :gutter="12">
        <el-col :span="8">
          <el-card shadow="never" style="border-color: #ffffff">
            <el-row style="margin-top: 8px; ">
              <el-avatar  :src="articleInfo.avatar_url" @click.native="router.push('/userProfile/'.concat(articleInfo.uid))" size="large"/>
              <div style="margin-left: 5px">
                <span style="margin-right: 8px">
                  {{articleInfo.nickname}}
                </span>
                <div v-if="!articleInfo.is_my_article">
                  <el-button style="border-radius: 25px" v-if="!articleInfo.followed" type="primary" @click="followUser">关注</el-button>
                  <el-button style="border-radius: 25px" v-else @click="cancelFollowUser">已关注</el-button>
                </div>

                <br>
                <el-text>
                  {{articleInfo.signature}}
                </el-text>
              </div>
            </el-row>
<!--            <el-button @click="router.push('/userProfile/'.concat(articleInfo.uid))">去ta的主页</el-button>-->
          </el-card>
        </el-col>
        <el-col :span="16">
          <el-card shadow="never" style="border-color: #ffffff">
              <h1 style="font-size: 30px;margin: 0">
                {{articleInfo.title}}
              </h1>
              <el-text size="">
                {{articleInfo.release_time}}
              </el-text>
          </el-card>
        </el-col>
<!--        <el-col :span="8">-->
<!--          <el-button v-if="!articleInfo.stared" :icon="Star" @click="starArticle">收藏文章</el-button>-->
<!--          <el-button v-else :icon="Star" type="primary" @click="cancelStaredArticle">已收藏文章</el-button>-->
<!--        </el-col>-->
      </el-row>

<!--      <el-card>-->
<!--        {{articleInfo.content}}-->
<!--      </el-card>-->

<!--      <v-md-editor v-model="markdownText">-->

<!--      </v-md-editor>-->

      <v-md-preview :text="articleInfo.content"></v-md-preview>
<!--      <el-card>-->
<!--      </el-card>-->

<!--      这里的sticky必须让overflow是visible,不然sticky效果会失效-->
      <div style="position: sticky;bottom: 0;z-index: 6;background-color: white">

        <el-row :gutter="6" style="margin: 10px">
          <el-col :span="4.5">

          </el-col>
          <el-col :span="1.5">
            <font-awesome-icon :icon="['far', 'message']" fade/>
            <el-text style="margin-left: 5px">评论数:{{ articleInfo.comment_number }}</el-text>
          </el-col>
          <el-col :span="6">
<!--            用span前一定要加:-->
            <el-row :gutter="8">
              <el-col :span="18">

                <el-input type="textarea" autosize resize="none" v-model="commentInput" placeholder="show me the code" style="">
                </el-input>
              </el-col>
              <el-col :span="6">

                <el-button type="primary" style="border-radius: 24px" @click="commentArticle" :disabled="commentInput==0">
                  发布
                </el-button>
              </el-col>
            </el-row>
<!--            <div style="">-->
<!--              <div style="">-->
<!--                <el-input type="textarea" autosize resize="none" v-model="test" placeholder="show me the code">-->
<!--                </el-input>-->
<!--              </div>-->
<!--              <div style="">-->
<!--                <el-button type="primary" style="border-radius: 24px">-->
<!--                  发布-->
<!--                </el-button>-->
<!--              </div>-->
<!--            </div>-->
          </el-col>
          <el-col span="4.5">

          </el-col>
          <el-col :span="1">
            <!--            <el-button type="primary" :icon="StarFilled">收藏</el-button>-->

            <!--            动画用Transition组件和animate.css来做-->
            <Transition mode="out-in"
                        enter-active-class="animate__animated animate__bounceIn"
                        leave-active-class="animate__animated animate__bounceOut"
            >
              <font-awesome-icon v-if="!articleInfo.stared" :icon="['far', 'star']" style="" @click="starArticle"/>
              <font-awesome-icon v-else :icon="['fas', 'star']" style="color: rgb(200,204,39);" @click="cancelStaredArticle"/>
            </Transition>
            <!--            <el-button type="info" >-->
            <!--              <font-awesome-icon :icon="['fas', 'star']" style="color: rgba(247,252,7,0.91);" />-->
            <!--              收藏-->
            <!--            </el-button>-->
          </el-col>
          <el-col :span="1">
            <el-tooltip content="修改文章" placement="top" v-if="articleInfo.is_my_article" >
              <font-awesome-icon :icon="['fas', 'pen-to-square']" @click="router.push('/updateArticle/'.concat(route.params.articleID))"/>
            </el-tooltip>
          </el-col>
          <el-col :span="1">
            <el-tooltip content="删除文章" placement="top" v-if="articleInfo.is_my_article">
              <font-awesome-icon :icon="['fas', 'trash']" />
            </el-tooltip>
          </el-col>
          <el-col :span="4.5">

          </el-col>
        </el-row>
      </div>

<!--      <el-button @click="commentDialog=!commentDialog">评论文章</el-button>-->
<!--      <el-dialog v-model="commentDialog" title="评论文章">-->
<!--        <el-input v-model="commentInput" placeholder="comment"/>-->
<!--        <el-button type="primary" @click="commentArticle">评论</el-button>-->
<!--      </el-dialog>-->

<!--      <el-text>评论数:{{ articleInfo.comment_number }}</el-text>-->

<!--      <div v-if="articleInfo.is_my_article">-->
<!--        <el-button @click="router.push('/updateArticle/'.concat(route.params.articleID))">修改文章</el-button>-->
<!--        <el-button >删除文章</el-button>-->
<!--      </div>-->


<!--      <el-container style="height: 600px">-->
<!--        <el-main>-->
<!--          <el-card>-->
<!--            {{articleInfo.content}}-->
<!--          </el-card>-->
<!--        </el-main>-->
<!--        <el-footer style="position: sticky;bottom: 0;z-index: 6">-->
<!--          <span>-->
<!--            footer-->
<!--          </span>-->
<!--        </el-footer>-->
<!--      </el-container>-->

<!--      这里评论功能用弹窗实现-->

      <el-empty v-if="comments.length==0" description="评论为空哦"/>

      <div style="padding: 10px;margin: 10px">
        <el-card v-for="comment in comments" style="border-radius: 20px;padding: 5px; margin: 8px; border-color: #ffffff;background-color: #ffffff;" shadow="always">
          <el-avatar :src="comment.avatar_url" @click="router.push('/userProfile/'.concat(comment.uid))"/>
          <el-link @click="router.push('/userProfile/'.concat(comment.uid))" style="vertical-align: 12px">
            <span style="font-size: 20px">
              {{comment.nickname}}
            </span>
          </el-link>
          <br/>
          <el-text style="" tag="b" size="large">
            {{comment.content}}
          </el-text>
          <br/>
          <el-text tag="i">
            {{comment.release_time}}
          </el-text>

          <font-awesome-icon :icon="['fas', 'comment']" style="margin-left: 10px" @click="!loginStore.isLogin?showLoginDialog=true:replyDialog=!replyDialog;toUserNickname='';toUid=0;replyType=COMMENTTYPE;commentID=comment.comment_id;commentOffset=comment.offset"/>

<!--          <el-button @click="toUserNickname='';toUid=0;replyType=COMMENTTYPE;commentID=comment.comment_id;replyDialog=!replyDialog;commentOffset=comment.offset">回复</el-button>-->

<!--          <el-button @click="router.push('/userProfile/'.concat(comment.uid))">查看用户</el-button>-->

          <el-card v-for="reply in comment.replys"  style="border-radius: 20px;padding: 5px; margin: 8px; border-color: #ffffff;background-color: #ffffff;" shadow="always">
            <el-avatar :src="reply.from_avatar_url"/>
            <el-link @click="router.push('/userProfile/'.concat(reply.from_uid))" style="vertical-align: 12px">
            <span style="font-size: 20px">
              {{reply.from_user_nickname}}
            </span>
            </el-link>
            <br/>
            <el-text v-if="reply.reply_type==REPLYTYPE">回复@ {{reply.to_user_nickname}}</el-text>
            {{reply.content}}
            <br/>
            <el-text tag="i">
              {{reply.release_time}}
            </el-text>

            <font-awesome-icon :icon="['fas', 'comment']" style="margin-left: 10px" @click="!loginStore.isLogin?showLoginDialog=true:replyDialog=!replyDialog;toUserNickname=reply.from_user_nickname;toUid=reply.from_uid;replyType=REPLYTYPE;commentID=comment.comment_id;commentOffset=comment.offset"/>
<!--            <el-button @click="toUserNickname=reply.from_user_nickname;toUid=reply.from_uid;replyType=REPLYTYPE;commentID=comment.comment_id;replyDialog=!replyDialog;commentOffset=comment.offset">回复</el-button>-->
<!--            <el-button @click="router.push('/userProfile/'.concat(reply.from_uid))">查看用户</el-button>-->

          </el-card>


        </el-card>
      </div>

    </div>

  <el-dialog v-model="replyDialog" title="REPLY" @closed="replyInput=''">
    <el-input v-model="replyInput" placeholder="请输入想回复的内容" style="--el-input-border-radius: 23px; --el-input-bg-color: #f5f7fa;"/>
<!--    <el-button type="primary" @click="replyComment">回复</el-button>-->
    <template #footer>
      <el-button type="danger" @click="replyDialog=false;replyInput=''" style="border-radius: 25px">
        没有,我再想想
      </el-button>
      <el-button type="primary" @click="replyDialog=false;replyComment();replyInput=''" style="border-radius: 25px">
        确认
      </el-button>
    </template>
  </el-dialog>
</template>