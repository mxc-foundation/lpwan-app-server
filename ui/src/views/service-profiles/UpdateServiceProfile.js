import React, { Component } from "react";
import { withRouter } from 'react-router-dom';

import i18n, { packageNS } from '../../i18n';
import ServiceProfileStore from "../../stores/ServiceProfileStore";
import ServiceProfileForm from "./ServiceProfileForm";


class UpdateServiceProfile extends Component {
  constructor(props) {
    super(props);
    this.state = {};

    this.onSubmit = this.onSubmit.bind(this);

  }

  onSubmit(serviceProfile) {
    ServiceProfileStore.update(serviceProfile, resp => {
      this.props.history.push(`/organizations/${this.props.match.params.organizationID}/service-profiles`);
    });
  }

  render() {
    return (<div className="position-relative">
      <ServiceProfileForm
        submitLabel={i18n.t(`${packageNS}:tr000066`)}
        object={this.props.serviceProfile}
        onSubmit={this.onSubmit}
        update={true}
        match={this.props.match}
        disabled={!this.props.admin}
      ></ServiceProfileForm>
    </div>
    );
  }
}

export default withRouter(UpdateServiceProfile);
