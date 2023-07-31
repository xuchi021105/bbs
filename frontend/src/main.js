import { createApp } from 'vue'
import App from './App.vue'
import ElementPlus from 'element-plus'
import router from './router'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { createPinia } from 'pinia'
import piniaPersist from 'pinia-plugin-persist'

import VMdEditor from '@kangc/v-md-editor/lib/codemirror-editor';
import '@kangc/v-md-editor/lib/style/codemirror-editor.css';
import '@kangc/v-md-editor/lib/theme/style/github.css';

import VMdPreview from '@kangc/v-md-editor/lib/preview';
import '@kangc/v-md-editor/lib/style/preview.css';
import githubTheme from '@kangc/v-md-editor/lib/theme/github.js';
import '@kangc/v-md-editor/lib/theme/style/github.css';

// highlightjs
import hljs from 'highlight.js';

// codemirror 编辑器的相关资源
import Codemirror from 'codemirror';
// mode
import 'codemirror/mode/markdown/markdown';
import 'codemirror/mode/javascript/javascript';
import 'codemirror/mode/css/css';
import 'codemirror/mode/htmlmixed/htmlmixed';
import 'codemirror/mode/vue/vue';
// edit
import 'codemirror/addon/edit/closebrackets';
import 'codemirror/addon/edit/closetag';
import 'codemirror/addon/edit/matchbrackets';
// placeholder
import 'codemirror/addon/display/placeholder';
// active-line
import 'codemirror/addon/selection/active-line';
// scrollbar
import 'codemirror/addon/scroll/simplescrollbars';
import 'codemirror/addon/scroll/simplescrollbars.css';
// style
import 'codemirror/lib/codemirror.css';

/* import the fontawesome core */
import { library } from '@fortawesome/fontawesome-svg-core'

/* import font awesome icon component */
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

/* import specific icons */
import { faUserSecret ,faTrash, faPenToSquare ,faComment, faUpload} from '@fortawesome/free-solid-svg-icons'
import { faStar as fasStar} from '@fortawesome/free-solid-svg-icons'
import { faStar as farStar} from '@fortawesome/free-regular-svg-icons'
import {faMessage} from "@fortawesome/free-regular-svg-icons";
import {faGithub} from "@fortawesome/free-brands-svg-icons";

import 'animate.css'

/* add icons to the library */
library.add(faUserSecret, fasStar, faMessage, farStar, faGithub, faTrash, faPenToSquare, faComment, faUpload)


VMdEditor.Codemirror = Codemirror;

VMdEditor.use(githubTheme, {
  Hljs: hljs,
});

VMdPreview.use(githubTheme, {
  Hljs: hljs,
});

const pinia = createPinia()
pinia.use(piniaPersist)

const app = createApp(App)

for (const [key, component] of Object.entries(ElementPlusIconsVue)) { // 导入icon图标
  app.component(key, component)
}

// use函数表示使用中间件
app.use(ElementPlus) // element-ui plus
app.use(router) // vue router
app.use(pinia) // pinia
app.use(VMdEditor) // markdown editor
app.use(VMdPreview)

app.component('font-awesome-icon', FontAwesomeIcon).mount('#app')