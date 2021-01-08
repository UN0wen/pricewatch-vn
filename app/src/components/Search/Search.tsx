import {
  Button,
  CircularProgress,
  Grid,
  Theme,
  Typography,
} from '@material-ui/core'
import React, { useEffect, useState } from 'react'
import { search } from '../../api/item'
import { ItemWithPrice } from '../../api/models'
import ItemCard from '../ItemCard'
import { makeStyles, createStyles } from '@material-ui/core/styles'
import { useHistory, useLocation } from 'react-router-dom'

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    grow: {
      flexGrow: 1,
    },
    paper: {
      marginTop: theme.spacing(1),
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      padding: theme.spacing(2),
    },
    returnButton: {
      margin: theme.spacing(1)
    }
  })
)

// A custom hook that builds on useLocation to parse
// the query string for you.
function useQuery() {
  return new URLSearchParams(useLocation().search)
}

export default function Home() {
  const classes = useStyles()
  const [items, setItems] = useState<ItemWithPrice[]>([])
  const [loading, setLoading] = useState(true)
  const query = useQuery()
  const history = useHistory()

  const queryString = query.get('q')

  const dataAvailable = queryString != '' && items.length != 0
  // Get items
  useEffect(() => {
    const fetchData = async () => {
      setLoading(true)

      const result = await search(query)
      setItems(result || [])
      setLoading(false)
    }

    fetchData()
  }, [])

  console.log(dataAvailable)
  return (
    <div className={classes.grow}>
      <div className={classes.paper}>
        {!dataAvailable ? (
          <React.Fragment>
            <Typography >No items found.</Typography>

            <Button
              onClick={() => {
                history.push('/')
              }}
              variant="outlined"
              className={classes.returnButton}
            >
              Return to Homepage
            </Button>
          </React.Fragment>
        ) : (
          <Grid
            container
            direction="row"
            justify="center"
            alignItems="stretch"
            spacing={3}
          >
            {loading ? (
              <CircularProgress />
            ) : (
              items.map((val, idx) => (
                <Grid item xs={3} key={idx}>
                  <ItemCard key={idx} {...val} />{' '}
                </Grid>
              ))
            )}
          </Grid>
        )}
      </div>
    </div>
  )
}
