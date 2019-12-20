import React, { Component } from "react";
import { Route, Switch, Link, withRouter } from "react-router-dom";

import { withStyles } from "@material-ui/core/styles";
import Grid from '@material-ui/core/Grid';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';

import i18n, { packageNS } from '../../../i18n';
import TitleBar from "../../../components/TitleBar";
import TitleBarTitle from "../../../components/TitleBarTitle";
import Spinner from "../../../components/ScaleLoader"
import SuperNodeEthAccount from "./EthAccount";

import topupStore from "../../../stores/TopupStore";
import styles from "./HistoryStyle";
import { SUPER_ADMIN } from "../../../util/M2mUtil";



class SupernodeHistory extends Component {
  constructor(props) {
    super(props);
    this.state = {
      tab: 0,
      loading: false,
      admin: false,
      income:0
    };

    this.onChangeTab = this.onChangeTab.bind(this);
    this.locationToTab = this.locationToTab.bind(this);
  }

  componentDidMount() {
    this.setState({loading:true});
    this.locationToTab();
    this.setState({loading:false});
    this.getIncome();
  }

  componentDidUpdate(oldProps) {
    if (this.props == oldProps) {
      return;
    }

    this.locationToTab();
  }

  getIncome(){
    topupStore.getIncome(0, resp => {
      this.setState({income:resp.amount});
    });
  }

  onChangeTab(e, v) {
    this.setState({
      tab: v,
    });
  }

  locationToTab() {
    let tab = 0;
    if (window.location.href.endsWith("/eth-account")) {
      tab = 1;
    } else if (window.location.href.endsWith("/network-activity")) {
      tab = 2;
    }  
    
    this.setState({
      tab,
    });
  }

  render() {
      
    return(
      <Grid container alignContent={'center'} spacing={24}>
        <Spinner on={this.state.loading}/>
        <Grid item xs={12} md={12} lg={12} className={this.props.classes.divider}>
          <div className={this.props.classes.TitleBar}>
                <TitleBar className={this.props.classes.padding}>
                  <TitleBarTitle title={i18n.t(`${packageNS}:menu.history.history`)} />
                </TitleBar>
                {/* <Divider light={true}/>
                <div className={this.props.classes.breadcrumb}>
                <TitleBar>
                  <TitleBarTitle component={Link} to="#" title="M2M Wallet" className={this.props.classes.link}/> 
                  <TitleBarTitle title="/" className={this.props.classes.navText}/>
                  <TitleBarTitle component={Link} to="#" title="History" className={this.props.classes.link}/>
                </TitleBar>
                </div> */}
            </div>
        </Grid>

        <Grid item container alignContent={'center'} xs={12} md={12} lg={12} justify="space-between" className={this.props.classes.tabsBlock}>
        <Tabs
            value={this.state.tab}
            onChange={this.onChangeTab}
            indicatorColor="primary"
            className={this.props.classes.tabs}
            variant="scrollable"
            scrollButtons="auto"
            textColor="primary"
          >
            <Tab label={i18n.t(`${packageNS}:menu.history.eth_account`)} component={Link} to={`/control-panel/history/`} />
          </Tabs>

            <Grid container justify="space-between" alignItems="center" className={this.props.classes.card}>
               <Grid item>{i18n.t(`${packageNS}:menu.history.last_income`)}</Grid>
              <Grid item align="right"><b>{this.state.income}MXC</b></Grid>
            </Grid>
        
        </Grid>

        <Grid item alignItems={'center'} xs={12} md={12} lg={12} >
          <Switch>
            <Route exact path={`${this.props.match.path}/`} render={props => <SuperNodeEthAccount organizationID={SUPER_ADMIN} {...props} />} />
            {/* <Redirect to={`/history/${organizationID}/transactions`} /> */}
          </Switch>
        </Grid>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(SupernodeHistory));
/* const _ExportedHistory = withStyles(styles)(withRouter(SupernodeHistory));
export default _ExportedHistory; */