import Grid from '@material-ui/core/Grid';
import { withStyles } from "@material-ui/core/styles";
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Delete from "mdi-material-ui/Delete";
import React, { Component } from "react";
import { Link, Route, Switch, withRouter } from "react-router-dom";
import DeviceAdmin from "../../components/DeviceAdmin";
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";
import TitleBarTitle from "../../components/TitleBarTitle";
import i18n, { packageNS } from '../../i18n';
import MulticastGroupStore from "../../stores/MulticastGroupStore";
import SessionStore from "../../stores/SessionStore";
import theme from "../../theme";
import ListMulticastGroupDevices from "./ListMulticastGroupDevices";
import UpdateMulticastGroup from "./UpdateMulticastGroup";






const styles = {
  tabs: {
    borderBottom: "1px solid " + theme.palette.divider,
    height: "48px",
    overflow: "visible",
  },
};


class MulticastGroupLayout extends Component {
  constructor() {
    super();
    this.state = {
      tab: 0,
      admin: false,
    };

    this.locationToTab = this.locationToTab.bind(this);
    this.onChangeTab = this.onChangeTab.bind(this);
    this.deleteMulticastGroup = this.deleteMulticastGroup.bind(this);
    this.setIsAdmin = this.setIsAdmin.bind(this);
  }

  componentDidMount() {
    MulticastGroupStore.get(this.props.match.params.multicastGroupID, resp => {
      this.setState({
        multicastGroup: resp,
      });
    });

    SessionStore.on("change", this.setIsAdmin);
    this.setIsAdmin();
  }

  componentWillUnmount() {
    SessionStore.removeListener("change", this.setIsAdmin);
  }

  componentDidUpdate(oldProps) {
    if (this.props === oldProps) {
      return;
    }
    this.locationToTab();
  }

  setIsAdmin() {
    this.setState({
      admin: SessionStore.isAdmin() || SessionStore.isOrganizationDeviceAdmin(this.props.match.params.organizationID),
    });
  }

  locationToTab() {
    let tab = 0;

    if (window.location.href.endsWith("/edit")) {
      tab = 1;
    }

    this.setState({
      tab: tab,
    });
  }

  onChangeTab(e, v) {
    this.setState({
      tab: v,
    });
  }

  deleteMulticastGroup() {
    if (window.confirm("Are you sure you want to delete this multicast-group?")) {
      MulticastGroupStore.delete(this.props.match.params.multicastGroupID, resp => {
        this.props.history.push(`/organizations/${this.props.match.params.organizationID}/multicast-groups`);
      });
    }
  }

  render() {
    if (this.state.multicastGroup === undefined) {
      return null;
    }

    return(
      <Grid container spacing={4}>
        <TitleBar
          buttons={
            <DeviceAdmin organizationID={this.props.match.params.organizationID}>
              <TitleBarButton
                label={i18n.t(`${packageNS}:tr000061`)}
                icon={<Delete />}
                onClick={this.deleteMulticastGroup}
              />
            </DeviceAdmin>
          }
        >
          <TitleBarTitle to={`/organizations/${this.props.match.params.organizationID}/multicast-groups`} title={i18n.t(`${packageNS}:tr000083`)} />
          <TitleBarTitle title="/" />
          <TitleBarTitle title={this.state.multicastGroup.multicastGroup.name} />
        </TitleBar>

        <Grid item xs={12}>
          <Tabs
            value={this.state.tab}
            onChange={this.onChangeTab}
            indicatorColor="primary"
            className={this.props.classes.tabs}
          >
            <Tab label={i18n.t(`${packageNS}:tr000278`)} component={Link} to={`/organizations/${this.props.match.params.organizationID}/multicast-groups/${this.props.match.params.multicastGroupID}`} />
            {this.state.admin && <Tab label={i18n.t(`${packageNS}:tr000298`)} component={Link} to={`/organizations/${this.props.match.params.organizationID}/multicast-groups/${this.props.match.params.multicastGroupID}/edit`} />}
          </Tabs>
        </Grid>

        <Grid item xs={12}>
          <Switch>
            <Route exact path={`${this.props.match.path}/edit`} render={props => <UpdateMulticastGroup multicastGroup={this.state.multicastGroup.multicastGroup} {...props} />} />
            <Route exact path={`${this.props.match.path}`} render={props => <ListMulticastGroupDevices multicastGroup={this.state.multicastGroup.multicastGroup} {...props} />} />
          </Switch>
        </Grid>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(MulticastGroupLayout));
