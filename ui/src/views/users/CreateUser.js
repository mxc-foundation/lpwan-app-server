import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';

import { Breadcrumb, BreadcrumbItem } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";
import Grid from '@material-ui/core/Grid';

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import UserForm from "./UserForm";
import UserStore from "../../stores/UserStore";

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


class CreateUser extends Component {
  onSubmit = (user) => {
    console.log('user', user);
    return false;
    UserStore.create(user, user.password, [], resp => {
      this.props.history.push("/users");
    });
  }

  render() {
    const { classes } = this.props;

    return (
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <TitleBar noButtons>
            <Breadcrumb className={classes.breadcrumb}>
              <BreadcrumbItem>
                <Link
                  className={classes.breadcrumbItemLink}
                  to={`/users`}>
                  {i18n.t(`${packageNS}:tr000036`)}
                </Link>
              </BreadcrumbItem>
              <BreadcrumbItem active>Create Profile</BreadcrumbItem>
            </Breadcrumb>
          </TitleBar>
        
          <UserForm
            submitLabel={i18n.t(`${packageNS}:tr000277`)}
            onSubmit={this.onSubmit}
          />
        </Grid>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(CreateUser));
