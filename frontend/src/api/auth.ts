import request from './request'

export interface LoginData {
  username: string
  password: string
}

export interface LoginResult {
  token: string
  user: {
    id: number
    username: string
    nickname: string
    email: string
    role: string
  }
}

export interface RegisterData {
  username: string
  password: string
  email?: string
  nickname?: string
}

export function login(data: LoginData) {
  return request.post<LoginResult>('/auth/login', data)
}

export function register(data: RegisterData) {
  return request.post('/auth/register', data)
}

export function getProfile() {
  return request.get('/auth/profile')
}
