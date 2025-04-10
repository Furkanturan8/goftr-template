<script lang="ts" setup>
import VerticalNavLink from '@layouts/components/VerticalNavLink.vue'
import VerticalNavSectionTitle from '@layouts/components/VerticalNavSectionTitle.vue'
</script>

<template>
  <!-- 👉 Home -->
  <VerticalNavLink
    :item="{
      title: 'Anasayfa',
      icon: 'bx-home-smile',
      to: '/dashboard',
    }"
  />
</template>
