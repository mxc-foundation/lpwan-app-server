import Grid from '@material-ui/core/Grid';
import { withStyles } from "@material-ui/core/styles";
import React, { Component } from "react";
import { Route, Switch, withRouter } from "react-router-dom";
import ApplicationStore from "../../stores/ApplicationStore";
import SessionStore from "../../stores/SessionStore";
import theme from "../../theme";
import ListDevices from "../devices/ListDevices";
import ListFUOTADeploymentsForApplication from "../fuota/ListFUOTADeploymentsForApplication";
import ApplicationDevices from "./ApplicationDevices";
import CreateIntegration from "./CreateIntegration";
import ListIntegrations from "./ListIntegrations";
import UpdateApplication from "./UpdateApplication";
import UpdateIntegration from "./UpdateIntegration";





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
      openModal: false
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

  openModal = () => {
    this.setState({openModal:true});
  }

  deleteApplication = (applicationId) => {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    ApplicationStore.delete(applicationId, resp => {
      this.props.history.push(`/organizations/${currentOrgID}/applications`);
    });
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
