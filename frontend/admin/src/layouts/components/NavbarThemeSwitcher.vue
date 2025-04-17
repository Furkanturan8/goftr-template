<script setup lang="ts">
import type { ThemeSwitcherTheme } from '@layouts/types'

const themes: ThemeSwitcherTheme[] = [
  {
    name: 'light',
    icon: 'bx bx-sun',
  },
  {
    name: 'dark',
    icon: 'bx bx-moon',
  },
]
</script>

<template>
  <ThemeSwitcher :themes="themes" />
</template>
