import { Backdrop, Dialog, Fade, Theme } from '@material-ui/core'
import React, { useEffect } from 'react'
import { useState } from 'react'
import { useHistory } from 'react-router-dom'
import Routes from '../../../utils/routes'
import { makeStyles, createStyles } from '@material-ui/core/styles'

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

export default function SignUpModal(props) {
  const classes = useStyles()
  const history = useHistory()
  const [open, setOpen] = useState(false)

  // change state on props change
  useEffect(() => {
    setOpen(props.open)
  }, [props.open])

  const handleEnter = () => {
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
            <h2 id="transition-modal-title">Successfully signed up</h2>
            <p id="transition-modal-description">
              You will be redirected to the home page shortly...
            </p>
          </div>
        </Fade>
      </Dialog>
    </div>
  )
}
