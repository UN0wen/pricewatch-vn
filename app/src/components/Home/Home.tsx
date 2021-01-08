import { CircularProgress, Grid,  Theme } from '@material-ui/core'
import React, { useEffect, useState } from 'react'
import { getAllItems } from '../../api/item'
import { ItemWithPrice } from '../../api/models'
import ItemCard from '../ItemCard'
import { makeStyles,  createStyles} from '@material-ui/core/styles'

const useStyles = makeStyles((theme: Theme) => createStyles({
  grow: {
    flexGrow: 1,
  },
  paper: {
    marginTop: theme.spacing(1),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    padding: theme.spacing(2)
  },
}))

export default function Home() {
  const classes = useStyles()
  const [items, setItems] = useState<ItemWithPrice[]>([])
  const [loading, setLoading] = useState(true)

  // Get items
  useEffect(() => {
    const fetchData = async () => {
      setLoading(true)

      const result = await getAllItems()
      setItems(result || [])
      console.log(result)
      setLoading(false)
    }

    fetchData()
  }, [])
  
  return (
    <div className={classes.grow}>
      <div className={classes.paper}>
        <Grid
          container
          direction="row"
          justify="center"
          alignItems="stretch"
          alignContent="center"
          spacing={3}
        >
          <Grid item xs>
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
                  <Grid item xs={3}>
                    <ItemCard key={idx} {...val} />{' '}
                  </Grid>
                ))
              )}
            </Grid>
          </Grid>
        </Grid>
      </div>
    </div>
  )
}
