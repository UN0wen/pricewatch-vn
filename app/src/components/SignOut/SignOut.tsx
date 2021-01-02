import React, { useState } from 'react'
import { makeStyles, Theme, createStyles } from '@material-ui/core/styles'
import Backdrop from '@material-ui/core/Backdrop'
import Fade from '@material-ui/core/Fade'
import { useHistory } from 'react-router-dom'
import Routes from '../../utils/routes'
import { Dialog, LinearProgress } from '@material-ui/core'
import { logout } from '../../contexts/actions'
import { useAuthDispatch } from '../../contexts/context'

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    modal: {
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
    },
    paper: {
      backgroundColor: theme.palette.background.paper,
      boxShadow: theme.shadows[5],
      padding: theme.spacing(2, 4, 3),
    },
  })
)

export default function SignOut() {
  const classes = useStyles()
  const dispatch = useAuthDispatch()
  const history = useHistory()
  const [open, setOpen] = useState(true)

  const handleEnter = () => {
    logout(dispatch)
    setTimeout(() => {
      setOpen(false)
    }, 1000)
  }

  const handleExit = () => {
    history.push(Routes.HOME)
  }

  return (
    <div>
      <Dialog
        aria-labelledby="transition-modal-title"
        aria-describedby="transition-modal-description"
        className={classes.modal}
        open={open}
        closeAfterTransition
        BackdropComponent={Backdrop}
        BackdropProps={{
          timeout: 500,
        }}
      >
        <Fade in={open} onEnter={handleEnter} onExited={handleExit}>
          <div className={classes.paper}>
            <h2 id="transition-modal-title">Successfully logged out</h2>
            <p id="transition-modal-description">
              You will be redirected to the home page shortly...
            </p>
            <LinearProgress/>
          </div>
        </Fade>
      </Dialog>
    </div>
  )
}
