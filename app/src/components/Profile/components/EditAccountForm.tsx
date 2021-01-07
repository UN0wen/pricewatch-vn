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
import { updateUsername } from '../../../api/user'

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
  username: string
}

export default function EditAccountForm({callback}) {
  const userAuth = useAuthState()
  const classes = useStyles()
  const {
    register,
    handleSubmit,
    watch,
    errors,
    setValue,
  } = useForm<FormData>()
  const dispatch = useAuthDispatch()

  if (!userAuth.user) {
    return <div></div>
  }

  const user = userAuth.user
  const password = React.useRef({})
  password.current = watch('password', '')

  setValue('username', user.username, { shouldValidate: true })

  const onSubmit = handleSubmit(async (data) => {
    try {
      const user = {
        username: data.username,
      }

      console.log(data)
      const response = await updateUsername(dispatch, { user })
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
                autoComplete="username"
                name="username"
                variant="outlined"
                required
                fullWidth
                id="username"
                label="Username"
                error={errors.username ? true : false}
                helperText={errors.username ? errors.username.message : ''}
                autoFocus
                defaultValue={user.username || ''}
                inputRef={register({
                  required: 'A new username is required',
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
