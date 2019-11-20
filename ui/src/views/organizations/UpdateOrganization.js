import React, { Component } from "react";
import { withRouter } from 'react-router-dom';

import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import CardContent from "@material-ui/core/CardContent";

import i18n, { packageNS } from '../../i18n';
import OrganzationStore from "../../stores/OrganizationStore";
import OrganizationForm from "./OrganizationForm";


class UpdateOrganization extends Component {
  constructor() {
    super();
    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit(organization) {
    OrganzationStore.update(organization, resp => {
      this.props.history.push("/organizations");
    });
  }

  render() {
    return(
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <OrganizationForm
                submitLabel={i18n.t(`${packageNS}:tr000066`)}
                object={this.props.organization.organization}
                onSubmit={this.onSubmit}
              />
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    );
  }
}

export default withRouter(UpdateOrganization);
