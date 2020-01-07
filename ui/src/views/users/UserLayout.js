import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";

import Delete from "mdi-material-ui/Delete";
import KeyVariant from "mdi-material-ui/KeyVariant";
import { Breadcrumb, BreadcrumbItem, Row, Button } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import UserStore from "../../stores/UserStore";
import UpdateUser from "./UpdateUser";

const styles = theme => ({
  [theme.breakpoints.down('sm')]: {
    breadcrumb: {
      fontSize: "1.1rem",
      margin: "0rem",
      padding: "0rem"
    },
  },
  [theme.breakpoints.up('sm')]: {
    breadcrumb: {
      fontSize: "1.25rem",
      margin: "0rem",
      padding: "0rem"
    },
  },
  breadcrumbItemLink: {
    color: "#71b6f9 !important"
  }
});

class UserLayout extends Component {
  constructor() {
    super();

    this.state = {
      loading: true,
    };
  }

  componentDidMount() {
    UserStore.get(this.props.match.params.userID, resp => {
      this.setState({
        user: resp.user,
        loading: false
      });
    });
  }

  deleteUser = () => {
    if (window.confirm("Are you sure you want to delete this user?")) {
      UserStore.delete(this.props.match.params.userID, () => {
        this.props.history.push("/users");
      });
    }
  }

  changePassword = () => {
    this.props.history.push(`/users/${this.props.match.params.userID}/password`);
  }

  render() {
    const { loading, user } = this.state;
    const { classes } = this.props;

    if (user === undefined) {
      return (<div></div>);
    }
    const isDisabled = (this.state.user.username === process.env.REACT_APP_DEMO_USER)
      ? true
      : false;

    return (
      <React.Fragment>
        <TitleBar
          buttons={[
            <Button color="danger"
              key={1}
              onClick={this.changePassword}
              disabled={isDisabled}
              className=""><i className="mdi mdi-key-change"></i>{' '}{i18n.t(`${packageNS}:tr000038`)}
            </Button>,
          ]}
        >
          <Breadcrumb className={classes.breadcrumb}>
            <BreadcrumbItem>
              <Link
                className={classes.breadcrumbItemLink}
                to={`/users`}
              >
                  {i18n.t(`${packageNS}:tr000036`)}
              </Link>
            </BreadcrumbItem>
            <BreadcrumbItem>
              <Link
                className={classes.breadcrumbItemLink}
                to={`/users/${user.id}`}
              >
                {user.username}
              </Link>
            </BreadcrumbItem>
            <BreadcrumbItem active>Edit</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>
        <Row xs={12} lg={12}>
          <UpdateUser
            user={user}
            loading={loading}
          />
        </Row>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withRouter(UserLayout));
