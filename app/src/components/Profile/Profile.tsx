import {
  CircularProgress,
  Divider,
  Grid,
  makeStyles,
  Theme,
  Typography,
} from '@material-ui/core'
import React, { useEffect, useState } from 'react'
import { useHistory } from 'react-router-dom'
import { Item } from '../../api/models'
import { getAllUserItems } from '../../api/userItems'
import { useAuthState } from '../../contexts/context'
import Routes from '../../utils/routes'
import ItemCard from '../ItemCard'
import Account from './components/Account'

const useStyles = makeStyles((theme: Theme) => ({
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
  section: {
    height: '100%',
  },
  followed: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    marginBottom: theme.spacing(8),
  },
}))

export default function Profile() {
  const classes = useStyles()
  const userAuth = useAuthState()
  const history = useHistory()
  const [items, setItems] = useState<Item[]>([])
  const [loading, setLoading] = useState(true)
  if (!userAuth.user) {
    history.push(Routes.SIGNIN)
  }

  // Get user items
  useEffect(() => {
    const fetchData = async () => {
      setLoading(true)

      const result = await getAllUserItems()
      setItems(result || [])
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
          <Grid item xs={3} className={classes.section}>
            <Account />
          </Grid>
          <Divider orientation="vertical" flexItem />
          <Grid item xs>
            <Typography
              component="h1"
              variant="h3"
              className={classes.followed}
            >
              Followed Items
            </Typography>
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
                  <Grid item xs={4}>
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
