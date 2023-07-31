<template>
  <swiper
      class="swiper"
      :modules="modules"
      :slides-per-view="1"
      :space-between="30"
      :pagination="{ clickable: true }"
      :navigation="true"
      :autoplay="{
      delay: 2500,
      disableOnInteraction: false
    }"
  >
    <swiper-slide class="slide" v-for="index in 3">
      <img style="width: 100%" :src="'/images/swiper/'.concat(index).concat('.jpg')">
    </swiper-slide>
  </swiper>
</template>

<script setup>
import {Pagination, Navigation, Autoplay} from 'swiper'
import { Swiper, SwiperSlide } from 'swiper/vue'
import 'swiper/css'
import 'swiper/css/pagination'
import 'swiper/css/navigation'

const modules = [Pagination, Navigation,  Autoplay]

</script>

<style lang="scss" scoped>
.swiper {
  width: 100%;
  .slide {
    display: flex;
    justify-content: center;
    align-items: center;
  }
}
</style>