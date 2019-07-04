import React, { Component } from "react";
import { Link, withRouter } from "react-router-dom";

import { withStyles } from "@material-ui/core/styles";
import Drawer from '@material-ui/core/Drawer';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import Typography from '@material-ui/core/Typography';

import Card from '@material-ui/core/Card';
import CardContent from "@material-ui/core/CardContent";

import Divider from '@material-ui/core/Divider';
import Domain from "mdi-material-ui/Domain";
import Account from "mdi-material-ui/Account";
import Server from "mdi-material-ui/Server";
import Apps from "mdi-material-ui/Apps";
import RadioTower from "mdi-material-ui/RadioTower";
import Tune from "mdi-material-ui/Tune";
import Settings from "mdi-material-ui/Settings";
import Rss from "mdi-material-ui/Rss";
import Wallet from "mdi-material-ui/Wallet";

import AccessPoint from "mdi-material-ui/AccessPoint";
import Repeat from "mdi-material-ui/Repeat";
import CalendarCheckOutline from "mdi-material-ui/CalendarCheckOutline";
import CreditCard from "mdi-material-ui/CreditCard";
import ArrowExpandLeft from "mdi-material-ui/ArrowExpandLeft";


//import ModifyEthAccount from "mdi-material-ui/Card-bulleted-settings-outline"
//import History from "mdi-material-ui/History"
//import Topup from "mdi-material-ui/Bank-transfer-in"
//import Withdraw from "mdi-material-ui/Cash-multiple"
import AccountDetails from "mdi-material-ui/AccountDetails";

import AutocompleteSelect from "./AutocompleteSelect";
import SessionStore from "../stores/SessionStore";
import OrganizationStore from "../stores/OrganizationStore";
import Admin from "./Admin";

import theme from "../theme";


