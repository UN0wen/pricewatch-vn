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

function Copyright() {
  return (
    <Typography variant="body2" color="textSecondary" align="center">
      {'Copyright © '}
      <Link color="inherit" href="https://material-ui.com/">
        Your Website
      </Link>{' '}
      {new Date().getFullYear()}
      {'.'}
    </Typography>
  );
}

function App() {
  return (
    <div>
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
      <Copyright/>
    </div>
  );
}

export default App;
 