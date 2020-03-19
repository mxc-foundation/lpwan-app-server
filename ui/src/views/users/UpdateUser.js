import React, { Component } from "react";
import { withRouter } from 'react-router-dom';

import i18n, { packageNS } from '../../i18n';
import UserStore from "../../stores/UserStore";
import UserForm from "./UserForm";
import User2FA from "./User2FA";
import PasswordReset from "./PasswordReset";


class UpdateUser extends Component {
  onSubmit = (user) => {
    UserStore.update(user, resp => {
      this.props.history.push("/users");
    });
  }

  render() {
    const { loading, user } = this.props;

    return (<React.Fragment>
      <UserForm
        submitLabel={i18n.t(`${packageNS}:tr000066`)}
        loading={loading}
        object={user}
        onSubmit={this.onSubmit}
        update={true}
      />

      <User2FA loading={loading}
        object={user}
        update={true} />

      <PasswordReset loading={loading}
        object={user}
        update={true} />

      </React.Fragment>
    );
  }
}

export default withRouter(UpdateUser);
