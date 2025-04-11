<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { authService } from '@/services/ApiService'
import JwtService from "@/services/JwtService";
import {useUserStore} from "@/store/user";

const router = useRouter()
const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
const rememberMe = ref(false)

const handleLogin = async () => {
  try {
    loading.value = true
    error.value = ''
    
    const response = await authService.login({
      email: email.value,
      password: password.value,
    })

    await useUserStore().login(response.data.access_token, response.data.refresh_token)

    if (rememberMe.value) {
      localStorage.setItem('remembered_email', email.value)
    } else {
      localStorage.removeItem('remembered_email')
    }
    
    // Ana sayfaya yönlendir
    await router.push('/dashboard')
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Giriş yapılırken bir hata oluştu'
  } finally {
    loading.value = false
  }
}

// Remember me özelliği için email'i yükle
const loadRememberedEmail = () => {
  const rememberedEmail = localStorage.getItem('remembered_email')
  if (rememberedEmail) {
    email.value = rememberedEmail
    rememberMe.value = true
  }
}

// Sayfa yüklendiğinde hatırlanan email'i getir
loadRememberedEmail()
</script>

<template>
  <div class="auth-wrapper d-flex align-center justify-center pa-4">
    <VCard class="auth-card pa-4 pt-7">
      <VCardItem class="justify-center">
        <VImg
          class="mb-4"
          src="@images/logo.png"
          max-width="40"
        />
      </VCardItem>

      <VCardText class="pt-2">
        <h5 class="text-h5 mb-1">
          Hoş Geldiniz! 👋
        </h5>
        <p class="mb-0">
          Lütfen hesabınıza giriş yapın
        </p>
      </VCardText>

      <VCardText>
        <VForm @submit.prevent="handleLogin">
          <VRow>
            <!-- Email -->
            <VCol cols="12">
              <VTextField
                v-model="email"
                label="Email"
                type="email"
                required
              />
            </VCol>

            <!-- Password -->
            <VCol cols="12">
              <VTextField
                v-model="password"
                label="Şifre"
                type="password"
                required
              />
            </VCol>

            <!-- Remember me -->
            <VCol
              cols="12"
              class="d-flex justify-space-between flex-wrap gap-3"
            >
              <VCheckbox
                v-model="rememberMe"
                label="Beni hatırla"
              />
              
              <RouterLink
                to="/forgot-password"
                class="text-primary ms-2 mb-1"
              >
                Şifremi unuttum
              </RouterLink>
            </VCol>

            <!-- Error -->
            <VCol
              v-if="error"
              cols="12"
            >
              <VAlert
                color="error"
                variant="tonal"
              >
                {{ error }}
              </VAlert>
            </VCol>

            <!-- Submit -->
            <VCol cols="12">
              <VBtn
                block
                type="submit"
                :loading="loading"
              >
                Giriş Yap
              </VBtn>
            </VCol>

            <!-- Register Link -->
            <VCol
              cols="12"
              class="text-center"
            >
              <span>Hesabınız yok mu?</span>
              <RouterLink
                class="text-primary ms-2"
                to="/register"
              >
                Kayıt ol
              </RouterLink>
            </VCol>
          </VRow>
        </VForm>
      </VCardText>
    </VCard>

  </div>
</template>

<style lang="scss">
@use "@core/scss/template/pages/page-auth.scss";
</style>
