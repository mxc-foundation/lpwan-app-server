import Grid from '@material-ui/core/Grid';
import { withStyles } from "@material-ui/core/styles";
import React, { Component } from "react";
import { Route, Switch, withRouter } from "react-router-dom";
import ApplicationStore from "../../stores/ApplicationStore";
import DeviceProfileStore from "../../stores/DeviceProfileStore";
import DeviceStore from "../../stores/DeviceStore";
import OrganizationStore from "../../stores/OrganizationStore";
import SessionStore from "../../stores/SessionStore";
import theme from "../../theme";
import DeviceDetailsDevicesTabs from "../../views/applications/DeviceDetailsDevicesTabs";
import ListFUOTADeploymentsForDevice from "../../views/fuota/ListFUOTADeploymentsForDevice";
import DeviceActivation from "./DeviceActivation";
import DeviceData from "./DeviceData";
import DeviceDetails from "./DeviceDetails";
import DeviceFrames from "./DeviceFrames";
import DeviceKeys from "./DeviceKeys";
import UpdateDevice from "./UpdateDevice";







const styles = {
  tabs: {
    borderBottom: "1px solid " + theme.palette.divider,
    height: "49px",
  },
};


class DeviceLayout extends Component {
  constructor() {
    super();
    this.state = {
      tab: 0,
      admin: false,
    };
  }

  componentDidMount() {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const currentApplicationID = this.props.applicationID || this.props.match.params.applicationID;

    if (currentApplicationID) {
      ApplicationStore.get(this.props.match.params.applicationID, resp => {
        this.setState({
          application: resp,
        });
      });
    }

    this.loadOrganization(currentOrgID);

    DeviceStore.on("update", this.getDevice);
    SessionStore.on("change", this.setIsAdmin);

    this.getMainTabDeviceIndexFromLocation();
    this.setIsAdmin();
    this.getDevice();
  }

  loadOrganization = async (id) => {
    let resp = await OrganizationStore.get(id);  
    this.setState({
      organization: resp.organization,
      loading: false
    });
  }

  componentWillUnmount() {
    SessionStore.removeListener("change", this.setIsAdmin);
    DeviceStore.removeListener("update", this.getDevice);
  }

  componentDidUpdate(oldProps) {
    if (this.props === oldProps) {
      return;
    }

    this.getMainTabDeviceIndexFromLocation();
  }

  getDevice = () => {
    DeviceStore.get(this.props.match.params.devEUI, resp => {
      this.setState({
        device: resp,
      });

      if (resp.device.deviceProfileID) {
        DeviceProfileStore.get(resp.device.deviceProfileID, resp => {
          this.setState({
            deviceProfile: resp,
          });
        });
      }

    });
  }

  setIsAdmin = () => {
    this.setState({
      admin: SessionStore.isAdmin() || SessionStore.isOrganizationDeviceAdmin(this.props.match.params.organizationID),
    }, () => {
      // we need to update the tab index, as for non-admins, some tabs are hidden
      this.getMainTabDeviceIndexFromLocation();
    });
  }

  onChangeTab = (e, v) => {
    this.setState({
      tab: v,
    });
  }

  getMainTabDeviceIndexFromLocation = () => {
    let tab = 0;

    if (window.location.href.endsWith("/edit")) {
      tab = 1;
    } else if (window.location.href.endsWith("/keys")) {
      tab = 2;
    } else if (window.location.href.endsWith("/activation")) {
      tab = 3;
    } else if (window.location.href.endsWith("/data")) {
      tab = 4;
    } else if (window.location.href.endsWith("/frames")) {
      tab = 5;
    } else if (window.location.href.endsWith("/fuota-deployments")) {
      tab = 6;
    }

    if (tab > 1 && !this.state.admin) {
      tab = tab - 1;
    }

    this.setState({
      tab: tab,
    });
  }

