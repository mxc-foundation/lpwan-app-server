import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';

import Admin from "../../components/Admin";
import { Breadcrumb, BreadcrumbItem } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";
import Grid from '@material-ui/core/Grid';

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import UserForm from "./UserForm";
import UserStore from "../../stores/UserStore";

import breadcrumbStyles from "../common/BreadcrumbStyles";

const localStyles = {};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};


class CreateUser extends Component {
  onSubmit = (newUserObject) => {
    UserStore.create(newUserObject, resp => {
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
              <Admin>
                <BreadcrumbItem className={classes.breadcrumbItem}>Control Panel</BreadcrumbItem>
              </Admin>
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
