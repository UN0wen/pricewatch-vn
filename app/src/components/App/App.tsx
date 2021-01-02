import React from 'react'
import NavBar from '../NavBar'
import { Route, Switch } from 'react-router-dom'
import ItemPage from '../ItemPage'
import NotFound from '../NotFound'
import Routes from '../../utils/routes'
import SignIn from '../SignIn'
import SignUp from '../SignUp'
import Profile from '../Profile'
import Home from '../Home'
import SignOut from '../SignOut'
import { AxiosInstance } from '../../utils/axios'
import { useAuthDispatch } from '../../contexts/context'
import { logout } from '../../contexts/actions'
import { CookieWrapper } from '../../utils/storage'
import { CssBaseline, ThemeProvider, useMediaQuery } from '@material-ui/core'
import { darkTheme, lightTheme } from '../../theme'

function App() {
  const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)')

  // Get current cookie
  const currentJWT = CookieWrapper.getCookie("jwt")
  AxiosInstance.defaults.headers.common['Authorization'] = currentJWT
    ? 'Bearer ' + currentJWT
    : ''

  // Initial check for current session
  const dispatch = useAuthDispatch()
  AxiosInstance.get('/user').catch(() => {
    console.log("error during user initialization, deleting cookies...")
    logout(dispatch)
  })

  const theme = prefersDarkMode ? darkTheme : lightTheme

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <NavBar />
      <Switch>
        <Route exact path={Routes.HOME} component={Home} />
        <Route path={Routes.PROFILE} component={Profile} />
        <Route path={Routes.SIGNIN} component={SignIn} />
        <Route path={Routes.SIGNUP} component={SignUp} />
        <Route path={Routes.ITEM} component={ItemPage} />
        <Route path={Routes.SIGNOUT} component={SignOut} />
        <Route component={NotFound} />
      </Switch>
    </ThemeProvider>
  )
}

export default App
