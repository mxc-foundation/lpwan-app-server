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

import MenuIcon from "mdi-material-ui/Menu";
import Backburger from "mdi-material-ui/Backburger";
import Wallet from "mdi-material-ui/Wallet";
import AccountCircle from "mdi-material-ui/AccountCircle";
import Magnify from "mdi-material-ui/Magnify";
import HelpCicle from "mdi-material-ui/HelpCircle";

import SessionStore from "../stores/SessionStore";
import theme from "../theme";


const styles = {
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
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
  },
  logo: {
    height: 32,
  },
  search: {
    marginRight: 3 * theme.spacing.unit,
    color: theme.palette.common.white,
    background: blue[400],
    width: 450,
    padding: 5,
    borderRadius: 3,
  },
  avatar: {
    background: blue[600],
    color: theme.palette.common.white,
  },
  chip: {
    background: blue[600],
    color: theme.palette.common.white,
    marginRight: theme.spacing.unit,
    "&:hover": {
      background: blue[400],
    },
    "&:active": {
      background: blue[400],
    },
  },
  iconButton: {
    color: theme.palette.common.white,
    marginRight: theme.spacing.unit,
  },
};


class TopNav extends Component {
  constructor() {
    super();

    this.state = {
      menuAnchor: null,
      search: "",
    };

    this.handleDrawerToggle = this.handleDrawerToggle.bind(this);
    this.onMenuOpen = this.onMenuOpen.bind(this);
    this.onMenuClose = this.onMenuClose.bind(this);
    this.onLogout = this.onLogout.bind(this);
    this.onSearchChange = this.onSearchChange.bind(this);
    this.onSearchSubmit = this.onSearchSubmit.bind(this);
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

  handleDrawerToggle() {
    this.props.setDrawerOpen(!this.props.drawerOpen);
    if(!this.props.drawerOpen){
      this.props.history.push("/");
    }else{
      this.props.history.push("/wallet");
    }
    
  }

  onSearchChange(e) {
    this.setState({
      search: e.target.value,
    });
  }

  onSearchSubmit(e) {
    e.preventDefault();
    this.props.history.push(`/search?search=${encodeURIComponent(this.state.search)}`);
  }

  render() {
    let drawerIcon;
    let logoIcon;
    if (!this.props.drawerOpen) {
      drawerIcon = <Wallet />;
      logoIcon = <img src="data:image/png;base64,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" className={this.props.classes.logo} alt="m2m-wallet" />;
    } else {
      drawerIcon = <MenuIcon />;
      logoIcon = <img src="/logo/logo.png" className={this.props.classes.logo} alt="LoRa Server" />
    }

    const open = Boolean(this.state.menuAnchor);

    return(
      <AppBar className={this.props.classes.appBar}>
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="toggle drawer"
            onClick={this.handleDrawerToggle}
            className={this.props.classes.menuButton}
          >
            {drawerIcon}
          </IconButton>

          <div className={this.props.classes.flex}>
            {logoIcon}
          </div>

          <form onSubmit={this.onSearchSubmit}>
            <Input
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
          </form>

          <a href="https://www.loraserver.io/lora-app-server/" target="loraserver-doc">
            <IconButton className={this.props.classes.iconButton}>
              <HelpCicle />
            </IconButton>
          </a>

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
            <MenuItem component={Link} to={`/users/${this.props.user.id}/password`}>Change password</MenuItem>
            <MenuItem onClick={this.onLogout}>Logout</MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>
    );
  }
}

export default withStyles(styles)(withRouter(TopNav));
