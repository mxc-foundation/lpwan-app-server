import React, { Component } from "react";
import { withRouter } from 'react-router-dom';

import i18n, { packageNS } from '../../i18n';
import GatewayStore from "../../stores/GatewayStore";
import GatewayForm from "./GatewayForm";


class UpdateGateway extends Component {
  constructor() {
    super();
    this.onSubmit = this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit(gateway) {
    GatewayStore.update(gateway, resp => {
      this.props.history.push(`/organizations/${this.props.match.params.organizationID}/gateways`);
    });
  }

  render() {
    return (<GatewayForm
        submitLabel={i18n.t(`${packageNS}:tr000066`)}
        object={this.props.gateway}
        onSubmit={this.onSubmit}
        update={true}
        match={this.props.match}
      ></GatewayForm>
    );
  }
}

export default withRouter(UpdateGateway);
