import React, { Component } from "react";
import { withRouter } from "react-router-dom";

import Grid from '@material-ui/core/Grid';

import Delete from "mdi-material-ui/Delete";
import KeyVariant from "mdi-material-ui/KeyVariant";

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import TitleBarButton from "../../components/TitleBarButton";
import UserStore from "../../stores/UserStore";
import UpdateUser from "./UpdateUser";

class UserLayout extends Component {
  constructor() {
    super();
    this.state = {
    };

    this.deleteUser = this.deleteUser.bind(this);
  }

  componentDidMount() {
    UserStore.get(this.props.match.params.userID, resp => {
      this.setState({
        user: resp,
      });
    });
  }

  deleteUser() {
    if (window.confirm("Are you sure you want to delete this user?")) {
      UserStore.delete(this.props.match.params.userID, () => {
        this.props.history.push("/users");
      });
    }
  }

  render() {
    if (this.state.user === undefined) {
      return(<div></div>);
    }
    const isDisabled = (this.state.user.user.username === process.env.REACT_APP_DEMO_USER)
                        ?true
                        :false; 
                        
    return(
      <Grid container spacing={4}>
        <TitleBar
          buttons={[
            <TitleBarButton
              key={1}
              label={i18n.t(`${packageNS}:tr000038`)}
              icon={<KeyVariant />}
              to={`/users/${this.props.match.params.userID}/password`}
              disabled={isDisabled}
            />,
            <TitleBarButton
              key={2}
              label={i18n.t(`${packageNS}:tr000061`)}
              icon={<Delete />}
              onClick={this.deleteUser}
            />,
          ]}
        >
          <TitleBarTitle to="/users" title={i18n.t(`${packageNS}:tr000036`)} />
          <TitleBarTitle title="/" />
          <TitleBarTitle title={this.state.user.user.username} />
        </TitleBar>

        <Grid item xs={12}>
          <UpdateUser user={this.state.user.user} />
        </Grid>
      </Grid>
    );
  }
}

export default withRouter(UserLayout);
