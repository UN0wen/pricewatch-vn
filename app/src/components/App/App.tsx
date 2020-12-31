import { Box, Container, Link, Typography } from '@material-ui/core';
import React from 'react';
import './App.css';
import NavBar from '../NavBar'
import { Route, Switch } from 'react-router-dom';
import HomePage from '../HomePage';
import ProfilePage from '../ProfilePage';
import SignIn from '../SignIn';

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
          <Route path='/signup' component={SignIn}/>
        </Switch>
      </main>
      <Copyright/>
    </div>
  );
}

export default App;
 