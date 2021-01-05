import React from 'react'
import { Theme, createStyles, makeStyles } from '@material-ui/core/styles'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardMedia from '@material-ui/core/CardMedia'
import Typography from '@material-ui/core/Typography'
import empty from './assets/empty.jpg'
import {
  Button,
  CardActionArea,
  CardActions,
  CardHeader,
  Divider,
} from '@material-ui/core'
import { useHistory } from 'react-router-dom'

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      display: 'flex',
      flexDirection: 'column',
      margin: theme.spacing(1),
      flexGrow: 1,
      maxWidth: 345,
    },
    content: {
      display: 'flex',
      flex: '1 0 auto',
      flexDirection: 'row',
      justifyContent: 'space-between',
      width: '100%',
    },
    cover: {
      height: 0,
      paddingTop: '56.25%', // 16:9
    },
    button: {
      height: theme.spacing(8),
    },
    cardAction: {
      padding: 0,
    },
    text: {
      display: 'flex',
      justifyContent: 'flex-end',
      flexGrow: 1,
      marginRight: theme.spacing(1)
    },
    vnd: {
      display: 'flex',
      justifyContent: 'flex-end',
      alignItems:'flex-end'
    },
  })
)

export default function ItemCard(props) {
  const classes = useStyles()
  const history = useHistory()

  const title = props.title || 'No title'
  const imgURL = props.image_url || empty
  const url = props.url || '/'
  const id = props.id || ''
  const updated = props.updated || 'Unknown'
  const price: number = props.price || 0
  const onClickStore = () => {
    window.location = url
  }

  const onClickItem = () => {
    history.push(`/item/${id}`)
  }

  return (
    <Card className={classes.root}>
      <CardActionArea onClick={onClickItem}>
        <CardHeader title={title} subheader={`Last updated: ${updated}`} />
        <CardMedia className={classes.cover} image={imgURL} title={title} />

        <div>
          <CardContent className={classes.content}>
            <Typography variant="h3" className={classes.text}>
              {price.toLocaleString()}
            </Typography>
            <Typography
              variant="h6"
              color="textSecondary"
              className={classes.vnd}
              align="right"
            >
              VND
            </Typography>
          </CardContent>
        </div>
      </CardActionArea>
      <Divider orientation="horizontal" />
      <CardActions className={classes.cardAction} disableSpacing>
        <Button onClick={onClickStore} className={classes.button} fullWidth>
          To store page
        </Button>
        <Divider orientation="vertical" flexItem/>
        <Button onClick={onClickItem} className={classes.button} fullWidth>
          Price History
        </Button>
      </CardActions>
    </Card>
  )
}
