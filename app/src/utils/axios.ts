import axios from 'axios'

export const AxiosInstance = axios.create({
  baseURL: 'http://172.20.204.69:8080/api',
})
