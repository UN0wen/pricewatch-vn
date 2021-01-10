import axios from 'axios'

export const AxiosInstance = axios.create({
  baseURL: 'http://pricewatch-vn.herokuapp.com/api',
})


