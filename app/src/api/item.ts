import { AxiosInstance } from '../utils/axios'
import { CookieWrapper } from '../utils/storage'

type CreateUserPayload = {
  user: {
    username: string
    email: string
    password: string
  }
}

export async function getAllItems(
  createUserPayload: CreateUserPayload
) {
  const requestOptions = {
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(createUserPayload),
  }
  try {
    const response = await AxiosInstance.post(
      `/signup`,
      createUserPayload,
      requestOptions
    )
    const data = response.data
    if (data.user) {
      let expires = new Date()
      expires = new Date(expires.getTime() + 1000 * 60 * 60 * 24 * 7)
      CookieWrapper.setCookie('userAuth', data.user, expires)
      CookieWrapper.setCookie('jwt', data.jwt, expires)
      // setup axios to use auth from now
      AxiosInstance.defaults.headers.common['Authorization'] =
        'Bearer ' + data.jwt
      return data
    }
    return
  } catch (error) {
    console.log(error.data.message)
  }
}
