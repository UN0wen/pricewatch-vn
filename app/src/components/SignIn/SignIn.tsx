import React, { useState } from 'react'
import LockOutlinedIcon from '@material-ui/icons/LockOutlined'
import { makeStyles } from '@material-ui/core/styles'
import { useForm } from 'react-hook-form'

import {
  Avatar,
  Button,
  CssBaseline,
  TextField,
  Link,
  Grid,
  Typography,
  Container,
} from '@material-ui/core'
import Routes from '../../utils/routes'
import { loginUser } from '../../api/user'
import { useAuthDispatch } from '../../contexts/context'
import SignInModal from './components/SignInModal'

const useStyles = makeStyles((theme) => ({
  paper: {
    marginTop: theme.spacing(8),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
  },
  avatar: {
    margin: theme.spacing(1),
    backgroundColor: theme.palette.secondary.main,
  },
  form: {
    width: '100%', // Fix IE 11 issue.
    marginTop: theme.spacing(1),
  },
  submit: {
    margin: theme.spacing(3, 0, 2),
  },
}))

type FormData = {
  email: string
  password: string
}

export default function SignIn() {
  const classes = useStyles()
  const { register, handleSubmit, setError , errors} = useForm<FormData>()
  const dispatch = useAuthDispatch()
  const [open, setOpen] = useState(false)

  const onSubmit = handleSubmit(async ({ email, password }) => {
    try {
      const user = {
        email,
        password,
      }
      const response = await loginUser(dispatch, { user })
      if (!response) return // TODO: error handling
      setOpen(true)
    } catch (error) {
      setError('email', {
        type: 'required',
      })
      setError('password', {
        type: 'required',
        message: 'Account does not exist or password does not match.',
      })
    }
  })

  return (
    <Container component="main" maxWidth="xs">
      <CssBaseline />

      <SignInModal open={open} />
      <div className={classes.paper}>
        <Avatar className={classes.avatar}>
          <LockOutlinedIcon />
        </Avatar>
        <Typography component="h1" variant="h5">
          Sign in
        </Typography>
        <form className={classes.form} noValidate onSubmit={onSubmit}>
          <TextField
            variant="outlined"
            margin="normal"
            required
            fullWidth
            id="email"
            label="Email Address"
            name="email"
            autoComplete="email"
            inputRef={register}
            error={errors.email ? true : false}
            helperText={errors.email ? errors.email.message : ''}
            autoFocus
          />
          <TextField
            variant="outlined"
            margin="normal"
            required
            fullWidth
            name="password"
            label="Password"
            type="password"
            id="password"
            autoComplete="current-password"
            error={errors.password ? true : false}
            helperText={errors.password ? errors.password.message : ''}
            inputRef={register}
          />
          <Button
            type="submit"
            fullWidth
            variant="contained"
            color="primary"
            className={classes.submit}
          >
            Sign In
          </Button>
          <Grid container>
            <Grid item>
              <Link href={Routes.SIGNUP} variant="body2">
                {"Don't have an account? Sign Up"}
              </Link>
            </Grid>
          </Grid>
        </form>
      </div>
    </Container>
  )
}
