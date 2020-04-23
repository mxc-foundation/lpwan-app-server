import React, { Component } from "react";
import { withRouter } from 'react-router-dom';
import { Card, CardBody } from 'reactstrap';
import i18n, { packageNS } from '../../i18n';
import OrganizationStore from "../../stores/OrganizationStore";
import OrganizationUserForm from "./OrganizationUserForm";




class UpdateOrganizationUser extends Component {
  constructor() {
    super();
    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit = async (organizationUser) => {
    const res = await OrganizationStore.updateUser(organizationUser);
    this.props.history.push(`/organizations/${organizationUser.organizationID}/users`);
  }

  render() {
    return (
      <React.Fragment>
        <Card>
          <CardBody>
            <OrganizationUserForm
              submitLabel={i18n.t(`${packageNS}:tr000066`)}
              object={this.props.organizationUser}
              update={true}
              onSubmit={this.onSubmit}
            />
          </CardBody>
        </Card>
      </React.Fragment>
    );
  }
}

export default withRouter(UpdateOrganizationUser);