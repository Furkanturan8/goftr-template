<script lang="ts" setup>
import { ref } from 'vue'
import { useUserStore } from '@/store/user.ts'

const route = useRoute()
const userStore = useUserStore()

const activeTab = ref(route.params.tab || 'account')
const loading = ref(false)

// Form verilerini tanımlama
const accountData = ref({
  email: userStore.user?.email || '',
  first_name: userStore.user?.first_name || '',
  last_name: userStore.user?.last_name || '',
  role: userStore.user?.role || '',
})

const passwordData = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: '',
})

// tabs
const tabs = [
  { title: 'Hesap', icon: 'bx-user', tab: 'account' },
  { title: 'Güvenlik', icon: 'bx-lock-alt', tab: 'security' },
]

// Form gönderme işlemleri
const onAccountSubmit = async () => {
  loading.value = true
  try {
    // Hesap bilgilerini güncelleme işlemi
    await userStore.updateProfile(accountData.value)
    // Başarılı mesajı göster
  } catch (error) {
    // Hata mesajı göster
  } finally {
    loading.value = false
  }
}

const onPasswordSubmit = async () => {
  loading.value = true
  try {
    // Şifre değiştirme işlemi
    await userStore.changePassword(passwordData.value)
    passwordData.value = {
      currentPassword: '',
      newPassword: '',
      confirmPassword: '',
    }
    // Başarılı mesajı göster
  } catch (error) {
    // Hata mesajı göster
  } finally {
    loading.value = false
  }
}

const formattedRole = computed({
  get: () => accountData.value.role.charAt(0).toUpperCase() + accountData.value.role.slice(1),
  set: (value) => accountData.value.role = value.toLowerCase()
})
</script>

<template>
  <div>
    <VTabs
      v-model="activeTab"
      show-arrows
      class="v-tabs-pill"
    >
      <VTab
        v-for="item in tabs"
        :key="item.icon"
        :value="item.tab"
      >
        <VIcon
          size="20"
          start
          :icon="item.icon"
        />
        {{ item.title }}
      </VTab>
    </VTabs>

    <VWindow
      v-model="activeTab"
      class="mt-5 disable-tab-transition"
    >
      <!-- Hesap -->
      <VWindowItem value="account">
        <VCard>
          <VCardText>
            <VForm @submit.prevent="onAccountSubmit">
              <VRow>
                <VCol cols="12" md="6">
                  <VTextField
                    v-model="accountData.email"
                    label="E-posta"
                    placeholder="E-posta adresinizi girin"
                    :rules="[
                      v => !!v || 'E-posta gerekli',
                      v => /.+@.+\..+/.test(v) || 'Geçerli bir e-posta adresi girin'
                    ]"
                  />
                </VCol>

                <VCol cols="12" md="6">
                  <VTextField
                    v-model="accountData.first_name"
                    label="Ad"
                    placeholder="Adınızı girin"
                  />
                </VCol>
                <VCol cols="12" md="6">
                  <VTextField
                    v-model="accountData.last_name"
                    label="Soyad"
                    placeholder="Soyadınızı girin"
                  />
                </VCol>
                <VCol cols="12" md="6">
                  <VTextField
                    v-model="formattedRole"
                    label="Role"
                    disabled
                    placeholder="Soyadınızı girin"
                  />
                </VCol>

                <VCol cols="12" class="text-right">
                  <VBtn
                    type="submit"
                    :loading="loading"
                    color="primary"
                  >
                    Değişiklikleri Kaydet
                  </VBtn>
                </VCol>
              </VRow>
            </VForm>
          </VCardText>
        </VCard>
      </VWindowItem>

      <!-- Güvenlik -->
      <VWindowItem value="security">
        <VCard>
          <VCardText>
            <VForm @submit.prevent="onPasswordSubmit">
              <VRow>
                <VCol cols="12" md="6">
                  <VTextField
                    v-model="passwordData.currentPassword"
                    label="Mevcut Şifre"
                    type="password"
                    placeholder="Mevcut şifrenizi girin"
                    :rules="[v => !!v || 'Mevcut şifre gerekli']"
                  />
                </VCol>

                <VCol cols="12" md="6">
                  <VTextField
                    v-model="passwordData.newPassword"
                    label="Yeni Şifre"
                    type="password"
                    placeholder="Yeni şifrenizi girin"
                    :rules="[
                      v => !!v || 'Yeni şifre gerekli',
                      v => v.length >= 8 || 'Şifre en az 8 karakter olmalı'
                    ]"
                  />
                </VCol>

                <VCol cols="12" md="6">
                  <VTextField
                    v-model="passwordData.confirmPassword"
                    label="Şifre Tekrar"
                    type="password"
                    placeholder="Yeni şifrenizi tekrar girin"
                    :rules="[
                      v => !!v || 'Şifre tekrarı gerekli',
                      v => v === passwordData.newPassword || 'Şifreler eşleşmiyor'
                    ]"
                  />
                </VCol>

                <VCol cols="12" class="text-right">
                  <VBtn
                    type="submit"
                    :loading="loading"
                    color="primary"
                  >
                    Şifreyi Değiştir
                  </VBtn>
                </VCol>
              </VRow>
            </VForm>
          </VCardText>
        </VCard>
      </VWindowItem>
    </VWindow>
  </div>
</template>

<style scoped>
.v-tabs-pill {
  margin-bottom: 1rem;
}
</style>
