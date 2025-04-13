import { defineStore } from 'pinia'
import JwtService from "@/services/JwtService";
import {ApiService} from "@/services/ApiService";
import {User} from "@/types";

export type UserRole = 'admin' | 'normal'

export const useUserStore = defineStore('UserStore', {
  state: (): { user: User; isAuthenticated: boolean } => ({
    user: localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user') as string) : {} as User,
    isAuthenticated: !!localStorage.getItem('access_token'),
  }),

  getters: {
    isUserAuthenticated: state => state.isAuthenticated,
    hasRole: state => (roles?: UserRole[]) => !roles || roles.includes(state.user.role as UserRole), // role paramteresi verilmediyse true döndürür
    getRole: state => () => state.user.role,
  },

  actions: {
    async login(access_token: string, refresh_token: string) {
      JwtService.saveTokens(access_token, refresh_token)
      this.isAuthenticated = true
      await this.getProfile()
    },
    async logout() {
      this.isAuthenticated = false
      this.user = {} as User
      JwtService.destroyTokens()
      localStorage.removeItem('user')

      const redirect = window.location.pathname + window.location.search
      const urlParams = new URLSearchParams()

      urlParams.set('redirect', redirect)
      document.location.href = `/auth/login?${urlParams.toString()}`
    },
    async getProfile() {
      if (!this.isAuthenticated) return

      try {
        const response = await ApiService.get('users/me')
        const userData = response.data.data  // 👈 Buraya dikkat, içteki data!

        this.user = {
          id: userData.id,
          email: userData.email,
          first_name: userData.first_name,
          last_name: userData.last_name,
          role: userData.role as UserRole,
        }

        localStorage.setItem('user', JSON.stringify(this.user))
      } catch (error) {
        console.error(error)
      }
    }
  },
})
