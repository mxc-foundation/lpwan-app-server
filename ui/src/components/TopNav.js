import { IconButton } from "@material-ui/core";
import AppBar from "@material-ui/core/AppBar";
import Avatar from '@material-ui/core/Avatar';
import Chip from '@material-ui/core/Chip';
import Input from "@material-ui/core/Input";
import InputAdornment from "@material-ui/core/InputAdornment";
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import { withStyles } from "@material-ui/core/styles";
import Toolbar from "@material-ui/core/Toolbar";
import Typography from '@material-ui/core/Typography';
import AccountCircle from "mdi-material-ui/AccountCircle";
import HelpCircle from "mdi-material-ui/HelpCircle";
import Magnify from "mdi-material-ui/Magnify";
import Wallet from "mdi-material-ui/WalletOutline";
import React, { Component } from "react";
import { Link, withRouter } from 'react-router-dom';
import i18n, { packageNS } from '../i18n';
import OrganizationStore from "../stores/OrganizationStore";
import SessionStore from "../stores/SessionStore";
import WalletStore from "../stores/WalletStore";
import theme from "../theme";
import { openM2M } from "../util/Util";
import DropdownMenuLanguage from "./DropdownMenuLanguage";





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
  brandLogo: {
    height: 32,
    marginLeft: 20,
  },
  search: {
    marginRight: 3 * theme.spacing.unit,
    color: theme.palette.textPrimary.main,
    backgroundColor: '#08005C',
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
      brandingLogo: null,
    };

    this.onLogout = this.onLogout.bind(this);
    this.onMenuOpen = this.onMenuOpen.bind(this);
    this.onMenuClose = this.onMenuClose.bind(this);
    this.onSearchChange = this.onSearchChange.bind(this);
    this.onSearchSubmit = this.onSearchSubmit.bind(this);
  }
  
  componentDidMount() {
    this.loadData();
    SessionStore.getBranding(resp => {
      if (resp.logoPath !== "") {
        this.setState({
          brandingLogo: resp.logoPath,
        });
      }
    });

    SessionStore.on("organization.change", () => {
      this.loadData();
    });
  }

  loadData = async () => {
    try {
      let organizationId = SessionStore.getOrganizationID();

      //var result = await getWalletBalance(organizationId);
      
      //this.setState({ balance: result.balance });
    } catch (error) {
      console.error(error);
      this.setState({ error });
    }
  }

  onChangeLanguage = (newLanguageState) => {
    this.props.onChangeLanguage(newLanguageState);
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
      
    });
  }

  onSearchChange(e) {
    this.setState({
      search: e.target.value,
    });
  }

  handlingExtLink = () => {
    const resp = SessionStore.getProfile();
    resp.then((res) => {
      let orgId = this.props.location.pathname.split('/')[2];
      const isBelongToOrg = res.body.organizations.some(e => e.organizationID === SessionStore.getOrganizationID());

      OrganizationStore.get(orgId, resp => {
        openM2M(resp.organization, isBelongToOrg, '/withdraw');
      });
    })
  }

  onSearchSubmit(e) {
    e.preventDefault();
    this.props.history.push(`/search?search=${encodeURIComponent(this.state.search)}`);
  }

  render() {
    //let drawerIcon;
    let logoIcon;
    let searchbar;
    let brandingLogo;
    if (!this.props.drawerOpen) {
      //drawerIcon = <Wallet />;
      logoIcon = <Typography type="body2" style={{ color: '#FFFFFF', fontFamily: 'Montserrat', fontSize: '22px' }} >M2M Wallet</Typography>
    } else {
      //drawerIcon = <MenuIcon />;
      logoIcon = <img src="/logo/logo_LP.png" className={this.props.classes.logo} alt="LPWAN Server" />
      searchbar = <Input
                    placeholder={i18n.t(`${packageNS}:tr000033`)}
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

    if (this.state.brandingLogo != null) {
      brandingLogo = <img src={this.state.brandingLogo} className={this.props.classes.brandLogo}/>;
    }else {
      brandingLogo = <div></div>
    }
    console.log("#######", brandingLogo);

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
            {brandingLogo}
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

          <DropdownMenuLanguage onChangeLanguage={this.onChangeLanguage} />

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
            <MenuItem disabled={isDisabled} component={Link} to={`/users/${this.props.user.id}/password`}>{i18n.t(`${packageNS}:tr000038`)}</MenuItem>
            <MenuItem onClick={this.onLogout}>{i18n.t(`${packageNS}:tr000035`)}</MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>
    );
  }
}

export default withStyles(styles)(withRouter(TopNav));
