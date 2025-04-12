export interface User {
  id: number
  email: string
  first_name: string
  last_name: string
  role: 'admin' | 'user'
  status: 'active' | 'inactive' | 'banned'
  last_login: string
  created_at: string
  updated_at: string
}

export interface UpdateUserRequest {
  first_name?: string
  last_name?: string
  email?: string
}

export interface UserState {
  user: User | null
  loading: boolean
  error: string | null
} 