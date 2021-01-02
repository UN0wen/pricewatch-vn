import { makeStyles, Theme, Typography } from '@material-ui/core'
import React from 'react'

const useStyles = makeStyles((theme: Theme) => ({
  grow: {
    flexGrow: 1,
  },
  paper: {
    marginTop: theme.spacing(4),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
  },
}))

export default function Account() {
  const classes = useStyles()
  return (
    <div className={classes.grow}>
      <div className={classes.paper}>
        <Typography component="h1" variant="h5">
          Account details
        </Typography>
      </div>
    </div>
  )
}
