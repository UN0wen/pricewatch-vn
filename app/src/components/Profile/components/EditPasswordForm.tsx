import React from 'react'
import { makeStyles } from '@material-ui/core/styles'

import {
  Button,
  CssBaseline,
  TextField,
  Grid,
  Container,
} from '@material-ui/core'
import { useForm } from 'react-hook-form'
import { useAuthDispatch, useAuthState } from '../../../contexts/context'
import { updatePassword } from '../../../api/user'

const useStyles = makeStyles((theme) => ({
  paper: {
    marginTop: theme.spacing(2),
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
  current_password: string
  password: string
  repeat_password: string
}

export default function EditPasswordForm({callback}) {
  const userAuth = useAuthState()
  const classes = useStyles()
  const {
    register,
    handleSubmit,
    watch,
    errors,
  } = useForm<FormData>()
  const dispatch = useAuthDispatch()

  if (!userAuth.user) {
    return <div></div>
  }

  const user = userAuth.user
  const password = React.useRef({})
  password.current = watch('password', '')

  const onSubmit = handleSubmit(async (data) => {
    try {
      const validate = {
        email: String(user.email),
        password: data.current_password
      }
      const userPayload = {
        password: data.password,
      }

      console.log(data)
      const response = await updatePassword(dispatch, { "user": validate }, {"user": userPayload})
      if (!response) return // TODO: error handling
    } catch (error) {
      console.log(error)
    }
  })

  const handleCancel = () => {
    callback()
  }

  return (
    <Container component="main" maxWidth="xs">
      <CssBaseline />

      <div className={classes.paper}>
        <form className={classes.form} noValidate onSubmit={onSubmit}>
          <Grid container spacing={2}>
          <Grid item xs={12}>
              <TextField
                variant="outlined"
                required
                fullWidth
                name="current_password"
                label="Password"
                type="password"
                id="current_password"
                error={errors.current_password ? true : false}
                helperText={errors.current_password ? errors.current_password.message : ''}
                defaultValue=""
                inputRef={register({
                  required: 'The current password is required',
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
                helperText={errors.password ? errors.password.message : ''}
                defaultValue=""
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
                defaultValue=""
                error={errors.repeat_password ? true : false}
                helperText={
                  errors.repeat_password ? errors.repeat_password.message : ''
                }
                inputRef={register({
                  validate: (value) =>
                    value === password.current || 'The passwords do not match',
                })}
              />
            </Grid>

            <Grid item xs={3}>
              <Button
                type="submit"
                fullWidth
                variant="contained"
                color="primary"
                className={classes.submit}
              >
                Save
              </Button>
            </Grid>
            <Grid item xs={3}>
              <Button
                type="button"
                fullWidth
                variant="outlined"
                color="primary"
                className={classes.submit}
                onClick={handleCancel}
              >
                Cancel
              </Button>
            </Grid>
          </Grid>
        </form>
      </div>
    </Container>
  )
}
