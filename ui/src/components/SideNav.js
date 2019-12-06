import React, { Component } from "react";
import { Link, withRouter } from "react-router-dom";

import { withStyles } from "@material-ui/core/styles";
import Drawer from '@material-ui/core/Drawer';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';

import Divider from '@material-ui/core/Divider';
import Domain from "mdi-material-ui/Domain";
import Account from "mdi-material-ui/Account";
import Server from "mdi-material-ui/Server";
import Apps from "mdi-material-ui/Apps";
import RadioTower from "mdi-material-ui/RadioTower";
import Tune from "mdi-material-ui/Tune";
import Settings from "mdi-material-ui/Settings";
import Rss from "mdi-material-ui/Rss";
import Wallet from "mdi-material-ui/WalletOutline";
import AccessPoint from "mdi-material-ui/AccessPoint";

import AccountDetails from "mdi-material-ui/AccountDetails";
import ServerInfoStore from "../stores/ServerInfoStore"
import AutocompleteSelect from "./AutocompleteSelect";
import SessionStore from "../stores/SessionStore";
import OrganizationStore from "../stores/OrganizationStore";
import Admin from "./Admin";

import theme from "../theme";
import { openM2M } from "../util/Util";
import i18n, { packageNS } from '../i18n';

const styles = {
  drawerPaper: {
    position: "fixed",
    width: 270,
    paddingTop: theme.spacing.unit * 10,
    paddingRight: 0,
    paddingLeft: 0,
    backgroundColor: theme.palette.secondary.secondary,
    color: theme.palette.secondary.secondary,
    boxShadow: '1px 1px 5px 0px rgba(29, 30, 33, 0.5)',
  },
  select: {
    paddingTop: theme.spacing(1),
    paddingLeft: theme.spacing(3),
    paddingRight: theme.spacing(3),
    paddingBottom: theme.spacing(1),
  },
  selected: {
    //fontSize: 'larger', 
    color: theme.palette.primary.white,
  },
  static: {
    position: 'static'
  },
  iconStyle: {
    color: theme.palette.common.white,
  },
  divider: {
    padding: 10,
  },
  autocompleteSelect: {
    color: theme.palette.secondary.main,
  },
};

function loadServerVersion() {
  return new Promise((resolve, reject) => {
    ServerInfoStore.getAppserverVersion(data=>{
      resolve(data);
    });
  });
} 

class SideNav extends Component {
  constructor() {
    super();

    this.state = {
      open: true,
      organization: null,
      cacheCounter: 0,
      version: '1.0.0'
    };


    this.onChange = this.onChange.bind(this);
    this.getOrganizationOption = this.getOrganizationOption.bind(this);
    this.getOrganizationOptions = this.getOrganizationOptions.bind(this);
    this.getOrganizationFromLocation = this.getOrganizationFromLocation.bind(this);
  }

  loadData = async () => {
    try {
      const organizationID = SessionStore.getOrganizationID();
      var data = await loadServerVersion();
      const serverInfo = JSON.parse(data);
      
      this.setState({
        organizationID,
        version: serverInfo.version
      })

      this.setState({loading: true})
      
    } catch (error) {
      this.setState({loading: false})
      console.error(error);
      this.setState({ error });
    }
  }

  componentDidMount() {
    this.loadData();
    SessionStore.on("organization.change", () => {
      OrganizationStore.get(SessionStore.getOrganizationID(), resp => {
        this.setState({
          organization: resp.organization,
        });
      });
    });

    OrganizationStore.on("create", () => {
      this.setState({
        cacheCounter: this.state.cacheCounter + 1,
      });
    });

    OrganizationStore.on("change", (org) => {
      if (this.state.organization !== null && this.state.organization.id === org.id) {
        this.setState({
          organization: org,
        });
      }

      this.setState({
        cacheCounter: this.state.cacheCounter + 1,
      });
    });

    OrganizationStore.on("delete", id => {
      if (this.state.organization !== null && this.state.organization.id === id) {
        this.setState({
          organization: null
        });
      }

      this.setState({
        cacheCounter: this.state.cacheCounter + 1,
      });
    });

    /* if (SessionStore.getOrganizationID() !== null) {
      OrganizationStore.get(SessionStore.getOrganizationID(), resp => {
        this.setState({
          organization: resp.organization,
        });
      });
    } */

    this.getOrganizationFromLocation();
  }

