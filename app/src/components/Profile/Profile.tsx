import { Divider, Grid, makeStyles, Theme, Typography } from '@material-ui/core'
import React from 'react'
import { useHistory } from 'react-router-dom'
import { useAuthState } from '../../contexts/context'
import Routes from '../../utils/routes'
import Account from './components/Account'

const useStyles = makeStyles((theme: Theme) => ({
  grow: {
    flexGrow: 1,
    height: '100%',
  },
  paper: {
    marginTop: theme.spacing(2),
    padding: theme.spacing(2),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
  },
  section: {
    height: '100%',
  },
}))

export default function Profile() {
  const classes = useStyles()
  const userAuth = useAuthState()
  const history = useHistory()

  if (!userAuth.user) {
    history.push(Routes.SIGNIN)
  }

  return (
    <div className={classes.grow}>
      <div className={classes.paper}>
        <Grid
          container
          direction="row"
          justify="center"
          alignItems="stretch"
          alignContent="center"
          spacing={3}
        >
          <Grid item xs={3} className={classes.section}>
            <Account />
          </Grid>
          <Divider orientation="vertical" flexItem/>
          <Grid item xs>
            <Typography component="h1" variant="h5" className={classes.section}>
              Followed Items
            </Typography>
          </Grid>
        </Grid>
      </div>
    </div>
  )
}
