import React, { Component } from "react";
import { Route, Switch, Link, withRouter } from "react-router-dom";

import { withStyles } from "@material-ui/core/styles";
import Grid from '@material-ui/core/Grid';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import TitleBarButton from "../../components/TitleBarButton";
import Admin from "../../components/Admin";

import ApplicationStore from "../../stores/ApplicationStore";
import SessionStore from "../../stores/SessionStore";
import ListDevices from "../devices/ListDevices";
import UpdateApplication from "./UpdateApplication";
import ListIntegrations from "./ListIntegrations";
import CreateIntegration from "./CreateIntegration";
import UpdateIntegration from "./UpdateIntegration";
import ListFUOTADeploymentsForApplication from "../fuota/ListFUOTADeploymentsForApplication";
import OrganizationDevices from "../devices/OrganizationDevices";
import ApplicationDevices from "./ApplicationDevices";

import theme from "../../theme";


const styles = {
  tabs: {
    borderBottom: "1px solid " + theme.palette.divider,
    height: "48px",
    overflow: "visible",
  },
};


class ApplicationLayout extends Component {
  constructor() {
    super();
    this.state = {
      tab: 0,
      admin: false,
    };
  }

  componentDidMount() {
    const currentApplicationID = this.props.applicationID || this.props.match.params.applicationID;

    ApplicationStore.get(currentApplicationID, resp => {
      this.setState({
        application: resp,
      });
    });

    SessionStore.on("change", this.setIsAdmin);

    this.setIsAdmin();
    this.getMainTabAppIndexFromLocation();
  }

  componentWillUnmount() {
    SessionStore.removeListener("change", this.setIsAdmin);
  }

  componentDidUpdate(oldProps) {
    if (this.props === oldProps) {
      return;
    }

    this.getMainTabAppIndexFromLocation();
  }

  setIsAdmin = () => {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    this.setState({
      admin: SessionStore.isAdmin() || SessionStore.isOrganizationAdmin(currentOrgID),
    });
  }

  deleteApplication = () => {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    if (window.confirm("Are you sure you want to delete this application? This will also delete all devices part of this application.")) {
      ApplicationStore.delete(currentOrgID, resp => {
        this.props.history.push(`/organizations/${currentOrgID}/applications`);
      });
    }
  }

  getMainTabAppIndexFromLocation() {
    let tab = 0; // Devices

    if (window.location.href.search("/edit") !== -1) {
      tab = 1;
    } else if (window.location.href.search("/integrations") !== -1) {
      tab = 2;
    } else if (window.location.href.search("/fuota-deployments") !== -1) {
      tab = 3;
    }

    this.setState({
      tab: tab,
    });
  }

  render() {
    const { admin, application, tab } = this.state;
    const { children } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    if (application === undefined) {
      return(<div></div>);
    }

    return(
      <Grid container spacing={4}>
        {/* <OrganizationDevices
          mainTabIndex={1}
          organizationID={currentOrgID}
        > */}
          <ApplicationDevices
            {...this.props}
            admin={admin}
            application={application}
            deleteApplication={this.deleteApplication}
            mainTabAppIndex={tab}
            organizationID={currentOrgID}
          >
            {children}
            <Switch>
              <Route exact path={`${this.props.match.path}`} render={props =>
                <ListDevices application={this.state.application.application} {...props} />} />
              <Route exact path={`${this.props.match.path}/edit`} render={props =>
                <UpdateApplication application={this.state.application.application} {...props} />} />
              <Route exact path={`${this.props.match.path}/integrations/create`} render={props =>
                <CreateIntegration application={this.state.application.application} {...props} /> } />
              <Route exact path={`${this.props.match.path}/integrations/:kind`} render={props =>
                <UpdateIntegration application={this.state.application.application} {...props} /> } />
              <Route exact path={`${this.props.match.path}/integrations`} render={props =>
                <ListIntegrations application={this.state.application.application} {...props} />} />
              <Route exact path={`${this.props.match.path}/fuota-deployments`} render={props =>
                <ListFUOTADeploymentsForApplication application={this.state.application.application} {...props} /> } />
            </Switch>
          </ApplicationDevices>
        {/* </OrganizationDevices> */}
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(ApplicationLayout));
