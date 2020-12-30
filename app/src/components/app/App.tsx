import { Box, Container, Link, Typography } from '@material-ui/core';
import React from 'react';
import './App.css';
import ProTip from './components/ProTips';

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
    <Container maxWidth="sm">
    <Box my={4}>
      <Typography variant="h4" component="h1" gutterBottom>
        Create React App v4-beta example with TypeScript
      </Typography>
      <ProTip />
      <Copyright />
    </Box>
  </Container>
  );
}

export default App;
