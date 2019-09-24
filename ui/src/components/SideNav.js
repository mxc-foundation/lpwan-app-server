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

import AccountDetails from "mdi-material-ui/AccountDetails";

import AutocompleteSelect from "./AutocompleteSelect";
import SessionStore from "../stores/SessionStore";
import OrganizationStore from "../stores/OrganizationStore";
import Admin from "./Admin";

import theme from "../theme";
import { getM2MLink } from "../util/Util";

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
/*   card: {
    width: '100%',
    height: 200,
    position: 'absolute',
    bottom: 0,
    backgroundColor: '#09006E',
    color: '#FFFFFF',
    marginTop: -20,
  }, */
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

class SideNav extends Component {
  constructor() {
    super();

    this.state = {
      open: true,
      organization: null,
      cacheCounter: 0,
    };


    this.onChange = this.onChange.bind(this);
    this.getOrganizationOption = this.getOrganizationOption.bind(this);
    this.getOrganizationOptions = this.getOrganizationOptions.bind(this);
    this.getOrganizationFromLocation = this.getOrganizationFromLocation.bind(this);
  }

  componentDidMount() {
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

    if (SessionStore.getOrganizationID() !== null) {
      OrganizationStore.get(SessionStore.getOrganizationID(), resp => {
        this.setState({
          organization: resp.organization,
        });
      });
    }

    this.getOrganizationFromLocation();
  }

  componentDidUpdate(prevProps) {
    if (this.props === prevProps) {
      return;
    }

    this.getOrganizationFromLocation();
  }

  onChange(e) {
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

  handleOpenM2M = () => {
    let orgId = this.state.organization.id;
    let orgName = '';
    if(!orgId){
      return false;
    }
    const user = SessionStore.getUser();  
    const org = SessionStore.getOrganizations(); 
    
    if(user.isAdmin){
      orgId = '0';
      orgName = 'Super_admin';
    }else{
      if(org.length > 0){
        orgName = org[0].organizationName;
      }else{
        orgName = '';
      }
    }
    
    const data = {
      jwt: window.localStorage.getItem("jwt"),
      path: `/withdraw/${orgId}`,
      orgId,
      orgName,
      username: user.username,
      loraHostUrl: window.location.origin
    };
    
    const dataString = encodeURIComponent(JSON.stringify(data));

    const host = getM2MLink();

    // for new tab, see: https://stackoverflow.com/questions/427479/programmatically-open-new-pages-on-tabs
    window.location.replace(host + `/#/j/${dataString}`);
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
              <ListItemText classes={selected('/network-servers')} primary="Network-servers" />
            </ListItem>
            <ListItem selected={active('/gateway-profiles')} button component={Link} to="/gateway-profiles">
              <ListItemIcon>
                <RadioTower />
              </ListItemIcon>
              <ListItemText classes={selected('/gateway-profiles')} primary="Gateway-profiles" />
            </ListItem>
            <Divider />
            <ListItem selected={active('/organizations')} button component={Link} to="/organizations">
            <ListItemIcon>
                <Domain />
              </ListItemIcon>
              <ListItemText classes={selected('/organizations')} primary="Organizations" />
            </ListItem>
            <ListItem selected={active('/users')} button component={Link} to="/users">
              <ListItemIcon>
                <Account />
              </ListItemIcon>
              <ListItemText classes={selected('/users')} primary="All users" />
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
            placeHolder="Change Organization"
          />
        </div>
        <Divider />
        {this.state.organization && <>
        <List className={this.props.classes.static}>
           <Admin>
            <ListItem selected={active(`edit`)} button component={Link} to={`/organizations/${this.state.organization.id}/edit`}>
              <ListItemIcon>
                <Settings />
              </ListItemIcon>
              <ListItemText classes={selected(`edit`)} primary="Org. settings" />
            </ListItem>
          </Admin>
          <Admin organizationID={this.state.organization.id}>
            <ListItem selected={active(`users`)} button component={Link} to={`/organizations/${this.state.organization.id}/users`}>
              <ListItemIcon>
                <Account />
              </ListItemIcon>
              <ListItemText classes={selected(`users`)} primary="Org. users" />
            </ListItem>
          </Admin>
          <ListItem selected={active(`service-profiles`)} button component={Link} to={`/organizations/${this.state.organization.id}/service-profiles`}>
            <ListItemIcon>
              <AccountDetails />
            </ListItemIcon>
            <ListItemText classes={selected(`service-profiles`)} primary="Service-profiles" />
          </ListItem>
          <ListItem selected={active(`device-profiles`)} button component={Link} to={`/organizations/${this.state.organization.id}/device-profiles`}>
            <ListItemIcon>
              <Tune />
            </ListItemIcon>
            <ListItemText classes={selected(`device-profiles`)} primary="Device-profiles" />
          </ListItem>
          {this.state.organization.canHaveGateways && <ListItem selected={active(`gateways`)} button component={Link} to={`/organizations/${this.state.organization.id}/gateways`}>
            <ListItemIcon>
              <RadioTower />
            </ListItemIcon>
            <ListItemText classes={selected(`gateways`)} primary="Gateways" />
          </ListItem>}
          <ListItem selected={active(`applications`)} button component={Link} to={`/organizations/${this.state.organization.id}/applications`}>
            <ListItemIcon>
              <Apps />
            </ListItemIcon>
            <ListItemText classes={selected(`applications`)} primary="Applications" />
          </ListItem>
          <ListItem selected={active(`multicast-groups`)} button component={Link} to={`/organizations/${this.state.organization.id}/multicast-groups`}>
            <ListItemIcon>
              <Rss />
            </ListItemIcon>
            <ListItemText classes={selected(`multicast-groups`)} primary="Multicast-groups" />
          </ListItem>
        </List>
        <Divider />
              <List className={this.props.classes.static}>
                <ListItem button onClick={this.handleOpenM2M} >
                  <ListItemIcon>
                    <Wallet />
                  </ListItemIcon>
                  <ListItemText primary="M2M Wallet" />
                </ListItem>
                <ListItem>
                  <ListItemText primary="Powered by" />
                  <ListItemIcon>
                    <img src="/logo/mxc_logo.png" className="iconStyle" alt="LoRa Server" onClick={this.handleMXC} />
                  </ListItemIcon>
                </ListItem>
              </List>

        </>}
      </Drawer>
    );
  }
}

export default withRouter(withStyles(styles)(SideNav));
