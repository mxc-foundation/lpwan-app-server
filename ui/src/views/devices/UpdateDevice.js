import React, { Component } from "react";
import { withRouter } from 'react-router-dom';

import { withStyles } from "@material-ui/core/styles";

import i18n, { packageNS } from '../../i18n';
import DeviceStore from "../../stores/DeviceStore";
import DeviceForm from "./DeviceForm";


const styles = {
  card: {
    overflow: "visible",
  },
};


class UpdateDevice extends Component {
  constructor() {
    super();
    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit(device) {
    const currentApplicationID = this.props.applicationID || this.props.match.params.applicationID;
    const isApplication = currentApplicationID && currentApplicationID !== "0"; 

    DeviceStore.update(device, resp => {
      isApplication
      ? this.props.history.push(`/organizations/${this.props.match.params.organizationID}/applications/${currentApplicationID}/devices/${this.props.match.params.devEUI}`)
      : this.props.history.push(`/organizations/${this.props.match.params.organizationID}/devices/${this.props.match.params.devEUI}`);
    });
  }

  render() {
    return(
      <DeviceForm
        submitLabel={i18n.t(`${packageNS}:tr000066`)}
        object={this.props.device}
        onSubmit={this.onSubmit}
        match={this.props.match}
        update={true}
        disabled={!this.props.admin}
      />
    );
  }
}

export default withStyles(styles)(withRouter(UpdateDevice));
