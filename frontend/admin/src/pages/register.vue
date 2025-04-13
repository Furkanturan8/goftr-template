<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { authService } from '@/services/ApiService'

const router = useRouter()
const loading = ref(false)
const error = ref('')

const formData = ref({
  name: '',
  email: '',
  password: '',
  passwordConfirm: '',
  terms: false,
})

const handleRegister = async () => {
  try {
    if (formData.value.password !== formData.value.passwordConfirm) {
      error.value = 'Şifreler eşleşmiyor'
      return
    }

    if (!formData.value.terms) {
      error.value = 'Kullanım koşullarını kabul etmelisiniz'
      return
    }

    loading.value = true
    error.value = ''

    const { passwordConfirm, terms, ...registerData } = formData.value
    
    await authService.register(registerData)

    // Kayıt başarılı, login sayfasına yönlendir
    router.push('/login')
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Kayıt olurken bir hata oluştu'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-wrapper d-flex align-center justify-center pa-4">
    <VCard class="auth-card pa-4 pt-7">
      <VCardText class="pt-2">
        <h5 class="text-h5 mb-1">
          Macera burada başlıyor 🚀
        </h5>
        <p class="mb-0">
          Hesabınızı oluşturmak çok kolay!
        </p>
      </VCardText>

      <VCardText>
        <VForm @submit.prevent="handleRegister">
          <VRow>
            <!-- Name -->
            <VCol cols="12">
              <VTextField
                v-model="formData.name"
                label="Ad Soyad"
                required
              />
            </VCol>

            <!-- Email -->
            <VCol cols="12">
              <VTextField
                v-model="formData.email"
                label="Email"
                type="email"
                required
              />
            </VCol>

            <!-- Password -->
            <VCol cols="12">
              <VTextField
                v-model="formData.password"
                label="Şifre"
                type="password"
                required
              />
            </VCol>

            <!-- Password Confirm -->
            <VCol cols="12">
              <VTextField
                v-model="formData.passwordConfirm"
                label="Şifre (Tekrar)"
                type="password"
                required
              />
            </VCol>

            <!-- Terms -->
            <VCol cols="12">
              <VCheckbox
                v-model="formData.terms"
                label="Kullanım koşullarını kabul ediyorum"
                required
              />
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
                Kayıt Ol
              </VBtn>
            </VCol>

            <!-- Login Link -->
            <VCol
              cols="12"
              class="text-center"
            >
              <span>Zaten hesabınız var mı?</span>
              <RouterLink
                class="text-primary ms-2"
                to="/login"
              >
                Giriş yap
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
