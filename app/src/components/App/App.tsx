import React from 'react'
import './App.css'
import NavBar from '../NavBar'
import { Route, Switch } from 'react-router-dom'
import ItemPage from '../ItemPage'
import NotFound from '../NotFound'
import axios from 'axios'
import { AuthProvider } from '../../contexts/context'
import Routes from '../../utils/routes'
import SignIn from '../SignIn'
import SignUp from '../SignUp'
import Profile from '../Profile'
import Home from '../Home'
import SignOut from '../SignOut'

function App() {
  axios.defaults.baseURL = 'http://172.20.204.69:8080/api'

  return (
    <AuthProvider>
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
    </AuthProvider>
  )
}

export default App
