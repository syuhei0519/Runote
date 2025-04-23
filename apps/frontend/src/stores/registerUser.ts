import { defineStore } from 'pinia'

export const useRegisterUserStore = defineStore('registerUser', {
    state: () => ({
      username: '',
      password: ''
    }),
  })