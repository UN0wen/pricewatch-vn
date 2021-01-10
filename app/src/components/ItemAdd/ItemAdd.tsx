import React from 'react'
import { makeStyles, Theme, createStyles } from '@material-ui/core/styles'
import Stepper from '@material-ui/core/Stepper'
import Step from '@material-ui/core/Step'
import StepLabel from '@material-ui/core/StepLabel'
import StepContent from '@material-ui/core/StepContent'
import Button from '@material-ui/core/Button'
import Typography from '@material-ui/core/Typography'
import { useHistory } from 'react-router-dom'
import { Item } from '../../api/models'
import ItemForm from './components/itemForm'
import { createItem, getItemFromURL } from '../../api/item'
import { Paper, TextField } from '@material-ui/core'
import { useForm } from 'react-hook-form'
import { isValidHttpUrl } from '../../utils/validate'
import { useAuthState } from '../../contexts/context'
import Routes from '../../utils/routes'

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      margin: theme.spacing(2)
    },
    button: {
      marginTop: theme.spacing(1),
      marginRight: theme.spacing(1),
    },
    actionsContainer: {
      marginBottom: theme.spacing(2),
    },
    resetContainer: {
      padding: theme.spacing(3),
    },
  })
)

function getSteps() {
  return ['Choose URL', 'Confirm Item', 'Finish']
}

function getStepContent(step: number) {
  switch (step) {
    case 0:
      return `Choose an URL to track prices for. We currently support: Tiki.vn, Shopee.vn .`
    case 1:
      return 'Check the item found by the server.'
    case 2:
      return `Your item has been successfully tracked. Click the button below to go to the item's page.`
    default:
      return 'Unknown step'
  }
}

type FormData = {
  url: string
}

export default function ItemAdd() {
  const classes = useStyles()
  const history = useHistory()
  const userAuth = useAuthState()

  const { register, handleSubmit, setError, errors } = useForm<FormData>()

  const [activeStep, setActiveStep] = React.useState(0)
  const [returnedItem, setItem] = React.useState<Item>({} as any)
  const steps = getSteps()

  if (!userAuth.user) {
    // TODO: error modal
    history.push(Routes.SIGNIN)
  }

  const handleClick = async () => {
    const responseItem = await createItem(returnedItem)
    if (!responseItem) {
      // set error
    } else if (responseItem == null) {
      return
    } else {
      setItem(responseItem)
      handleNext()
    }
  }

  const getStepButton = (step) => {
    switch (step) {
      case 0:
        return (
          <Button
            variant="contained"
            color="primary"
            type="submit"
            form="url-form"
            className={classes.button}
          >
            Submit URL
          </Button>
        )
      case 1:
        return (
          <Button
            variant="contained"
            color="primary"
            type="submit"
            onClick={handleClick}
            className={classes.button}
          >
            Confirm
          </Button>
        )
      case 2:
        return (
          <Button
            variant="contained"
            color="primary"
            onClick={handleNext}
            className={classes.button}
          >
            {activeStep === steps.length - 1 ? 'Finish' : 'Next'}
          </Button>
        )
      default:
        return <div>Unknown step</div>
    }
  }

  const handleNext = () => {
    if (activeStep === steps.length - 1) {
      history.push(`/item/${returnedItem.id}`)
      return
    }
    setActiveStep((prevActiveStep) => prevActiveStep + 1)
  }

  const handleBack = () => {
    setActiveStep((prevActiveStep) => prevActiveStep - 1)
  }

  const onSubmit = handleSubmit(async ({ url }) => {
    const item = {
      url,
    }
    if (!isValidHttpUrl(url)) {
      setError('url', {
        type: 'required',
        message: 'The URL must be valid.',
      })
      return
    }

    const responseItem = await getItemFromURL({ item })
    console.log(responseItem)

    if (typeof responseItem === 'number') {
      if (responseItem === 501) {
        setError('url', {
          type: 'required',
          message: 'This site is not currently supported.',
        })
      } else {
        setError('url', {
          type: 'required',
          message: 'This item cannot be found.',
        })
      }
    } else if (responseItem == null) {
      return
    } else {
      console.log(responseItem)
      setItem(responseItem)
      handleNext()
    }
  })

  const handleSteps = (step) => {
    switch (step) {
      case 0:
        return (
          <form noValidate onSubmit={onSubmit} id="url-form">
            <TextField
              variant="outlined"
              margin="normal"
              required
              fullWidth
              id="url"
              label="URL"
              name="url"
              autoComplete="url"
              inputRef={register({
                required: 'An URL is required',
              })}
              error={errors.url ? true : false}
              helperText={errors.url ? errors.url.message : ''}
              autoFocus
            />
          </form>
        )
      case 1:
        return (
          <React.Fragment>
            {returnedItem ? (
              <ItemForm item={returnedItem} />
            ) : (
              <div>An error has occured. Please go back to the first step.</div>
            )}
          </React.Fragment>
        )
      default:
        break
    }
  }
  return (
    <Paper className={classes.root}>
      <Stepper activeStep={activeStep} orientation="vertical">
        {steps.map((label, index) => (
          <Step key={label}>
            <StepLabel>{label}</StepLabel>
            <StepContent>
              <Typography>{getStepContent(index)}</Typography>
              {handleSteps(activeStep)}
              <div className={classes.actionsContainer}>
                <div>
                  <Button
                    disabled={activeStep === 0}
                    onClick={handleBack}
                    className={classes.button}
                  >
                    Back
                  </Button>
                  {getStepButton(activeStep)}
                </div>
              </div>
            </StepContent>
          </Step>
        ))}
      </Stepper>
    </Paper>
  )
}
