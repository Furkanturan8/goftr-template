<script lang="ts" setup>
import VerticalNavLink from '@layouts/components/VerticalNavLink.vue'
import VerticalNavSectionTitle from '@layouts/components/VerticalNavSectionTitle.vue'
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
    :item="{
      title: 'Kullanıcı İşlemleri',
      icon: 'bx bx-user',
      to: '/users',
    }"
  />

</template>
