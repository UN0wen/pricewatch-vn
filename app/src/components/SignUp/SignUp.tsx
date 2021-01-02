import React from 'react'
import LockOutlinedIcon from '@material-ui/icons/LockOutlined'
import { makeStyles } from '@material-ui/core/styles'

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
import { useForm } from 'react-hook-form'
import { useAuthDispatch } from '../../contexts/context'
import { createUser } from '../../contexts/actions'
import SignUpModal from './components/SignUpModal'

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
  username: string
  email: string
  password: string
  repeat_password: string
}

const defaultValues: FormData = {
  username: '',
  email: '',
  password: '',
  repeat_password: '',
}

export default function SignUp() {
  const classes = useStyles()
  const { register, handleSubmit, watch, errors, setError } = useForm<FormData>({
    defaultValues,
  })
  const dispatch = useAuthDispatch()
  const [open, setOpen] = React.useState(false)
  const password = React.useRef({})
  password.current = watch('password', '')

  const onSubmit = handleSubmit(async (data) => {
    try {
      console.log('here')
      const user = {
        username: data.username,
        email: data.email,
        password: data.password,
      }

      console.log(data)
      const response = await createUser(dispatch, { user })
      if (!response.user) return // TODO: error handling
      setOpen(!open)
    } catch (error) {
      setError("email", {
        type: "notMatch", 
        message: "Email already in use. Please sign in or choose another email address."
      });

    }
  })

  return (
    <Container component="main" maxWidth="xs">
      <CssBaseline />

      <SignUpModal open={open} />
      <div className={classes.paper}>
        <Avatar className={classes.avatar}>
          <LockOutlinedIcon />
        </Avatar>
        <Typography component="h1" variant="h5">
          Sign up
        </Typography>
        <form className={classes.form} noValidate onSubmit={onSubmit}>
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <TextField
                autoComplete="username"
                name="username"
                variant="outlined"
                required
                fullWidth
                id="username"
                label="Username"
                error={errors.username ? true : false}
                helperText={errors.username ? errors.username.message : ""}
                autoFocus
                inputRef={register({
                  required: 'A username is required',
                })}
              />

              
            </Grid>
            <Grid item xs={12}>
              <TextField
                variant="outlined"
                required
                fullWidth
                id="email"
                label="Email Address"
                name="email"
                autoComplete="email"
                error={errors.email ? true : false}
                helperText={errors.email ? errors.email.message : ""}
                inputRef={register({
                  required: "An email is required"
                })}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                variant="outlined"
                required
                fullWidth
                name="password"
                label="Password"
                type="password"
                id="password"
                error={errors.password ? true : false}
                helperText={errors.password ? errors.password.message : ""}
                inputRef={register({
                  required: 'A password is required',
                  minLength: {
                    value: 6,
                    message: 'Password must have at least 6 characters',
                  },
                })}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                variant="outlined"
                required
                fullWidth
                name="repeat_password"
                label="Repeat Password"
                type="password"
                id="repeat_password"
                error={errors.repeat_password ? true : false}
                helperText={errors.repeat_password ? errors.repeat_password.message : ""}
                inputRef={register({
                  validate: (value) =>
                    value === password.current || 'The passwords do not match',
                })}
              />
            </Grid>
          </Grid>
          <Button
            type="submit"
            fullWidth
            variant="contained"
            color="primary"
            className={classes.submit}
          >
            Sign Up
          </Button>
          <Grid container justify="flex-end">
            <Grid item>
              <Link href={Routes.SIGNIN} variant="body2">
                Already have an account? Sign in
              </Link>
            </Grid>
          </Grid>
        </form>
      </div>
    </Container>
  )
}
