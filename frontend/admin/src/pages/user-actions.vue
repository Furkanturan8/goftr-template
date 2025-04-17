<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {userService} from "@/services/ApiService";

interface User {
  id: number
  first_name: string
  last_name: string
  email: string
  role: string
  status: 'active' | 'inactive'
}

const users = ref<User[]>([])
const loading = ref(false)
const dialog = ref(false)
const editedIndex = ref(-1)
const search = ref('')
const { width } = useWindowSize()

const editedItem = ref<User>({
  id: 0,
  first_name: '',
  last_name: '',
  email: '',
  role: '',
  status: 'active',
})

const defaultItem: User = {
  id: 0,
  first_name: '',
  last_name: '',
  email: '',
  role: '',
  status: 'active',
}

const headers = [
  { title: '#', key: 'index', sortable: false }, // 👈 bu eklendi
  { title: 'İsim', key: 'first_name', sortable: true },
  { title: 'Soyisim', key: 'last_name', sortable: true },
  { title: 'E-posta', key: 'email', sortable: true },
  { title: 'Rol', key: 'role', sortable: true },
  { title: 'Durum', key: 'status', sortable: true },
  { title: 'İşlemler', key: 'actions', sortable: false },
]

const fetchUsers = async () => {
  loading.value = true
  try {
    const data = await userService.listUsers()
    users.value = data.data.data
  } catch (error) {
    console.error('Kullanıcılar yüklenirken hata oluştu:', error)
  } finally {
    loading.value = false
  }
}

const editItem = (id: number) => {
  // id'ye göre kullanıcıyı bul
  const index = users.value.findIndex(user => user.id === id)
  if (index !== -1) {
    editedIndex.value = index
    // editedItem'a seçilen kullanıcının bilgilerini kopyala
    editedItem.value = { ...users.value[index] }
    dialog.value = true
  }
}

const deleteItem = async (id: number) => {
  if (confirm('Bu kullanıcıyı silmek istediğinizden emin misiniz?')) {
    try {
      await userService.deleteUser(id)
      const index = users.value.findIndex(user => user.id === id)
      if (index !== -1) {
        users.value.splice(index, 1)
      }
    } catch (error) {
      console.error('Kullanıcı silinirken hata oluştu:', error)
    }
  }
  await fetchUsers()
}


const save = async () => {
  try {
    if (editedIndex.value > -1) {
      // Güncelleme işlemi
      const response = await userService.updateUser(editedItem.value.id, editedItem.value)
      console.log(response.data.message)
    } else {
      // Yeni kullanıcı ekleme
      const response = await userService.createUser(editedItem.value)
      console.log(response.data.message)
    }
    close()
  } catch (error) {
    console.error('Kullanıcı kaydedilirken hata oluştu:', error)
  }
  await fetchUsers()
}


const close = () => {
  dialog.value = false
  // Bir sonraki açılışta yeni kullanıcı oluşturma modunda olacak
  editedIndex.value = -1
  // Form alanlarını temizle
  editedItem.value = { ...defaultItem }
}

onMounted(() => {
  fetchUsers()
})
</script>

<template>
  <VCard>
    <VCardTitle class="pa-6">
      <!-- 1. Satır: Başlık -->
      <div class="text-h5 mb-4">Kullanıcı Yönetimi</div>

      <!-- 2. Satır: Arama ve Buton -->
      <div class="d-flex flex-column flex-sm-row align-stretch gap-2">
        <VTextField
          v-model="search"
          prepend-inner-icon="bx bx-search"
          label="Ara"
          single-line
          hide-details
          density="compact"
          :style="{ width: width <= 600 ? '100%' : '250px' }"
        />
        <VBtn
          color="primary"
          prepend-icon="bx bx-plus"
          @click="dialog = true"
          :style="{ width: width <= 600 ? '100%' : '250px' }"
        >
          {{ width <= 600 ? 'Ekle' : 'Yeni Kullanıcı' }}
        </VBtn>
      </div>
    </VCardTitle>

    <VCardText>
      <VDataTable
        v-model:search="search"
        :headers="headers"
        :items="users"
        :loading="loading"
        hover
      >
        <!-- Sıra numarası -->
        <template #item.index="{ index }">
          {{ index + 1 }}
        </template>
        <!-- Status sütunu için özel render -->
        <template #item.status="{ item }">
          <VChip
            :color="item.status === 'active' ? 'success' : 'error'"
            size="small"
          >
            {{ item.status === 'active' ? 'Aktif' : 'Pasif' }}
          </VChip>
        </template>

        <!-- İşlemler sütunu için özel render -->
        <template #item.actions="{ item }">
          <VBtn
            icon
            variant="text"
            color="primary"
            size="small"
            @click="editItem(item.id)"
          >
            <VIcon size="20">
              bx bx-edit
            </VIcon>
          </VBtn>

          <VBtn
            icon
            variant="text"
            color="error"
            size="small"
            @click="deleteItem(item.id)"
          >
            <VIcon size="20">
              bx bx-trash
            </VIcon>
          </VBtn>
        </template>
      </VDataTable>
    </VCardText>

    <!-- Kullanıcı Ekleme/Düzenleme Dialog -->
    <VDialog
      v-model="dialog"
      max-width="500px"
    >
      <VCard>
        <VCardTitle>
          <span class="text-h5">{{ editedIndex === -1 ? 'Yeni Kullanıcı' : 'Kullanıcıyı Düzenle' }}</span>
        </VCardTitle>

        <VCardText>
          <VContainer>
            <VRow>
              <VCol
                cols="12"
                sm="6"
              >
                <VTextField
                  v-model="editedItem.first_name"
                  label="İsim"
                />
              </VCol>
              <VCol
                cols="12"
                sm="6"
              >
                <VTextField
                  v-model="editedItem.last_name"
                  label="Soyisim"
                />
              </VCol>
              <VCol
                cols="12"
                sm="12"
              >
                <VTextField
                  v-model="editedItem.email"
                  label="E-posta"
                  type="email"
                  :rules="[
                      v => !!v || 'E-posta gerekli',
                      v => /.+@.+\..+/.test(v) || 'Geçerli bir e-posta adresi girin'
                    ]"
                />
              </VCol>
              <VCol
                cols="12"
                sm="6"
              >

                <VSelect
                  v-model="editedItem.status"
                  :items="[
                    { title: 'Aktif', value: 'active' },
                    { title: 'Pasif', value: 'inactive' }
                  ]"
                  item-title="title"
                  item-value="value"
                  label="Durum"
                />
              </VCol>
              <VCol
                cols="12"
                sm="6"
              >
                <VSelect
                  v-model="editedItem.role"
                  :items="['admin', 'user']"
                  label="Rol"
                />
              </VCol>
            </VRow>
          </VContainer>
        </VCardText>

        <VCardActions>
          <VSpacer />
          <VBtn
            color="error"
            variant="text"
            @click="close"
          >
            İptal
          </VBtn>
          <VBtn
            color="primary"
            variant="text"
            @click="save"
          >
            Kaydet
          </VBtn>
        </VCardActions>
      </VCard>
    </VDialog>
  </VCard>
</template>

<style scoped lang="scss">
.v-data-table {
  .v-data-table-header {
    th {
      font-weight: 600;
      text-transform: uppercase;
      white-space: nowrap;
    }
  }
}
</style>
