import React, { Component } from "react";
import { Route, Switch, Link, withRouter } from "react-router-dom";

import { withStyles } from "@material-ui/core/styles";
import Grid from '@material-ui/core/Grid';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
//import Card from '@material-ui/core/Card';
//import CardContent from "@material-ui/core/CardContent";

import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import Divider from '@material-ui/core/Divider';
//import TitleBarButton from "../../components/TitleBarButton";

import SessionStore from "../../stores/SessionStore";

import Transactions from "./Transactions";
import EthAccount from "./EthAccount";
import Subsriptions from "./Subsriptions";

import theme from "../../theme";


const styles = {
  tabs: {
    borderBottom: "1px solid " + theme.palette.divider,
    height: "49px",
  },
  navText: {
    fontSize: 14,
  },
  TitleBar: {
    height: 115,
    width: '50%',
    light: true,
    display: 'flex',
    flexDirection: 'column'
  },
  card: {
    minWidth: 180,
    width: 220,
    backgroundColor: "#0C0270",
  },
  divider: {
    padding: 0,
    color: '#FFFFFF',
    width: '100%',
  },
  padding: {
    padding: 0,
  },
};


class HistoryLayout extends Component {
  constructor() {
    super();
    this.state = {
      tab: 0,
      admin: false,
    };

    this.onChangeTab = this.onChangeTab.bind(this);
    this.locationToTab = this.locationToTab.bind(this);
  }

  componentDidMount() {
    this.locationToTab();
  }

  componentDidUpdate(oldProps) {
    console.log(oldProps)
    if (this.props === oldProps) {
      return;
    }

    this.locationToTab();
  }

  onChangeTab(e, v) {
    this.setState({
      tab: v,
    });
  }

  locationToTab() {
    let tab = 0;
    console.log("window.location.href");
console.log(window.location.href);
    if (window.location.href.endsWith("/edit")) {
      tab = 1;
    } else if (window.location.href.endsWith("/keys")) {
      tab = 2;
    } 

    if (tab > 1 && !this.state.admin) {
      tab = tab - 1;
    }

    this.setState({
      tab: tab,
    });
  }

  render() {
    const organizationID = SessionStore.getOrganizationID();
  
    console.log('props', this.props)
    return(
      <Grid container spacing={24}>
        <Grid item xs={12} className={this.props.classes.divider}>
          <div className={this.props.classes.TitleBar}>
              <TitleBar className={this.props.classes.padding}>
                <TitleBarTitle title="History" />
              </TitleBar>
              <Divider light={true}/>
              <TitleBar>
                <TitleBarTitle title="M2M Wallet" className={this.props.classes.navText}/>
                <TitleBarTitle title="/" className={this.props.classes.navText}/>
                <TitleBarTitle title="History" className={this.props.classes.navText}/>
              </TitleBar>
          </div>
        </Grid>

        <Grid item xs={12}>
          <Tabs
            value={this.state.tab}
            onChange={this.onChangeTab}
            indicatorColor="primary"
            className={this.props.classes.tabs}
            variant="scrollable"
            scrollButtons="auto"
            textColor="primary"
          >
            <Tab label="Transactions" component={Link} to={`/history/${organizationID}`} />
            <Tab label="ETH Account" component={Link} to={`/history/${organizationID}/edit`} />
            <Tab label="Subsriptions" component={Link} to={`/history/${organizationID}/keys`} />
          </Tabs>
        </Grid>

        <Grid item xs={12}>
          <Switch>
            <Route exact path={`${this.props.match.path}`} render={props => <Transactions {...props} />} />
            <Route exact path={`${this.props.match.path}/${organizationID}/edit`} render={props => <EthAccount {...props} />} />
            <Route exact path={`${this.props.match.path}/${organizationID}/keys`} render={props => <Subsriptions {...props} />} />
          </Switch>
        </Grid>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(HistoryLayout));
