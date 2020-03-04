import React, { Component } from "react";
import { withRouter } from 'react-router-dom';

import i18n, { packageNS } from '../../i18n';

import OrganzationStore from "../../stores/OrganizationStore";
import Loader from "../../components/Loader";
import OrganizationForm from "./OrganizationForm";


class UpdateOrganization extends Component {
  constructor() {
    super();
    this.state = {
      loading: false
    };

    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit(organization) {
    this.setState({ loading: true });
    OrganzationStore.update(organization, resp => {
      this.setState({ loading: false });
      this.props.history.push("/organizations");
    }, error => { this.setState({ loading: false }) });
  }

  render() {
    return(<div className="position-relative">
      {this.state.loading && <Loader />}
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