  componentDidUpdate(prevProps) {
    if (this.props === prevProps) {
      return;
    }

    this.getOrganizationFromLocation();
  }

  onChange(e) {
    SessionStore.setOrganizationID(e.target.value);

    this.props.history.push(`/organizations/${e.target.value}/applications`);
  }

  getOrganizationFromLocation() {
    const organizationRe = /\/organizations\/(\d+)/g;
    const match = organizationRe.exec(this.props.history.location.pathname);

    if (match !== null && (this.state.organization === null || this.state.organization.id !== match[1])) {
      SessionStore.setOrganizationID(match[1]);
    }
  }

  getOrganizationOption(id, callbackFunc) {
    OrganizationStore.get(id, resp => {
      callbackFunc({label: resp.organization.name, value: resp.organization.id, color:"black"});
    });
  }

  getOrganizationOptions(search, callbackFunc) {
    OrganizationStore.list(search, 10, 0, resp => {
      const options = resp.result.map((o, i) => {return {label: o.name, value: o.id, color:'black'}});
      callbackFunc(options);
    });
  }

  handlingExtLink = () => {
    const resp = SessionStore.getProfile();
    resp.then((res) => {
      let orgId = SessionStore.getOrganizationID();
      const isBelongToOrg = res.body.organizations.some(e => e.organizationID === SessionStore.getOrganizationID());

      OrganizationStore.get(orgId, resp => {
        openM2M(resp.organization, isBelongToOrg, '/modify-account');
      });
    })
  }

