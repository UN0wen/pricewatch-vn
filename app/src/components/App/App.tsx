import React from 'react'
import './App.css'
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

function App() {
  // Get current cookie
  const currentUser = CookieWrapper.getCookie()
  AxiosInstance.defaults.headers.common['Authorization'] = currentUser
    ? "Bearer " + currentUser.jwt
    : ''

  // Initial check for current session
  const dispatch = useAuthDispatch()
  AxiosInstance.get('/user')
    .catch((error) => {
      console.log(error)
      logout(dispatch)
    })

  return (
    <div>
      <NavBar />
      <main style={{ margin: '0 15px' }}>
        <Switch>
          <Route exact path={Routes.HOME} component={Home} />
          <Route path={Routes.PROFILE} component={Profile} />
          <Route path={Routes.SIGNIN} component={SignIn} />
          <Route path={Routes.SIGNUP} component={SignUp} />
          <Route path={Routes.ITEM} component={ItemPage} />
          <Route path={Routes.SIGNOUT} component={SignOut} />
          <Route component={NotFound} />
        </Switch>
      </main>
    </div>
  )
}

export default App
