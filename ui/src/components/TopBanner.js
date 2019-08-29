import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';

import AppBar from "@material-ui/core/AppBar";
import { withStyles } from "@material-ui/core/styles";
import Input from "@material-ui/core/Input";
import InputAdornment from "@material-ui/core/InputAdornment";
import Typography from '@material-ui/core/Typography';

//import MenuIcon from "mdi-material-ui/Menu";
//import Backburger from "mdi-material-ui/Backburger";
//import Wallet from "mdi-material-ui/Wallet";
import AccountCircle from "mdi-material-ui/AccountCircle";
import Magnify from "mdi-material-ui/Magnify";
import HelpCircle from "mdi-material-ui/HelpCircle";

import SessionStore from "../stores/SessionStore";
import theme from "../theme";


const styles = {
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
    backgroundColor: '#09006E50',
    marginTop: '96vh',
    opacity: 25,
  },
  menuButton: {
    marginLeft: -12,
    marginRight: 10,
  },
  hidden: {
    display: "none",
  },
  flex: {
    display: 'flex', 
    justifyContent: 'center',
    color: '#00FFD9',
    fontWeight: '900',
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
    marginRight: theme.spacing.unit,
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
  iconButton: {
    color: theme.palette.common.white,
    marginRight: theme.spacing.unit,
  },
};


class TopBanner extends Component {
  constructor() {
    super();

    this.state = {
      menuAnchor: null,
      search: "",
    };

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

    const open = Boolean(this.state.menuAnchor);
    const isDisabled = (this.props.user.username === process.env.REACT_APP_DEMO_USER)
                        ?true
                        :false;
    return(
      <AppBar className={this.props.classes.appBar}>
            <div className={this.props.classes.flex}>
                <span>THIS IS A DEMO VERSION.</span>
            </div>
      </AppBar>
    );
  }
}

export default withStyles(styles)(withRouter(TopBanner));