  render() {
    let organizationID = "";
    if (this.state.organization) {
      organizationID = this.state.organization.id;
    }
    const { pathname } = this.props.location;
    const pathLastName = pathname.split('/').pop();
    
    const active = (sideNavName) => Boolean(pathLastName.match(sideNavName));
    const selected = (sideNavName) => {
      if(Boolean(pathLastName.match(sideNavName))){
        return { primary: this.props.classes.selected };
      }else{
        return {};
      }
    }

    return(
      <Drawer
        variant="persistent"
        anchor="left"
        open={this.props.open}
        classes={{paper: this.props.classes.drawerPaper}}
      >
        <Admin>
          <List>
            <ListItem selected={active('/network-servers')} button component={Link} to="/network-servers">
              <ListItemIcon>
                <Server />
              </ListItemIcon>
              <ListItemText classes={selected('/network-servers')} primary={i18n.t(`${packageNS}:tr000040`)} />
            </ListItem>
            <ListItem selected={active('/gateway-profiles')} button component={Link} to="/gateway-profiles">
              <ListItemIcon>
                <RadioTower />
              </ListItemIcon>
              <ListItemText classes={selected('/gateway-profiles')} primary={i18n.t(`${packageNS}:tr000046`)} />
            </ListItem>
            <ListItem selected={active('/organizations')} button component={Link} to="/organizations">
            <ListItemIcon>
                <Domain />
              </ListItemIcon>
              <ListItemText classes={selected('/organizations')} primary={i18n.t(`${packageNS}:tr000049`)} />
            </ListItem>
            <ListItem selected={active('/users')} button component={Link} to="/users">
              <ListItemIcon>
                <Account />
              </ListItemIcon>
              <ListItemText classes={selected('/users')} primary={i18n.t(`${packageNS}:tr000055`)} />
            </ListItem>
          </List>
        </Admin>
        <div>
          <AutocompleteSelect
            id="organizationID"
            margin="none"
            value={organizationID}
            onChange={this.onChange}
            getOption={this.getOrganizationOption}
            getOptions={this.getOrganizationOptions}
            className={this.props.classes.select}
            triggerReload={this.state.cacheCounter}
            placeHolder={i18n.t(`${packageNS}:tr000358`)}
          />
        </div>
        {this.state.organization && <>
        <List className={this.props.classes.static}>
           <Admin>
            <ListItem selected={active(`edit`)} button component={Link} to={`/organizations/${this.state.organization.id}/edit`}>
              <ListItemIcon>
                <Settings />
              </ListItemIcon>
              <ListItemText classes={selected(`edit`)} primary={i18n.t(`${packageNS}:tr000060`)} />
            </ListItem>
          </Admin>
          <Admin organizationID={this.state.organization.id}>
            <ListItem selected={active(`users`)} button component={Link} to={`/organizations/${this.state.organization.id}/users`}>
              <ListItemIcon>
                <Account />
              </ListItemIcon>
              <ListItemText classes={selected(`users`)} primary={i18n.t(`${packageNS}:tr000067`)} />
            </ListItem>
          </Admin>
          <ListItem selected={active(`service-profiles`)} button component={Link} to={`/organizations/${this.state.organization.id}/service-profiles`}>
            <ListItemIcon>
              <AccountDetails />
            </ListItemIcon>
            <ListItemText classes={selected(`service-profiles`)} primary={i18n.t(`${packageNS}:tr000069`)} />
          </ListItem>
          <ListItem selected={active(`device-profiles`)} button component={Link} to={`/organizations/${this.state.organization.id}/device-profiles`}>
            <ListItemIcon>
              <Tune />
            </ListItemIcon>
            <ListItemText classes={selected(`device-profiles`)} primary={i18n.t(`${packageNS}:tr000070`)} />
          </ListItem>
          {this.state.organization.canHaveGateways && <ListItem selected={active(`gateways`)} button component={Link} to={`/organizations/${this.state.organization.id}/gateways`}>
            <ListItemIcon>
              <RadioTower />
            </ListItemIcon>
            <ListItemText classes={selected(`gateways`)} primary={i18n.t(`${packageNS}:tr000072`)} />
          </ListItem>}
          <ListItem selected={active(`applications`)} button component={Link} to={`/organizations/${this.state.organization.id}/applications`}>
            <ListItemIcon>
              <Apps />
            </ListItemIcon>
            <ListItemText classes={selected(`applications`)} primary={i18n.t(`${packageNS}:tr000076`)} />
          </ListItem>
          <ListItem selected={active(`multicast-groups`)} button component={Link} to={`/organizations/${this.state.organization.id}/multicast-groups`}>
            <ListItemIcon>
              <Rss />
            </ListItemIcon>
            <ListItemText classes={selected(`multicast-groups`)} primary={i18n.t(`${packageNS}:tr000083`)} />
          </ListItem>
        </List>
        <Divider />
              <List className={this.props.classes.static}>
                <ListItem button onClick={this.handlingExtLink} >
                  <ListItemIcon>
                    <Wallet />
                  </ListItemIcon>
                  <ListItemText primary={i18n.t(`${packageNS}:tr000084`)} />
                </ListItem>
                <ListItem button className={this.props.classes.static}>  
                  <ListItemIcon>
                    <AccessPoint />
                  </ListItemIcon>
                  <ListItemText primary={i18n.t(`${packageNS}:tr000085`)} />
                </ListItem>
                <ListItem>
                  <ListItemText primary={i18n.t(`${packageNS}:tr000086`)} />
                  <ListItemIcon>
                    <img src="/logo/mxc_logo.png" className="iconStyle" alt={i18n.t(`${packageNS}:tr000051`)} onClick={this.handleMXC} />
                  </ListItemIcon>
                </ListItem>
                <ListItem>
                  <ListItemText secondary={`${i18n.t(`${packageNS}:tr000087`)} ${this.state.version}`} />
                </ListItem>
              </List>
        </>}
      </Drawer>
    );
  }
}

export default withRouter(withStyles(styles)(SideNav));
