import { Box, Container, Link, Typography } from '@material-ui/core';
import React from 'react';
import './App.css';
import NavBar from '../NavBar'
import { Route, Switch } from 'react-router-dom';
import HomePage from '../Home';
import ProfilePage from '../Profile';
import SignUpPage from '../SignUp';
import LoginPage from '../Login';
import ItemPage from '../ItemPage';
import NotFound from '../NotFound';
import { SessionContext, UserContext, UserCtx } from '../../utils/sessions';
import axios from 'axios';


function App() {
  axios.defaults.baseURL = 'http://172.20.204.69:8080/api';
  const jwt = SessionContext.getJWT()
  axios.get('/user', {

  }).then((response) => {
    console.log(response.data);
    console.log(response.status);
    console.log(response.statusText);
    console.log(response.headers);
    console.log(response.config);
  })

  return (
    <div>
      <UserContext.Provider value={new UserCtx("", jwt, false)}>
      <NavBar />
      <main style={{margin: '0 15px'}}>
        <Switch>
          <Route exact path='/'   component={HomePage}/>
          <Route path='/profile' component={ProfilePage}/>
          <Route path='/login' component={LoginPage}/>
          <Route path='/signup' component={SignUpPage}/>
          <Route path='/item/:itemid' component={ItemPage}/>
          <Route component={NotFound} />
        </Switch>
      </main>
      </UserContext.Provider>
    </div>
  );
}

export default App;
 