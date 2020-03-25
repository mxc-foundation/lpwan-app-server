import Grid from '@material-ui/core/Grid';
import { withStyles } from "@material-ui/core/styles";
import React, { Component } from "react";
import { Link, withRouter } from 'react-router-dom';
import { Breadcrumb, BreadcrumbItem } from 'reactstrap';
import TitleBar from "../../components/TitleBar";
import i18n, { packageNS } from '../../i18n';
import UserStore from "../../stores/UserStore";
import breadcrumbStyles from "../common/BreadcrumbStyles";
import UserForm from "./UserForm";




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
              <BreadcrumbItem className={classes.breadcrumbItem}>Control Panel</BreadcrumbItem>
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
