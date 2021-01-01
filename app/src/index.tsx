import React from 'react'
import ReactDOM from 'react-dom'
import './index.css'
import App from './components/App'
import { CssBaseline, ThemeProvider } from '@material-ui/core'
import theme from './theme'
import { BrowserRouter } from 'react-router-dom'
import { AuthProvider } from './contexts/context'

ReactDOM.render(
  <ThemeProvider theme={theme}>
    {/* CssBaseline kickstart an elegant, consistent, and simple baseline to build upon. */}
    <CssBaseline />
    <AuthProvider>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </AuthProvider>
  </ThemeProvider>,
  document.querySelector('#root')
)
