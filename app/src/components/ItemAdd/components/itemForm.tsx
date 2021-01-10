import {
  Card,
  CardMedia,
  CardContent,
  Typography,
  Divider,
  CardActions,
  Button,
} from '@material-ui/core'
import { createStyles, makeStyles, Theme } from '@material-ui/core/styles'
import React from 'react'
import empty from '../../../images/empty.jpg'
import { Item } from '../../../api/models'

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    grow: {
      flexGrow: 1,
    },
    root: {
      display: 'flex',
      flexDirection: 'row',
      margin: theme.spacing(1),
      flexGrow: 1,
    },
    titleContentPrice: {
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'space-between',
      alignContent: 'center',
      width: '75%'
    },
    buttonArea: {
      flexGrow: 1,
    },
    titleContent: {
      display: 'flex',
      flexGrow: 5,
      flexDirection: 'column',
      justifyContent: 'flex-start',
      alignContent: 'center',
    },
    title: {
      wordWrap: 'break-word',
      margin: theme.spacing(1),
      justifyContent: 'center',
      display: 'flex',
    },
    cover: {
      display: 'block',
      width: 400,
      height: 400
    },
    button: {
      height: '100%',
    },
    cardAction: {
      padding: 0,
      height: '100%',
    },
    text: {
      display: 'flex',
      justifyContent: 'flex-end',
      alignItems: 'flex-end',
      flexGrow: 1,
      marginRight: theme.spacing(1),
    },
    description: {
      display: 'flex',
      margin: theme.spacing(1, 0),
    },
  })
)

export default function ItemForm(props) {
  const classes = useStyles()
  let item: Item

  if (props.item) {
    item = props.item
  } else {
    return <div>Choose an item!</div>
  }
  const title = item?.name || 'No title'
  const imgURL = item?.image_url || empty
  const url = item?.url || '/'
  const description = item?.description || 'No description'
  const onClickStore = () => {
    window.open(url)
  }

  return (
    <div className={classes.grow}>
        <Card className={classes.root} variant="outlined">
          <CardMedia className={classes.cover} image={imgURL} title={title} />

          <div className={classes.titleContentPrice}>
            <CardContent className={classes.titleContent}>
                <Typography variant="h5" className={classes.title}>
                  {title}
                </Typography>
              <Divider
                orientation="horizontal"
                className={classes.description}
              />

              <Typography className={classes.description}>
                {description}
              </Typography>
            </CardContent>
            <div className={classes.buttonArea}>
              <Divider orientation="horizontal" />
              <CardActions className={classes.cardAction} disableSpacing>
                <Button
                  onClick={onClickStore}
                  className={classes.button}
                  fullWidth
                >
                  To store page
                </Button>
              </CardActions>
            </div>
          </div>
        </Card>
    </div>
  )
}
