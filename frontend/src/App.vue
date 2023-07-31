<script setup>
import CommonHeader from './components/CommonHeader.vue';
</script>

<template>

  <div style="">
    <el-container>

      <el-header style="position: sticky;top: 0;z-index: 7;padding: 0;">
        <common-header/>
      </el-header>

      <el-main style="overflow: visible;padding:0">
        <router-view :key="$route.fullPath"/>
      </el-main>

    </el-container>
  </div>
</template>

<style>

/*记得设置全局样式*/
.html, body {
  margin: 0;
  padding: 0;
}

</style>