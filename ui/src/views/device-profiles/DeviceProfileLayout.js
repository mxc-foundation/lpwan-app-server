import React, { Component } from "react";
import { Link, withRouter } from "react-router-dom";

import { Breadcrumb, BreadcrumbItem } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";

import Grid from '@material-ui/core/Grid';

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";
import DeviceProfileStore from "../../stores/DeviceProfileStore";
import SessionStore from "../../stores/SessionStore";
import OrganizationDevices from "../devices/OrganizationDevices";
import UpdateDeviceProfile from "./UpdateDeviceProfile";

import breadcrumbStyles from "../common/BreadcrumbStyles";

const localStyles = {};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

class DeviceProfileLayout extends Component {
  constructor() {
    super();
    this.state = {
      admin: false,
    };
  }

  componentDidMount() {
    DeviceProfileStore.get(this.props.match.params.deviceProfileID, resp => {
      this.setState({
        deviceProfile: resp,
      });
    });

    SessionStore.on("change", this.setIsAdmin);
    this.setIsAdmin();
  }

  componentWillUpdate() {
    SessionStore.removeListener("change", this.setIsAdmin);
  }

  setIsAdmin = () => {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    this.setState({
      admin: SessionStore.isAdmin() || SessionStore.isOrganizationDeviceAdmin(currentOrgID),
    });
  }

  deleteDeviceProfile = () => {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    if (window.confirm("Are you sure you want to delete this device-profile?")) {
      DeviceProfileStore.delete(this.props.match.params.deviceProfileID, resp => {
        this.props.history.push(`/organizations/${currentOrgID}/device-profiles`);
      });
    }
  }

  render() {
    const { classes } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    if (this.state.deviceProfile === undefined) {
      return(<div></div>);
    }

    let buttons = [];
    if (this.state.admin) {
      buttons = [
        <TitleBarButton
          key={1}
          label={i18n.t(`${packageNS}:tr000061`)}
          icon={<i className="mdi mdi-delete mr-1 align-middle"></i>}
          color="danger"
          onClick={this.deleteDeviceProfile}
        />,
      ];
    }

    return(
      <Grid container spacing={4}>
        <OrganizationDevices
          mainTabIndex={2}
          organizationID={currentOrgID}
        >
          <TitleBar
            buttons={buttons}
          >
            <Breadcrumb className={classes.breadcrumb}>
              <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                `/organizations/${currentOrgID}/device-profiles`
              }>{i18n.t(`${packageNS}:tr000070`)}</Link></BreadcrumbItem>
              <BreadcrumbItem active>{this.state.deviceProfile.deviceProfile.name}</BreadcrumbItem>
            </Breadcrumb>
          </TitleBar>

          <Grid item xs={12}>
            <UpdateDeviceProfile deviceProfile={this.state.deviceProfile.deviceProfile} admin={this.state.admin} />
          </Grid>
        </OrganizationDevices>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(DeviceProfileLayout));
