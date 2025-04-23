<script lang="ts" setup>
import VerticalNavLink from '@layouts/components/VerticalNavLink.vue'
import { useUserStore } from "@/store/user";

const userStore = useUserStore()
</script>


<template>
  <!-- 👉 Home -->
  <VerticalNavLink
    :item="{
      title: 'Anasayfa',
      icon: 'bx bx-home-smile',
      to: '/dashboard',
    }"
  />
  <VerticalNavLink
    v-if="userStore.user.role === 'admin'"
    :item="{
      title: 'Kullanıcı İşlemleri',
      icon: 'bx bx-user',
      to: '/users',
    }"
  />

</template>
