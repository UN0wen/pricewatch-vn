import {
  Card,
  CardMedia,
  CardContent,
  Typography,
  Divider,
  CardActions,
  Button,
  Paper,
} from '@material-ui/core'
import { createStyles, makeStyles, Theme } from '@material-ui/core/styles'
import React from 'react'
import empty from '../../images/empty.jpg'
import { Item } from '../../../api/models'

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    grow: {
      flexGrow: 1,
      height: '100%',
    },
    paper: {
      marginTop: theme.spacing(2),
      padding: theme.spacing(2),
      display: 'flex',
      flexDirection: 'column',
      flexGrow: 1,
      alignItems: 'center',
      margin: theme.spacing(2),
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
    },
    buttonArea: {
      flexGrow: 1,
    },
    titleContent: {
      display: 'flex',
      flexGrow: 5,
      flexDirection: 'column',
      justifyContent: 'space-between',
      alignContent: 'center',
    },
    content: {
      display: 'flex',
      flex: '1 0 auto',
      flexDirection: 'row',
      justifyContent: 'space-between',
      width: '100%',
    },
    title: {
      wordWrap: 'break-word',
      margin: theme.spacing(1),
      justifyContent: 'center',
      display: 'flex',
    },
    titleSubtile: {
      display: 'flex',
      alignItems: 'center',
      flexDirection: 'column',
    },
    cover: {
      width: 400,
      height: 400,
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
    vnd: {
      display: 'flex',
      justifyContent: 'flex-end',
      alignItems: 'flex-end',
    },
    chart: {
      display: 'flex',
      width: '90%',
    },
    description: {
      display: 'flex',
      margin: theme.spacing(1, 0),
    },
  })
)

export default function ItemPage(props) {
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
    window.location.href = url
  }

  return (
    <div className={classes.grow}>
      <Paper className={classes.paper}>
        <Card className={classes.root}>
          <CardMedia className={classes.cover} image={imgURL} title={title} />

          <div className={classes.titleContentPrice}>
            <CardContent className={classes.titleContent}>
              <div className={classes.titleSubtile}>
                <Typography variant="h5" className={classes.title}>
                  {title}
                </Typography>
              </div>
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
      </Paper>
    </div>
  )
}
