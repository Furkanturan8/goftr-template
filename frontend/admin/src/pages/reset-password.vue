<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import {authService} from "@/services/ApiService";

const route = useRoute()
const router = useRouter()

const token = route.query.token as string || ''
const newPassword = ref('')
const confirmPassword = ref('')
const errorMessage = ref('')
const successMessage = ref('')
const loading = ref(false)

const resetPassword = async () => {
  errorMessage.value = ''
  successMessage.value = ''

  if (!token) {
    errorMessage.value = 'Token bulunamadı.'
    return
  }

  if (newPassword.value !== confirmPassword.value) {
    errorMessage.value = 'Şifreler uyuşmuyor.'
    return
  }

  if (newPassword.value.length < 6) {
    errorMessage.value = 'Şifre en az 6 karakter olmalı.'
    return
  }

  loading.value = true

  try {
    const data = {
      new_password: newPassword.value,
      token: token,
    };

    await authService.resetPassword(data);

    successMessage.value = 'Şifreniz başarıyla güncellendi. Giriş sayfasına yönlendiriliyorsunuz...'

    setTimeout(() => {
      router.push('/login')
    }, 2000)
  } catch (error: any) {
    errorMessage.value = error.response?.data?.message || 'Bir hata oluştu.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="reset-password">
    <h1>Şifre Sıfırlama</h1>

    <form @submit.prevent="resetPassword">
      <div>
        <label>Yeni Şifre</label>
        <input v-model="newPassword" type="password" required />
      </div>

      <div>
        <label>Yeni Şifre (Tekrar)</label>
        <input v-model="confirmPassword" type="password" required />
      </div>

      <button type="submit" :disabled="loading">
        {{ loading ? 'Gönderiliyor...' : 'Şifreyi Sıfırla' }}
      </button>

      <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
      <p v-if="successMessage" class="success">{{ successMessage }}</p>
    </form>
  </div>
</template>

<style scoped lang="scss">
.reset-password {
  max-width: 400px;
  margin: 50px auto;
  padding: 2rem;
  border: 1px solid #ccc;
  border-radius: 8px;

  h1 {
    margin-bottom: 1.5rem;
    text-align: center;
  }

  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: bold;
  }

  input {
    width: 100%;
    padding: 0.5rem;
    margin-bottom: 1.2rem;
    border: 1px solid #ccc;
    border-radius: 4px;
  }

  button {
    width: 100%;
    padding: 0.7rem;
    background-color: #2d8cf0;
    color: white;
    border: none;
    border-radius: 4px;
    font-weight: bold;
    cursor: pointer;

    &:disabled {
      background-color: #bbb;
      cursor: not-allowed;
    }
  }

  .error {
    color: red;
    margin-top: 1rem;
    text-align: center;
  }

  .success {
    color: green;
    margin-top: 1rem;
    text-align: center;
  }
}
</style>
