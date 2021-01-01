import { createContext } from 'react'
import Cookies from 'universal-cookie'

const cookies = new Cookies()

export const SessionContext = (function () {
  const getJWT = function (): string {
    return cookies.get('jwt')
  }

  const setJWT = function (jwt: string, expire: Date) {
    cookies.set('jwt', jwt, { path: '/', expires: expire })
  }

  return {
    getJWT: getJWT,
    setJWT: setJWT,
  }
})()

export class UserCtx {
  username = ''
  jwt = ''
  auth = false
  constructor(username: string, jwt: string, auth: boolean) {
    this.username = username
    this.jwt = jwt
    this.auth = auth
  }
}

export const UserContext = createContext<UserCtx>(
  new UserCtx('', SessionContext.getJWT(), false)
)
