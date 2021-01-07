import {
  Card,
  CardHeader,
  CardMedia,
  CardContent,
  Typography,
  Divider,
  CardActions,
  Button,
  Paper,
} from '@material-ui/core'
import { createStyles, makeStyles, Theme } from '@material-ui/core/styles'
import { formatDistanceToNow, parseISO } from 'date-fns'
import React, { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { getItem, getItemPrices } from '../../api/item'
import { ItemPrice, ItemWithPrice } from '../../api/models'
import empty from '../../images/empty.jpg'

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
      alignItems: 'center',
    },
    root: {
      display: 'flex',
      flexDirection: 'column',
      margin: theme.spacing(1),
      flexGrow: 1,
      maxWidth: 400,
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
      marginRight: theme.spacing(1),
    },
    vnd: {
      display: 'flex',
      justifyContent: 'flex-end',
      alignItems: 'flex-end',
    },
  })
)
interface ParamTypes {
  itemID: string
}
export default function ItemPage() {
  const classes = useStyles()

  const { itemID } = useParams<ParamTypes>()
  const [item, setItem] = useState<ItemWithPrice>({} as any)
  const [itemPrices, setItemPrices] = useState<ItemPrice[]>([])
  const [loading, setLoading] = useState(true)
  const [loadingPrices, setLoadingPrices] = useState(true)
  useEffect(() => {
    const fetchData = async () => {
      setLoading(true)

      const result = await getItem(itemID)
      setItem(result || ({} as any))
      setLoading(false)
    }

    fetchData()
  }, [itemID])

  useEffect(() => {
    const fetchData = async () => {
      setLoadingPrices(true)

      const result = await getItemPrices(itemID)
      setItemPrices(result || [])
      setLoadingPrices(false)
    }

    fetchData()
  }, [itemID])

  let title = item?.name || 'No title'
  let imgURL = item?.image_url || empty
  let url = item?.url || '/'
  let updated = item?.time ? parseISO(item?.time) : new Date()
  let price: number = item?.price || 0
  const onClickStore = () => {
    window.location.href = url
  }

  useEffect(() => {
    title = item?.name || 'No title'
    imgURL = item?.image_url || empty
    url = item?.url || '/'
    updated = item?.time ? parseISO(item?.time) : new Date()
    price = item?.price || 0
  }, [item])

  console.log(item)
  return (
    <div className={classes.grow}>
      <Paper className={classes.paper}>
        <Card className={classes.root}>
          <CardHeader
            title={title}
            subheader={`Last updated: ${formatDistanceToNow(updated, {
              addSuffix: true,
            })}`}
          />
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
          <Divider orientation="horizontal" />
          <CardActions className={classes.cardAction} disableSpacing>
            <Button onClick={onClickStore} className={classes.button} fullWidth>
              To store page
            </Button>
          </CardActions>
        </Card>
      </Paper>
    </div>
  )
}
