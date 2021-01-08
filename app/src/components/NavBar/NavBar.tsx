import React from 'react'
import { fade, makeStyles, Theme, createStyles } from '@material-ui/core/styles'
import HomeIcon from '@material-ui/icons/Home'
import SearchIcon from '@material-ui/icons/Search'
import AccountCircle from '@material-ui/icons/AccountCircle'
import KeyboardArrowUp from '@material-ui/icons/KeyboardArrowUp'
import {
  AppBar,
  Toolbar,
  IconButton,
  Typography,
  InputBase,
  MenuItem,
  Menu,
  Divider,
  Button,
  Fab,
} from '@material-ui/core'
import { useHistory } from 'react-router-dom'
import Routes from '../../utils/routes'
import { useAuthState } from '../../contexts/context'
import ScrollTop from './components/ScrollTop'
import { AddCircleOutlined } from '@material-ui/icons'

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    grow: {
      flexGrow: 1,
    },
    homeButton: {
      marginRight: theme.spacing(2),
    },
    title: {
      display: 'none',
      [theme.breakpoints.up('sm')]: {
        display: 'block',
      },
    },
    root: {
      margin: theme.spacing(1, 1, 1, 4),
      padding: theme.spacing(1),
      display: 'flex',
      alignItems: 'center',
      width: '75%',
      height: 40,
      background: fade(theme.palette.common.white, 0.15),
    },
    input: {
      marginLeft: theme.spacing(1),
      flex: 1,
    },
    iconButton: {
      padding: 10,
    },
    section: {
      display: 'none',
      [theme.breakpoints.up('md')]: {
        display: 'flex',
      },
    },
    username: {
      padding: theme.spacing(0, 2),
    },
  })
)

function isValidHttpUrl(string: string) {
  let url

  if (!string.startsWith('http')) {
    string = 'https://' + string
  }

  try {
    url = new URL(string)
  } catch (_) {
    return false
  }

  return url.protocol === 'http:' || url.protocol === 'https:'
}

export default function NavBar() {
  const history = useHistory()
  const classes = useStyles()
  const userAuth = useAuthState()
  const [value, setValue] = React.useState('')
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null)

  const isMenuOpen = Boolean(anchorEl)

  const handleProfileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleMenuClose = () => {
    setAnchorEl(null)
  }

  const onSubmit = (e) => {
    e.preventDefault()
    if (isValidHttpUrl(value)) {
      history.push(`/item/add?url=${value}`)
    }
    history.push(`/search?q=${value}`)
    setValue('')
  }

  const menuId = 'primary-search-account-menu'
  const renderMenu = (
    <Menu
      anchorEl={anchorEl}
      getContentAnchorEl={null}
      anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      id={menuId}
      keepMounted
      transformOrigin={{ vertical: 'top', horizontal: 'center' }}
      open={isMenuOpen}
      onClose={handleMenuClose}
    >
      <MenuItem
        onClick={() => {
          history.push(Routes.PROFILE)
          handleMenuClose()
        }}
      >
        Signed in as {userAuth.user.username}
      </MenuItem>
      <Divider />
      <MenuItem
        onClick={() => {
          history.push(Routes.PROFILE)
          handleMenuClose()
        }}
      >
        Profile
      </MenuItem>
      <MenuItem
        onClick={() => {
          history.push(Routes.ADDITEM)
          handleMenuClose()
        }}
      >
        Add New Item
      </MenuItem>
      <MenuItem
        onClick={() => {
          history.push(Routes.SIGNOUT)
          handleMenuClose()
        }}
      >
        Sign Out
      </MenuItem>
    </Menu>
  )

  return (
    <div className={classes.grow}>
      <AppBar position="sticky" color="default">
        <Toolbar id="back-to-top-anchor">
          <IconButton
            edge="start"
            className={classes.homeButton}
            color="inherit"
            aria-label="home"
            onClick={() => {
              history.push('/')
            }}
          >
            <HomeIcon />
          </IconButton>
          <Typography className={classes.title} variant="h6" noWrap>
            PriceWatch-VN
          </Typography>
          <form onSubmit={onSubmit} className={classes.root}>
            <InputBase
              className={classes.input}
              placeholder="Search or URL to add new item (e.g. tiki.vn, shopee.vn)"
              inputProps={{ 'aria-label': 'search' }}
              onChange={(event) => {
                setValue(event.target.value)
              }}
            />
            <IconButton
              type="submit"
              className={classes.iconButton}
              aria-label="search"
            >
              <SearchIcon />
            </IconButton>
          </form>
          <div className={classes.grow} />
          <div className={classes.section}>
            {userAuth.user ? (
              <div>
                <IconButton
                  edge="end"
                  aria-label="add new item"
                  aria-controls={menuId}
                  aria-haspopup="true"
                  onClick={() => {history.push(Routes.ADDITEM)}}
                  color="inherit"
                >
                  <AddCircleOutlined />
                </IconButton>
                <IconButton
                  edge="end"
                  aria-label="account of current user"
                  aria-controls={menuId}
                  aria-haspopup="true"
                  onClick={handleProfileMenuOpen}
                  color="inherit"
                >
                  <AccountCircle />
                </IconButton>
              </div>
            ) : (
              <div>
                <Button
                  color="inherit"
                  onClick={() => {
                    history.push(Routes.SIGNUP)
                    handleMenuClose()
                  }}
                >
                  Sign Up
                </Button>
                <Button
                  color="inherit"
                  onClick={() => {
                    history.push(Routes.SIGNIN)
                    handleMenuClose()
                  }}
                >
                  Sign In
                </Button>
              </div>
            )}
          </div>
        </Toolbar>
      </AppBar>
      {renderMenu}
      <ScrollTop>
        <Fab color="secondary" size="small" aria-label="scroll back to top">
          <KeyboardArrowUp />
        </Fab>
      </ScrollTop>
    </div>
  )
}
