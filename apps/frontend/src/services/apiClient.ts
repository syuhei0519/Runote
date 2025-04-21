import axios from 'axios'

export const apiClient = axios.create({
  baseURL: '/api', // Nginx が /api/* を API Gateway に中継
  withCredentials: true, // cookie送信が必要なら（JWT等）
})

apiClient.interceptors.response.use(
    (res) => res,
    (err) => {
      const status = err.response?.status
      if (status === 401) {
        console.warn('認証エラー。ログインが必要です。')
        // router.push('/login') なども可能
      } else {
        console.error('APIエラー:', err.message)
      }
      return Promise.reject(err)
    }
)