  deleteDevice = () => {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const currentApplicationID = this.props.applicationID || this.props.match.params.applicationID;
    const isApplication = currentApplicationID && currentApplicationID !== "0"; 

    if (window.confirm("Are you sure you want to delete this device?")) {
      DeviceStore.delete(this.props.match.params.devEUI, resp => {
        isApplication
        ? this.props.history.push(`/organizations/${currentOrgID}/applications/${currentApplicationID}`)
        : this.props.history.push(`/organizations/${currentOrgID}`);
      });
    }
  }

  render() {
    const { admin, application, device, deviceProfile, organization, tab } = this.state;
    const { children, match } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    // Note: Must prefix with `\\d+` and `\\w` here instead of `\d+` and `\w` as we've done in App.js
    const urlPrefixDeviceNoApp = `/organizations/:organizationID(\\d+)/devices/:devEUI([\\w]{16})`;

    // if (application === undefined || device === undefined|| deviceProfile === undefined) {
    if (device === undefined) {
      return(<div></div>);
    }

    return(
      <Grid container spacing={4}>
        <DeviceDetailsDevicesTabs
          {...this.props}
          admin={admin}
          application={application}
          deleteDevice={this.deleteDevice}
          device={device}
          deviceProfile={deviceProfile}
          mainTabDeviceIndex={tab}
          organizationID={currentOrgID}
          organization={organization}
        >
          {children}
          <Switch>
            <Route exact path={`${urlPrefixDeviceNoApp}/edit`} render={props => <UpdateDevice device={device.device} admin={admin} {...props} />} />
            <Route exact path={`${urlPrefixDeviceNoApp}/keys`} render={props => <DeviceKeys devEUI={this.props.match.params.devEUI} device={device.device} admin={admin} deviceProfile={deviceProfile && deviceProfile.deviceProfile} {...props} />} />
            <Route exact path={`${urlPrefixDeviceNoApp}/activation`} render={props => <DeviceActivation device={device.device} admin={admin} deviceProfile={deviceProfile && deviceProfile.deviceProfile} {...props} />} />
            <Route exact path={`${urlPrefixDeviceNoApp}/data`} render={props => <DeviceData device={device.device} admin={admin} {...props} />} />
            <Route exact path={`${urlPrefixDeviceNoApp}/frames`} render={props => <DeviceFrames device={device.device} admin={admin} {...props} />} />
            <Route exact path={`${urlPrefixDeviceNoApp}/fuota-deployments`} render={props => <ListFUOTADeploymentsForDevice device={device.device} admin={admin} {...props} /> } />
            <Route exact path={`${urlPrefixDeviceNoApp}`} render={props => <DeviceDetails device={device} deviceProfile={deviceProfile} admin={admin} {...props} />} />
            
            <Route exact path={`${match.path}/edit`} render={props => <UpdateDevice device={device.device} admin={admin} {...props} />} />
            <Route exact path={`${match.path}/keys`} render={props => <DeviceKeys device={device.device} admin={admin} deviceProfile={deviceProfile && deviceProfile.deviceProfile} {...props} />} />
            <Route exact path={`${match.path}/activation`} render={props => <DeviceActivation device={device.device} admin={admin} deviceProfile={deviceProfile && deviceProfile.deviceProfile} {...props} />} />
            <Route exact path={`${match.path}/data`} render={props => <DeviceData device={device.device} admin={admin} {...props} />} />
            <Route exact path={`${match.path}/frames`} render={props => <DeviceFrames device={device.device} admin={admin} {...props} />} />
            <Route exact path={`${match.path}/fuota-deployments`} render={props => <ListFUOTADeploymentsForDevice device={device.device} admin={admin} {...props} /> } />
            <Route exact path={`${match.path}`} render={props => <DeviceDetails device={device} deviceProfile={deviceProfile} admin={admin} {...props} />} />
          </Switch>
        </DeviceDetailsDevicesTabs>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(DeviceLayout));
