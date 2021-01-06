import { AuthDispatch, AuthReducerActions } from '../contexts/reducer'
import { CookieWrapper } from '../utils/storage'
import { AxiosInstance } from '../utils/axios'

type LoginPayload = {
  user: {
    email: string
    password: string
  }
}

type CreateUserPayload = {
  user: {
    username: string
    email: string
    password: string
  }
}

export async function createUser(
  dispatch: React.Dispatch<AuthDispatch>,
  createUserPayload: CreateUserPayload
) {
  const requestOptions = {
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(createUserPayload),
  }
  try {
    dispatch({
      type: AuthReducerActions.REQUEST_LOGIN,
      payload: {} as any,
      error: {} as any,
    })
    const response = await AxiosInstance.post(
      `/signup`,
      createUserPayload,
      requestOptions
    )
    const data = response.data
    if (data.user) {
      dispatch({
        type: AuthReducerActions.LOGIN_SUCCESS,
        payload: data,
        error: {} as any,
      })
      let expires = new Date()
      expires = new Date(expires.getTime() + 1000 * 60 * 60 * 24 * 7)
      CookieWrapper.setCookie("userAuth", data.user, expires)
      CookieWrapper.setCookie("jwt", data.jwt, expires)
      // setup axios to use auth from now
      AxiosInstance.defaults.headers.common['Authorization'] =
        'Bearer ' + data.jwt
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
    const response = await AxiosInstance.post(
      `/login`,
      loginPayload,
      requestOptions
    )
    const data = response.data
    if (data.user) {
      dispatch({
        type: AuthReducerActions.LOGIN_SUCCESS,
        payload: data,
        error: {} as any,
      })
      let expires = new Date()
      expires = new Date(expires.getTime() + 1000 * 60 * 60 * 24 * 7)
      CookieWrapper.setCookie("userAuth", data.user, expires)
      CookieWrapper.setCookie("jwt", data.jwt, expires)

      // setup axios to use auth from now
      AxiosInstance.defaults.headers.common['Authorization'] =
        'Bearer ' + data.jwt
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

type UpdateUsernamePayload = {
  user: {
    username: string
  }
}

type UpdatePasswordPayload = {
  user: {
    password: string
  }
}

export async function updatePassword(
  dispatch: React.Dispatch<AuthDispatch>,
  loginPayload: LoginPayload,
  updatePayload: UpdatePasswordPayload
) {
  const requestOptions = (payload) => {
    return {
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    }
  }

  try {
    const loginResponse = await AxiosInstance.post(
      `/login`,
      loginPayload,
      requestOptions(loginPayload)
    )
    const loginData = loginResponse.data
    if (loginData.user) {
      const response = await AxiosInstance.put(
        `/user`,
        updatePayload,
        requestOptions(updatePayload)
      )
      const data = response.data
      if (data.user) {
        dispatch({
          type: AuthReducerActions.LOGIN_SUCCESS,
          payload: data,
          error: {} as any,
        })
        let expires = new Date()
        expires = new Date(expires.getTime() + 1000 * 60 * 60 * 24 * 7)
        CookieWrapper.setCookie("userAuth", data.user, expires)
        return data
      }
    }

    return
  } catch (error) {
    console.log(error)
  }
}

export async function updateUsername(
  dispatch: React.Dispatch<AuthDispatch>,
  updatePayload: UpdateUsernamePayload
) {
  const requestOptions = {
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(updatePayload),
  }

  try {
    const response = await AxiosInstance.put(
      `/user`,
      updatePayload,
      requestOptions
    )
    const data = response.data
    if (data.user) {
      dispatch({
        type: AuthReducerActions.LOGIN_SUCCESS,
        payload: data,
        error: {} as any,
      })
      let expires = new Date()
      expires = new Date(expires.getTime() + 1000 * 60 * 60 * 24 * 7)
      CookieWrapper.setCookie("userAuth", data.user, expires)

      return data
    }

    return
  } catch (error) {
    console.log(error)
  }
}
