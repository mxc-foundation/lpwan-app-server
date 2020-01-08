import React, { Component } from "react";
import { withRouter } from 'react-router-dom';

import { withStyles } from "@material-ui/core/styles";
import Grid from '@material-ui/core/Grid';

import i18n, { packageNS } from '../../i18n';
import DeviceProfileStore from "../../stores/DeviceProfileStore";
import DeviceProfileForm from "./DeviceProfileForm";

// FIXME - this isn't being used, and we can also remove `withStyles` here
const styles = {
  card: {
    overflow: "visible",
  },
};


class UpdateDeviceProfile extends Component {
  constructor() {
    super();
    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit(deviceProfile) {
    DeviceProfileStore.update(deviceProfile, resp => {
      this.props.history.push(`/organizations/${this.props.match.params.organizationID}/device-profiles`);
    });
  }

  render() {
    return(
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <DeviceProfileForm
            submitLabel={i18n.t(`${packageNS}:tr000066`)}
            object={this.props.deviceProfile}
            onSubmit={this.onSubmit}
            match={this.props.match}
            disabled={!this.props.admin}
            update={true}
          />
        </Grid>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(UpdateDeviceProfile));
