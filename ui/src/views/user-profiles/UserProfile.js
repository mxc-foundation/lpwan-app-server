import React, { Component } from 'react';
import { Link, withRouter } from 'react-router-dom';
import { Breadcrumb, BreadcrumbItem, Card, UncontrolledDropdown, DropdownMenu, DropdownItem, DropdownToggle } from 'reactstrap';

import SessionStore from '../../stores/SessionStore';
import i18n, { packageNS } from '../../i18n';
import TitleBar from '../../components/TitleBar';
import UserProfileForm from './UserProfileForm';

class UserProfile extends Component {
  constructor() {
    super();

    this.state = {};
  }

  componentDidMount() {
    const orgId = SessionStore.getOrganizationID();
    const user = SessionStore.getUser();
    this.setState({
      orgId,
      user
    });
  }

  onSubmit = (user) => {
    const { orgId } = this.state;

    // TODO - implement functions to update session and user profile in
    // stores and backend
    SessionStore.update(user, resp => {
      this.props.history.push(`/modify-account/${orgId}/users/${user.id}/user-profile`);
    });
  }

  render() {
    if (this.state.user === undefined) {
      return(<div></div>);
    }

    return(
      <React.Fragment>
        <TitleBar>
          <Breadcrumb style={{ fontSize: "1.25rem", margin: "0rem", padding: "0rem" }}>
            <BreadcrumbItem>{i18n.t(`${packageNS}:tr000430`)}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>
        <Card className="card-box shadow-sm" style={{ minWidth: "25rem" }}>
          <UserProfileForm
            onSubmit={this.onSubmit}
            object={this.state.user}
            orgId={this.state.orgId}
            submitLabel={i18n.t(`${packageNS}:tr000066`)}
          />
        </Card>
      </React.Fragment>
    );
  }
}

export default withRouter(UserProfile);
