import React, { useContext } from "react";
import {
  fade,
  makeStyles,
  Theme,
  createStyles,
} from "@material-ui/core/styles";
import HomeIcon from "@material-ui/icons/Home";
import SearchIcon from "@material-ui/icons/Search";
import AccountCircle from "@material-ui/icons/AccountCircle";

import {
  AppBar,
  Toolbar,
  IconButton,
  Typography,
  InputBase,
  Badge,
  MenuItem,
  Menu,
  ListItem,
  Divider,
  Button,
} from "@material-ui/core";
import { useHistory } from "react-router-dom";
import Routes from "../../utils/routes";
import { UserContext } from "../../utils/sessions";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    grow: {
      flexGrow: 1,
    },
    homeButton: {
      marginRight: theme.spacing(2),
    },
    title: {
      display: "none",
      [theme.breakpoints.up("sm")]: {
        display: "block",
      },
    },
    search: {
      position: "relative",
      marginRight: theme.spacing(2),
      marginLeft: theme.spacing(1),
      borderRadius: 2,
      background: fade(theme.palette.common.white, 0.15),
      "&:hover": {
        background: fade(theme.palette.common.white, 0.25),
      },
      "& $inputInput": {
        transition: theme.transitions.create("width"),
        width: 120,
        [theme.breakpoints.down("sm")]: {
          width: "100%",
        },
        "&:focus": {
          width: 400,
          [theme.breakpoints.down("sm")]: {
            width: "100%",
          },
        },
      },
      [theme.breakpoints.down("sm")]: {
        marginLeft: 0,
        marginRight: 0,
      },
    },
    searchIcon: {
      padding: theme.spacing(0, 2),
      height: "100%",
      position: "absolute",
      pointerEvents: "none",
      display: "flex",
      alignItems: "center",
      justifyContent: "center",
    },
    inputRoot: {
      color: "inherit",
    },
    inputInput: {
      padding: theme.spacing(1, 1, 1, 0),
      // vertical padding + font size from searchIcon
      paddingLeft: `calc(1em + ${theme.spacing(4)}px)`,
      transition: theme.transitions.create("width"),
      width: "100%",
      [theme.breakpoints.up("md")]: {
        width: "20ch",
      },
    },
    section: {
      display: "none",
      [theme.breakpoints.up("md")]: {
        display: "flex",
      },
    },
    username: {
      padding: theme.spacing(0, 2),
    },
  })
);

export default function NavBar() {
  const history = useHistory();
  const classes = useStyles();
  const userCtx = useContext(UserContext);
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

  const isMenuOpen = Boolean(anchorEl);

  const handleProfileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const menuId = "primary-search-account-menu";
  const renderMenu = (
    <Menu
      anchorEl={anchorEl}
      getContentAnchorEl={null}
      anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
      id={menuId}
      keepMounted
      transformOrigin={{ vertical: "top", horizontal: "center" }}
      open={isMenuOpen}
      onClose={handleMenuClose}
    >
      <MenuItem
        onClick={() => {
          history.push(Routes.Profile);
          handleMenuClose();
        }}
      >
        Signed in as {}
      </MenuItem>
      <Divider />
      <MenuItem
        onClick={() => {
          history.push(Routes.Profile);
          handleMenuClose();
        }}
      >
        Profile
      </MenuItem>
      <MenuItem
        onClick={() => {
          history.push(Routes.Account);
          handleMenuClose();
        }}
      >
        My Account
      </MenuItem>
      <MenuItem
        onClick={() => {
          history.push(Routes.Logout);
          handleMenuClose();
        }}
      >
        Logout
      </MenuItem>
    </Menu>
  );

  return (
    <div className={classes.grow}>
      <AppBar position="static">
        <Toolbar>
          <IconButton
            edge="start"
            className={classes.homeButton}
            color="inherit"
            aria-label="home"
            onClick={() => {
              history.push("/");
            }}
          >
            <HomeIcon />
          </IconButton>
          <Typography className={classes.title} variant="h6" noWrap>
            PriceWatch-VN
          </Typography>
          <div className={classes.search}>
            <div className={classes.searchIcon}>
              <SearchIcon />
            </div>
            <InputBase
              placeholder="Search…"
              classes={{
                root: classes.inputRoot,
                input: classes.inputInput,
              }}
              inputProps={{ "aria-label": "search" }}
            />
          </div>
          <div className={classes.grow} />
          <div className={classes.section}>
            {userCtx.auth ? (
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
            ) : (
              <div>
                <Button color="inherit">Sign Up</Button>
                <Button color="inherit">Login</Button>
              </div>
            )}
          </div>
        </Toolbar>
      </AppBar>
      {renderMenu}
    </div>
  );
}
