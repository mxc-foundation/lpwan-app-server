import Grid from '@material-ui/core/Grid';
import { withStyles } from "@material-ui/core/styles";
import React, { Component } from "react";
import { Route, Switch } from "react-router-dom";
import ApplicationStore from "../../stores/ApplicationStore";
import FUOTADeploymentStore from "../../stores/FUOTADeploymentStore";
import theme from "../../theme";
import ApplicationFUOTADeploymentTabs from "../applications/ApplicationFUOTADeploymentTabs";
import FUOTADeploymentDetails from "./FUOTADeploymentDetails";
import ListFUOTADeploymentDevices from "./ListFUOTADeploymentDevices";





const styles = {
  tabs: {
    borderBottom: "1px solid " + theme.palette.divider,
    height: "48px",
    overflow: "visible",
  },
};


class FUOTADeploymentLayout extends Component {
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

    FUOTADeploymentStore.on("reload", this.getFuotaDeployment);


    this.getFuotaDeployment();
    this.getMainTabAppIndexFromLocation();
  }

  componentWillUnmount() {
    FUOTADeploymentStore.removeListener("reload", this.getFuotaDeployment);
  }

  getFuotaDeployment = () => {
    FUOTADeploymentStore.get(this.props.match.params.fuotaDeploymentID, resp => {
      this.setState({
        fuotaDeployment: resp
      });
    });
  }

  getMainTabAppIndexFromLocation = () => {
    let tab = 0; // Information

    if (window.location.href.endsWith("/devices")) {
      tab = 1;
    }

    this.setState({
      tab: tab,
    });
  }


  render() {
    const { admin, application, fuotaDeployment, tab } = this.state;
    const { children } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    if (application === undefined || fuotaDeployment === undefined) {
      return null;
    }

    return(
      <Grid container spacing={4}>
        {/* <OrganizationDevices
          mainTabIndex={1}
          organizationID={currentOrgID}
        > */}
          <ApplicationFUOTADeploymentTabs
            {...this.props}
            admin={admin}
            application={application}
            deleteApplication={this.deleteApplication}
            fuotaDeployment={fuotaDeployment}
            mainTabAppIndex={tab}
            organizationID={currentOrgID}
          >
            {children}
            <Switch>
              <Route exact path={`${this.props.match.path}`} render={props =>
                <FUOTADeploymentDetails fuotaDeployment={fuotaDeployment} {...props} />} />
              <Route exact path={`${this.props.match.path}/devices`} render={props =>
                <ListFUOTADeploymentDevices fuotaDeployment={fuotaDeployment} {...props} />} />
            </Switch>
          </ApplicationFUOTADeploymentTabs>
        {/* </OrganizationDevices> */}
      </Grid>
    );
  }
}

export default withStyles(styles)(FUOTADeploymentLayout);

