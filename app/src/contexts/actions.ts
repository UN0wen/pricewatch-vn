import { AuthDispatch, AuthReducerActions } from './reducer'
import { CookieWrapper } from '../utils/storage'
import { AxiosInstance } from '../utils/axios'

export interface LoginPayload {
  user: {
    email: string
    password: string
  }
  
}

export async function loginUser(
  dispatch: React.Dispatch<AuthDispatch>,
  loginPayload: LoginPayload
) {
  const requestOptions = {
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(loginPayload),
  }

  try {
    dispatch({
      type: AuthReducerActions.REQUEST_LOGIN,
      payload: {} as any,
      error: {} as any,
    })
    const response = await AxiosInstance.post(`/login`, loginPayload, requestOptions)
    const data = response.data
    if (data.user) {
      dispatch({
        type: AuthReducerActions.LOGIN_SUCCESS,
        payload: data,
        error: {} as any,
      })
      let expires = new Date()
      expires = new Date(expires.getTime() + 1000 * 60 * 60 * 24 * 7)
      CookieWrapper.setCookie(data, expires)

      // setup axios to use auth from now
      AxiosInstance.defaults.headers.common['Authorization'] = "Bearer " + data.jwt
      return data
    }

    dispatch({
      type: AuthReducerActions.LOGIN_ERROR,
      payload: {} as any,
      error: data.errors[0],
    })
    return
  } catch (error) {
    dispatch({
      type: AuthReducerActions.LOGIN_ERROR,
      payload: {} as any,
      error: error,
    })
  }
}

export async function logout(dispatch: React.Dispatch<AuthDispatch>) {
  dispatch({
    type: AuthReducerActions.LOGOUT,
    payload: {} as any,
    error: {} as any,
  })

  // setup axios to stop using auth
  AxiosInstance.defaults.headers.common['Authorization'] = ''
  CookieWrapper.removeCookie()
}
