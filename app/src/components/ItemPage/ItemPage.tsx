import {
  ArgumentScale,
  EventTracker,
  ValueScale,
} from '@devexpress/dx-react-chart'
import {
  ArgumentAxis,
  Chart,
  LineSeries,
  ScatterSeries,
  Title,
  Tooltip,
  ValueAxis,
} from '@devexpress/dx-react-chart-material-ui'
import {
  Card,
  CardMedia,
  CardContent,
  Typography,
  Divider,
  CardActions,
  Button,
  Paper,
  CircularProgress,
} from '@material-ui/core'
import { createStyles, makeStyles, Theme } from '@material-ui/core/styles'
import { format, formatDistanceToNow, fromUnixTime, parseISO } from 'date-fns'
import React, { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { getItem, getItemPrices } from '../../api/item'
import { ItemPrice, ItemWithPrice } from '../../api/models'
import empty from '../../images/empty.jpg'
import { line, curveStepAfter } from 'd3-shape'

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
      flexDirection: 'column'
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
  })
)

interface ParamTypes {
  itemID: string
}

const Line = (props) => (
  <LineSeries.Path
    {...props}
    path={line()
      .x(({ arg }) => arg)
      .y(({ val }) => val)
      .curve(curveStepAfter)}
  />
)

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
      if (result) {
        result.sort((a, b) => (a.time > b.time ? 1 : -1))
        result.forEach((ip) => {
          ip.time = new Date(ip.time)
        })
      }
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

  const TooltipContent = ({ targetItem }) => (
    <div>{itemPrices[targetItem.point].price.toLocaleString()} VND</div>
  )

  const ArgumentLabel = ({ x, y, dy, text, textAnchor }) => {
    const res = text.replace(/[ ,.]/g, '')
    return (
      <Chart.Label x={x} y={y} dy={dy} textAnchor={textAnchor}>
        {format(fromUnixTime(res / 1000), 'h a, MMM Do')}
      </Chart.Label>
    )
  }
  return (
    <div className={classes.grow}>
      <Paper className={classes.paper}>
        {loading ? (
          <CircularProgress />
        ) : (
          <Card className={classes.root}>
            <CardMedia className={classes.cover} image={imgURL} title={title} />

            <div className={classes.titleContentPrice}>
              <CardContent className={classes.titleContent}>
                <div className={classes.titleSubtile}>
                  <Typography variant="h5" className={classes.title}>
                    {title}
                  </Typography>
                  <Typography variant="body1" color="textSecondary">
                    {`Last updated: ${formatDistanceToNow(updated, {
                      addSuffix: true,
                    })}`}
                  </Typography>
                </div>

                <div className={classes.content}>
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
                </div>
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
        )}

        {loadingPrices ? (
          <CircularProgress />
        ) : (
          <div className={classes.chart}>
            <Chart data={itemPrices} width={5000}>
              <ValueScale name="price" />
              <ArgumentScale name="time" />
              <ArgumentAxis labelComponent={ArgumentLabel} />
              <ValueAxis scaleName="price" showLine showTicks />
              <LineSeries
                name="Price over time"
                valueField="price"
                argumentField="time"
                scaleName="price"
                seriesComponent={Line}
              />
              <ScatterSeries valueField="price" argumentField="time" />{' '}
              <EventTracker />
              <Tooltip contentComponent={TooltipContent} />
              <Title text="Price over time" />
            </Chart>
          </div>
        )}
      </Paper>
    </div>
  )
}
