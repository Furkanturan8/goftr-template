export const required = (v: string) => !!v || 'Bu alan gerekli'
export const emailRule = (v: string) =>
  /.+@.+\..+/.test(v) || 'Geçerli bir e-posta adresi girin'

export const passwordVAL = [
  (v: string) => !!v || 'Şifre gerekli',
  (v: string) => v.length >= 8 || 'Şifre en az 8 karakter olmalı'
]

export function passwordAgainVAL(newPassword: string) {
  return [
    (v: string) => !!v || 'Şifre tekrarı gerekli',
    (v: string) => v === newPassword || 'Şifreler uyuşmuyor',
  ]
}
