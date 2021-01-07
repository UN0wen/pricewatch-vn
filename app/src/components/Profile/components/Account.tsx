import {
  Button,
  Grid,
  Theme,
  Typography,
} from '@material-ui/core'
import { AccountCircle, Email } from '@material-ui/icons'
import React, { useState } from 'react'
import { useAuthState } from '../../../contexts/context'
import EditPasswordForm from './EditPasswordForm'
import EditAccountForm from './EditAccountForm'
import { makeStyles,  createStyles} from '@material-ui/core/styles'

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    avatar: {
      width: 200,
      height: 200,
      margin: 'auto',
    },
    centeredTop: {
      textAlign: 'center',
      padding: theme.spacing(1),
    },
    editButton: {
      margin: theme.spacing(1),
    },
    emailLine: {
      display: 'flex',
      alignItems: 'center',
      flexWrap: 'wrap',
      justifyContent: 'center',
      margin: theme.spacing(1)
    },
    emailString: {
      marginLeft: theme.spacing(1),
    },
  })
)

export default function Account() {
  const userAuth = useAuthState()
  const classes = useStyles()
  const [edit, setEdit] = useState(0)

  if (!userAuth.user) {
    return <div></div>
  }
  const user = userAuth.user

  const onCancel = () => {
    setEdit(0)
  }

  const onClick = (value) => {
    setEdit(value)
  }

  return (
    <div>
      <div className={classes.centeredTop}>
        <AccountCircle className={classes.avatar} />
        <Typography component="h1" align="center" variant="h5">
          {user.username}
        </Typography>

        <Typography className={classes.emailLine}>
          <Email />
          <span className={classes.emailString}>{user.email}</span>
        </Typography>

        <Grid container direction="row" spacing={2}>
          <Grid item xs={6}>
            <Button
              variant="outlined"
              className={classes.editButton}
              onClick={() => {
                onClick(1)
              }}
            >
              Edit username
            </Button>
          </Grid>
          <Grid item xs={6}>
            <Button
              variant="outlined"
              className={classes.editButton}
              onClick={() => {
                onClick(2)
              }}
            >
              Edit password
            </Button>
          </Grid>
        </Grid>

        {edit == 1 ? <EditAccountForm callback={onCancel} /> : <div />}
        {edit == 2 ? <EditPasswordForm callback={onCancel}/> : <div />}
      </div>
    </div>
  )
}
