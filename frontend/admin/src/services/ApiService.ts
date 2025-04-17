import axios from 'axios'

const API_URL = 'http://localhost:3005/api/v1'

const ApiService = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor - token ekleme
ApiService.interceptors.request.use(
  config => {
    const token = localStorage.getItem('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => Promise.reject(error),
)

// Response interceptor - hata yönetimi
ApiService.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config

    // Token süresi dolmuşsa ve yenileme denemesi yapılmamışsa
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      const refreshToken = localStorage.getItem('refresh_token')

      if (refreshToken) {
        try {
          // Token yenileme isteği
          const response = await axios.post(`${API_URL}/auth/refresh`, {
            refresh_token: refreshToken,
          })

          const { access_token: accessToken, refresh_token: newRefreshToken } = response.data.data

          // Yeni token'ları kaydet
          localStorage.setItem('access_token', accessToken)
          localStorage.setItem('refresh_token', newRefreshToken)

          // Tüm sonraki istekler için header'ı güncelle
          ApiService.defaults.headers.common.Authorization = `Bearer ${accessToken}`

          // Başarısız olan orijinal isteği yeni token ile tekrar dene
          originalRequest.headers.Authorization = `Bearer ${accessToken}`
          return ApiService(originalRequest)
        } catch (refreshError) {
          // Token yenileme başarısız olursa
          localStorage.removeItem('access_token')
          localStorage.removeItem('refresh_token')
          
          // Kullanıcıyı login sayfasına yönlendir
          window.location.href = '/login'
          return Promise.reject(refreshError)
        }
      } else {
        // Refresh token yoksa login sayfasına yönlendir
        window.location.href = '/login'
        return Promise.reject(error)
      }
    }

    return Promise.reject(error)
  },
)

// Auth servisleri
export const authService = {
  register: (data: any) => ApiService.post('/auth/register', data),
  login: (data: any) => ApiService.post('/auth/login', data),
  logout: () => ApiService.post('/auth/logout'),
  forgotPassword: (data: any) => ApiService.post('/auth/forgot-password', data),
  resetPassword: (data: any) => ApiService.post('/auth/reset-password', data),
}

// Kullanıcı servisleri
export const userService = {
  getProfile: () => ApiService.get('/users/me'),
  updateProfile: (data: any) => ApiService.put('/users/me', data),

  // Admin only routes
  listUsers: () => ApiService.get('/users'),
  createUser: (data: any) => ApiService.post('/users', data),
  getUserById: (id: string) => ApiService.get(`/users/${id}`),
  updateUser: (id: number, data: any) => ApiService.put(`/users/${id}`, data),
  deleteUser: (id: number) => ApiService.delete(`/users/${id}`),
}

export { ApiService }
