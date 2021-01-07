import React from 'react'
import { User } from '../api/models'
import { AuthDispatch, AuthReducer, initialState } from './reducer'

const anything = {} as any
const AuthStateContext = React.createContext<UserAuth>(initialState)
const AuthDispatchContext = React.createContext<React.Dispatch<AuthDispatch>>(
  anything
)

// interface for userauth objects stored in cookies
export interface UserAuth {
  user: User
  token: string
  loading: boolean
  errorMessage: any
}

export function useAuthState() {
  const context = React.useContext(AuthStateContext)
  if (context === undefined) {
    throw new Error('useAuthState must be used within a AuthProvider')
  }

  return context
}

export function useAuthDispatch() {
  const context = React.useContext(AuthDispatchContext)
  if (context === undefined) {
    throw new Error('useAuthDispatch must be used within a AuthProvider')
  }

  return context
}

export const AuthProvider = ({ children }) => {
  const [user, dispatch] = React.useReducer(AuthReducer, initialState)

  return (
    <AuthStateContext.Provider value={user}>
      <AuthDispatchContext.Provider value={dispatch}>
        {children}
      </AuthDispatchContext.Provider>
    </AuthStateContext.Provider>
  )
}
