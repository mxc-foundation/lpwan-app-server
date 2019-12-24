import React, { Component } from "react";
import { withRouter } from 'react-router-dom';

import i18n, { packageNS } from '../../i18n';

import OrganzationStore from "../../stores/OrganizationStore";
import OrganizationForm from "./OrganizationForm";


class UpdateOrganization extends Component {
  constructor() {
    super();
    this.state = {};

    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit(organization) {
    OrganzationStore.update(organization, resp => {
      this.props.history.push("/organizations");
    });
  }

  render() {
    return(<div className="position-relative">
        <OrganizationForm
            submitLabel={i18n.t(`${packageNS}:tr000066`)}
            object={this.props.organization.organization}
            onSubmit={this.onSubmit}
            update={true}
            match={this.props.match}
        ></OrganizationForm>
      </div>
    );
  }
}

export default withRouter(UpdateOrganization);
