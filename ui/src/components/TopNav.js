import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';

import AppBar from "@material-ui/core/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import { withStyles } from "@material-ui/core/styles";
import { IconButton } from "@material-ui/core";
import MenuItem from '@material-ui/core/MenuItem';
import Menu from '@material-ui/core/Menu';
import Input from "@material-ui/core/Input";
import InputAdornment from "@material-ui/core/InputAdornment";
import blue from "@material-ui/core/colors/blue";
import Avatar from '@material-ui/core/Avatar';
import Chip from '@material-ui/core/Chip';
import Typography from '@material-ui/core/Typography';

//import MenuIcon from "mdi-material-ui/Menu";
//import Backburger from "mdi-material-ui/Backburger";
//import Wallet from "mdi-material-ui/Wallet";
import AccountCircle from "mdi-material-ui/AccountCircle";
import Magnify from "mdi-material-ui/Magnify";
import HelpCircle from "mdi-material-ui/HelpCircle";

import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import Wallet from "mdi-material-ui/WalletOutline";
import { openM2M } from "../util/Util";

import SessionStore from "../stores/SessionStore";
import WalletStore from "../stores/WalletStore";
import theme from "../theme";


const styles = {
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
    backgroundColor: theme.palette.secondary.main,
  },
  menuButton: {
    marginLeft: -12,
    marginRight: 10,
  },
  hidden: {
    display: "none",
  },
  flex: {
    flex: 1,
    paddingLeft: 40,
  },
  logo: {
    height: 32,
    marginLeft: -45,
  },
  search: {
    marginRight: 3 * theme.spacing.unit,
    color: theme.palette.textPrimary.main,
    backgroundColor: theme.palette.primary.secondary,
    width: 480,
    padding: 5,
    borderRadius: 3,
  },
  avatar: {
    background: theme.palette.secondary.main,
    color: theme.palette.common.white,
  },
  chip: {
    background: theme.palette.secondary.main,
    color: theme.palette.common.white,
    marginRight: theme.spacing(1),
    "&:hover": {
      background: theme.palette.primary.secondary,
    },
    "&:active": {
      background: theme.palette.primary.main,
    },
    "&:visited": {
      background: theme.palette.primary.main,
    },
  },
  iconStyle: {
    color: theme.palette.primary.main,
  },
  iconButton: {
    color: theme.palette.common.white,
    marginRight: theme.spacing(1),
  },
  noPadding: {
    "&:hover": {
      color: theme.palette.primary.main,
      cursor: 'pointer'
    }
  }
};

function getWalletBalance(organizationId) {
  if (!organizationId) {
    return 0;
  }
  
  return new Promise((resolve, reject) => {
    
    WalletStore.getWalletBalance(organizationId, resp => {
      return resolve(resp);
    });
  });
}

class TopNav extends Component {
  constructor() {
    super();

    this.state = {
      menuAnchor: null,
      balance: null,
      organizationId: SessionStore.getOrganizationID(),
      search: "",
    };

    this.onMenuOpen = this.onMenuOpen.bind(this);
    this.onMenuClose = this.onMenuClose.bind(this);
    this.onLogout = this.onLogout.bind(this);
    this.onSearchChange = this.onSearchChange.bind(this);
    this.onSearchSubmit = this.onSearchSubmit.bind(this);
  }
  
  componentDidMount() {
    this.loadData();
  }

  loadData = async () => {
    try {
      let organizationId = null;

      if(this.state.organizationId !== null){
        organizationId = this.state.organizationId;
      }else{
        if(SessionStore.getOrganizations().length > 0){
          organizationId = SessionStore.getOrganizations()[0].organizationID;
        }else{
          organizationId = null;
        }
      }

      var result = await getWalletBalance(organizationId);
      this.setState({ balance: result.balance });
    } catch (error) {
      console.error(error);
      this.setState({ error });
    }
  }

  onMenuOpen(e) {
    this.setState({
      menuAnchor: e.currentTarget,
    });
  }

  onMenuClose() {
    this.setState({
      menuAnchor: null,
    });
  }

  onLogout() {
    SessionStore.logout(() => {
      this.props.history.push("/login");
    });
  }

  onSearchChange(e) {
    this.setState({
      search: e.target.value,
    });
  }

  handlingExtLink = () => {
    const orgId = this.props.location.pathname.split('/')[2];
    openM2M(orgId, '/withdraw');
  }

  onSearchSubmit(e) {
    e.preventDefault();
    this.props.history.push(`/search?search=${encodeURIComponent(this.state.search)}`);
  }

  render() {
    //let drawerIcon;
    let logoIcon;
    let searchbar;
    if (!this.props.drawerOpen) {
      //drawerIcon = <Wallet />;
      logoIcon = <Typography type="body2" style={{ color: '#FFFFFF', fontFamily: 'Montserrat', fontSize: '22px' }} >M2M Wallet</Typography>
    } else {
      //drawerIcon = <MenuIcon />;
      logoIcon = <img src="/logo/logo.png" className={this.props.classes.logo} alt="LoRa Server" />
      searchbar = <Input
                    placeholder="Search organization, application, gateway or device"
                    className={this.props.classes.search}
                    disableUnderline={true}
                    value={this.state.search || ""}
                    onChange={this.onSearchChange}
                    startAdornment={
                      <InputAdornment position="start">
                      <Magnify />
                      </InputAdornment>
                    }
                  />
    }
    const { balance } = this.state;
    
    const balanceEl = balance === null ? 
      <span className="color-gray">(no org selected)</span> : 
      balance + " MXC";

    const open = Boolean(this.state.menuAnchor);
    const isDisabled = (this.props.user.username === process.env.REACT_APP_DEMO_USER)
                        ?true
                        :false;
    return(
      <AppBar className={this.props.classes.appBar}>
        <Toolbar>
          {/* <IconButton
            color="inherit"
            aria-label="toggle drawer"
            onClick={this.handleDrawerToggle}
            className={this.props.classes.menuButton}
          >
            {drawerIcon}
          </IconButton> */}

          <div className={this.props.classes.flex}>
            {logoIcon}
          </div>

          <form onSubmit={this.onSearchSubmit}>
            { searchbar }
          </form>

          <List>
            <ListItem>
              <ListItemIcon >
                <Wallet className={this.props.classes.iconStyle} />
              </ListItemIcon>
              <ListItemText primary={ balanceEl } classes={{ primary: this.props.classes.noPadding }} onClick={this.handlingExtLink}/>
            </ListItem>
          </List>

          <Chip
            avatar={
              <Avatar>
                <AccountCircle />
              </Avatar>
            }
            label={this.props.user.username}
            onClick={this.onMenuOpen}
            classes={{
              avatar: this.props.classes.avatar,
              root: this.props.classes.chip,
            }}
          />
          <a href="https://www.mxc.org/support" target="mxc-support">
            <IconButton className={this.props.classes.iconButton}>
              <HelpCircle />
            </IconButton>
          </a>

          <Menu
            id="menu-appbar"
            anchorEl={this.state.menuAnchor}
            anchorOrigin={{
              vertical: "top",
              horizontal: "right",
            }}
            transformOrigin={{
              vertical: "top",
              horizontal: "right",
            }}
            open={open}
            onClose={this.onMenuClose}
          >
            <MenuItem disabled={isDisabled} component={Link} to={`/users/${this.props.user.id}/password`}>Change password</MenuItem> :
            <MenuItem onClick={this.onLogout}>Logout</MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>
    );
  }
}

export default withStyles(styles)(withRouter(TopNav));
