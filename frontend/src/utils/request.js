import axios from "axios";

const http = axios.create({
    baseURL: 'http://127.0.0.1:8080/', // 用于本机测试的后端接口地址
    timeout: 5000,
})
export default http