const styles = {
  drawerPaper: {
    position: "fixed",
    width: 270,
    paddingTop: theme.spacing.unit * 9,
    backgroundColor: '#09006E',
    color: '#FFFFFF',
  },
  select: {
    paddingTop: theme.spacing.unit,
    paddingLeft: theme.spacing.unit * 3,
    paddingRight: theme.spacing.unit * 3,
    paddingBottom: theme.spacing.unit * 1,
  },
  card: {
    width: '100%',
    height: 200,
    position: 'absolute',
    bottom: 0,
    backgroundColor: '#09006E',
    color: '#FFFFFF',
  },
  static: {
    position: 'static'
  },
  iconStyle: {
    color: theme.palette.common.white,
  }
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
        console.log('org organization.change', resp.organization);
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
        console.log('org change', org);
        this.setState({
          organization: org,
        });
      }
    });

    OrganizationStore.on("delete", id => {
      if (this.state.organization !== null && this.state.organization.id === id) {
        console.log('org delete');
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
        console.log('org componentDidMount', resp.organization);
        this.setState({
          organization: resp.organization,
        });
      });
    }

    this.getOrganizationFromLocation();
  }

  componentWillUnmount() {
    console.log('SideNav.componentWillUnmount');
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
      callbackFunc({label: resp.organization.name, value: resp.organization.id});
    });
  }

  getOrganizationOptions(search, callbackFunc) {
    OrganizationStore.list(search, 10, 0, resp => {
      const options = resp.result.map((o, i) => {return {label: o.name, value: o.id}});
      callbackFunc(options);
    });
  }

  handleOpenM2M = () => {
    //this.props.setDrawerOpen(false);
    //this.props.history.push(`/withdraw/${this.state.organization.id}`);
    
    const data = {
      jwt: window.localStorage.getItem("jwt"),
      path: `/withdraw/${this.state.organization.id}`,
      org_id: `${this.state.organization.id}`
    };
    const dataString = encodeURIComponent(JSON.stringify(data));
    
    // for new tab, see: https://stackoverflow.com/questions/427479/programmatically-open-new-pages-on-tabs
    window.location.replace(`http://localhost:3001/#/j/${dataString}`);
  }

  handleOpenLora = () => {
    this.props.setDrawerOpen(true);
    this.props.history.push(`/`);
  }  

  render() {
    let organizationID = "";
    if (this.state.organization) {
      organizationID = this.state.organization.id;
    }
   
    return(
      <>
      <Drawer
        variant="persistent"
        anchor="left"
        open={this.props.open}
        classes={{paper: this.props.classes.drawerPaper}}
      >
        <Admin>
          <List>
            <ListItem button component={Link} to="/network-servers">
              <ListItemIcon>
                <Server />
              </ListItemIcon>
              <ListItemText disableTypography
        primary={<Typography type="body2" className="default-text">Network-servers</Typography>} />
            </ListItem>
            <ListItem button component={Link} to="/gateway-profiles">
              <ListItemIcon>
                <RadioTower />
              </ListItemIcon>
              <ListItemText disableTypography
        primary={<Typography type="body2" className="default-text">Gateway-profiles</Typography>} />
            </ListItem>
            <ListItem button component={Link} to="/organizations">
            <ListItemIcon>
                <Domain />
              </ListItemIcon>
              <ListItemText disableTypography
        primary={<Typography type="body2" className="default-text">Organizations</Typography>} />
            </ListItem>
            <ListItem button component={Link} to="/users">
              <ListItemIcon>
                <Account />
              </ListItemIcon>
              <ListItemText disableTypography
        primary={<Typography type="body2" className="default-text">All users</Typography>} />
            </ListItem>
          </List>
          <Divider />
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
          />
        </div>

        {this.state.organization && <>
        <List className={this.props.classes.static}>
          <Admin>
            <ListItem button component={Link} to={`/organizations/${this.state.organization.id}/edit`}>
              <ListItemIcon>
                <Settings />
              </ListItemIcon>
              <ListItemText disableTypography
        primary={<Typography type="body2" className="default-text" >Org. settings</Typography>} />
            </ListItem>
          </Admin>
          <Admin organizationID={this.state.organization.id}>
            <ListItem button component={Link} to={`/organizations/${this.state.organization.id}/users`}>
              <ListItemIcon>
                <Account />
              </ListItemIcon>
              <ListItemText disableTypography
        primary={<Typography type="body2" className="default-text" >Org. users</Typography>} />
            </ListItem>
          </Admin>
          <ListItem button component={Link} to={`/organizations/${this.state.organization.id}/service-profiles`}>
            <ListItemIcon>
              <AccountDetails />
            </ListItemIcon>
            <ListItemText disableTypography
        primary={<Typography type="body2" className="default-text" >Service-profiles</Typography>} />
          </ListItem>
          <ListItem button component={Link} to={`/organizations/${this.state.organization.id}/device-profiles`}>
            <ListItemIcon>
              <Tune />
            </ListItemIcon>
            <ListItemText disableTypography
        primary={<Typography type="body2" className="default-text" >Device-profiles</Typography>} />
          </ListItem>
          {this.state.organization.canHaveGateways && <ListItem button component={Link} to={`/organizations/${this.state.organization.id}/gateways`}>
            <ListItemIcon>
              <RadioTower />
            </ListItemIcon>
            <ListItemText disableTypography
        primary={<Typography type="body2" className="default-text" >Gateways</Typography>} />
          </ListItem>}
          <ListItem button component={Link} to={`/organizations/${this.state.organization.id}/applications`}>
            <ListItemIcon>
              <Apps />
            </ListItemIcon>
            <ListItemText disableTypography
        primary={<Typography type="body2" className="default-text" >Applications</Typography>} />
          </ListItem>
          <ListItem button component={Link} to={`/organizations/${this.state.organization.id}/multicast-groups`}>
            <ListItemIcon>
              <Rss />
            </ListItemIcon>
            <ListItemText disableTypography
        primary={<Typography type="body2" className="default-text" >Multicast-groups</Typography>} />
          </ListItem>
        </List>

        <Card className={this.props.classes.card}>
            <CardContent>
              <List className={this.props.classes.static}>
                <ListItem button  onClick={this.handleOpenLora}>
                  <ListItemIcon>
                    <AccessPoint />
                  </ListItemIcon>
                  <ListItemText disableTypography primary={<Typography type="body2" className="default-text" >Lora</Typography>} />
                </ListItem>
                <ListItem button onClick={this.handleOpenM2M} >
                  <ListItemIcon>
                    <Wallet />
                  </ListItemIcon>
                  <ListItemText disableTypography primary={<Typography type="body2" className="default-text" >M2M Wallet</Typography>} />
                </ListItem>
                <ListItem button  onClick={this.handleOpenLora}>
                  <ListItemText disableTypography primary={<Typography type="body2" className="default-text" >Account name</Typography>} />
                  <ListItemIcon>
                    <Settings />
                  </ListItemIcon>
                </ListItem>
                <ListItem button onClick={this.handleOpenM2M} >
                  <ListItemText disableTypography primary={<Typography type="body2" className="default-text" >Change Account</Typography>} />
                  <ListItemIcon>
                    <Repeat />
                  </ListItemIcon>
                </ListItem>
              </List>
            </CardContent>
          </Card>
        </>}
      </Drawer>
      <Drawer 
        variant="persistent"
        anchor="left"
        open={!this.props.open}
        classes={{paper: this.props.classes.drawerPaper}}
      >
        {this.state.organization && <List className={this.props.classes.static}>
        
          <ListItem button component={Link} to={`/withdraw/${this.state.organization.id}`}>
            <ListItemIcon className={this.props.classes.iconStyle}>
              <ArrowExpandLeft />
            </ListItemIcon>
            <ListItemText disableTypography
        primary={<Typography type="body2" style={{ color: '#FFFFFF', fontFamily: 'Montserrat' }} >Withdraw</Typography>} />
          </ListItem>
          <ListItem button component={Link} to={`/history`}>
            <ListItemIcon>
              <CalendarCheckOutline />
            </ListItemIcon>
            <ListItemText disableTypography
        primary={<Typography type="body2" className="default-text" >History</Typography>} />
          </ListItem>
          <ListItem button component={Link} to={`/modify-account`}>
            <ListItemIcon>
              <CreditCard />
            </ListItemIcon>
            <ListItemText disableTypography
        primary={<Typography type="body2" className="default-text" >ModifyEthAccount</Typography>} />
          </ListItem>
          <Card className={this.props.classes.card}>
            <CardContent>
              <List className={this.props.classes.static}>
                <ListItem button  onClick={this.handleOpenLora}>
                  <ListItemIcon>
                    <AccessPoint />
                  </ListItemIcon>
                  <ListItemText disableTypography primary={<Typography type="body2" className="default-text" >Lora</Typography>} />
                </ListItem>
                <ListItem button onClick={this.handleOpenM2M} >
                  <ListItemIcon>
                    <Wallet />
                  </ListItemIcon>
                  <ListItemText disableTypography primary={<Typography type="body2" className="default-text" >M2M Wallet</Typography>} />
                </ListItem>
                <ListItem button  onClick={this.handleOpenLora}>
                  <ListItemText disableTypography primary={<Typography type="body2" className="default-text" >Account name</Typography>} />
                  <ListItemIcon>
                    <Settings />
                  </ListItemIcon>
                </ListItem>
                <ListItem button onClick={this.handleOpenM2M} >
                  <ListItemText disableTypography primary={<Typography type="body2" className="default-text" >Change Account</Typography>} />
                  <ListItemIcon>
                    <Repeat />
                  </ListItemIcon>
                </ListItem>
              </List>
            </CardContent>
          </Card>
        </List>}
      </Drawer>
      </>
    );
  }
}

export default withRouter(withStyles(styles)(SideNav));
